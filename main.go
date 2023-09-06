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
		"--exclude",
		"--debug",
	},
}

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	Run:   runPlot,
	ValidArgs: []string{
		"--output-filename",
		"--ext",
		"--monthly",
		"--weekday",
		"--timeofday",
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
	var author []string
	out, err := cmd.Output()
	if err == nil && len(out) == 0 {
		author = []string{string(out)}
		fmt.Printf("Using author: %s\n", string(out))
	}

	reportCmd.
		Flags().
		StringArrayVarP(&config.Authors, "author", "a", author, "The author to filter git logs")

	reportCmd.
		Flags().
		StringVarP(&config.OutputDir, "output-dir", "o", "", "The output path for the report")

	reportCmd.
		Flags().
		StringArrayVarP(&config.Exclude, "exclude", "e", []string{}, "The directories to exclude from the report")

	reportCmd.
		Flags().
		BoolVar(&config.Debug, "debug", false, "Generate a debug report, which includes the raw git logs")
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
		BoolVar(&plotConfig.Monthly, "monthly", false, "Plot the data by month.")

	plotCmd.
		Flags().
		BoolVar(&plotConfig.Weekday, "weekday", false, "Plot the data by week.")

	plotCmd.
		Flags().
		BoolVar(&plotConfig.TimeOfDay, "timeofday", false, "Plot the data by time of day.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.StartDate, "start-date", "s", time.Now().AddDate(-2, 0, 0).Format("2006-01-02"), "The start date of the plot. The default is '1970-01-01'.")

	plotCmd.
		Flags().
		StringVarP(&plotConfig.FinishDate, "finish-date", "f", time.Now().Format("2006-01-02"), "The end date of the plot. The default is '3000-01-01'.")

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
}
