package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	cmdutil "github.com/christian-gama/produgit/internal/util/cmd"
)

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

func GetLog(repoPath string, exclude []string) ([]string, error) {
	if err := checkGitExists(); err != nil {
		return nil, err
	}

	args, err := createBaseArgs(repoPath, exclude)
	if err != nil {
		return nil, err
	}

	args = append(args, createExcludeArgs(exclude)...)

	cmd, err := cmdutil.Run("git", args...)
	if err != nil && !strings.Contains(err.Error(), "does not have any commits yet") {
		return nil, fmt.Errorf("Could not run git log: %s", err)
	}

	return strings.Split(cmd.String(), "\n"), nil
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

func createBaseArgs(repoPath string, exclude []string) ([]string, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("Could not convert to absolute path: %s", err)
	}

	return []string{
		"-C", absRepoPath, "log",
		"--pretty=format:%ad,%ae,%an",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
	}, nil
}

func checkGitExists() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("Error: %s\n", err)
	}
	return nil
}
