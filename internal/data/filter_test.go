package data

import (
	reflect "reflect"
	"testing"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestFilter_WithDate(t *testing.T) {
	type args struct {
		startTime time.Time
		endTime   time.Time
	}

	tests := []struct {
		name     string
		logs     []*Log
		args     args
		expected []*Log
		wantErr  bool
	}{
		{
			name: "empty logs",
			logs: []*Log{},
			args: args{
				startTime: time.Now().Add(-24 * time.Hour),
				endTime:   time.Now(),
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "filter by date",
			logs: []*Log{
				{Date: timestamppb.New(time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC))},
				{Date: timestamppb.New(time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC))},
			},
			args: args{
				startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				endTime:   time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			expected: []*Log{
				{Date: timestamppb.New(time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC))},
				{Date: timestamppb.New(time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC))},
			},
			wantErr: false,
		},
		{
			name: "filter by date with no logs",
			logs: []*Log{
				{Date: timestamppb.New(time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC))},
				{Date: timestamppb.New(time.Date(2020, 1, 2, 12, 0, 0, 0, time.UTC))},
			},
			args: args{
				startTime: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				endTime:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.logs, WithDate(tt.args.startTime, tt.args.endTime))
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.expected[i]) {
					t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestFilter_WithAuthors(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*Log
		authors  []string
		expected []*Log
		wantErr  bool
	}{
		{
			name: "match authors",
			logs: []*Log{
				{Author: "John"},
				{Author: "Jane"},
				{Author: "Doe"},
			},
			authors: []string{"John", "Jane"},
			expected: []*Log{
				{Author: "John"},
				{Author: "Jane"},
			},
			wantErr: false,
		},
		{
			name: "case insensitive match",
			logs: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			authors: []string{"JoHn"},
			expected: []*Log{
				{Author: "john"},
			},
			wantErr: false,
		},
		{
			name: "no match",
			logs: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			authors:  []string{"Alan"},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "regex match",
			logs: []*Log{
				{Author: "john"},
				{Author: "johndoe"},
				{Author: "Jane"},
			},
			authors: []string{"john.*"},
			expected: []*Log{
				{Author: "john"},
				{Author: "johndoe"},
			},
			wantErr: false,
		},
		{
			name: "invalid regex",
			logs: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			authors:  []string{"john("},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.logs, WithAuthors(tt.authors))
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i].Author != tt.expected[i].Author {
					t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestFilter_WithMergeAuthors(t *testing.T) {
	tests := []struct {
		name     string
		logs     []*Log
		authors  []string
		expected []*Log
		wantErr  bool
	}{
		{
			name: "merge authors",
			logs: []*Log{
				{Author: "John"},
				{Author: "John Doe"},
				{Author: "Jane"},
			},
			authors: []string{"John"},
			expected: []*Log{
				{Author: "John"},
				{Author: "John"},
				{Author: "Jane"},
			},
			wantErr: false,
		},
		{
			name: "case insensitive merge",
			logs: []*Log{
				{Author: "john"},
				{Author: "JoHn Doe"},
				{Author: "Jane"},
			},
			authors: []string{"jOhN"},
			expected: []*Log{
				{Author: "jOhN"},
				{Author: "jOhN"},
				{Author: "Jane"},
			},
			wantErr: false,
		},
		{
			name: "no merge",
			logs: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			authors: []string{"Alan"},
			expected: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			wantErr: false,
		},
		{
			name: "multiple author merge",
			logs: []*Log{
				{Author: "john"},
				{Author: "johndoe"},
				{Author: "Jane"},
			},
			authors: []string{"john", "Jane"},
			expected: []*Log{
				{Author: "john"},
				{Author: "john"},
				{Author: "Jane"},
			},
			wantErr: false,
		},
		{
			name: "invalid regex",
			logs: []*Log{
				{Author: "john"},
				{Author: "Jane"},
			},
			authors:  []string{"john("},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.logs, WithMergeAuthors(tt.authors))
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i].Author != tt.expected[i].Author {
					t.Errorf("Filter() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
