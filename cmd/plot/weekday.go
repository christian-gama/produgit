package plot

import (
	"github.com/christian-gama/produgit/pkg/plot"
	"github.com/christian-gama/produgit/pkg/utils/date"
	"github.com/spf13/cobra"
)

var weekdayCmd = &cobra.Command{
	Use:   "weekday",
	Short: "Plot the weekday data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.Weekday(
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
