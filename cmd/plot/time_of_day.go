package plot

import (
	"github.com/christian-gama/produgit/internal/plot"
	"github.com/spf13/cobra"
)

var timeOfDayCmd = &cobra.Command{
	Use:   "time_of_day",
	Short: "Plot the time of day data from the report command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return plot.NewTimeOfDay(logs, cfg).Plot()
	},
}
