package components

import (
	"strings"
	"testing"
	"time"
)

func TestGetValueLocalizesDateTimes(t *testing.T) {
	date := time.Date(2026, time.November, 26, 13, 45, 0, 0, time.UTC)
	got := GetValue(&date)
	if !strings.Contains(got, `data-localized-date-time datetime="2026-11-26T13:45:00Z"`) || !strings.Contains(got, "26 Nov 26 13:45 UTC") {
		t.Fatalf("GetValue() = %q, want localized datetime element", got)
	}
	if got := GetValue((*time.Time)(nil)); got != "" {
		t.Fatalf("nil date = %q, want empty", got)
	}
}
