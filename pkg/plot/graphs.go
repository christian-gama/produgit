package plot

import (
	"github.com/christian-gama/produgit/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Monthly(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	monthLabel := createMonthLabels("01-2006", data)
	formattedData := formatDataByDate("01-2006", data)

	bar := charts.NewBar()
	createBarHeader("Monthly Report", bar, filterConfig)

	bar.SetXAxis(monthLabel)
	generateBarSeries(bar, monthLabel, filterConfig.Authors, formattedData)

	save("monthly", config, filterConfig, bar)
}

func TimeOfDay(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	formattedData := formatDataByTimeOfDay(data)

	bar := charts.NewBar()
	createBarHeader("Time of Day Report", bar, filterConfig)

	bar.SetXAxis(timeOfDay)
	generateBarSeries(bar, timeOfDay, filterConfig.Authors, formattedData)

	save("time-of-day", config, filterConfig, bar)
}

func TopLanguages(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	langLabel := createLanguageLabel(formatDataByLanguage(data))
	formattedData := formatDataByLanguage(data)

	bar := charts.NewBar()
	createBarHeader("Top Languages Report", bar, filterConfig)

	bar.SetXAxis(langLabel)
	generateBarSeries(bar, langLabel, filterConfig.Authors, formattedData)

	save("top-languages", config, filterConfig, bar)
}

func Weekday(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	formattedData := formatDataByWeekday(data)

	bar := charts.NewBar()
	createBarHeader("Weekday Report", bar, filterConfig)

	bar.SetXAxis(daysOfWeek)
	generateBarSeries(bar, daysOfWeek, filterConfig.Authors, formattedData)

	save("weekday", config, filterConfig, bar)
}

func TopAuthors(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)

	formattedData := formatDataByAuthor(data)

	pie := charts.NewPie()
	createPieHeader("Top Authors Report", pie, filterConfig)

	pieData := generatePieData(formattedData)
	pie.AddSeries("Authors", pieData, charts.WithLabelOpts(opts.Label{
		Show:      true,
		Formatter: "{b}: {d}%",
	}),
	)

	save("top-authors", config, filterConfig, pie)
}
