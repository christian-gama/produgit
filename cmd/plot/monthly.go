package plot

import (
	"github.com/christian-gama/produgit/internal/plot"
	"github.com/spf13/cobra"
)

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Plot the monthly data from the report command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return plot.NewMonthly(logs, cfg).Plot()
	},
}
