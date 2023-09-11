package plot

import (
	"fmt"
	"strings"
	"time"

	"github.com/christian-gama/produgit/internal/data"
)

type Config struct {
	startDate time.Time
	endDate   time.Time
	authors   []string
	period    string
	output    string
}

func NewConfig(
	startDate time.Time,
	endDate time.Time,
	authors []string,
	period string,
	output string,
) (*Config, error) {
	if len(authors) == 0 {
		return nil, fmt.Errorf("At least one author must be provided")
	}

	if len(authors) > 3 {
		return nil, fmt.Errorf("Only up to 3 authors are supported")
	}

	for _, author := range authors {
		if strings.TrimSpace(author) == "" {
			return nil, fmt.Errorf("Author cannot be empty")
		}
	}

	if period != "" && (!startDate.IsZero() || !endDate.IsZero()) {
		return nil, fmt.Errorf("Period cannot be used with start date and end date")
	}

	maxDuration := 30 * 30 * 24 * time.Hour // 2.5 years
	if strings.TrimSpace(period) == "" {
		if startDate.IsZero() {
			startDate = time.Now().Add(-maxDuration + 1*time.Second)
		}
		if endDate.IsZero() {
			endDate = time.Now()
		}
	} else if p, err := data.Period(period, time.Now()); err == nil {
		startDate = p.StartDate
		endDate = p.EndDate
	} else {
		return nil, fmt.Errorf("Invalid period")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("Start date cannot be after end date")
	}

	if startDate.Equal(endDate) {
		return nil, fmt.Errorf("Start date cannot be equal to end date")
	}

	diff := endDate.Sub(startDate)
	if diff > maxDuration {
		return nil, fmt.Errorf(
			"The difference between start and end date cannot be greater than 2.5 years",
		)
	}

	cfg := &Config{
		startDate: startDate,
		endDate:   endDate,
		authors:   authors,
		period:    period,
		output:    output,
	}

	return cfg, nil
}
