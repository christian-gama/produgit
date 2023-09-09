package report

import (
	"github.com/christian-gama/productivity/pkg/report"
	"github.com/spf13/cobra"
)

var (
	authors []string
	dir     []string
	output  string
	exclude []string
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your productivity using a specialized git log",
	ValidArgs: []string{
		"--dir",
		"--output",
		"--exclude",
		"-d",
		"-o",
		"-e",
	},
	Run: func(cmd *cobra.Command, args []string) {
		report.Generate(
			report.NewConfig(
				dir,
				authors,
				output,
				exclude,
			),
		)
	},
}

func init() {
	ReportCmd.
		Flags().
		StringArrayVarP(&dir, "dir", "d", []string{}, "The starting directory to search for .git repositories")

	ReportCmd.
		Flags().
		StringVarP(&output, "output", "o", "", "The output path for the report")

	ReportCmd.
		Flags().
		StringArrayVarP(&exclude, "exclude", "e", []string{}, "The directories to exclude from the report")
}
