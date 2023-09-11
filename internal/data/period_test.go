package data

import (
	"testing"
	"time"
)

func TestPeriod(t *testing.T) {
	tests := []struct {
		key       string
		wantError bool
	}{
		// Test for valid period keys
		{"today", false},
		{"24h", false},
		{"this_week", false},
		{"7d", false},
		{"this_month", false},
		{"30d", false},
		{"this_year", false},
		{"1y", false},

		// Test for an invalid period key
		{"invalid_key", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			_, err := Period(tt.key, time.Now())

			if tt.wantError && err == nil {
				t.Fatalf("expected an error but got none")
			} else if !tt.wantError && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}

			if tt.wantError {
				return
			}
		})
	}
}

func TestPeriodRanges(t *testing.T) {
	now := time.Date(2023, 9, 10, 15, 30, 0, 0, time.UTC)

	expectedPeriods := map[string]*PeriodRange{
		"today": {
			time.Date(2023, 9, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 10, 23, 59, 59, 0, time.UTC),
		},
		"24h": {now.AddDate(0, 0, -1), now},
		"this_week": {
			time.Date(2023, 9, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 16, 23, 59, 59, 0, time.UTC),
		},
		"7d": {now.AddDate(0, 0, -7), now},
		"this_month": {
			time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 10, 15, 30, 0, 0, time.UTC),
		},
		"30d": {now.AddDate(0, 0, -30), now},
		"this_year": {
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 10, 15, 30, 0, 0, time.UTC),
		},
		"1y": {now.AddDate(-1, 0, 0), now},
	}

	for key, expected := range expectedPeriods {
		t.Run(key, func(t *testing.T) {
			got, err := Period(key, now)
			if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			if !got.StartDate.Equal(expected.StartDate) || !got.EndDate.Equal(expected.EndDate) {
				t.Errorf("for %s, got %v, want %v", key, got, expected)
			}
		})
	}
}
