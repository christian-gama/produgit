package plot

import (
	"sort"
	"time"

	"github.com/christian-gama/produgit/internal/data"
)

// monthly is a struct that represents the monthly plot.
type monthly struct {
	bar *bar
}

// NewMonthly creates a new Monthly plot.
func NewMonthly(
	logs *data.Logs,
	config *Config,
) *monthly {
	return &monthly{
		bar: newBar(
			logs,
			"monthly",
			config,
		),
	}
}

// Plot generates the monthly chart.
func (m *monthly) Plot() error {
	logs, err := data.Filter(
		m.bar.logs,
		data.WithDate(m.bar.startDate, m.bar.endDate),
		data.WithAuthors(m.bar.authors),
		data.WithMergeAuthors(m.bar.authors),
	)
	if err != nil {
		return err
	}

	monthLabel := m.createLabels(logs)
	formattedData := m.bar.generateDataMap(
		func(l *data.Log) dataValueMap { return dataValueMap{l.GetAuthor(): l.GetPlus()} },
	)

	bar := m.bar.renderer
	bar.SetGlobalOptions(m.bar.defaultGlobalOpts("Monthly report")...)

	bar.SetXAxis(monthLabel)
	m.bar.generateSeries(monthLabel, formattedData)

	return m.bar.save()
}

// createLabels creates the labels for the monthly plot.
func (m *monthly) createLabels(logs *data.Logs) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, log := range logs.Logs {
		if log.Date.AsTime().Before(m.bar.startDate) {
			continue
		}

		monthYear := log.Date.AsTime().Format(m.bar.dateFmt)
		if _, exists := seen[monthYear]; !exists {
			result = append(result, monthYear)
			seen[monthYear] = struct{}{}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(m.bar.dateFmt, result[i])
		t2, _ := time.Parse(m.bar.dateFmt, result[j])
		return t1.Before(t2)
	})

	return result
}
