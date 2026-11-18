package services

import (
	"testing"
	"time"
)

func TestAPIKeyExpiration(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	future := now.Add(time.Minute)
	past := now.Add(-time.Minute)
	equal := now

	for _, test := range []struct {
		name     string
		expireAt *time.Time
		want     bool
	}{
		{name: "no expiry", want: false},
		{name: "future", expireAt: &future, want: false},
		{name: "past", expireAt: &past, want: true},
		{name: "expires now", expireAt: &equal, want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := apiKeyIsExpired(test.expireAt, now); got != test.want {
				t.Fatalf("apiKeyIsExpired() = %t, want %t", got, test.want)
			}
		})
	}
}
