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
)

type PlotConfig struct {
	OutputFile string
	Extension  string
	Type       string
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
		fmt.Printf(
			"Invalid extension: %s. Valid extensions are png, svg, pdf, tif.\n",
			plotConfig.Extension,
		)
		return
	}

	byteValue, err := os.ReadFile("output.json")
	if err != nil {
		panic(err)
	}

	var data []GitLog
	if err := json.Unmarshal(byteValue, &data); err != nil {
		panic(err)
	}

	// Convert date strings to time.Time and sort
	var formattedData []GitLog
	for _, entry := range data {
		dateStr := entry.Date
		date, err := time.Parse("2006-01-02 15:04", dateStr)
		if err != nil {
			panic(err)
		}
		formattedData = append(
			formattedData,
			GitLog{Date: date.Format("2006-01-02 15:04"), Plus: entry.Plus},
		)
	}

	sort.Slice(formattedData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02 15:04", formattedData[i].Date)
		dateJ, _ := time.Parse("2006-01-02 15:04", formattedData[j].Date)
		return dateI.Before(dateJ)
	})

	switch plotConfig.Type {
	case "monthly":
		monthly(formattedData)
	case "weekday":
		weekday(formattedData)
	case "time":
		timeOfDay(formattedData)
	default:
		fmt.Printf(
			"Invalid plot type: %s. Valid types are monthly, weekday, time.\n",
			plotConfig.Type,
		)
	}
}

func monthly(formattedData []GitLog) {
	// Organize data by year and month
	monthlyPlus := make(map[string]int)
	for _, entry := range formattedData {
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

	// trim to a max of 24 months
	if len(sortedKeys) > 24 {
		sortedKeys = sortedKeys[len(sortedKeys)-24:]
	}

	// Extract sorted labels and values for the graph
	var sortedLabels []string
	var sortedValues plotter.Values
	for _, k := range sortedKeys {
		sortedLabels = append(sortedLabels, k)
		sortedValues = append(sortedValues, float64(monthlyPlus[k]))
	}
	p := plot.New()
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Mês"
	p.X.Label.Text = "Mês"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	w := vg.Points(22)

	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 72, G: 61, B: 139, A: 255}

	p.Add(bars)
	p.NominalX(sortedLabels...)

	filename := fmt.Sprintf("%s.%s", plotConfig.OutputFile, plotConfig.Extension)
	if err := p.Save(50*vg.Centimeter, 30*vg.Centimeter, filename); err != nil {
		panic(err)
	}
}

func weekday(formattedData []GitLog) {
	weekdayPlus := make(map[string]int)
	weekdayLabels := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}

	for _, entry := range formattedData {
		date, _ := time.Parse("2006-01-02 15:04", entry.Date)
		day := weekdayLabels[int(date.Weekday())]
		weekdayPlus[day] += entry.Plus
	}

	var sortedValues plotter.Values
	for _, day := range weekdayLabels {
		sortedValues = append(sortedValues, float64(weekdayPlus[day]))
	}

	p := plot.New()
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Dia da Semana"
	p.X.Label.Text = "Dia da Semana"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	w := vg.Points(120)

	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 0, G: 0, B: 128, A: 255}

	p.Add(bars)
	p.NominalX(weekdayLabels...)

	filename := fmt.Sprintf("%s.%s", plotConfig.OutputFile, plotConfig.Extension)
	if err := p.Save(50*vg.Centimeter, 30*vg.Centimeter, filename); err != nil {
		panic(err)
	}
}

func timeOfDay(formattedData []GitLog) {
	timeOfDayPlus := make(map[string]int)
	timeLabels := []string{"Madrugada", "Manhã", "Tarde", "Noite"}

	for _, entry := range formattedData {
		date, _ := time.Parse("2006-01-02 15:04", entry.Date)
		hour := date.Hour()
		var period string

		if hour < 6 {
			period = "Madrugada"
		} else if hour < 12 {
			period = "Manhã"
		} else if hour < 18 {
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

	p := plot.New()
	p.Title.Text = "Quantidade de Linhas de Código Escritas por Período do Dia"
	p.X.Label.Text = "Período do Dia"
	p.Y.Label.Text = "Quantidade de Linhas de Código"

	w := vg.Points(210)

	bars, err := plotter.NewBarChart(sortedValues, w)
	if err != nil {
		panic(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}

	p.Add(bars)
	p.NominalX(timeLabels...)

	filename := fmt.Sprintf("%s.%s", plotConfig.OutputFile, plotConfig.Extension)
	if err := p.Save(50*vg.Centimeter, 30*vg.Centimeter, filename); err != nil {
		panic(err)
	}
}
