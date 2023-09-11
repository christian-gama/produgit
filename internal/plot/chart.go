package plot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/christian-gama/produgit/internal/data"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

// dataMap is a map with date as key and value as a map of any string as key to int32 as value.
type (
	dataValueMap map[string]int32
	dataMap      map[string]dataValueMap
)

// chart holds the configuration for a plot.
type chart[T render.Renderer] struct {
	*Config

	logs *data.Logs

	renderer  T
	chartName string
	dateFmt   string
}

// NewPlot creates a new plot.
func NewPlot[T render.Renderer](
	render T,
	chartName string,
	config *Config,
	logs *data.Logs,
) *chart[T] {
	return &chart[T]{
		Config: config,

		logs: logs,

		renderer:  render,
		chartName: chartName,
		dateFmt:   "2006-01",
	}
}

// save saves the plot to a HTML file.
func (p *chart[T]) save() error {
	fileName, err := p.createFileName()
	if err != nil {
		return err
	}

	if _, err := os.Stat(fileName); err == nil {
		err = os.Remove(fileName)
		if err != nil {
			return fmt.Errorf("Removing existing file %s failed: %v", fileName, err)
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %v", fileName, err)
	}

	err = p.renderer.Render(f)
	if err != nil {
		return fmt.Errorf("Could not render file %s: %v", fileName, err)
	}

	return nil
}

// createFileName creates the file name for the plot.
func (p *chart[T]) createFileName() (string, error) {
	output := p.output

	output = strings.ReplaceAll(output, "<authors>", strings.Join(p.authors, "_"))
	output = strings.ReplaceAll(
		output,
		"<date>",
		fmt.Sprintf("%s_%s", p.startDate.Format("200601021504"), p.endDate.Format("200601021504")),
	)
	output = strings.ReplaceAll(output, "<chart>", p.chartName)
	output = strings.ReplaceAll(output, "<timestamp>", time.Now().Format("20060102150405"))
	output = strings.ReplaceAll(output, "<start_date>", p.startDate.Format("200601021504"))
	output = strings.ReplaceAll(output, "<end_date>", p.endDate.Format("200601021504"))

	if ext := filepath.Ext(output); ext != "" {
		output = fmt.Sprintf("%s%s", strings.TrimSuffix(output, ext), ext)
	} else {
		output = fmt.Sprintf("%s.html", output)
	}

	return output, nil
}

// defaultGlobalOpts creates the default global options for a plot.
func (p *chart[T]) defaultGlobalOpts(
	title string,
	globalOpts ...charts.GlobalOpts,
) []charts.GlobalOpts {
	subtitle := ""
	if p.startDate.IsZero() {
		subtitle = fmt.Sprintf(
			"From the beginning to %s",
			p.endDate.Format("2006-01-02"),
		)
	} else {
		subtitle = fmt.Sprintf(
			"From %s to %s",
			p.startDate.Format("2006-01-02"),
			p.endDate.Format("2006-01-02"),
		)
	}

	return append(
		globalOpts,
		charts.WithTitleOpts(opts.Title{Title: title, Subtitle: subtitle}),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:     "1280px",
				Height:    "720px",
				PageTitle: fmt.Sprintf("%s - %s", title, subtitle),
			},
		),
	)
}

// generateDataMap generates a map of data from a list of logs.
func (p *chart[T]) generateDataMap(
	createKey func(c *chart[T], l *data.Log) string,
	createData func(c *chart[T], l *data.Log) dataValueMap,
) dataMap {
	accumulatedLogs := make(dataMap)

	for _, log := range p.logs.Logs {
		rootKey := createKey(p, log)

		if _, ok := accumulatedLogs[rootKey]; !ok {
			accumulatedLogs[rootKey] = make(dataValueMap)
		}

		for key, value := range createData(p, log) {
			accumulatedLogs[rootKey][key] += value
		}
	}

	return accumulatedLogs
}
