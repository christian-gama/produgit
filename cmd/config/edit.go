package config

import (
	"fmt"
	"os"

	"github.com/christian-gama/produgit/config"
	cmdutil "github.com/christian-gama/produgit/internal/util/cmd"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:       "edit",
	Short:     "Edit the produgit configuration",
	ValidArgs: []string{},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.Exists() {
			if err := config.Reset(); err != nil {
				return fmt.Errorf("Failed to reset config: %v", err)
			}
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		configPath, err := config.DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("Failed to get config path: %v", err)
		}

		if _, err := cmdutil.Run(editor, configPath); err != nil {
			return fmt.Errorf("Failed to run editor: %v", err)
		}

		return nil
	},
}
