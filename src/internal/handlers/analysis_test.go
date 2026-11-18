package handlers

import (
	"testing"
	"time"
)

func TestAnalysisDateRange(t *testing.T) {
	now := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.Local)
	start, end, err := analysisDateRange("last_7_days", "", "", now)
	if err != nil || !start.Equal(now.AddDate(0, 0, -7)) || !end.Equal(now) {
		t.Fatalf("last_7_days range = (%v, %v, %v)", start, end, err)
	}

	start, end, err = analysisDateRange("custom", "2026-01-08", "2026-01-09", now)
	if err != nil || start.Format(time.DateOnly) != "2026-01-08" || end.Format(time.DateOnly) != "2026-01-09" || end.Hour() != 23 {
		t.Fatalf("custom range = (%v, %v, %v)", start, end, err)
	}
	if _, _, err = analysisDateRange("custom", "2026-01-10", "2026-01-09", now); err == nil {
		t.Fatal("expected reversed custom range to fail")
	}
}
