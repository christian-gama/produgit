package plot

import (
	"time"

	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/spf13/cobra"
)

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Plot the monthly data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.Monthly(
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
