package data

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	dateutil "github.com/christian-gama/produgit/internal/util/date"
)

// FilterOption represents a filter option.
type FilterOption func(logs []*Log) ([]*Log, error)

// Filter filters the logs.
func Filter(logs *Logs, options ...FilterOption) (*Logs, error) {
	result := &Logs{Logs: make([]*Log, 0)}

	currentLogs := logs.Logs
	for _, option := range options {
		tempResults, err := option(currentLogs)
		if err != nil {
			return nil, err
		}

		if len(result.Logs) == 0 {
			result.Logs = append(result.Logs, tempResults...)
		} else {
			result.Logs = tempResults
		}

		currentLogs = result.Logs
	}

	return result, nil
}

// WithDate filters logs by date.
func WithDate(startTime, endTime time.Time) FilterOption {
	return func(logs []*Log) ([]*Log, error) {
		var filteredLogs []*Log

		if startTime.IsZero() {
			startTime = time.Now().AddDate(-2, -6, 0)
		}

		if endTime.IsZero() {
			endTime = time.Now()
		}

		for _, log := range logs {
			logTime := log.GetDate().AsTime()

			if logTime.After(startTime) && logTime.Before(endTime) {
				filteredLogs = append(filteredLogs, log)
			}
		}

		if len(filteredLogs) == 0 {
			return nil, fmt.Errorf(
				"No logs found between %s and %s.",
				dateutil.ToString(startTime),
				dateutil.ToString(endTime),
			)
		}

		return filteredLogs, nil
	}
}

// WithAuthors filters logs by author.
func WithAuthors(authors []string) FilterOption {
	return func(logs []*Log) ([]*Log, error) {
		var filteredLogs []*Log

		for _, author := range authors {
			r, err := regexp.Compile(fmt.Sprintf("(?i)%s", author))
			if err != nil {
				return nil, fmt.Errorf("Author expected to be a valid regex: %s.", author)
			}

			for _, log := range logs {
				if r.MatchString(log.GetAuthor()) {
					filteredLogs = append(filteredLogs, log)
				}
			}
		}

		if len(filteredLogs) == 0 {
			return nil, fmt.Errorf("No logs found for authors %s.", authors)
		}

		return filteredLogs, nil
	}
}

// WithMergeAuthors merges authors.
func WithMergeAuthors(authors []string) FilterOption {
	return func(logs []*Log) ([]*Log, error) {
		var filteredLogs []*Log

		r, err := regexp.Compile(fmt.Sprintf("(?i)(%s)", strings.Join(authors, "|")))
		if err != nil {
			return nil, fmt.Errorf(
				"Author expected to be a valid regex: %s.",
				strings.Join(authors, ","),
			)
		}

		for _, log := range logs {
			matches := r.FindStringSubmatch(log.GetAuthor())
			if len(matches) > 0 {
				newLog := log
				for _, author := range authors {
					if strings.EqualFold(author, matches[1]) {
						newLog.Author = author
						break
					}
				}
				filteredLogs = append(filteredLogs, newLog)
			} else {
				filteredLogs = append(filteredLogs, log)
			}
		}

		return filteredLogs, nil
	}
}
