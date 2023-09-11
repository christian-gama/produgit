package plot

import (
	"github.com/christian-gama/produgit/pkg/plot"
	"github.com/christian-gama/produgit/pkg/utils/date"
	"github.com/spf13/cobra"
)

var topLanguages = &cobra.Command{
	Use:   "top-languages",
	Short: "Plot the top languages data from the report command",
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
