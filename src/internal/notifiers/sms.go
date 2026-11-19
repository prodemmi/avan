package notifiers

import (
	"avan/internal/db/sqlc"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	MethodGET  int64 = 0
	MethodPOST int64 = 1
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Replace placeholders in URL or body with dynamic data
func replacePlaceholders(template string, replacements map[string]string) string {
	for key, value := range replacements {
		template = strings.ReplaceAll(template, fmt.Sprintf("$%s", key), value)
	}
	return template
}

func buildSmsPostBody(bodySchema string, message sqlc.Message, authKey, authValue *string, includeAuth bool) ([]byte, error) {
	var payload any
	if message.Payload != nil {
		if err := json.Unmarshal([]byte(*message.Payload), &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}
	}
	var body any
	if bodySchema == "" {
		if payload == nil {
			return nil, errors.New("payload is nil")
		}
		body = payload
	} else {
		if err := json.Unmarshal([]byte(bodySchema), &body); err != nil {
			return nil, fmt.Errorf("invalid SMS POST body schema: %w", err)
		}
		var replace func(any) any
		replace = func(value any) any {
			switch value := value.(type) {
			case map[string]any:
				for key, nested := range value {
					value[key] = replace(nested)
				}
				return value
			case []any:
				for i, nested := range value {
					value[i] = replace(nested)
				}
				return value
			case string:
				if value == "$payload" {
					return payload
				}
				if value == "$retry" {
					return message.Retry
				}
				if value == "$provider_id" {
					return message.ProviderID
				}
				if value == "$workspace_id" {
					return message.WorkspaceID
				}
				if value == "$template_id" {
					return message.TemplateID
				}
				if value == "$schedule_at" {
					if message.ScheduleAt == nil {
						return nil
					}
					return message.ScheduleAt.Format(time.RFC3339)
				}
				value = strings.ReplaceAll(value, "$message", message.Content)
				value = strings.ReplaceAll(value, "$send_to", message.Recipient)
				value = strings.ReplaceAll(value, "$retry", strconv.FormatBool(message.Retry))
				value = strings.ReplaceAll(value, "$provider_id", strconv.FormatInt(message.ProviderID, 10))
				value = strings.ReplaceAll(value, "$workspace_id", strconv.FormatInt(message.WorkspaceID, 10))
				if message.TemplateID != nil {
					value = strings.ReplaceAll(value, "$template_id", strconv.FormatInt(*message.TemplateID, 10))
				}
				if message.ScheduleAt != nil {
					value = strings.ReplaceAll(value, "$schedule_at", message.ScheduleAt.Format(time.RFC3339))
				}
				return value
			default:
				return value
			}
		}
		body = replace(body)
	}
	object, ok := body.(map[string]any)
	if !ok {
		return nil, errors.New("SMS POST body schema must be a JSON object")
	}
	if includeAuth && authKey != nil && authValue != nil {
		object[*authKey] = *authValue
	}
	encoded, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return encoded, nil
}

// HandleSmsNotifier processes and sends SMS notifications
func (n Notifier) HandleSmsNotifier(provider sqlc.Provider, message sqlc.Message) error {
	smsConfig, err := n.ProviderService.GetSmsConfig(provider.ID)
	if err != nil {
		return fmt.Errorf("failed to get sms config: %w", err)
	}
	inHeader := smsConfig.IsAuthInHeader != nil && *smsConfig.IsAuthInHeader

	var req *http.Request
	if smsConfig.Method == MethodGET {
		replaces := map[string]string{
			"message": message.Content,
			"send_to": message.Recipient,
		}

		encodedUrl := replacePlaceholders(smsConfig.Url, replaces)

		if !inHeader && smsConfig.AuthTokenKey != nil && smsConfig.AuthTokenValue != nil {
			encodedUrl = replacePlaceholders(encodedUrl, map[string]string{*smsConfig.AuthTokenKey: *smsConfig.AuthTokenValue})

			if encodedUrl == smsConfig.Url {
				u, err := url.Parse(encodedUrl)
				if err != nil {
					return fmt.Errorf("failed to parse URL: %w", err)
				}
				query := u.Query()
				query.Set(*smsConfig.AuthTokenKey, *smsConfig.AuthTokenValue)
				u.RawQuery = query.Encode()
				encodedUrl = u.String()
			}
		}

		req, err = http.NewRequest("GET", encodedUrl, nil)
		if err != nil {
			return fmt.Errorf("failed to create GET request: %w", err)
		}
	} else if smsConfig.Method == MethodPOST {
		bodySchema := ""
		if smsConfig.BodySchema != nil {
			bodySchema = *smsConfig.BodySchema
		}
		payload, err := buildSmsPostBody(bodySchema, message, smsConfig.AuthTokenKey, smsConfig.AuthTokenValue, !inHeader)
		if err != nil {
			return err
		}

		req, err = http.NewRequest("POST", smsConfig.Url, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to create POST request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		return errors.New("unsupported method type")
	}

	if inHeader && smsConfig.AuthTokenKey != nil && smsConfig.AuthTokenValue != nil {
		req.Header.Set(*smsConfig.AuthTokenKey, *smsConfig.AuthTokenValue)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS, status: %d, response: %s", resp.StatusCode, string(body))
	}

	fmt.Println("SMS sent successfully")
	return nil
}
