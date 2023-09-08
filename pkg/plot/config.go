package plot

import (
	"fmt"
	"os"
)

type Config struct {
	Output string
	Input  string
}

func NewConfig(input, output string) *Config {
	if output == "" {
		output = "./"
	} else {
		if isFile(output) {
			panic("Output should be a directory")
		}

		output = fmt.Sprintf("%s/", output)

	}

	if input == "" {
		input = "produgit_report.json"
	}

	return &Config{
		Output: output,
		Input:  input,
	}
}

func isFile(output string) bool {
	info, err := os.Stat(output)
	if err != nil {
		panic(err)
	}

	return !info.IsDir()
}
