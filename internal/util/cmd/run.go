package cmdutil

import (
	"fmt"
	"os"
	"os/exec"
)

// Run executes a command and returns the command object
func Run(cmd string, args ...string) (*exec.Cmd, error) {
	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("Command failed with error: %v", err)
	}

	return command, nil
}
