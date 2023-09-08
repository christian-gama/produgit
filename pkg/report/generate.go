package report

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/christian-gama/productivity/pkg/gitlog"
)

func Generate(config *Config) {
	logs := processDir(config)
	saveReport(config, logs)
}

func processDir(config *Config) []*gitlog.Log {
	var logs []*gitlog.Log
	var rawLogs []string

	err := filepath.WalkDir(
		config.Dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if strings.Contains(err.Error(), "permission denied") {
					return nil
				}

				return err
			}

			if d.IsDir() && d.Name() == ".git" {
				fmt.Printf("Processing %s\n", strings.TrimSuffix(path, "/.git"))
				rawLogs = gitlog.GetLog(filepath.Dir(path), &gitlog.Config{
					Exclude: config.Exclude,
				})
				if err != nil {
					return err
				}

				logs = append(logs, gitlog.Parse(rawLogs)...)

				return filepath.SkipDir
			}

			return nil
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Walking the directory tree failed: %s\n", err))
	}

	return logs
}

func saveReport(config *Config, logs []*gitlog.Log) {
	data, err := json.Marshal(logs)
	if err != nil {
		panic(fmt.Sprintf("JSON marshaling failed: %s\n", err))
	}

	if _, err := os.Stat(config.Output); err == nil {
		err = os.Remove(config.Output)
		if err != nil {
			panic(fmt.Sprintf("Removing existing JSON file failed: %s\n", err))
		}
	}

	err = os.WriteFile(config.Output, data, 0600)
	if err != nil {
		panic(fmt.Sprintf("Writing to JSON file failed: %s\n", err))
	}
}
