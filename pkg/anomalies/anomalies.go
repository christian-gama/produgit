package anomalies

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/christian-gama/produgit/pkg/gitlog"
	"github.com/christian-gama/produgit/pkg/report"
)

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

type Config struct {
	StartDate time.Time
	EndDate   time.Time
	Quantity  int
	Input     string
	Authors   []string
}

func NewConfig(startDate, endDate time.Time, quantity int, input string, authors []string) *Config {
	if endDate.IsZero() {
		endDate = time.Now()
	}

	if quantity <= 0 {
		panic("Quantity must be greater than 0")
	}

	if len(authors) == 0 {
		panic("No authors provided")
	}

	if strings.TrimSpace(input) == "" {
		input = "produgit_report.json"
	}

	return &Config{
		StartDate: startDate,
		EndDate:   endDate,
		Quantity:  quantity,
		Input:     input,
		Authors:   authors,
	}
}

func Anomalies(config *Config) {
	data := report.Read(config.Input)
	data = filterByDate(data, config.StartDate, config.EndDate)
	data = filterByAuthor(data, config.Authors)

	found := false
	for _, log := range data {
		if log.Plus > config.Quantity {
			if !found {
				fmt.Printf("%-6s - %-15s - %s\n", "Plus", "Author", "Path")
			}
			fmt.Printf("%-6d - %-15s - %s\n", log.Plus, fmt.Sprintf("%.15s", log.Author), log.Path)
			found = true
		}
	}

	if !found {
		fmt.Println("No anomalies found")
	}
}
