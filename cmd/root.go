package main

import (
	"fmt"
	"os"

	"github.com/christian-gama/produgit/cmd/anomalies"
	"github.com/christian-gama/produgit/cmd/config"
	"github.com/christian-gama/produgit/cmd/plot"
	"github.com/christian-gama/produgit/cmd/report"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "produgit", Short: "A produgit tool for git"}

func init() {
	rootCmd.AddCommand(plot.PlotCmd)
	rootCmd.AddCommand(report.ReportCmd)
	rootCmd.AddCommand(anomalies.AnomaliesCmd)
	rootCmd.AddCommand(config.ConfigCmd)
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
