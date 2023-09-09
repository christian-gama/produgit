package plot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/christian-gama/productivity/pkg/gitlog"
)

type FilterConfig struct {
	StartDate time.Time
	EndDate   time.Time
	Period    string
	Authors   []string
}

func NewFilterConfig(startDate, endDate time.Time, period string, authors []string) *FilterConfig {
	if len(authors) == 0 {
		panic("At least one author must be provided")
	}

	if len(authors) > 3 {
		panic("Only up to 3 authors are supported")
	}

	for _, author := range authors {
		if strings.TrimSpace(author) == "" {
			panic("Author cannot be empty")
		}
	}

	maxDuration := 30 * 30 * 24 * time.Hour
	if strings.TrimSpace(period) == "" {
		if startDate.IsZero() {
			startDate = time.Now().Add(-maxDuration + 1*time.Second)
		}
		if endDate.IsZero() {
			endDate = time.Now()
		}
	} else if p, ok := periods[period]; ok {
		startDate = p[0]
		endDate = p[1]
	} else {
		panic("Invalid period")
	}

	if startDate.After(endDate) {
		panic("Start date cannot be after end date")
	}

	if startDate.Equal(endDate) {
		panic("Start date cannot be equal to end date")
	}

	diff := endDate.Sub(startDate)
	if diff > maxDuration {
		panic("The difference between start and end date cannot be greater than 2.5 years")
	}

	if period != "" && !startDate.IsZero() && !endDate.IsZero() {
		panic("Period cannot be used with start date and end date")
	}

	return &FilterConfig{
		StartDate: startDate,
		EndDate:   endDate,
		Authors:   authors,
	}
}

func filterByDate(logs []*gitlog.Log, start, end time.Time) []*gitlog.Log {
	var result []*gitlog.Log
	for _, log := range logs {
		if (log.Date.After(start) || log.Date.Equal(start)) && log.Date.Before(end) {
			result = append(result, log)
		}
	}
	return result
}

func filterByAuthor(logs []*gitlog.Log, authors []string) []*gitlog.Log {
	var result []*gitlog.Log

	authorMap := make(map[string]string)
	for _, author := range authors {
		r, err := regexp.Compile(fmt.Sprintf("(?i)%s", author))
		if err != nil {
			panic(err)
		}
		for _, log := range logs {
			if r.MatchString(log.Author) {
				authorMap[log.Author] = author
			}
		}
	}

	for _, log := range logs {
		if replacement, ok := authorMap[log.Author]; ok {
			newLog := *log
			newLog.Author = replacement
			result = append(result, &newLog)
		}
	}

	return result
}

func mergeAuthors(logs []*gitlog.Log, authors []string) []*gitlog.Log {
	var result []*gitlog.Log
	r, err := regexp.Compile(fmt.Sprintf("(?i)(%s)", strings.Join(authors, "|")))
	if err != nil {
		panic(err)
	}

	for _, log := range logs {
		matches := r.FindStringSubmatch(log.Author)
		if len(matches) > 0 {
			newLog := *log
			newLog.Author = matches[1]
			result = append(result, &newLog)
		} else {
			result = append(result, log)
		}
	}
	return result
}

func filter(config *FilterConfig, logs []*gitlog.Log) []*gitlog.Log {
	logs = filterByDate(logs, config.StartDate, config.EndDate)
	logs = filterByAuthor(logs, config.Authors)
	logs = mergeAuthors(logs, config.Authors)
	return logs
}
