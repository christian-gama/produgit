package plot

import (
	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/christian-gama/productivity/pkg/utils/date"
	"github.com/spf13/cobra"
)

var topAuthors = &cobra.Command{
	Use:   "top-authors",
	Short: "Plot the top authors data from the report command",
	Run: func(cmd *cobra.Command, args []string) {
		plot.TopAuthors(
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
