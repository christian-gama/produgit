package list

import (
	"github.com/spf13/cobra"
)

var dir []string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List general information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	ValidArgs: []string{
		"--dir",
		"-d",
	},
}

func Init() {
	ListCmd.AddCommand(authorCmd)

	ListCmd.
		PersistentFlags().
		StringArrayVarP(&dir, "dir", "d", []string{"."}, "The starting directory to search for .git repositories")
}
