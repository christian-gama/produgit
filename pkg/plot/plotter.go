package plot

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type renderer interface {
	Render(io.Writer) error
}

func save(graph string, config *Config, filterConfig *FilterConfig, render renderer) {
	config.Output = createFileName(graph, config, filterConfig)

	if _, err := os.Stat(config.Output); err == nil {
		err = os.Remove(config.Output)
		if err != nil {
			panic(fmt.Sprintf("Removing existing file %s failed: %v", config.Output, err))
		}
	}

	f, err := os.Create(config.Output)
	if err != nil {
		panic(fmt.Sprintf("Could not create file %s: %v", config.Output, err))
	}

	err = render.Render(f)
	if err != nil {
		panic(fmt.Sprintf("Could not render file %s: %v", config.Output, err))
	}
}

func createFileName(graph string, config *Config, filterConfig *FilterConfig) string {
	fileName := config.Output
	fileName = fmt.Sprintf("%s%s_%s", fileName, graph, strings.Join(filterConfig.Authors, "_"))

	if filterConfig.StartDate.IsZero() {
		fileName = fmt.Sprintf("%s_%s", fileName, filterConfig.EndDate.Format("2006-01-02"))
	}

	fileName = fmt.Sprintf(
		"%s_%s_%s",
		fileName,
		filterConfig.StartDate.Format("2006-01-02"),
		filterConfig.EndDate.Format("2006-01-02"),
	)

	return fmt.Sprintf("%s.html", fileName)
}

func createBarHeader(
	title string,
	bar *charts.Bar,
	filterConfig *FilterConfig,
	o ...charts.GlobalOpts,
) {
	subtitle := ""
	if filterConfig.StartDate.IsZero() {
		subtitle = fmt.Sprintf(
			"From the beginning to %s",
			filterConfig.EndDate.Format("2006-01-02"),
		)
	} else {
		subtitle = fmt.Sprintf(
			"From %s to %s",
			filterConfig.StartDate.Format("2006-01-02"),
			filterConfig.EndDate.Format("2006-01-02"),
		)
	}

	bar.SetGlobalOptions(
		append(o, charts.WithTitleOpts(opts.Title{Title: title, Subtitle: subtitle}))...,
	)
}

func createPieHeader(
	title string,
	pie *charts.Pie,
	filterConfig *FilterConfig,
	o ...charts.GlobalOpts,
) {
	subtitle := ""
	if filterConfig.StartDate.IsZero() {
		subtitle = fmt.Sprintf(
			"From the beginning to %s",
			filterConfig.EndDate.Format("2006-01-02"),
		)
	} else {
		subtitle = fmt.Sprintf(
			"From %s to %s",
			filterConfig.StartDate.Format("2006-01-02"),
			filterConfig.EndDate.Format("2006-01-02"),
		)
	}

	pie.SetGlobalOptions(
		append(o, charts.WithTitleOpts(opts.Title{Title: title, Subtitle: subtitle}))...,
	)
}
