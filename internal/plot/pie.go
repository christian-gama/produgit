package plot

import (
	"github.com/christian-gama/produgit/internal/data"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// pie is a struct that contains the pie chart.
type pie struct {
	*chart[*charts.Pie]
}

// newPie returns a new pie chart.
func newPie(
	logs *data.Logs,
	chartName string,
	config *Config,
) *pie {
	return &pie{
		NewPlot[*charts.Pie](
			charts.NewPie(),
			chartName,
			config,
			logs,
		),
	}
}

// generateData generates the data for the bar chart.
func (p *pie) generateData(
	labels []string,
	data dataMap,
) []opts.PieData {
	var pieData []opts.PieData

	for _, label := range labels {
		for author, authorData := range data[label] {
			pieData = append(
				pieData,
				opts.PieData{
					Name:  author,
					Value: authorData,
				},
			)
		}
	}

	return pieData
}

// generateSeries generates the series for the bar chart.
func (p *pie) generateSeries(
	labels []string,
	data dataMap,
) {
	p.renderer.AddSeries(
		"Authors",
		p.generateData(labels, data),
		charts.WithLabelOpts(opts.Label{
			Show:      true,
			Formatter: "{b}: {c} ({d}%)",
		}),
	)
}

// setGlobalOptions sets the global options for the bar chart.
func (p *pie) setGlobalOptions(title string) {
	p.renderer.SetGlobalOptions(
		append(
			p.defaultGlobalOpts(title),
			charts.WithTooltipOpts(opts.Tooltip{
				Show:      true,
				Trigger:   "item",
				Formatter: "{a} <br/>{b} : {c} ({d}%)",
			}),
		)...,
	)
}
