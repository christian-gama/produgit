package plot

import "testing"

func TestGetTimeOfDay(t *testing.T) {
	tests := []struct {
		hour int
		want string
	}{
		{0, "midnight"},
		{1, "midnight"},
		{2, "midnight"},
		{3, "midnight"},
		{4, "midnight"},
		{5, "midnight"},
		{6, "morning"},
		{7, "morning"},
		{8, "morning"},
		{9, "morning"},
		{10, "morning"},
		{11, "morning"},
		{12, "afternoon"},
		{13, "afternoon"},
		{14, "afternoon"},
		{15, "afternoon"},
		{16, "afternoon"},
		{17, "afternoon"},
		{18, "night"},
		{19, "night"},
		{20, "night"},
		{21, "night"},
		{22, "night"},
		{23, "night"},
		{24, "midnight"},
		{25, "invalid hour"},
		{-1, "invalid hour"},
	}

	for _, test := range tests {
		got := getTimeOfDay(test.hour)
		if got != test.want {
			t.Errorf("getTimeOfDay(%d) = %s; want %s", test.hour, got, test.want)
		}
	}
}
