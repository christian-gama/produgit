package plot

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func getTimeOfDay(hour int) string {
	if hour >= 0 && hour < 6 {
		return "midnight"
	} else if hour >= 6 && hour < 12 {
		return "morning"
	} else if hour >= 12 && hour < 18 {
		return "afternoon"
	} else if hour >= 18 && hour < 24 {
		return "night"
	} else {
		return "midnight"
	}
}

func getStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getDayOfWeek(day int) string {
	return time.Weekday(day).String()
}

func identifyLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) > 1 {
		if lang, ok := languages[ext[1:]]; ok {
			return lang
		}
	}
	return "Others"
}
