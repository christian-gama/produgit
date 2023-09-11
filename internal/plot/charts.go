package plot

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/christian-gama/produgit/internal/data"
	"github.com/go-echarts/go-echarts/v2/charts"
)

// monthly is a struct that represents the monthly plot.
type monthly struct {
	bar *bar
}

// NewMonthly creates a new Monthly plot.
func NewMonthly(
	logs *data.Logs,
	config *Config,
) *monthly {
	return &monthly{
		bar: newBar(
			logs,
			"monthly",
			config,
		),
	}
}

// Plot generates the monthly chart.
func (m *monthly) Plot() error {
	logs, err := data.Filter(
		m.bar.logs,
		data.WithDate(m.bar.startDate, m.bar.endDate),
		data.WithAuthors(m.bar.authors),
		data.WithMergeAuthors(m.bar.authors),
	)
	if err != nil {
		return err
	}

	monthLabel := m.createLabels(logs)
	formattedData := m.bar.generateDataMap(
		func(c *chart[*charts.Bar], l *data.Log) string {
			return l.GetDate().AsTime().Format(c.dateFmt)
		},
		func(_ *chart[*charts.Bar], l *data.Log) dataValueMap {
			return dataValueMap{l.GetAuthor(): l.GetPlus()}
		},
	)

	m.bar.setGlobalOptions("Monthly Report")
	bar := m.bar.renderer

	bar.SetXAxis(monthLabel)
	m.bar.generateSeries(monthLabel, formattedData)

	return m.bar.save()
}

// createLabels creates the labels for the monthly plot.
func (m *monthly) createLabels(logs *data.Logs) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, log := range logs.Logs {
		if log.GetDate().AsTime().Before(m.bar.startDate) {
			continue
		}

		monthYear := log.GetDate().AsTime().Format(m.bar.dateFmt)
		if _, exists := seen[monthYear]; !exists {
			result = append(result, monthYear)
			seen[monthYear] = struct{}{}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(m.bar.dateFmt, result[i])
		t2, _ := time.Parse(m.bar.dateFmt, result[j])
		return t1.Before(t2)
	})

	return result
}

// timeOfDay is a struct that represents the time of day plot.
type timeOfDay struct {
	bar *bar
}

// NewTimeOfDay creates a new Time of day plot.
func NewTimeOfDay(
	logs *data.Logs,
	config *Config,
) *timeOfDay {
	return &timeOfDay{
		bar: newBar(
			logs,
			"time_of_day",
			config,
		),
	}
}

// Plot generates the time of day chart.
func (t *timeOfDay) Plot() error {
	logs, err := data.Filter(
		t.bar.logs,
		data.WithDate(t.bar.startDate, t.bar.endDate),
		data.WithAuthors(t.bar.authors),
		data.WithMergeAuthors(t.bar.authors),
	)
	if err != nil {
		return err
	}

	timeLabel := t.createLabels(logs)
	formattedData := t.bar.generateDataMap(
		func(c *chart[*charts.Bar], l *data.Log) string {
			return t.getTimeOfDay(l.GetDate().AsTime().Hour())
		},
		func(_ *chart[*charts.Bar], l *data.Log) dataValueMap {
			return dataValueMap{l.GetAuthor(): l.GetPlus()}
		},
	)

	t.bar.setGlobalOptions("Time of Day Report")
	bar := t.bar.renderer

	bar.SetXAxis(timeLabel)
	t.bar.generateSeries(timeLabel, formattedData)

	return t.bar.save()
}

// createLabels creates the labels for the time of day plot.
func (t *timeOfDay) createLabels(logs *data.Logs) []string {
	return []string{
		"Midnight",
		"Morning",
		"Afternoon",
		"Night",
	}
}

// getTimeOfDay receives a hour and return a time of day.
func (t *timeOfDay) getTimeOfDay(hour int) string {
	if hour >= 0 && hour < 6 {
		return "Midnight"
	} else if hour >= 6 && hour < 12 {
		return "Morning"
	} else if hour >= 12 && hour < 19 {
		return "Afternoon"
	} else if hour >= 19 && hour < 24 {
		return "Night"
	} else {
		return "Midnight"
	}
}

// topLanguages is a struct that represents the top languages plot.
type topLanguages struct {
	bar *bar
}

// NewTopLanguages creates a new Top languages plot.
func NewTopLanguages(
	logs *data.Logs,
	config *Config,
) *topLanguages {
	return &topLanguages{
		bar: newBar(
			logs,
			"top_languages",
			config,
		),
	}
}

// Plot generates the top languages chart.
func (t *topLanguages) Plot() error {
	logs, err := data.Filter(
		t.bar.logs,
		data.WithDate(t.bar.startDate, t.bar.endDate),
		data.WithAuthors(t.bar.authors),
		data.WithMergeAuthors(t.bar.authors),
	)
	if err != nil {
		return err
	}

	languageLabel := t.createLabels(logs)
	formattedData := t.bar.generateDataMap(
		func(c *chart[*charts.Bar], l *data.Log) string {
			return t.identifyLanguage(l.GetPath())
		},
		func(_ *chart[*charts.Bar], l *data.Log) dataValueMap {
			return dataValueMap{
				l.GetAuthor(): l.GetPlus(),
			}
		},
	)

	t.bar.setGlobalOptions("Top Languages Report")
	bar := t.bar.renderer

	bar.SetXAxis(languageLabel)
	t.bar.generateSeries(languageLabel, formattedData)

	return t.bar.save()
}

// languages returns a map of languages and its names.
func (t *topLanguages) languages() map[string]string {
	return map[string]string{
		"go":     "Go",
		"py":     "Python",
		"js":     "JavaScript",
		"ts":     "TypeScript",
		"rs":     "Rust",
		"html":   "HTML",
		"css":    "CSS",
		"sh":     "Shell",
		"sql":    "SQL",
		"c":      "C",
		"cpp":    "C++",
		"h":      "C header",
		"hpp":    "C++ header",
		"java":   "Java",
		"cs":     "C#",
		"rb":     "Ruby",
		"php":    "PHP",
		"pl":     "Perl",
		"m":      "Objective-C",
		"swift":  "Swift",
		"kt":     "Kotlin",
		"lua":    "Lua",
		"r":      "R",
		"f":      "Fortran",
		"f90":    "Fortran 90",
		"p":      "Pascal",
		"pas":    "Pascal",
		"jl":     "Julia",
		"dart":   "Dart",
		"scala":  "Scala",
		"groovy": "Groovy",
		"clj":    "Clojure",
		"cljs":   "ClojureScript",
		"el":     "Emacs Lisp",
		"hs":     "Haskell",
		"asm":    "Assembly",
		"erl":    "Erlang",
		"ex":     "Elixir",
		"cob":    "COBOL",
		"vb":     "Visual Basic",
		"fs":     "F#",
		"ml":     "OCaml",
		"pro":    "Prolog",
		"tcl":    "Tcl",
	}
}

// identifyLanguage receives a path and returns the language.
func (t *topLanguages) identifyLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) > 1 {
		if lang, ok := t.languages()[ext[1:]]; ok {
			return lang
		}
	}
	return "Others"
}

// createLabels creates the labels for the language plot.
func (t *topLanguages) createLabels(logs *data.Logs) []string {
	var result []string
	for _, log := range logs.Logs {
		lang := t.identifyLanguage(log.GetPath())
		if !contains(result, lang) {
			result = append(result, lang)
		}
	}

	sort.Strings(result)
	return result
}

// weekday is a struct that represents the weekday plot.
type weekday struct {
	bar *bar
}

// NewWeekday creates a new weekday plot.
func NewWeekday(
	logs *data.Logs,
	config *Config,
) *weekday {
	return &weekday{
		bar: newBar(
			logs,
			"weekday",
			config,
		),
	}
}

// Plot generates the weekday chart.
func (w *weekday) Plot() error {
	logs, err := data.Filter(
		w.bar.logs,
		data.WithDate(w.bar.startDate, w.bar.endDate),
		data.WithAuthors(w.bar.authors),
		data.WithMergeAuthors(w.bar.authors),
	)
	if err != nil {
		return err
	}

	weekdayLabel := w.createLabels(logs)
	formattedData := w.bar.generateDataMap(
		func(c *chart[*charts.Bar], l *data.Log) string {
			return w.identifyWeekday(l.GetDate().AsTime())
		},
		func(_ *chart[*charts.Bar], l *data.Log) dataValueMap {
			return dataValueMap{
				l.GetAuthor(): l.GetPlus(),
			}
		},
	)

	w.bar.setGlobalOptions("Weekday Report")
	bar := w.bar.renderer

	bar.SetXAxis(weekdayLabel)
	w.bar.generateSeries(weekdayLabel, formattedData)

	return w.bar.save()
}

// identifyWeekday receives a date and returns the weekday.
func (w *weekday) identifyWeekday(date time.Time) string {
	return date.Weekday().String()
}

// createLabels creates the labels for the weekday plot.
func (w *weekday) createLabels(logs *data.Logs) []string {
	return []string{
		time.Sunday.String(),
		time.Monday.String(),
		time.Tuesday.String(),
		time.Wednesday.String(),
		time.Thursday.String(),
		time.Friday.String(),
		time.Saturday.String(),
	}
}

// contains is a helper function to check if a slice contains a string.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
