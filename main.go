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
)

type GitLog struct {
	Date string `json:"date"`
	Plus int    `json:"plus"`
}

func main() {
	var allLogs []*GitLog
	var allRawLogs []string

	startingDir := "."

	err := filepath.WalkDir(startingDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			parentDir := filepath.Dir(path)
			fmt.Println("Found a .git repository:", parentDir)

			logs, rawLogs, err := getGitLogs(parentDir)
			if err != nil {
				return err
			}
			allLogs = append(allLogs, logs...)
			allRawLogs = append(allRawLogs, rawLogs...)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Walking the directory tree failed: %v\n", err)
		return
	}

	jsonData, err := json.Marshal(allLogs)
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
		return
	}

	if err := writeJSONToFile("output.json", jsonData); err != nil {
		fmt.Printf("Writing to JSON file failed: %s", err)
	}

	if err := os.WriteFile("output.txt", []byte(strings.Join(allRawLogs, "\n")), 0o600); err != nil {
		fmt.Printf("Writing to text file failed: %s", err)
	}
}

func getGitLogs(repoPath string) ([]*GitLog, []string, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not convert to absolute path: %v", err)
	}

	args := []string{
		"-C", absRepoPath, "log",
		"--author=Christian",
		"--pretty=format:'%ad,%s'",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
		":(exclude)node_modules/*",
		":(exclude)*/node_modules/*",
		":(exclude)vendor/*",
		":(exclude)*/vendor/*",
		":(exclude)go.sum",
		":(exclude)yarn.lock",
		":(exclude)package-lock.json",
		":(exclude)pnp-lock.yaml",
		":(exclude)pnpm-lock.yaml",
		":(exclude)dist/*",
		":(exclude)*/dist/*",
		":(exclude)build/*",
		":(exclude)*/build/*",
		":(exclude).git/*",
		":(exclude).idea/*",
		":(exclude)mocks/*",
		":(exclude)*/mocks/*",
		":(exclude)*.csv",
		":(exclude)*.pdf",
		":(exclude)*.doc",
		":(exclude)*.docx",
		":(exclude)*.json",
		":(exclude)*.png",
		":(exclude)*.jpg",
		":(exclude)*.jpeg",
		":(exclude)*.gif",
		":(exclude)*.svg",
		":(exclude)*.ico",
		":(exclude)*.woff",
		":(exclude)*.woff2",
		":(exclude)*.ttf",
		":(exclude)*.eot",
	}

	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Command failed with error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
		return nil, nil, err
	}

	logs := make([]*GitLog, 0)
	lines := strings.Split(stdout.String(), "\n")
	plusRegex := regexp.MustCompile(`(\d+)\t(\d+)\t`)

	var currentGitLog *GitLog

	for _, line := range lines {
		if strings.Contains(line, ",") {
			if currentGitLog != nil {
				logs = append(logs, currentGitLog)
			}
			dateRegex := regexp.MustCompile(`'(\d{4}-\d{2}-\d{2} \d{2}:\d{2})'`)
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

	return logs, lines, nil
}

func writeJSONToFile(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return err
	}
	return nil
}
