package main

import (
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your productivity",
	Run:   runReport,
	ValidArgs: []string{
		"--author",
		"--dir",
		"--output",
		"--exclude",
	},
}

var (
	authors []string
	dir     string
	output  string
	exclude []string
)

func init() {
	reportCmd.
		Flags().
		StringVar(&dir, "dir", "", "The starting directory to search for .git repositories")

	reportCmd.
		Flags().
		StringArrayVar(&authors, "author", []string{}, "The author to filter git logs")

	reportCmd.
		Flags().
		StringVar(&output, "output", "", "The output path for the report")

	reportCmd.
		Flags().
		StringArrayVar(&exclude, "exclude", []string{}, "The directories to exclude from the report")
}

func runReport(cmd *cobra.Command, args []string) {
	report.Generate(
		report.NewConfig(
			dir,
			authors,
			output,
			exclude,
		),
	)
}
