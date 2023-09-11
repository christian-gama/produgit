package report

import (
	"github.com/christian-gama/produgit/config"
	"github.com/christian-gama/produgit/internal/report"
	"github.com/spf13/cobra"
)

var (
	dir     []string
	output  string
	exclude []string
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your produgit using a specialized git log",
	ValidArgs: []string{
		"--dir",
		"-d",
		"--output",
		"-o",
		"--exclude",
		"-e",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		report := report.NewReport(dir, exclude, output)
		return report.Generate()
	},
}

func init() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	ReportCmd.
		Flags().
		StringArrayVarP(&dir, "dir", "d", []string{"."}, "The starting directory to search for .git repositories")

	ReportCmd.
		Flags().
		StringVarP(&output, "output", "o", config.Report.Output, "The output path for the report")

	ReportCmd.
		Flags().
		StringArrayVarP(&exclude, "exclude", "e", config.Report.Exclude, "The directories to exclude from the report")
}
