package config

import (
	"github.com/christian-gama/produgit/config"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:       "reset",
	Short:     "Reset the produgit configuration to its default values",
	ValidArgs: []string{},
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Reset()
	},
}
