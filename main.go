package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "productivity", Short: "A productivity tool"}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your productivity",
	Run:   runReport,
	ValidArgs: []string{
		"--author",
		"--starting-dir",
	},
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	Run:   runPlot,
	ValidArgs: []string{
		"--output-file",
		"--extension",
		"--type",
	},
}

func init() {
	configReportCmd()
	configPlotCmd()
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(plotCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
}

func configReportCmd() {
	reportCmd.
		Flags().
		StringVarP(&config.StartingDir, "starting-dir", "d", ".", "The starting directory to search for .git repositories")

	reportCmd.
		Flags().
		StringVarP(&config.Author, "author", "a", "", "The author to filter git logs")

	err := reportCmd.MarkFlagRequired("author")
	if err != nil {
		panic(fmt.Sprintf("Marking flag as required failed: %s\n", err))
	}
}

func configPlotCmd() {
	plotCmd.
		Flags().
		StringVarP(&plotConfig.OutputFile, "output-file", "o", "chart", "The input file to plot. The default is 'chart'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.Extension, "extension", "e", "png", "The extension of the output file. The default is 'png'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.Type, "type", "t", "monthly", "The type of the plot. The default is 'monthly'.")

	if err := plotCmd.MarkFlagRequired("type"); err != nil {
		panic(fmt.Sprintf("Marking flag as required failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"extension",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"png",
				"svg",
				"pdf",
				"eps",
				"tif",
				"tiff",
				"jpg",
				"jpeg",
			}, cobra.ShellCompDirectiveDefault
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"type",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"monthly",
				"weekday",
				"time",
			}, cobra.ShellCompDirectiveDefault
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}
}
