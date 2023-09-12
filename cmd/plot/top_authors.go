package plot

import (
	"github.com/christian-gama/produgit/internal/plot"
	"github.com/spf13/cobra"
)

var topAuthorsCmd = &cobra.Command{
	Use:   "top_authors",
	Short: "Plot the top authors data from the report command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return plot.NewTopAuthors(logs, cfg).Plot()
	},
}
