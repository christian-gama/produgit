package plot

import (
	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
)

func TimeOfDay(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	accumulated := accumulateLogsByTimeOfDay(data)

	bar := charts.NewBar()
	createBarHeader("Time of Day Report", bar, filterConfig)

	bar.SetXAxis(timeOfDay)
	addSeries(bar, timeOfDay, filterConfig.Authors, accumulated)

	save("time-of-day", config, filterConfig, bar)
}

var timeOfDay = []string{
	"midnight",
	"morning",
	"afternoon",
	"night",
}

func getTimeOfDay(hour int) string {
	if hour >= 0 && hour < 6 {
		return "midnight"
	} else if hour >= 6 && hour < 12 {
		return "morning"
	} else if hour >= 12 && hour < 18 {
		return "afternoon"
	} else if hour >= 18 && hour < 24 {
		return "night"
	} else {
		return "midnight"
	}
}

func accumulateLogsByTimeOfDay(logs []*gitlog.Log) map[string][]*Data {
	accumulatedLogs := make(map[string][]*Data)

	for _, log := range logs {
		timeOfDayKey := getTimeOfDay(log.Date.Hour())

		if _, ok := accumulatedLogs[timeOfDayKey]; !ok {
			accumulatedLogs[timeOfDayKey] = []*Data{}
		}

		var found *Data
		for _, aLog := range accumulatedLogs[timeOfDayKey] {
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
			accumulatedLogs[timeOfDayKey] = append(accumulatedLogs[timeOfDayKey], newLog)
		} else {
			found.Plus += log.Plus
		}
	}

	return accumulatedLogs
}
