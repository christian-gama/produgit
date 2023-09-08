package plot

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/go-echarts/go-echarts/v2/charts"
)

func TopLanguages(config *Config, filterConfig *FilterConfig) {
	data := report.Read(config.Input)
	data = filter(filterConfig, data)
	langLabel := uniqueLanguages(accumulateLogsByLanguage(data))
	accumulated := accumulateLogsByLanguage(data)

	bar := charts.NewBar()
	createBarHeader("Top Languages Report", bar, filterConfig)

	bar.SetXAxis(langLabel)
	addSeries(bar, langLabel, filterConfig.Authors, accumulated)

	save("top-languages", config, filterConfig, bar)
}

var languageMap = map[string]string{
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

func identifyLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) > 1 {
		if lang, ok := languageMap[ext[1:]]; ok {
			return lang
		}
	}
	return "Others"
}

func accumulateLogsByLanguage(logs []*gitlog.Log) map[string][]*Data {
	accumulatedLogs := make(map[string][]*Data)

	for _, log := range logs {
		lang := identifyLanguage(log.Path)
		if _, ok := accumulatedLogs[lang]; !ok {
			accumulatedLogs[lang] = []*Data{}
		}

		var found *Data
		for _, aLog := range accumulatedLogs[lang] {
			if aLog.Author == log.Author {
				found = aLog
				break
			}
		}

		if found == nil {
			newLog := &Data{
				Author:   log.Author,
				Plus:     log.Plus,
				Language: lang,
			}
			accumulatedLogs[lang] = append(accumulatedLogs[lang], newLog)
		} else {
			found.Plus += log.Plus
		}
	}

	return accumulatedLogs
}

func uniqueLanguages(logs map[string][]*Data) []string {
	var result []string
	for lang := range logs {
		result = append(result, lang)
	}

	sort.Strings(result)
	return result
}
