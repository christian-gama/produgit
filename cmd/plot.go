package main

import (
	"encoding/json"
	"os"

	"github.com/christian-gama/productivity/pkg/gitlog"
	"github.com/christian-gama/productivity/pkg/plot"
	"github.com/spf13/cobra"
)

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plot the data from the report command",
	Run:   runPlot,
}

func runPlot(cmd *cobra.Command, args []string) {
	fileData, err := os.ReadFile("testando4.json")
	if err != nil {
		panic(err)
	}

	var data []*gitlog.Log
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		panic(err)
	}

	plot.Monthly(data)
}
