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

	toml "github.com/pelletier/go-toml"
)

type GitLog struct {
	Date string `json:"date"`
	Plus int    `json:"plus"`
}

type Config struct {
	StartingDir string   `toml:"starting_dir"`
	Author      string   `toml:"author"`
	Excludes    []string `toml:"excludes"`
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	config, err := loadConfig("config.toml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		return
	}

	var allLogs []*GitLog

	err = filepath.WalkDir(config.StartingDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			parentDir := filepath.Dir(path)
			fmt.Println("Found a .git repository:", parentDir)

			logs, err := getGitLogs(parentDir, config)
			if err != nil {
				return err
			}
			allLogs = append(allLogs, logs...)
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
}

func getGitLogs(repoPath string, config *Config) ([]*GitLog, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("Could not convert to absolute path: %v", err)
	}

	args := []string{
		"-C", absRepoPath, "log",
		"--author=" + config.Author,
		"--pretty=format:'%ad,%s'",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
	}

	// Add excludes from the TOML config
	for _, exclude := range config.Excludes {
		args = append(args, ":(exclude)"+exclude)
	}

	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Command failed with error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
		return nil, err
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

	return logs, nil
}

func writeJSONToFile(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return err
	}
	return nil
}
