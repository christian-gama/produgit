package gitlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestGetLog(t *testing.T) {
	config := &Config{
		Authors: []string{"Christian"},
		Exclude: []string{},
	}
	repoPath := getGitDir()

	logs := GetLog(repoPath, config)

	if logs == nil {
		t.Errorf("Logs are nil")
	}

	if len(logs) == 0 {
		t.Errorf("Logs are empty")
	}
}

func TestGetLog_NoCommitsYet(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "gitlogtest")
	if err != nil {
		t.Fatal("Could not create temporary directory:", err)
	}

	cmd := exec.Command("git", "init", tmpDir)
	_, err = cmd.Output()
	if err != nil {
		t.Fatal("Failed to initialize git repository:", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetLog panicked: %s", r)
		}
	}()

	config := &Config{
		Authors: []string{"Christian"},
		Exclude: []string{},
	}

	GetLog(tmpDir, config)
}

func TestGetLog_NoGit(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "gitlogtest")
	if err != nil {
		t.Fatal("Could not create temporary directory:", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetLog did not panic")
		}
	}()

	config := &Config{
		Authors: []string{"Christian"},
		Exclude: []string{},
	}

	GetLog(tmpDir, config)
}

func TestGetLog_UseGlobalAuthor(t *testing.T) {
	config := &Config{}
	checkAuthor(config)

	if len(config.Authors) != 1 {
		t.Errorf("Expected 1 author, got %d", len(config.Authors))
	}
}

func getGitDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Could not get working directory: %s", err))
	}

	return getGitDirRecursive(wd)
}

func getGitDirRecursive(dir string) string {
	gitDir := fmt.Sprintf("%s/.git", dir)
	if _, err := os.Stat(gitDir); err == nil {
		return gitDir
	}

	parentDir := fmt.Sprintf("%s/..", dir)
	if parentDir == dir {
		panic("Could not find .git directory")
	}

	return getGitDirRecursive(parentDir)
}
