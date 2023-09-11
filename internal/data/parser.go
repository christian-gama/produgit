package data

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/christian-gama/produgit/internal/logger"
	dateutil "github.com/christian-gama/produgit/internal/util/date"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	headerRegex  = regexp.MustCompile(`^'(\d{4}-\d{2}-\d{2} \d{2}:\d{2})',(.*),(.*)`)
	writtenRegex = regexp.MustCompile(`^(\d+)\t(\d+)\t`)
	pathRegex    = regexp.MustCompile(`^\d+\t\d+\t(.*)`)
)

// Parse parses the logs.
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

			if log.Plus > 3_000 || log.Minus > 3_000 {
				logger.Warn(
					"Possibly found an anomaly: Plus: %d, Minus: %d, Path: %s, Author: %s",
					log.Plus,
					log.Minus,
					log.Path,
					log.Author,
				)
			}

			logs = append(logs, log)
			log = &Log{
				Date:   log.Date,
				Author: log.Author,
			}
		}

	}

	sort.Stable(logSorter(logs))

	return logs, nil
}
