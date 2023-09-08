package plot

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type FilterConfig struct {
	StartDate time.Time
	EndDate   time.Time
	Period    string
	Authors   []string
}

type Data struct {
	Author   string
	Language string
	Plus     int
	FileExt  string
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

var (
	now        = time.Now()
	startOfDay = StartOfDay(now)
)

var periods = map[string][2]time.Time{
	"today":      {startOfDay, now},
	"24h":        {now.Add(-24 * time.Hour), now},
	"this-week":  {startOfDay.Add(-7 * 24 * time.Hour), now},
	"this-month": {time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), now},
	"this-year":  {time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), now},
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

	if startDate.After(endDate) {
		panic("Start date cannot be after end date")
	}

	if startDate.Equal(endDate) {
		panic("Start date cannot be equal to end date")
	}

	maxDuration := 30 * 30 * 24 * time.Hour
	if strings.TrimSpace(period) == "" {
		startDate = time.Now().Add(-maxDuration + 1*time.Second)
		endDate = time.Now()
	} else if p, ok := periods[period]; ok {
		startDate = p[0]
		endDate = p[1]
	} else {
		panic("Invalid period")
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
		r, err := regexp.Compile(author)
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
	r, err := regexp.Compile("(" + strings.Join(authors, "|") + ")")
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

func uniqueMonths(dateFmt string, logs []*gitlog.Log) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, log := range logs {
		monthYear := log.Date.Format(dateFmt)
		if _, exists := seen[monthYear]; !exists {
			result = append(result, monthYear)
			seen[monthYear] = struct{}{}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(dateFmt, result[i])
		t2, _ := time.Parse(dateFmt, result[j])
		return t1.Before(t2)
	})

	return result
}

func generateBarData(
	labels []string, author string,
	data map[string][]*Data,
) []opts.BarData {
	var result []opts.BarData

	for _, label := range labels {
		authorDataFound := false
		if authorLogs, ok := data[label]; ok {
			for _, log := range authorLogs {
				if log.Author == author {
					result = append(result, opts.BarData{Value: log.Plus})
					authorDataFound = true
					break
				}
			}
		}
		if !authorDataFound {
			result = append(result, opts.BarData{Value: 0})
		}
	}
	return result
}

func addSeries(
	bar *charts.Bar,
	labels []string,
	authors []string,
	data map[string][]*Data,
) {
	for _, author := range authors {
		found := false

		for _, label := range labels {
			authorLogs, ok := data[label]
			if !ok {
				continue
			}

			for _, log := range authorLogs {
				if log.Author == author {
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if found {
			bar.AddSeries(author, generateBarData(labels, author, data))
		} else {
			bar.AddSeries(author, []opts.BarData{})
		}
	}
}
