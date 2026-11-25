package providers

type ProviderTemplate struct {
	ID             string
	Name           string
	ConfigType     string
	Region         string
	Description    string
	Documentation  string
	RequiredFields []string
	URL            string
	Method         string
	AuthKey        string
	AuthType       string
	BodyEncoding   string
	BodySchema     string
	SMTPHost       string
	SMTPPort       string
	SMTPUsername   string
}

var providerTemplates = []ProviderTemplate{
	{
		ID: "kavenegar", Name: "Kavenegar", ConfigType: "sms", Region: "Local",
		Description:    "GET-based REST SMS gateway; enter your API key before sending.",
		Documentation:  "https://kavenegar.com/rest.html",
		RequiredFields: []string{"API key", "Receptor", "Message"},
		URL:            "https://api.kavenegar.com/v1/$API-KEY/sms/send.json?receptor=$send_to&message=$message",
		Method:         "0", AuthKey: "API-KEY",
	},
	{
		ID: "melipayamak", Name: "Melipayamak", ConfigType: "sms", Region: "Local",
		Description:    "Send SMS through the HTTPS form-urlencoded API.",
		Documentation:  "https://www.melipayamak.com/api/sendsimplesms2/",
		RequiredFields: []string{"API username", "Password/API key", "Sender number", "Recipient", "Message"},
		URL:            "https://rest.payamak-panel.com/api/SendSMS/SendSMS",
		Method:         "1", AuthKey: "password", AuthType: "token", BodyEncoding: "form",
		BodySchema: `{"username":"$username","password":"","to":"$send_to","from":"$sender","text":"$message"}`,
	},
	{
		ID: "twilio", Name: "Twilio SMS", ConfigType: "sms", Region: "Global",
		Description:    "Send SMS through Twilio's Messages API using API-key Basic Auth.",
		Documentation:  "https://www.twilio.com/docs/messaging/api/message-resource",
		RequiredFields: []string{"Account SID", "API Key SID", "API key secret", "From phone number"},
		URL:            "https://api.twilio.com/2010-04-01/Accounts/$account_sid/Messages.json",
		Method:         "1", AuthType: "basic", BodyEncoding: "form",
		BodySchema: `{"To":"$send_to","From":"$sender","Body":"$message"}`,
	},
	{
		ID: "sendgrid", Name: "SendGrid SMTP", ConfigType: "email", Region: "Global",
		Description:    "Send email through SendGrid's SMTP relay.",
		Documentation:  "https://www.twilio.com/docs/sendgrid/for-developers/sending-email/getting-started-smtp",
		RequiredFields: []string{"API key (SMTP password)", "Verified sender email"},
		SMTPHost:       "smtp.sendgrid.net", SMTPPort: "587", SMTPUsername: "apikey",
	},
	{
		ID: "slack", Name: "Slack Incoming Webhook", ConfigType: "slack", Region: "Global",
		Description:    "Post message text to a channel using a Slack incoming webhook.",
		Documentation:  "https://docs.slack.dev/messaging/sending-messages-using-incoming-webhooks",
		RequiredFields: []string{"Incoming webhook URL"},
	},
	{
		ID: "webhook", Name: "Custom Webhook", ConfigType: "webhook", Region: "Global",
		Description:    "Send the service's JSON message payload to your own endpoint.",
		RequiredFields: []string{"Endpoint URL"},
	},
}

func AllProviderTemplates() []ProviderTemplate {
	return providerTemplates
}

func FindProviderTemplate(id string) (ProviderTemplate, bool) {
	for _, template := range providerTemplates {
		if template.ID == id {
			return template, true
		}
	}
	return ProviderTemplate{}, false
}
