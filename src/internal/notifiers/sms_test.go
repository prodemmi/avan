package notifiers

import (
	"avan/internal/db/sqlc"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestBuildSmsPostBody(t *testing.T) {
	payload := `{"tag":"notice"}`
	message := sqlc.Message{Recipient: "+98123", Content: `hello "world"`, Payload: &payload, Retry: true, ProviderID: 7, WorkspaceID: 9}
	key, value := "token", "secret"
	body, err := buildSmsPostBody(`{"to":"$send_to","message":"$message","data":"$payload","retry":"$retry","provider_id":"$provider_id","workspace_id":"$workspace_id"}`, message, &key, &value, true)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["to"] != "+98123" || got["message"] != `hello "world"` || got["token"] != "secret" {
		t.Fatalf("unexpected body: %s", body)
	}
	if got["data"].(map[string]any)["tag"] != "notice" {
		t.Fatalf("payload placeholder was not inserted as JSON: %s", body)
	}
	if got["retry"] != true || got["provider_id"] != float64(7) || got["workspace_id"] != float64(9) {
		t.Fatalf("dynamic variables were not substituted: %s", body)
	}
}

func TestBuildTwilioSMSRequest(t *testing.T) {
	accountSID, apiKeySID, apiKeySecret, sender := "AC123", "SK456", "secret", "+15551234567"
	bodySchema := `{"To":"$send_to","From":"$sender","Body":"$message"}`
	config := sqlc.SmsConfig{
		Url:            "https://api.twilio.com/2010-04-01/Accounts/$account_sid/Messages.json",
		Method:         MethodPOST,
		BodyEncoding:   "form",
		AuthType:       "basic",
		AuthUsername:   &apiKeySID,
		AccountSid:     &accountSID,
		Sender:         &sender,
		AuthTokenValue: &apiKeySecret,
		BodySchema:     &bodySchema,
	}
	req, err := buildSmsRequest(config, sqlc.Message{Recipient: "+15558675310", Content: "hello & world"})
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "POST" || req.URL.Path != "/2010-04-01/Accounts/AC123/Messages.json" {
		t.Fatalf("request target = %s %s", req.Method, req.URL.Path)
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q", req.Header.Get("Content-Type"))
	}
	username, password, ok := req.BasicAuth()
	if !ok || username != apiKeySID || password != apiKeySecret {
		t.Fatal("Twilio Basic Auth credentials were not set")
	}
	values := readFormBody(t, req)
	if values.Get("To") != "+15558675310" || values.Get("From") != sender || values.Get("Body") != "hello & world" {
		t.Fatalf("unexpected Twilio form fields: %v", values)
	}
}

func TestBuildMelipayamakSMSRequest(t *testing.T) {
	username, password, sender := "api-user", "password-secret", "5000"
	bodySchema := `{"username":"$username","password":"","to":"$send_to","from":"$sender","text":"$message"}`
	isAuthInHeader := false
	config := sqlc.SmsConfig{
		Url:            "https://rest.payamak-panel.com/api/SendSMS/SendSMS",
		Method:         MethodPOST,
		BodyEncoding:   "form",
		AuthType:       "token",
		AuthUsername:   &username,
		Sender:         &sender,
		AuthTokenKey:   stringPtr("password"),
		AuthTokenValue: &password,
		IsAuthInHeader: &isAuthInHeader,
		BodySchema:     &bodySchema,
	}
	req, err := buildSmsRequest(config, sqlc.Message{Recipient: "09120000000", Content: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if req.URL.String() != "https://rest.payamak-panel.com/api/SendSMS/SendSMS" {
		t.Fatalf("request URL = %q", req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q", req.Header.Get("Content-Type"))
	}
	values := readFormBody(t, req)
	if values.Get("username") != username || values.Get("password") != password || values.Get("to") != "09120000000" || values.Get("from") != sender || values.Get("text") != "hello" {
		t.Fatalf("unexpected Melipayamak form fields: %v", values)
	}
}

func TestBuildLegacySMSRequests(t *testing.T) {
	t.Run("JSON and header auth", func(t *testing.T) {
		bodySchema := `{"to":"$send_to","message":"$message"}`
		authKey, authValue := "Authorization", "Bearer token"
		isAuthInHeader := true
		config := sqlc.SmsConfig{
			Url: "https://legacy.test/send", Method: MethodPOST,
			AuthTokenKey: &authKey, AuthTokenValue: &authValue, IsAuthInHeader: &isAuthInHeader,
			BodySchema: &bodySchema,
		}
		req, err := buildSmsRequest(config, sqlc.Message{Recipient: "+123", Content: "hello"})
		if err != nil {
			t.Fatal(err)
		}
		if req.Header.Get("Content-Type") != "application/json" || req.Header.Get(authKey) != authValue {
			t.Fatalf("legacy headers changed: %v", req.Header)
		}
		var fields map[string]string
		if err := json.NewDecoder(req.Body).Decode(&fields); err != nil || fields["to"] != "+123" || fields["message"] != "hello" {
			t.Fatalf("legacy JSON request changed: %v, error: %v", fields, err)
		}
	})

	t.Run("GET query placeholders and token", func(t *testing.T) {
		key, secret := "token", "secret"
		config := sqlc.SmsConfig{
			Url:    "https://legacy.test/send?to=$send_to&message=$message&key=$token",
			Method: MethodGET, AuthTokenKey: &key, AuthTokenValue: &secret,
		}
		req, err := buildSmsRequest(config, sqlc.Message{Recipient: "12345", Content: "hello"})
		if err != nil {
			t.Fatal(err)
		}
		query := req.URL.Query()
		if req.Method != "GET" || query.Get("to") != "12345" || query.Get("message") != "hello" || query.Get("key") != secret {
			t.Fatalf("legacy GET query changed: %v", query)
		}
	})
}

func TestSmsResponseError(t *testing.T) {
	t.Run("accepts Twilio resource creation", func(t *testing.T) {
		if err := smsResponseError("https://api.twilio.com/2010-04-01/Accounts/AC123/Messages.json", 201, []byte(`{"sid":"SM123"}`)); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("accepts Melipayamak message ID", func(t *testing.T) {
		if err := smsResponseError("https://rest.payamak-panel.com/api/SendSMS/SendSMS", 200, []byte("123456")); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("rejects Melipayamak application error", func(t *testing.T) {
		if err := smsResponseError("https://rest.payamak-panel.com/api/SendSMS/SendSMS", 200, []byte("-110")); err == nil {
			t.Fatal("expected a provider application error")
		}
	})
	t.Run("rejects HTTP errors", func(t *testing.T) {
		if err := smsResponseError("https://api.example.com/send", 503, []byte("unavailable")); err == nil {
			t.Fatal("expected an HTTP error")
		}
	})
}

func readFormBody(t *testing.T, req *http.Request) url.Values {
	t.Helper()
	body, err := req.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	values, err := url.ParseQuery(string(data))
	if err != nil {
		t.Fatal(err)
	}
	return values
}

func stringPtr(value string) *string { return &value }
