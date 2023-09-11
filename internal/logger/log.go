package logger

import (
	"fmt"

	"github.com/christian-gama/produgit/config"
)

// Print prints the given message to stdout if the quiet flag is not set.
func Print(msg string, args ...interface{}) {
	if !config.Config.Quiet {
		fmt.Printf(fmt.Sprintf("%s\n", msg), args...)
	}
}

// Warn prints the given message to stdout if the quiet flag is not set.
func Warn(msg string, args ...interface{}) {
	if !config.Config.Quiet {
		fmt.Printf(fmt.Sprintf("Warning: %s\n", msg), args...)
	}
}
