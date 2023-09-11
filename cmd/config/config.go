package config

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:       "config",
	Short:     "Manage the produgit configuration",
	ValidArgs: []string{},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Init() {
	ConfigCmd.AddCommand(editCmd)
	ConfigCmd.AddCommand(resetCmd)
}
