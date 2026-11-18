package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"avan/internal/db/sqlc"
)

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (n Notifier) HandleTelegramNotifier(provider sqlc.Provider, message sqlc.Message) error {
	telegramConfig, err := n.ProviderService.GetTelegramConfig(provider.ID)
	if err != nil {
		return err
	}

	token := telegramConfig.BotToken // فرض بر این که توکن نیز در message وجود دارد

	// ساختن payload پیام
	telegramMessage := TelegramMessage{
		ChatID: message.Recipient,
		Text:   message.Content,
	}
	payload, err := json.Marshal(telegramMessage)
	if err != nil {
		return err
	}

	// تنظیم URL API تلگرام
	url := fmt.Sprintf("https://telegram-api-proxy-anonymous.pages.dev/api/bot%s/sendMessage", token)

	// ارسال درخواست HTTP
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// بررسی وضعیت پاسخ
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %d", resp.StatusCode)
	}

	return nil
}
