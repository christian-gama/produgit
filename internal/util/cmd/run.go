package cmdutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunAndWait executes a command and returns the output as a string.
func RunAndWait(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("Command failed with error: %v", err)
	}

	return stdout.String(), nil
}

// Run executes a command and returns the command object.
func Run(cmd string, args ...string) (*exec.Cmd, error) {
	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("Command failed with error: %v", err)
	}

	return command, nil
}
