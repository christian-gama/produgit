package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type PlotConfig struct {
	OutputFile string
	Extension  string
	StartDate  string
	FinishDate string
	Monthly    bool
	Weekday    bool
	TimeOfDay  bool
}

var plotConfig PlotConfig

func runPlot(cmd *cobra.Command, args []string) {
	validExtensions := map[string]bool{
		"png":  true,
		"svg":  true,
		"pdf":  true,
		"tif":  true,
		"jpg":  true,
		"jpeg": true,
		"tiff": true,
	}
	if !validExtensions[plotConfig.Extension] {
		panic(
			fmt.Sprintf(
				"Invalid extension: %s. Valid extensions are png, svg, pdf, tif.\n",
				plotConfig.Extension,
			),
		)
	}

	data := getData()
	if len(data) == 0 {
		panic("No data to plot.")
	}

	if plotConfig.Monthly {
		monthly(data)
	}

	if plotConfig.Weekday {
		weekday(data)
	}

	if plotConfig.TimeOfDay {
		timeOfDay(data)
	}

	if !plotConfig.Monthly && !plotConfig.Weekday && !plotConfig.TimeOfDay {
		panic("No plot type specified, please specify one of --monthly, --weekday, or --timeofday.")
	}
}

func getData() []*GitLog {
	byteValue, err := os.ReadFile("productivity_report_output.json")
	if err != nil {
		panic(err)
	}

	var data []*GitLog
	if err := json.Unmarshal(byteValue, &data); err != nil {
		panic(err)
	}

	startDate, err := time.Parse("2006-01-02", plotConfig.StartDate)
	if err != nil {
		panic(err)
	}

	endDate, err := time.Parse("2006-01-02", plotConfig.FinishDate)
	if err != nil {
		panic(err)
	}

	if endDate.Before(startDate) {
		panic("Start date must be before end date.")
	}

	if plotConfig.Monthly && endDate.Sub(startDate).Hours() > 24*365*2.5 {
		panic("Monthly plot only supported for 2.5 years of data or less.")
	}

	// Convert date strings to time.Time and sort
	var formattedData []*GitLog
	for _, entry := range data {
		dateStr := entry.Date
		date, err := time.Parse("2006-01-02 15:04", dateStr)
		if err != nil {
			panic(err)
		}

		if date.Before(startDate) || date.After(endDate) {
			continue
		}

		formattedData = append(
			formattedData,
			&GitLog{Date: date.Format("2006-01-02 15:04"), Plus: entry.Plus},
		)
	}

	sort.Slice(formattedData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02 15:04", formattedData[i].Date)
		dateJ, _ := time.Parse("2006-01-02 15:04", formattedData[j].Date)
		return dateI.Before(dateJ)
	})

	return formattedData
}

func monthly(data []*GitLog) {
	// Organize data by year and month
	monthlyPlus := make(map[string]int)
	for _, entry := range data {
		date, _ := time.Parse("2006-01-02 15:04", entry.Date)
		yearMonth := date.Format("2006-01")
		monthlyPlus[yearMonth] += entry.Plus
	}

	// Sort keys of monthlyPlus for graphing
	var sortedKeys []string
	for k := range monthlyPlus {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var sortedLabels []string
	var sortedValues plotter.Values
	for _, k := range sortedKeys {
		sortedLabels = append(sortedLabels, k)
		sortedValues = append(sortedValues, float64(monthlyPlus[k]))
	}

	p, w := prettifyPlot(sortedLabels)
	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 72, G: 61, B: 139, A: 255}
	p.Add(bars)
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Mês"
	p.X.Label.Text = "Mês"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	save(p, "monthly")
}

func prettifyPlot(keys []string) (p *plot.Plot, w vg.Length) {
	p = plot.New()
	p.X.Padding = 1 * vg.Centimeter
	p.Y.Padding = 1 * vg.Centimeter
	p.Title.Padding = 2 * vg.Centimeter
	p.Title.TextStyle.Font.Size = vg.Points(24)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Max = float64(len(keys))
	p.X.Tick.Label.Rotation = 45
	p.X.Tick.Label.XAlign = draw.XCenter
	w = vg.Points(600 / float64(len(keys)))
	p.NominalX(keys...)
	return p, w
}

func save(p *plot.Plot, t string) {
	filename := fmt.Sprintf("%s.%s.%s", plotConfig.OutputFile, t, plotConfig.Extension)
	if err := p.Save(40*vg.Centimeter, 28*vg.Centimeter, filename); err != nil {
		panic(err)
	}
}

func weekday(data []*GitLog) {
	weekdayPlus := make(map[string]int)
	weekdayLabels := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}

	for _, entry := range data {
		date, _ := time.Parse("2006-01-02 15:04", entry.Date)
		day := weekdayLabels[int(date.Weekday())]
		weekdayPlus[day] += entry.Plus
	}

	var sortedValues plotter.Values
	for _, day := range weekdayLabels {
		sortedValues = append(sortedValues, float64(weekdayPlus[day]))
	}

	p, w := prettifyPlot(weekdayLabels)
	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}

	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 0, G: 0, B: 128, A: 255}
	p.Add(bars)
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Dia da Semana"
	p.X.Label.Text = "Dia da Semana"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	save(p, "weekday")
}

func timeOfDay(data []*GitLog) {
	timeOfDayPlus := make(map[string]int)
	timeLabels := []string{"Madrugada", "Manhã", "Tarde", "Noite"}

	for _, entry := range data {
		date, _ := time.Parse("2006-01-02 15:04", entry.Date)
		hour := date.Hour()
		var period string

		if hour >= 0 && hour < 7 {
			period = "Madrugada"
		} else if hour >= 7 && hour < 13 {
			period = "Manhã"
		} else if hour >= 13 && hour < 19 {
			period = "Tarde"
		} else {
			period = "Noite"
		}

		timeOfDayPlus[period] += entry.Plus
	}

	var sortedValues plotter.Values
	for _, period := range timeLabels {
		sortedValues = append(sortedValues, float64(timeOfDayPlus[period]))
	}

	p, w := prettifyPlot(timeLabels)
	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}

	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	p.Add(bars)
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Período do Dia"
	p.X.Label.Text = "Período do Dia"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	save(p, "timeofday")
}
