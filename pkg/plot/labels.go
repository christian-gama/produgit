package plot

import (
	"sort"
	"time"

	"github.com/christian-gama/productivity/pkg/gitlog"
)

func createLanguageLabel(logs map[string][]*Data) []string {
	var result []string
	for lang := range logs {
		result = append(result, lang)
	}

	sort.Strings(result)
	return result
}

func createMonthLabels(dateFmt string, logs []*gitlog.Log) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, log := range logs {
		monthYear := log.Date.Format(dateFmt)
		if _, exists := seen[monthYear]; !exists {
			result = append(result, monthYear)
			seen[monthYear] = struct{}{}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(dateFmt, result[i])
		t2, _ := time.Parse(dateFmt, result[j])
		return t1.Before(t2)
	})

	return result
}
