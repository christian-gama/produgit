package plot

import (
	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func TopAuthors(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)

	accumulated := accumulateLogsByAuthor(data)

	pie := charts.NewPie()
	createPieHeader("Top Authors Report", pie, filterConfig)

	pieData := generatePieData(accumulated)
	pie.AddSeries("Authors", pieData)

	save("top-authors", config, filterConfig, pie)
}

func accumulateLogsByAuthor(logs []*gitlog.Log) map[string]*Data {
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

func generatePieData(accumulated map[string]*Data) []opts.PieData {
	var pieData []opts.PieData

	for author, log := range accumulated {
		pieData = append(pieData, opts.PieData{
			Value: log.Plus,
			Name:  author,
		})
	}

	return pieData
}
