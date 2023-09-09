package report

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	Dir     []string
	Authors []string
	Output  string
	Exclude []string
}

func NewConfig(dir []string, authors []string, output string, exclude []string) *Config {
	if len(dir) == 0 {
		dir = []string{"."}
	}

	if output == "" {
		output = filepath.Join(".", "produgit_report.json")
	} else if filepath.Ext(output) != ".json" {
		output = fmt.Sprintf("%s.json", output)
	}

	return &Config{
		Dir:     dir,
		Authors: authors,
		Output:  output,
		Exclude: exclude,
	}
}
