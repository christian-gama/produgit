package report

import (
	"encoding/json"
	"os"

	"github.com/christian-gama/productivity/pkg/gitlog"
)

func Read(reportPath string) []*gitlog.Log {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		panic(err)
	}

	var logs []*gitlog.Log
	err = json.Unmarshal(data, &logs)
	if err != nil {
		panic(err)
	}

	return logs
}
