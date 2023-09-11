package gitlog

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/christian-gama/produgit/pkg/cmd"
)

type Config struct {
	Exclude []string
}

var excludeItems = []string{
	"**node_modules/*",
	"**vendor/*",
	"**go.sum",
	"**go.mod",
	"**yarn.lock",
	"**package-lock.json",
	"**requirements.txt",
	"**venv/*",
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

func GetLog(repoPath string, config *Config) []string {
	checkGitExists()

	args := createBaseArgs(repoPath, config)
	args = append(args, createExcludeArgs(config.Exclude)...)

	cmdOut, err := cmd.Run("git", args...)
	if err != nil && !strings.Contains(err.Error(), "does not have any commits yet") {
		panic(fmt.Sprintf("Could not run git log: %s", err))
	}

	return strings.Split(cmdOut, "\n")
}

func createExcludeArgs(exclude []string) []string {
	excludeArgs := make([]string, 0, len(exclude)+len(excludeItems))
	for _, item := range exclude {
		if item != "" {
			excludeArgs = append(excludeArgs, fmt.Sprintf(":(exclude)%s", item))
		}
	}

	for _, item := range excludeItems {
		if item != "" {
			excludeArgs = append(excludeArgs, fmt.Sprintf(":(exclude)%s", item))
		}
	}

	return excludeArgs
}

func createBaseArgs(repoPath string, config *Config) []string {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		panic(fmt.Sprintf("Could not convert to absolute path: %s", err))
	}

	return []string{
		"-C", absRepoPath, "log",
		"--pretty=format:%ad,%ae,%an",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
	}
}

func checkGitExists() {
	if _, err := exec.LookPath("git"); err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}
}
