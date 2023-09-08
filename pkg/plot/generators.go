package plot

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

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

func generateBarSeries(
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
