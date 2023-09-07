package plot

import (
	"fmt"
	"os"
	"sort"

	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Monthly(data []*gitlog.Log) {
	monthlyPlusBySpecificAuthor := make(map[string]int)
	monthlyPlusByAnotherAuthor := make(map[string]int)
	specificAuthor := "example@email.com"
	anotherAuthor := "example2@email.com"

	for _, entry := range data {
		yearMonth := entry.Date.Format("2006-01")
		if entry.Author == specificAuthor {
			monthlyPlusBySpecificAuthor[yearMonth] += entry.Plus
		}
		if entry.Author == anotherAuthor {
			monthlyPlusByAnotherAuthor[yearMonth] += entry.Plus
		}
	}

	var sortedLabels []string
	var sortedKeys []string
	keySet := make(map[string]bool)

	for k := range monthlyPlusBySpecificAuthor {
		keySet[k] = true
	}
	for k := range monthlyPlusByAnotherAuthor {
		keySet[k] = true
	}

	for k := range keySet {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	var sortedValuesBySpecificAuthor, sortedValuesByAnotherAuthor []opts.BarData
	for _, k := range sortedKeys {
		sortedLabels = append(sortedLabels, k)
		sortedValuesBySpecificAuthor = append(
			sortedValuesBySpecificAuthor,
			opts.BarData{Value: monthlyPlusBySpecificAuthor[k]},
		)
		sortedValuesByAnotherAuthor = append(
			sortedValuesByAnotherAuthor,
			opts.BarData{Value: monthlyPlusByAnotherAuthor[k]},
		)
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Code written per month",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Data: sortedLabels,
			Name: "Month",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Lines of code",
		}),
	)

	// Add the data
	bar.SetXAxis(sortedLabels).
		AddSeries(specificAuthor, sortedValuesBySpecificAuthor).
		AddSeries("fulano@gmail.com", sortedValuesByAnotherAuthor)

	// Save the chart
	filename := "plot.monthly.html"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	err = bar.Render(f)
	if err != nil {
		fmt.Println(err)
	}
}
