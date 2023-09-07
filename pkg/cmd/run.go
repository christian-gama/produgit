package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Run(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil || stderr.Len() > 0 {
		return "", fmt.Errorf("Command failed with error: %s\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}
