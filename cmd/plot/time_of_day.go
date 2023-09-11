package plot

import (
	"github.com/christian-gama/produgit/pkg/plot"
	"github.com/christian-gama/produgit/pkg/utils/date"
	"github.com/spf13/cobra"
)

var timeOfDay = &cobra.Command{
	Use:   "time-of-day",
	Short: "Plot the time of day data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.TimeOfDay(
			plot.NewConfig(input, output),
			plot.NewFilterConfig(
				date.Parse(startDate),
				date.Parse(endDate),
				period,
				authors,
			),
		)
	},
}
