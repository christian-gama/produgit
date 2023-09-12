package plot

import (
	"github.com/christian-gama/produgit/config"
	"github.com/christian-gama/produgit/internal/data"
	"github.com/christian-gama/produgit/internal/plot"
	dateutil "github.com/christian-gama/produgit/internal/util/date"
	"github.com/spf13/cobra"
)

var (
	logs *data.Logs
	cfg  *plot.Config
)

var PlotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		logs, err = data.Load(input)
		if err != nil {
			return err
		}

		start, err := dateutil.ToTime(startDate)
		if err != nil {
			return err
		}

		end, err := dateutil.ToTime(endDate)
		if err != nil {
			return err
		}

		cfg, err = plot.NewConfig(start, end, authors, period, output)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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
	output    string
	input     string
	startDate string
	endDate   string
	authors   []string
	period    string
)

func Init() {
	PlotCmd.AddCommand(monthlyCmd)
	PlotCmd.AddCommand(timeOfDayCmd)
	PlotCmd.AddCommand(topLanguagesCmd)
	PlotCmd.AddCommand(topAuthorsCmd)
	PlotCmd.AddCommand(weekdayCmd)

	PlotCmd.
		PersistentFlags().
		StringVarP(&startDate, "start-date", "s", "", "Start date")

	PlotCmd.
		PersistentFlags().
		StringVarP(&endDate, "end-date", "e", "", "End date")

	PlotCmd.
		PersistentFlags().
		StringSliceVarP(&authors, "author", "a", config.Config.Authors, "Authors")

	PlotCmd.
		PersistentFlags().
		StringVarP(&output, "output", "o", config.Config.Plot.Output, "Output file")

	PlotCmd.
		PersistentFlags().
		StringVarP(&input, "input", "i", config.Config.Report.Output, "Input file")

	PlotCmd.
		PersistentFlags().
		StringVarP(&period, "period", "p", "", "Period to plot")

	if err := PlotCmd.RegisterFlagCompletionFunc("period", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"today", "24h", "this_week", "7d", "this_month", "30d", "this_year", "1y"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(err)
	}
}
