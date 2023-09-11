package data

import (
	"fmt"
	"time"
)

// PeriodRange represents a range of time.
type PeriodRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// generatePeriods generates a map of periods.
func generatePeriods(now time.Time) map[string]*PeriodRange {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
	for startOfDay.Weekday() != time.Sunday {
		startOfDay = startOfDay.Add(-24 * time.Hour)
	}
	endOfWeek := startOfDay.AddDate(0, 0, 6).Add(24*time.Hour - time.Second)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	return map[string]*PeriodRange{
		"today":      {startOfDay, endOfDay},
		"24h":        {now.AddDate(0, 0, -1), now},
		"this_week":  {startOfDay, endOfWeek},
		"7d":         {now.AddDate(0, 0, -7), now},
		"this_month": {startOfMonth, now},
		"30d":        {now.AddDate(0, 0, -30), now},
		"this_year":  {startOfYear, now},
		"1y":         {now.AddDate(-1, 0, 0), now},
	}
}

// Period returns a PeriodRange for a given key.
func Period(key string, now time.Time) (*PeriodRange, error) {
	periods := generatePeriods(now)
	p, ok := periods[key]
	if !ok {
		return p, fmt.Errorf("The period is invalid, must be one of %v", periods)
	}
	return p, nil
}
