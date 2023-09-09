package date

import (
	"fmt"
	"time"
)

func Parse(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
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
			return date
		}
	}

	panic(fmt.Sprintf("Invalid date '%s'. Expected one of %v", dateStr, formats))
}
