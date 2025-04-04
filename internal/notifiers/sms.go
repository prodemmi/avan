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
		var payload []byte
		if message.Payload != nil {
			decodedPayload := []byte(*message.Payload) // Convert string back to []byte
			var data map[string]interface{}
			if err := json.Unmarshal(decodedPayload, &data); err != nil {
				return fmt.Errorf("failed to unmarshal payload: %w", err)
			}

			if !inHeader && smsConfig.AuthTokenKey != nil && smsConfig.AuthTokenValue != nil {
				data[*smsConfig.AuthTokenKey] = *smsConfig.AuthTokenValue
			}

			payload, err = json.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}
		} else {
			return errors.New("payload is nil")
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
