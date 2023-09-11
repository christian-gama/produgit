package dateutil

import (
	"fmt"
	"time"
)

func ToTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01",
		"2006",
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("Invalid date format: %s.", dateStr)
}

func ToString(date time.Time) string {
	return date.Format("2006-01-02 15:04")
}
