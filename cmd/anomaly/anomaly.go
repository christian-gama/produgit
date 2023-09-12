package anomaly

import (
	"github.com/christian-gama/produgit/config"
	"github.com/christian-gama/produgit/internal/anomaly"
	"github.com/christian-gama/produgit/internal/data"
	dateutil "github.com/christian-gama/produgit/internal/util/date"
	"github.com/spf13/cobra"
)

var (
	quantity  int32
	input     string
	startDate string
	endDate   string
	authors   []string
)

var AnomalyCmd = &cobra.Command{
	Use:   "anomaly",
	Short: "Check what commits have a quantity of lines that may be considered an anomaly",
	ValidArgs: []string{
		"--quantity",
		"-q",
		"--start-date",
		"-s",
		"--end-date",
		"-e",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		start, err := dateutil.ToTime(startDate)
		if err != nil {
			return err
		}

		end, err := dateutil.ToTime(endDate)
		if err != nil {
			return err
		}

		cfg, err := anomaly.NewConfig(
			start,
			end,
			quantity,
			input,
			authors,
		)
		if err != nil {
			return err
		}

		logs, err := data.Load(input)
		if err != nil {
			return err
		}

		return anomaly.Anomaly(logs, cfg)
	},
}

func Init() {
	AnomalyCmd.
		Flags().
		Int32VarP(&quantity, "quantity", "q", 3000, "Quantity of lines to be considered an anomaly")

	AnomalyCmd.
		PersistentFlags().
		StringVarP(&input, "input", "i", config.Config.Report.Output, "Input file")

	AnomalyCmd.
		PersistentFlags().
		StringVarP(&startDate, "start-date", "s", "", "Start date")

	AnomalyCmd.
		PersistentFlags().
		StringVarP(&endDate, "end-date", "e", "", "End date")

	AnomalyCmd.
		PersistentFlags().
		StringSliceVarP(&authors, "authors", "a", config.Config.Authors, "Authors to be considered")
}
