package plot

import (
	"github.com/christian-gama/produgit/internal/data"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// bar is a struct that contains the bar chart.
type bar struct {
	*chart[*charts.Bar]
}

// newBar returns a new bar chart.
func newBar(
	logs *data.Logs,
	chartName string,
	config *Config,
) *bar {
	return &bar{
		NewPlot[*charts.Bar](
			charts.NewBar(),
			chartName,
			config,
			logs,
		),
	}
}

// generateData generates the data for the bar chart.
func (b *bar) generateData(
	labels []string,
	author string,
	data dataMap,
) []opts.BarData {
	var result []opts.BarData

	for _, label := range labels {
		authorData, ok := data[label][author]
		if ok {
			result = append(result, opts.BarData{Value: authorData})
		} else {
			result = append(result, opts.BarData{Value: int32(0)})
		}
	}

	return result
}

// generateSeries generates the series for the bar chart.
func (b *bar) generateSeries(
	labels []string,
	data dataMap,
) {
	for _, author := range b.authors {
		found := false

		for _, label := range labels {
			if _, ok := data[label][author]; ok {
				found = true
				break
			}
		}

		if found {
			b.renderer.AddSeries(author, b.generateData(labels, author, data))
		} else {
			b.renderer.AddSeries(author, []opts.BarData{})
		}
	}
}

// setGlobalOptions sets the global options for the bar chart.
func (b *bar) setGlobalOptions(title string) {
	b.renderer.SetGlobalOptions(
		append(
			b.defaultGlobalOpts(title),
			charts.WithTooltipOpts(opts.Tooltip{
				Show:      true,
				Trigger:   "axis",
				Formatter: "{b} <br />{a} : {c}",
			}),
		)...,
	)
}
