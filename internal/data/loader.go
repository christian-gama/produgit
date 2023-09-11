package data

import (
	"os"

	"google.golang.org/protobuf/proto"
)

// Load loads the logs from the report.
func Load(reportPath string) (*Logs, error) {
	file, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, err
	}

	logs := &Logs{}
	if err := proto.Unmarshal(file, logs); err != nil {
		return nil, err
	}

	return logs, nil
}
