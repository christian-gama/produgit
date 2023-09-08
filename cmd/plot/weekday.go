package plot

import (
	"time"

	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/spf13/cobra"
)

var weekdayCmd = &cobra.Command{
	Use:   "weekday",
	Short: "Plot the weekday data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.Weekday(
			plot.NewConfig(input, output),
			plot.NewFilterConfig(
				parseDate(startDate, time.Time{}),
				parseDate(endDate, time.Now()),
				period,
				authors,
			),
		)
	},
}
