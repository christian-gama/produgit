package data

import (
	"sort"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestLogSorter(t *testing.T) {
	tests := []struct {
		name   string
		logs   []*Log
		sorted []*Log
	}{
		{
			name: "sort by author",
			logs: []*Log{
				{Author: "zach"},
				{Author: "adam"},
				{Author: "mike"},
			},
			sorted: []*Log{
				{Author: "adam"},
				{Author: "mike"},
				{Author: "zach"},
			},
		},
		{
			name: "sort by date when authors are same",
			logs: []*Log{
				{
					Author: "mike",
					Date:   timestamppb.New(time.Date(2023, 01, 02, 0, 0, 0, 0, time.UTC)),
				},
				{
					Author: "mike",
					Date:   timestamppb.New(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			sorted: []*Log{
				{
					Author: "mike",
					Date:   timestamppb.New(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
				{
					Author: "mike",
					Date:   timestamppb.New(time.Date(2023, 01, 02, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(logSorter(tt.logs))
			for i, v := range tt.logs {
				if v.Author != tt.sorted[i].Author ||
					v.Date.AsTime() != tt.sorted[i].Date.AsTime() ||
					v.Path != tt.sorted[i].Path ||
					v.Diff != tt.sorted[i].Diff {
					t.Fatalf("expected %+v but got %+v", tt.sorted[i], v)
				}
			}
		})
	}
}
