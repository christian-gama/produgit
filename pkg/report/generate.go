package report

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/christian-gama/produgit/pkg/gitlog"
)

func Generate(config *Config) {
	logs := processDir(config)
	saveReport(config, logs)
}

func processDir(config *Config) []*gitlog.Log {
	var logs []*gitlog.Log

	// Define a channel to collect logs from the goroutines
	results := make(chan []*gitlog.Log, len(config.Dir))

	// Define a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, dir := range config.Dir {
		err := filepath.WalkDir(
			dir,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					if strings.Contains(err.Error(), "permission denied") {
						return nil
					}

					return err
				}

				if d.IsDir() && d.Name() == ".git" {
					wg.Add(1) // Increment the wait group counter

					go func(path string) { // Start a goroutine
						defer wg.Done()

						var localLogs []*gitlog.Log
						fmt.Printf("Processing %s\n", strings.TrimSuffix(path, "/.git"))
						rawLogs := gitlog.GetLog(filepath.Dir(path), &gitlog.Config{
							Exclude: config.Exclude,
						})

						localLogs = append(localLogs, gitlog.Parse(rawLogs)...)

						results <- localLogs // Send the results to the channel
					}(path)

					return filepath.SkipDir
				}

				return nil
			},
		)
		if err != nil {
			panic(fmt.Sprintf("Walking the directory tree failed: %s\n", err))
		}
	}

	go func() { // Start another goroutine to close the results channel after all goroutines are done
		wg.Wait()
		close(results)
	}()

	for localLogs := range results {
		logs = append(logs, localLogs...) // Aggregate results from the channel
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
