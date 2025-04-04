package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"avan/internal/db/sqlc"
)

type WebhookMessage struct {
	Content string `json:"content"`
}

func (n Notifier) HandleWebhookNotifier(provider sqlc.Provider, message sqlc.Message) error {
	// دریافت تنظیمات Webhook از ProviderService
	webhookConfig, err := n.ProviderService.GetWebhookConfig(provider.ID)
	if err != nil {
		return fmt.Errorf("failed to get webhook config: %w", err)
	}

	// ساخت پیام برای ارسال به Webhook
	webhookMessage := WebhookMessage{
		Content: message.Content,
	}

	// تبدیل پیام به فرمت JSON
	payload, err := json.Marshal(webhookMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	// ارسال درخواست HTTP به Webhook
	req, err := http.NewRequest("POST", webhookConfig.Url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ارسال درخواست
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to webhook: %w", err)
	}
	defer resp.Body.Close()

	// بررسی وضعیت پاسخ
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to webhook, status: %d", resp.StatusCode)
	}

	fmt.Println("Message sent to webhook successfully")
	return nil
}
