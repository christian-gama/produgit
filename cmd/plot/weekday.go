package plot

import (
	"github.com/christian-gama/produgit/internal/plot"
	"github.com/spf13/cobra"
)

var weekdayCmd = &cobra.Command{
	Use:   "weekday",
	Short: "Plot the weekday data from the report command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return plot.NewWeekday(logs, cfg).Plot()
	},
}
