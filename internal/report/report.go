package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/christian-gama/produgit/internal/data"
	"github.com/christian-gama/produgit/internal/git"
	"github.com/christian-gama/produgit/internal/logger"
	"google.golang.org/protobuf/proto"
)

// Report is the configuration for the report command.
type Report struct {
	Dir     []string
	Exclude []string
	Output  string
}

// NewReport creates a new Report.
func NewReport(dir []string, exclude []string, output string) *Report {
	return &Report{
		Dir:     dir,
		Exclude: exclude,
		Output:  output,
	}
}

// Generate generates the report.
func (r *Report) Generate() error {
	logs, err := r.processDir()
	if err != nil {
		return fmt.Errorf("Processing directory failed: %w", err)
	}

	err = r.save(&data.Logs{Logs: logs})
	if err != nil {
		return fmt.Errorf("Saving report failed: %w", err)
	}

	return nil
}

// processGitDir processes a .git directory to fetch and parse logs.
func (r *Report) processGitDir(
	path string,
	results chan []*data.Log,
	errs chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	logger.Print("Processing %s", strings.TrimSuffix(path, "/.git"))

	var localLogs []*data.Log

	rawLogs, err := git.GetLog(filepath.Dir(path), r.Exclude)
	if err != nil {
		errs <- fmt.Errorf("Getting logs failed: %w", err)
		return
	}

	parsedLogs, err := data.Parse(rawLogs)
	if err != nil {
		errs <- fmt.Errorf("Parsing logs failed: %w", err)
		return
	}

	localLogs = append(localLogs, parsedLogs...)
	results <- localLogs
}

func (r *Report) processDir() ([]*data.Log, error) {
	var logs []*data.Log
	results := make(chan []*data.Log, len(r.Dir))
	errs := make(chan error, len(r.Dir))
	var wg sync.WaitGroup

	callback := func(path string) error {
		wg.Add(1)
		go r.processGitDir(path, results, errs, &wg)
		return nil
	}

	err := git.WalkDirs(r.Dir, callback)
	if err != nil {
		return nil, err
	}

	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	if len(errs) > 0 {
		var collectedErrors []string
		for err := range errs {
			collectedErrors = append(collectedErrors, err.Error())
		}
		return nil, fmt.Errorf("Errors encountered: %s", strings.Join(collectedErrors, "; "))
	}

	for localLogs := range results {
		logs = append(logs, localLogs...)
	}

	return logs, nil
}

// save saves the report.
func (r *Report) save(logs *data.Logs) error {
	data, err := proto.Marshal(logs)
	if err != nil {
		return fmt.Errorf("Marshaling logs failed: %w", err)
	}

	if _, err := os.Stat(r.Output); err == nil {
		err = os.Remove(r.Output)
		if err != nil {
			return fmt.Errorf("Removing file failed: %w", err)
		}
	}

	err = os.WriteFile(r.Output, data, 0600)
	if err != nil {
		return fmt.Errorf("Writing file failed: %w", err)
	}

	return nil
}
