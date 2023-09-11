package data

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	dateutil "github.com/christian-gama/produgit/internal/util/date"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	headerRegex  = regexp.MustCompile(`^'(\d{4}-\d{2}-\d{2} \d{2}:\d{2})',(.*),(.*)`)
	writtenRegex = regexp.MustCompile(`^(\d+)\t(\d+)\t`)
	pathRegex    = regexp.MustCompile(`^\d+\t\d+\t(.*)`)
)

func Parse(rawLogs []string) ([]*Log, error) {
	logs := make([]*Log, 0)
	log := &Log{}

	for _, line := range rawLogs {
		if matches := headerRegex.FindStringSubmatch(line); len(matches) == 4 {
			date, err := dateutil.ToTime(matches[1])
			if err != nil {
				return nil, err
			}
			log.Date = timestamppb.New(date)

			name := matches[3]
			if name == "" {
				name = "Unknown Name"
			}

			email := matches[2]
			if email == "" {
				email = "Unknown Email"
			}

			author := strings.TrimSpace(fmt.Sprintf("%s (%s)", name, email))
			log.Author = author

			continue
		}

		if matches := writtenRegex.FindStringSubmatch(line); len(matches) == 3 {
			var plus int32
			if _, err := fmt.Sscanf(matches[1], "%d", &plus); err == nil {
				log.Plus = plus
			}

			var minus int32
			if _, err := fmt.Sscanf(matches[2], "%d", &minus); err == nil {
				log.Minus = minus
			}

			log.Diff = log.Plus - log.Minus
		}

		if matches := pathRegex.FindStringSubmatch(line); len(matches) == 2 {
			log.Path = matches[1]
			logs = append(logs, log)
			log = &Log{
				Date:   log.Date,
				Author: log.Author,
			}
		}
	}

	sort.Stable(LogSlice(logs))

	return logs, nil
}

// LogSlice is a custom type to make sorting Logs easier.
type LogSlice []*Log

func (s LogSlice) Len() int {
	return len(s)
}

func (s LogSlice) Less(i, j int) bool {
	if s[i].Author < s[j].Author {
		return true
	} else if s[i].Author > s[j].Author {
		return false
	}

	// If Author Emails are equal, compare by Date
	if s[i].Date.AsTime().Before(s[j].Date.AsTime()) {
		return true
	} else if s[i].Date.AsTime().After(s[j].Date.AsTime()) {
		return false
	}

	// If Dates are equal, compare by Path
	if s[i].Path < s[j].Path {
		return true
	} else if s[i].Path > s[j].Path {
		return false
	}

	// If Dates and Paths are equal, compare by Diff
	return s[i].Diff < s[j].Diff
}

func (s LogSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
