package gitlog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	headerRegex  = regexp.MustCompile(`^'(\d{4}-\d{2}-\d{2} \d{2}:\d{2})',(.*),(.*)`)
	writtenRegex = regexp.MustCompile(`^(\d+)\t(\d+)\t`)
	pathRegex    = regexp.MustCompile(`^\d+\t\d+\t(.*)`)
)

type Log struct {
	Date   time.Time `json:"date"`
	Plus   int       `json:"plus"`
	Minus  int       `json:"minus"`
	Diff   int       `json:"diff"`
	Path   string    `json:"path"`
	Author string    `json:"author"`
}

func Parse(output []string) []*Log {
	logs := make([]*Log, 0)
	log := &Log{}

	for _, line := range output {
		if matches := headerRegex.FindStringSubmatch(line); len(matches) == 4 {
			date, err := time.Parse("2006-01-02 15:04", matches[1])
			if err != nil {
				panic(err)
			}
			log.Date = date

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
			plus := 0
			if _, err := fmt.Sscanf(matches[1], "%d", &plus); err == nil {
				log.Plus = plus
			}

			minus := 0
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

	return logs
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
	if s[i].Date.Before(s[j].Date) {
		return true
	} else if s[i].Date.After(s[j].Date) {
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
