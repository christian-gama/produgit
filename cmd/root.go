package main

import (
	"fmt"
	"os"

	"github.com/christian-gama/productivity/cmd/plot"
	"github.com/christian-gama/productivity/cmd/report"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "produgit", Short: "A productivity tool for git"}

func init() {
	rootCmd.AddCommand(plot.PlotCmd)
	rootCmd.AddCommand(report.ReportCmd)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
			os.Exit(1)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
}
