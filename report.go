package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type GitLog struct {
	Date string `json:"date"`
	Plus int    `json:"plus"`
}

type ReportConfig struct {
	StartingDir string
	Authors     []string
	OutputDir   string
	Exclude     []string
	Debug       bool
}

var config ReportConfig

func processDir(path string, d fs.DirEntry, formattedLogs *[]*GitLog, rawLogs *[]string) error {
	if d.IsDir() && d.Name() == ".git" {
		parentDir := filepath.Dir(path)
		fmt.Println("Processing directory: ", parentDir)

		logs, lines := getGitLogs(parentDir)

		*formattedLogs = append(*formattedLogs, logs...)
		*rawLogs = append(*rawLogs, lines...)

		return filepath.SkipDir
	}

	return nil
}

func runReport(cmd *cobra.Command, args []string) {
	if len(config.Authors) == 0 {
		if err := cmd.Help(); err != nil {
			panic(fmt.Sprintf("Error: %s", err))
		}
		panic("\nError: required flag(s) \"author\" not set")
	}

	var formattedLogs []*GitLog
	var rawLogs []string

	err := filepath.WalkDir(
		config.StartingDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if strings.Contains(err.Error(), "permission denied") {
					return nil
				}

				return err
			}
			return processDir(path, d, &formattedLogs, &rawLogs)
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Walking the directory tree failed: %s\n", err))
	}

	jsonData, err := json.Marshal(formattedLogs)
	if err != nil {
		panic(fmt.Sprintf("JSON marshaling failed: %s\n", err))
	}

	if config.OutputDir == "" {
		config.OutputDir = config.StartingDir
	}

	if _, err := os.Stat(filepath.Join(config.OutputDir, "productivity_report_output.json")); err == nil {
		err = os.Remove(filepath.Join(config.OutputDir, "productivity_report_output.json"))
		if err != nil {
			panic(fmt.Sprintf("Removing existing JSON file failed: %s\n", err))
		}
	}

	err = os.WriteFile(
		filepath.Join(config.OutputDir, "productivity_report_output.json"),
		jsonData,
		0600,
	)
	if err != nil {
		panic(fmt.Sprintf("Writing to JSON file failed: %s\n", err))
	}

	if config.Debug {
		if _, err := os.Stat(filepath.Join(config.OutputDir, "productivity_report_output_debug.txt")); err == nil {
			err = os.Remove(filepath.Join(config.OutputDir, "productivity_report_output_debug.txt"))
			if err != nil {
				panic(fmt.Sprintf("Removing existing debug file failed: %s\n", err))
			}
		}

		err = os.WriteFile(
			filepath.Join(config.OutputDir, "productivity_report_output_debug.txt"),
			[]byte(strings.Join(rawLogs, "\n")),
			0600,
		)
		if err != nil {
			panic(fmt.Sprintf("Writing to debug file failed: %s\n", err))
		}
	}
}

func getGitLogs(repoPath string) ([]*GitLog, []string) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		panic(fmt.Sprintf("Could not convert to absolute path: %s", err))
	}

	baseArgs := []string{
		"-C", absRepoPath, "log",
		fmt.Sprintf("--author=%s", strings.Join(config.Authors, "\\|")),
		"--pretty=format:'%ad,%s'",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
	}

	excludeItems := []string{
		"**node_modules/*",
		"**vendor/*",
		"**go.sum",
		"**go.mod",
		"**yarn.lock",
		"**package-lock.json",
		"**requirements.txt",
		"**pnpm-lock.yaml",
		"**dist/*",
		"**build/*",
		"**assets/*",
		"**.git/*",
		"**.idea/*",
		"**mocks/*",
		"**.vscode/*",
		"**/.pytest_cache/*",
		"**.pytest_cache/*",
		"**.next/*",
		".next/*",
		"**.cache/*",
		"**__pycache__/*",
		"**coverage/*",
		"**coverage.xml",
		"**coverage.html",
		"**coverage.lcov",
		"**LICENSE.md",
		"*.csv",
		"*.pdf",
		"*.doc",
		"*.docx",
		"*.json",
		"*.png",
		"*.jpg",
		"*.jpeg",
		"*.gif",
		"*.svg",
		"*.ico",
		"*.woff",
		"*.woff2",
		"*.ttf",
		"*.eot",
		"*.txt",
		"**.DS_Store",
		"**Thumbs.db",
		"*.log",
		"*.bak",
		"*.swp",
		"*.swo",
		"*.tmp",
		"*.temp",
		"*.o",
		"*.out",
		"*.jar",
		"*.war",
		"*.ear",
		"*.sqlite3",
		"android/*",
		"ios/*",
	}

	excludeItems = append(excludeItems, config.Exclude...)

	// Prefix each item with ":(exclude)"
	excludeArgs := make([]string, len(excludeItems))
	for i, item := range excludeItems {
		excludeArgs[i] = fmt.Sprintf(":(exclude)%s", item)
	}

	args := append(baseArgs, excludeArgs...)
	cmdOut, stdErr, err := runCommand("git", args)
	if err != nil {
		if strings.Contains(stdErr, "does not have any commits yet") {
			return make([]*GitLog, 0), make([]string, 0)
		}

		panic(fmt.Sprintf("Command failed with error: %s\n%s", err, stdErr))
	}

	return parseGitLogs(cmdOut)
}

func runCommand(cmd string, args []string) (string, string, error) {
	command := exec.Command(cmd, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return "", stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}

var (
	dateRegex = regexp.MustCompile(`'(\d{4}-\d{2}-\d{2} \d{2}:\d{2})'`)
	plusRegex = regexp.MustCompile(`(\d+)\t(\d+)\t`)
)

func parseGitLogs(output string) ([]*GitLog, []string) {
	logs := make([]*GitLog, 0)
	lines := strings.Split(output, "\n")

	var currentGitLog *GitLog
	for _, line := range lines {
		if strings.Contains(line, ",") {
			if currentGitLog != nil {
				logs = append(logs, currentGitLog)
			}

			matches := dateRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				currentGitLog = &GitLog{
					Date: matches[1],
				}
			} else {
				currentGitLog = nil
			}
		} else if currentGitLog != nil {
			matches := plusRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				plus := 0
				fmt.Sscanf(matches[1], "%d", &plus)
				currentGitLog.Plus += plus
			}
		}
	}

	if currentGitLog != nil {
		logs = append(logs, currentGitLog)
	}

	return logs, lines
}
