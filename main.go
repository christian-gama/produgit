package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "productivity", Short: "A productivity tool"}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your productivity",
	Run:   runReport,
	ValidArgs: []string{
		"--author",
		"--starting-dir",
		"--output-dir",
	},
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	Run:   runPlot,
	ValidArgs: []string{
		"--output-filename",
		"--ext",
		"--type",
		"--start-date",
		"--finish-date",
	},
}

func init() {
	if _, err := exec.LookPath("git"); err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}

	configReportCmd()
	configPlotCmd()
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(plotCmd)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
}

func configReportCmd() {
	reportCmd.
		Flags().
		StringVarP(&config.StartingDir, "starting-dir", "d", ".", "The starting directory to search for .git repositories")

	// get Author from git config user.name
	cmd := exec.Command("git", "config", "--global", "user.email")
	var author string
	out, err := cmd.Output()
	if err == nil && len(out) == 0 {
		author = string(out)
		fmt.Printf("Using author: %s\n", author)
	}

	reportCmd.
		Flags().
		StringVarP(&config.Author, "author", "a", author, "The author to filter git logs")

	reportCmd.
		Flags().
		StringVarP(&config.OutputDir, "output-dir", "o", "", "The output path for the report")
}

func configPlotCmd() {
	plotCmd.
		Flags().
		StringVarP(&plotConfig.OutputFile, "output-filename", "o", "chart", "The output filename of the plot. The default is 'chart'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.Extension, "ext", "e", "png", "The extension of the output file. The default is 'png'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.Type, "type", "t", "monthly", "The type of the plot. The default is 'monthly'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.StartDate, "start-date", "s", time.Now().AddDate(-2, 0, 0).Format("2006-01-02"), "The start date of the plot. The default is '1970-01-01'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.FinishDate, "finish-date", "f", time.Now().Format("2006-01-02"), "The end date of the plot. The default is '3000-01-01'.")

	if err := plotCmd.MarkFlagRequired("type"); err != nil {
		panic(fmt.Sprintf("Marking flag as required failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"ext",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"png",
				"svg",
				"pdf",
				"eps",
				"tif",
				"tiff",
				"jpg",
				"jpeg",
			}, cobra.ShellCompDirectiveDefault
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"start-date",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"finish-date",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"output-filename",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}

	if err := plotCmd.RegisterFlagCompletionFunc(
		"type",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"monthly",
				"weekday",
				"timeofday",
			}, cobra.ShellCompDirectiveDefault
		},
	); err != nil {
		panic(fmt.Sprintf("Registering flag completion function failed: %s\n", err))
	}
}
