package plot

import (
	"time"

	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
)

func Weekday(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	accumulated := accumulateLogsByWeekday(data)

	bar := charts.NewBar()
	createBarHeader("Weekday Report", bar, filterConfig)

	bar.SetXAxis(daysOfWeek)
	addSeries(bar, daysOfWeek, filterConfig.Authors, accumulated)

	save("weekday", config, filterConfig, bar)
}

var daysOfWeek = []string{
	time.Sunday.String(),
	time.Monday.String(),
	time.Tuesday.String(),
	time.Wednesday.String(),
	time.Thursday.String(),
	time.Friday.String(),
	time.Saturday.String(),
}

func getDayOfWeek(day int) string {
	return time.Weekday(day).String()
}

func accumulateLogsByWeekday(logs []*gitlog.Log) map[string][]*Data {
	accumulatedLogs := make(map[string][]*Data)

	for _, log := range logs {
		weekdayKey := getDayOfWeek(int(log.Date.Weekday()))

		if _, ok := accumulatedLogs[weekdayKey]; !ok {
			accumulatedLogs[weekdayKey] = []*Data{}
		}

		var found *Data
		for _, aLog := range accumulatedLogs[weekdayKey] {
			if aLog.Author == log.Author {
				found = aLog
				break
			}
		}

		if found == nil {
			newLog := &Data{
				Author: log.Author,
				Plus:   log.Plus,
			}
			accumulatedLogs[weekdayKey] = append(accumulatedLogs[weekdayKey], newLog)
		} else {
			found.Plus += log.Plus
		}
	}

	return accumulatedLogs
}
