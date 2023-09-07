package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "produgit", Short: "A productivity tool for git"}

func init() {
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(plotCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
}
