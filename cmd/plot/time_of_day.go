package plot

import (
	"time"

	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/spf13/cobra"
)

var timeOfDay = &cobra.Command{
	Use:   "timeOfDay",
	Short: "Plot the time of day data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.TimeOfDay(
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
