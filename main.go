package main

import (
	"fmt"
	"os"

	"github.com/christian-gama/produgit/cmd/anomaly"
	"github.com/christian-gama/produgit/cmd/config"
	"github.com/christian-gama/produgit/cmd/list"
	"github.com/christian-gama/produgit/cmd/plot"
	"github.com/christian-gama/produgit/cmd/report"
	"github.com/spf13/cobra"

	appconfig "github.com/christian-gama/produgit/config"
)

var rootCmd = &cobra.Command{
	Use:   "produgit",
	Short: "A tool to help you understand your git repositories",
}

func init() {
	err := appconfig.Load()
	if err != nil {
		panic(err)
	}

	plot.Init()
	report.Init()
	config.Init()
	list.Init()
	anomaly.Init()

	rootCmd.AddCommand(plot.PlotCmd)
	rootCmd.AddCommand(report.ReportCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(anomaly.AnomalyCmd)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
			os.Exit(1)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
