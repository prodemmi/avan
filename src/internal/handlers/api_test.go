package handlers

import (
	"avan/internal/db/sqlc"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMessageStatusName(t *testing.T) {
	scheduleAt := time.Now()
	for _, test := range []struct {
		name       string
		status     int64
		scheduleAt *time.Time
		want       string
	}{
		{name: "pending", status: 1, want: "pending"},
		{name: "scheduled", status: 1, scheduleAt: &scheduleAt, want: "scheduled"},
		{name: "sent", status: 0, want: "sent"},
		{name: "failed", status: 2, want: "failed"},
		{name: "delivered", status: 3, want: "delivered"},
		{name: "read", status: 4, want: "read"},
		{name: "unknown", status: 99, want: "unknown"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := messageStatusName(test.status, test.scheduleAt); got != test.want {
				t.Fatalf("messageStatusName() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestReportedMessageStatus(t *testing.T) {
	for _, test := range []struct {
		value string
		want  int64
		valid bool
	}{
		{value: " SENT ", want: 0, valid: true},
		{value: "delivered", want: 3, valid: true},
		{value: "failed", want: 2, valid: true},
		{value: "read", want: 4, valid: true},
		{value: "pending", valid: false},
	} {
		got, valid := reportedMessageStatus(test.value)
		if got != test.want || valid != test.valid {
			t.Errorf("reportedMessageStatus(%q) = (%d, %t), want (%d, %t)", test.value, got, valid, test.want, test.valid)
		}
	}
}

func TestDuplicateFailedSendReturnsItsCurrentState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	writeSendResult(context, sqlc.Message{ID: 42, Status: 2}, true)

	var result struct {
		Data struct {
			MessageID int64  `json:"message_id"`
			Status    string `json:"status"`
			Duplicate bool   `json:"duplicate"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusBadGateway || result.Data.MessageID != 42 || result.Data.Status != "failed" || !result.Data.Duplicate {
		t.Fatalf("duplicate response = %d %+v, want failed existing message", response.Code, result.Data)
	}
}

func TestIdempotencyKeyAndRequestFingerprint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest(http.MethodPost, "/api/send", nil)
	if _, ok := requestIdempotencyKey(context); ok {
		t.Fatal("missing idempotency key was accepted")
	}
	context.Request.Header.Set("Idempotency-Key", "  request-1  ")
	if key, ok := requestIdempotencyKey(context); !ok || key != "request-1" {
		t.Fatalf("normalized key = %q, valid = %t", key, ok)
	}

	body := SendBody{ProviderID: 2, Recipient: "+15555550100", Message: "one"}
	first := sendRequestHash(body, body.Message, nil)
	body.Message = "two"
	if first == sendRequestHash(body, body.Message, nil) {
		t.Fatal("different send bodies produced the same request fingerprint")
	}
}
