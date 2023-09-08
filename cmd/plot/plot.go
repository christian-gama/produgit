package plot

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var PlotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
	ValidArgs: []string{
		"--input",
		"-i",
		"--start-date",
		"-s",
		"--end-date",
		"-e",
		"--author",
		"-a",
		"--period",
		"-p",
	},
}

var (
	output string
	input  string

	startDate string
	endDate   string
	authors   []string
	period    string
)

func init() {
	PlotCmd.AddCommand(monthlyCmd)
	PlotCmd.AddCommand(timeOfDay)
	PlotCmd.AddCommand(weekdayCmd)
	PlotCmd.AddCommand(topAuthors)
	PlotCmd.AddCommand(topLanguages)

	PlotCmd.
		PersistentFlags().
		StringVarP(&startDate, "start-date", "s", "", "Start date")

	PlotCmd.
		PersistentFlags().
		StringVarP(&endDate, "end-date", "e", "", "End date")

	PlotCmd.
		PersistentFlags().
		StringSliceVarP(&authors, "author", "a", []string{}, "Authors")

	PlotCmd.
		PersistentFlags().
		StringVarP(&output, "output", "o", "", "Output file")

	PlotCmd.
		PersistentFlags().
		StringVarP(&input, "input", "i", "", "Input file")

	PlotCmd.
		PersistentFlags().
		StringVarP(&period, "period", "p", "", "Period to plot")

	if err := PlotCmd.RegisterFlagCompletionFunc("period", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"today", "24h", "this-week", "this-month", "this-year"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(err)
	}
}

func parseDate(dateStr string, defaultDate time.Time) time.Time {
	if dateStr == "" {
		return defaultDate
	}

	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01",
		"2006",
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date
		}
	}

	panic(fmt.Sprintf("Invalid date '%s'. Expected one of %v", dateStr, formats))
}
