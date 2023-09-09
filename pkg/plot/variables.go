package plot

import "time"

var (
	now        = time.Now()
	startOfDay = getStartOfDay(now)
)

var periods = map[string][2]time.Time{
	"today":      {startOfDay, now},
	"24h":        {now.Add(-24 * time.Hour), now},
	"this-week":  {startOfDay.Add(-7 * 24 * time.Hour), now},
	"this-month": {time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), now},
	"this-year":  {time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), now},
}

var daysOfWeek = []string{
	time.Sunday.String(),
	time.Monday.String(),
	time.Tuesday.String(),
	time.Wednesday.String(),
	time.Thursday.String(),
	time.Friday.String(),
	time.Saturday.String(),
}

var languages = map[string]string{
	"go":   "Go",
	"py":   "Python",
	"js":   "JavaScript",
	"ts":   "TypeScript",
	"rs":   "Rust",
	"html": "HTML",
	"css":  "CSS",
	"sh":   "Shell",
	"sql":  "SQL",
	"c":    "C",
	"cpp":  "C++",
	"h":    "C header",
	"hpp":  "C++ header",
	"java": "Java",
	"cs":   "C#", "rb": "Ruby",
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

var timeOfDay = []string{
	"Midnight",
	"Morning",
	"Afternoon",
	"Night",
}
