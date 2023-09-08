package plot

import (
	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
)

func Monthly(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	monthLabel := uniqueMonths("01-2006", data)
	accumulated := accumulateLogsByDate("01-2006", data)

	bar := charts.NewBar()
	createBarHeader("Monthly Report", bar, filterConfig)

	bar.SetXAxis(monthLabel)
	addSeries(bar, monthLabel, filterConfig.Authors, accumulated)

	save("monthly", config, filterConfig, bar)
}

func accumulateLogsByDate(dateFmt string, logs []*gitlog.Log) map[string][]*Data {
	accumulatedLogs := make(map[string][]*Data)

	for _, log := range logs {
		dateKey := log.Date.Format(dateFmt)

		if _, ok := accumulatedLogs[dateKey]; !ok {
			accumulatedLogs[dateKey] = []*Data{}
		}

		var found *Data
		for _, aLog := range accumulatedLogs[dateKey] {
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
			accumulatedLogs[dateKey] = append(accumulatedLogs[dateKey], newLog)
		} else {
			found.Plus += log.Plus
		}
	}

	return accumulatedLogs
}
