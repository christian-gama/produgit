package anomalies

import (
	"github.com/christian-gama/productivity/pkg/anomalies"
	"github.com/christian-gama/productivity/pkg/utils/date"
	"github.com/spf13/cobra"
)

var (
	quantity  int
	input     string
	startDate string
	endDate   string
	authors   []string
)

var AnomaliesCmd = &cobra.Command{
	Use:   "anomalies",
	Short: "Check what commits have a quantity of lines that may be considered an anomaly",
	ValidArgs: []string{
		"--quantity",
		"-q",
		"--start-date",
		"-s",
		"--end-date",
		"-e",
	},
	Run: func(cmd *cobra.Command, args []string) {
		anomalies.Anomalies(
			anomalies.NewConfig(
				date.Parse(startDate),
				date.Parse(endDate),
				quantity,
				input,
				authors,
			),
		)
	},
}

func init() {
	AnomaliesCmd.
		Flags().
		IntVarP(&quantity, "quantity", "q", 2000, "Quantity of lines to be considered an anomaly")

	AnomaliesCmd.
		PersistentFlags().
		StringVarP(&input, "input", "i", "", "Input file")

	AnomaliesCmd.
		PersistentFlags().
		StringVarP(&startDate, "start-date", "s", "", "Start date")

	AnomaliesCmd.
		PersistentFlags().
		StringVarP(&endDate, "end-date", "e", "", "End date")

	AnomaliesCmd.
		PersistentFlags().
		StringSliceVarP(&authors, "authors", "a", []string{}, "Authors to filter")
}
