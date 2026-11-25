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

func buildSmsPostFields(bodySchema string, message sqlc.Message, authKey, authValue *string, includeAuth bool, replacements map[string]string) (map[string]any, error) {
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
				for key, replacement := range replacements {
					value = strings.ReplaceAll(value, "$"+key, replacement)
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
	return object, nil
}

func buildSmsPostBody(bodySchema string, message sqlc.Message, authKey, authValue *string, includeAuth bool) ([]byte, error) {
	object, err := buildSmsPostFields(bodySchema, message, authKey, authValue, includeAuth, nil)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return encoded, nil
}

func buildSmsRequestBody(bodySchema string, message sqlc.Message, authKey, authValue *string, includeAuth bool, replacements map[string]string, encoding string) ([]byte, string, error) {
	object, err := buildSmsPostFields(bodySchema, message, authKey, authValue, includeAuth, replacements)
	if err != nil {
		return nil, "", err
	}
	if encoding == "" || encoding == "json" {
		body, err := json.Marshal(object)
		return body, "application/json", err
	}
	if encoding != "form" {
		return nil, "", errors.New("unsupported SMS POST body encoding")
	}
	values := make(url.Values, len(object))
	for key, value := range object {
		if value == nil {
			continue
		}
		if text, ok := value.(string); ok {
			values.Set(key, text)
			continue
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode form field %q: %w", key, err)
		}
		values.Set(key, string(encoded))
	}
	return []byte(values.Encode()), "application/x-www-form-urlencoded", nil
}

func buildSmsRequest(config sqlc.SmsConfig, message sqlc.Message) (*http.Request, error) {
	replacements := map[string]string{"message": message.Content, "send_to": message.Recipient}
	if config.AuthUsername != nil {
		replacements["username"] = *config.AuthUsername
	}
	if config.AccountSid != nil {
		replacements["account_sid"] = url.PathEscape(*config.AccountSid)
	}
	if config.Sender != nil {
		replacements["sender"] = *config.Sender
	}
	requestURL := replacePlaceholders(config.Url, replacements)
	authType := config.AuthType
	if authType == "" {
		authType = "token"
	}
	if authType != "token" && authType != "basic" {
		return nil, errors.New("unsupported SMS authentication type")
	}
	inHeader := config.IsAuthInHeader != nil && *config.IsAuthInHeader
	var req *http.Request
	var err error
	switch config.Method {
	case MethodGET:
		if authType == "token" && !inHeader && config.AuthTokenKey != nil && config.AuthTokenValue != nil {
			key, value := *config.AuthTokenKey, *config.AuthTokenValue
			if strings.Contains(config.Url, "$"+key) {
				requestURL = replacePlaceholders(requestURL, map[string]string{key: value})
			} else {
				u, parseErr := url.Parse(requestURL)
				if parseErr != nil {
					return nil, fmt.Errorf("failed to parse URL: %w", parseErr)
				}
				query := u.Query()
				query.Set(key, value)
				u.RawQuery = query.Encode()
				requestURL = u.String()
			}
		}
		req, err = http.NewRequest(http.MethodGet, requestURL, nil)
	case MethodPOST:
		encoding := config.BodyEncoding
		if encoding == "" {
			encoding = "json"
		}
		bodySchema := ""
		if config.BodySchema != nil {
			bodySchema = *config.BodySchema
		}
		includeAuth := authType == "token" && !inHeader
		var body []byte
		var contentType string
		body, contentType, err = buildSmsRequestBody(bodySchema, message, config.AuthTokenKey, config.AuthTokenValue, includeAuth, replacements, encoding)
		if err == nil {
			req, err = http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
			if err == nil {
				req.Header.Set("Content-Type", contentType)
			}
		}
	default:
		return nil, errors.New("unsupported SMS method type")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to build SMS request: %w", err)
	}
	if authType == "basic" {
		if config.AuthUsername == nil || config.AuthTokenValue == nil || *config.AuthUsername == "" || *config.AuthTokenValue == "" {
			return nil, errors.New("HTTP Basic username and password are required")
		}
		req.SetBasicAuth(*config.AuthUsername, *config.AuthTokenValue)
	} else if inHeader && config.AuthTokenKey != nil && config.AuthTokenValue != nil {
		req.Header.Set(*config.AuthTokenKey, *config.AuthTokenValue)
	}
	return req, nil
}

func smsResponseError(requestURL string, statusCode int, body []byte) error {
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to send SMS, status: %d, response: %s", statusCode, body)
	}
	u, err := url.Parse(requestURL)
	if err != nil || !strings.EqualFold(u.Hostname(), "rest.payamak-panel.com") || u.Path != "/api/SendSMS/SendSMS" {
		return nil
	}
	returnValue, err := strconv.Atoi(strings.Trim(strings.TrimSpace(string(body)), `"`))
	if err != nil || returnValue <= 0 {
		return fmt.Errorf("Melipayamak rejected SMS, response: %s", body)
	}
	return nil
}

// HandleSmsNotifier processes and sends SMS notifications
func (n Notifier) HandleSmsNotifier(provider sqlc.Provider, message sqlc.Message) error {
	smsConfig, err := n.ProviderService.GetSmsConfig(provider.ID)
	if err != nil {
		return fmt.Errorf("failed to get sms config: %w", err)
	}
	req, err := buildSmsRequest(smsConfig, message)
	if err != nil {
		return err
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

	if err := smsResponseError(req.URL.String(), resp.StatusCode, body); err != nil {
		return err
	}

	fmt.Println("SMS sent successfully")
	return nil
}
