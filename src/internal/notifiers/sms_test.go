package notifiers

import (
	"avan/internal/db/sqlc"
	"encoding/json"
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
