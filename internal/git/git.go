package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	cmdutil "github.com/christian-gama/produgit/internal/util/cmd"
)

// GetLog returns the git log for the given repoPath
func GetLog(repoPath string, exclude []string) ([]string, error) {
	if err := checkGitExists(); err != nil {
		return nil, err
	}

	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("Could not convert to absolute path: %s", err)
	}

	args := []string{
		"-C", absRepoPath, "log",
		"--pretty=format:%ad,%ae,%an",
		"--date=format:'%Y-%m-%d %H:%M'",
		"--numstat",
		"--",
		".",
	}

	args = appendExcludeArgs(args, exclude)

	output, err := cmdutil.RunAndWait("git", args...)
	if err != nil && !strings.Contains(err.Error(), "does not have any commits yet") {
		return nil, fmt.Errorf("Could not run git log: %s", err)
	}

	return strings.Split(output, "\n"), nil
}

// ListAllAuthors returns a list of all authors in the given repoPath.
func ListAllAuthors(repoPath string) ([]string, error) {
	if err := checkGitExists(); err != nil {
		return nil, err
	}

	args := []string{
		"-C", repoPath, "log",
		"--format=%an (%ae)",
	}

	// Running git command
	output, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("could not run git log: %s", err)
	}

	// Split, sort, and deduplicate the result
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	authorsMap := make(map[string]bool)

	for _, line := range lines {
		line = formatAuthor(line)
		authorsMap[line] = true
	}

	// Convert the map keys to a slice
	var authors []string
	for author := range authorsMap {
		authors = append(authors, author)
	}

	return authors, nil
}

func formatAuthor(input string) string {
	name := strings.Split(input, " (")[0]
	email := strings.TrimSuffix(strings.Split(input, " (")[1], ")")

	if name == "" {
		name = "Unknown Name"
	}

	if email == "" {
		email = "Unknown Email"
	}

	return fmt.Sprintf("%s (%s)", name, email)
}

func appendExcludeArgs(args []string, exclude []string) []string {
	excludeArgs := make([]string, 0, len(exclude))
	for _, item := range exclude {
		if item != "" {
			excludeArgs = append(excludeArgs, fmt.Sprintf(":(exclude)%s", item))
		}
	}

	return append(args, excludeArgs...)
}

// checkGitExists checks if git is installed.
func checkGitExists() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("Could not find git in your PATH: %s", err)
	}
	return nil
}
