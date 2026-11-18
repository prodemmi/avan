package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"avan/internal/db/sqlc"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func (n Notifier) HandleSlackNotifier(provider sqlc.Provider, message sqlc.Message) error {
	// دریافت تنظیمات Slack از ProviderService
	slackConfig, err := n.ProviderService.GetSlackConfig(provider.ID)
	if err != nil {
		return fmt.Errorf("failed to get Slack config: %w", err)
	}

	// ساخت پیام برای ارسال به Slack
	slackMessage := SlackMessage{
		Text: message.Content,
	}

	// تبدیل پیام به فرمت JSON
	payload, err := json.Marshal(slackMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	// ارسال درخواست HTTP به Webhook Slack
	req, err := http.NewRequest("POST", slackConfig.WebhookUrl, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ارسال درخواست
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to Slack: %w", err)
	}
	defer resp.Body.Close()

	// بررسی وضعیت پاسخ
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack, status: %d", resp.StatusCode)
	}

	fmt.Println("Message sent to Slack successfully")
	return nil
}
