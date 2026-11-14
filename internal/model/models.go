package model

import (
	"github.com/satori/go.uuid"
	"time"
)

// Workspace represents a user's project where messages are managed
type Workspace struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Unique workspace ID
	Name      string    `json:"name"`                                 // Workspace name
	OwnerID   uuid.UUID `json:"owner_id"`                             // User who owns the workspace
	CreatedAt time.Time `json:"created_at"`                           // Creation timestamp
	UpdatedAt time.Time `json:"updated_at"`                           // Last updated timestamp
}

// MessageTemplate represents a reusable message template (Markdown or HTML)
type MessageTemplate struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Unique template ID
	WorkspaceID uuid.UUID `json:"workspace_id"`                         // Workspace association
	Name        string    `json:"name"`                                 // Template name
	Format      string    `json:"format"`                               // Format: markdown or html
	Content     string    `json:"content"`                              // Template content
	CreatedAt   time.Time `json:"created_at"`                           // Timestamp when created
	UpdatedAt   time.Time `json:"updated_at"`                           // Last updated timestamp
}

// Message represents a scheduled or instant message
type Message struct {
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;"` // Unique message ID
	WorkspaceID uuid.UUID  `json:"workspace_id"`                         // Belongs to a workspace
	TemplateID  *uuid.UUID `json:"template_id,omitempty"`                // Optional template ID (if used)
	Type        string     `json:"type"`                                 // Message type: email, sms, telegram, slack
	Recipient   string     `json:"recipient"`                            // Recipient address/phone/user ID
	Content     string     `json:"content"`                              // Message body (could be from a template)
	Status      string     `json:"status"`                               // Status: pending, sent, failed
	ScheduleAt  *time.Time `json:"schedule_at,omitempty"`                // Optional scheduled time
	CreatedAt   time.Time  `json:"created_at"`                           // Timestamp when created
	SentAt      *time.Time `json:"sent_at,omitempty"`                    // Timestamp when sent
}

// EmailConfig represents email-specific configurations for a workspace
type EmailConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Config ID
	WorkspaceID uuid.UUID `json:"workspace_id"`                         // Workspace association
	SMTPHost    string    `json:"smtp_host"`                            // SMTP Server
	SMTPPort    int       `json:"smtp_port"`                            // SMTP Port
	Username    string    `json:"username"`                             // SMTP Username
	Password    string    `json:"password"`                             // SMTP Password
	FromEmail   string    `json:"from_email"`                           // Sender Email
}

// SMSConfig represents SMS provider configurations for a workspace
type SMSConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Config ID
	WorkspaceID uuid.UUID `json:"workspace_id"`                         // Workspace association
	Provider    string    `json:"provider"`                             // e.g., Twilio, Nexmo
	APIKey      string    `json:"api_key"`                              // API Key
	SenderID    string    `json:"sender_id"`                            // Sender Identifier
}

// TelegramConfig represents Telegram bot configurations for a workspace
type TelegramConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Config ID
	WorkspaceID uuid.UUID `json:"workspace_id"`                         // Workspace association
	BotToken    string    `json:"bot_token"`                            // Target Chat ID
}

// SlackConfig represents Slack webhook configurations for a workspace
type SlackConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Config ID
	WorkspaceID uuid.UUID `json:"workspace_id"`                         // Workspace association
	WebhookURL  string    `json:"webhook_url"`                          // Slack Webhook SmsURL
}

// NotificationLog stores logs of sent messages
type NotificationLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"` // Log ID
	MessageID uuid.UUID `json:"message_id"`                           // Reference to Message ID
	Status    string    `json:"status"`                               // Success or failure
	Error     *string   `json:"error,omitempty"`                      // Error message if failed
	CreatedAt time.Time `json:"created_at"`                           // Timestamp
}
