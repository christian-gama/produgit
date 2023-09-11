package plot

import (
	"github.com/christian-gama/produgit/internal/plot"
	"github.com/spf13/cobra"
)

var topLanguagesCmd = &cobra.Command{
	Use:   "top_languages",
	Short: "Plot the top languages data from the report command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return plot.NewTopLanguages(logs, cfg).Plot()
	},
}
