package plot

import (
	"time"

	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/spf13/cobra"
)

var topLanguages = &cobra.Command{
	Use:   "top-languages",
	Short: "Plot the top languages data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.TopAuthors(
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
