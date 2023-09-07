package gitlog

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/christian-gama/productivity/pkg/cmd"
)

type Config struct {
	Authors []string
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
	checkAuthor(config)

	args := createBaseArgs(repoPath, config)
	args = append(args, createExcludeArgs(excludeItems)...)

	cmdOut, err := cmd.Run("git", args...)
	if err != nil && !strings.Contains(err.Error(), "does not have any commits yet") {
		panic(fmt.Sprintf("Could not run git log: %s", err))
	}

	return strings.Split(cmdOut, "\n")
}

func createExcludeArgs(excludeItems []string) []string {
	excludeArgs := make([]string, len(excludeItems))
	for i, item := range excludeItems {
		excludeArgs[i] = fmt.Sprintf(":(exclude)%s", item)
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
		fmt.Sprintf("--author=%s", strings.Join(config.Authors, "\\|")),
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

func checkAuthor(config *Config) {
	if len(config.Authors) == 0 {
		cmdOut, err := cmd.Run("git", "config", "--global", "user.email")
		if err != nil {
			panic(fmt.Sprintf("Could not run git config: %s", err))
		}

		if cmdOut == "" {
			cmdOut, err = cmd.Run("git", "config", "--global", "user.name")
			if err != nil {
				panic(fmt.Sprintf("Could not run git config: %s", err))
			}
		}

		config.Authors = []string{strings.TrimSpace(cmdOut)}
	}

	if len(config.Authors) == 0 {
		panic("No author provided")
	}
}
