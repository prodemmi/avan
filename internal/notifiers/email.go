package notifiers

import (
	"avan/internal/db/sqlc"
	"fmt"
	"net/smtp"
)

func (n Notifier) HandleEmailNotifier(provider sqlc.Provider, message sqlc.Message) error {
	// دریافت تنظیمات ایمیل از ProviderService
	emailConfig, err := n.ProviderService.GetEmailConfig(provider.ID)
	if err != nil {
		return fmt.Errorf("failed to get email config: %w", err)
	}

	// بررسی مقدار Username و Password
	username := ""
	password := ""
	if emailConfig.Username != nil {
		username = *emailConfig.Username
	}
	if emailConfig.Password != nil {
		password = *emailConfig.Password
	}

	// احراز هویت SMTP (در صورت نیاز)
	var auth smtp.Auth
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, emailConfig.SmtpHost)
	}

	// قالب ایمیل
	subject := "\n"
	msg := []byte(subject + "Content-Type: text/html; charset=\"UTF-8\"\n\n" + message.Content)

	// ارسال ایمیل
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", emailConfig.SmtpHost, emailConfig.SmtpPort),
		auth,
		emailConfig.FromEmail,
		[]string{message.Recipient},
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email successfully sent to:", message.Recipient)
	return nil
}
