package plot

import "github.com/christian-gama/produgit/pkg/gitlog"

type Data struct {
	Author   string
	Language string
	Plus     int
	FileExt  string
}

func formatDataByDate(dateFmt string, logs []*gitlog.Log) map[string][]*Data {
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

func formatDataByTimeOfDay(logs []*gitlog.Log) map[string][]*Data {
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

func formatDataByLanguage(logs []*gitlog.Log) map[string][]*Data {
	accumulatedLogs := make(map[string][]*Data)

	for _, log := range logs {
		lang := identifyLanguage(log.Path)
		if _, ok := accumulatedLogs[lang]; !ok {
			accumulatedLogs[lang] = []*Data{}
		}

		var found *Data
		for _, aLog := range accumulatedLogs[lang] {
			if aLog.Author == log.Author {
				found = aLog
				break
			}
		}

		if found == nil {
			newLog := &Data{
				Author:   log.Author,
				Plus:     log.Plus,
				Language: lang,
			}
			accumulatedLogs[lang] = append(accumulatedLogs[lang], newLog)
		} else {
			found.Plus += log.Plus
		}
	}

	return accumulatedLogs
}

func formatDataByWeekday(logs []*gitlog.Log) map[string][]*Data {
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

func formatDataByAuthor(logs []*gitlog.Log) map[string]*Data {
	accumulatedLogs := make(map[string]*Data)

	for _, log := range logs {
		if _, ok := accumulatedLogs[log.Author]; !ok {
			accumulatedLogs[log.Author] = &Data{
				Author: log.Author,
				Plus:   0,
			}
		}

		accumulatedLogs[log.Author].Plus += log.Plus
	}

	return accumulatedLogs
}
