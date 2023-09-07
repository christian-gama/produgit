package gitlog

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var rawOutput string = `
'2021-07-30 07:23',example@email.com,Author Name
8	5	src/scripts/core/{confirmAccount.ts => confirm.ts}
103	10	public/css/style.css
1	1	public/css/style.css.map
5	22	src/scripts/DBConnection/accountDB.ts

'2022-01-21 19:25',,Author Name
0	3	src/domain/handler/handler.model.ts
0	6	tests/e2e/cypress/plugins/index.ts
0	1	tests/e2e/cypress/support/index.ts

'2022-01-20 01:46',,
0	0	{src/external/@types => @types}/env.d.ts
1	0	{src/external/@types => @types}/index.d.ts
931	30	{src/external/@types => @types}/vite-env.d.ts
`

func TestParse(t *testing.T) {
	sampleOutput := strings.Split(rawOutput, "\n")

	expectedLogs := []*Log{
		{
			Date:   time.Date(2021, 7, 30, 7, 23, 0, 0, time.UTC),
			Plus:   8,
			Minus:  5,
			Diff:   3,
			Path:   "src/scripts/core/{confirmAccount.ts => confirm.ts}",
			Author: "example@email.com",
		},
		{
			Date:   time.Date(2021, 7, 30, 7, 23, 0, 0, time.UTC),
			Plus:   103,
			Minus:  10,
			Diff:   93,
			Path:   "public/css/style.css",
			Author: "example@email.com",
		},
		{
			Date:   time.Date(2021, 7, 30, 7, 23, 0, 0, time.UTC),
			Plus:   1,
			Minus:  1,
			Diff:   0,
			Path:   "public/css/style.css.map",
			Author: "example@email.com",
		},
		{
			Date:   time.Date(2021, 7, 30, 7, 23, 0, 0, time.UTC),
			Plus:   5,
			Minus:  22,
			Diff:   -17,
			Path:   "src/scripts/DBConnection/accountDB.ts",
			Author: "example@email.com",
		},

		{
			Date:   time.Date(2022, 1, 21, 19, 25, 0, 0, time.UTC),
			Plus:   0,
			Minus:  3,
			Diff:   -3,
			Path:   "src/domain/handler/handler.model.ts",
			Author: "Author Name",
		},
		{
			Date:   time.Date(2022, 1, 21, 19, 25, 0, 0, time.UTC),
			Plus:   0,
			Minus:  6,
			Diff:   -6,
			Path:   "tests/e2e/cypress/plugins/index.ts",
			Author: "Author Name",
		},
		{
			Date:   time.Date(2022, 1, 21, 19, 25, 0, 0, time.UTC),
			Plus:   0,
			Minus:  1,
			Diff:   -1,
			Path:   "tests/e2e/cypress/support/index.ts",
			Author: "Author Name",
		},

		{
			Date:   time.Date(2022, 1, 20, 1, 46, 0, 0, time.UTC),
			Plus:   0,
			Minus:  0,
			Diff:   0,
			Path:   "{src/external/@types => @types}/env.d.ts",
			Author: "Unknown",
		},
		{
			Date:   time.Date(2022, 1, 20, 1, 46, 0, 0, time.UTC),
			Plus:   1,
			Minus:  0,
			Diff:   1,
			Path:   "{src/external/@types => @types}/index.d.ts",
			Author: "Unknown",
		},
		{
			Date:   time.Date(2022, 1, 20, 1, 46, 0, 0, time.UTC),
			Plus:   931,
			Minus:  30,
			Diff:   901,
			Path:   "{src/external/@types => @types}/vite-env.d.ts",
			Author: "Unknown",
		},
	}

	parsedLogs := Parse(sampleOutput)

	parsedLogsMap := make(map[string]*Log)
	for _, log := range parsedLogs {
		parsedLogsMap[fmt.Sprintf("%s%s", log.Date.String(), log.Path)] = log
	}

	for _, expectedLog := range expectedLogs {
		parsedLog, ok := parsedLogsMap[fmt.Sprintf("%s%s", expectedLog.Date.String(), expectedLog.Path)]
		if !ok {
			t.Errorf("Expected log with date %s, but not found", expectedLog.Date)
		} else {
			if parsedLog.Plus != expectedLog.Plus {
				t.Errorf("Expected log with date %s to have plus %d, but got %d",
					expectedLog.Date, expectedLog.Plus, parsedLog.Plus)
			}

			if parsedLog.Minus != expectedLog.Minus {
				t.Errorf("Expected log with date %s to have minus %d, but got %d",
					expectedLog.Date, expectedLog.Minus, parsedLog.Minus)
			}

			if parsedLog.Diff != expectedLog.Diff {
				t.Errorf("Expected log with date %s to have diff %d, but got %d",
					expectedLog.Date, expectedLog.Diff, parsedLog.Diff)
			}

			if parsedLog.Path != expectedLog.Path {
				t.Errorf("Expected log with date %s to have path %s, but got %s",
					expectedLog.Date, expectedLog.Path, parsedLog.Path)
			}

			if parsedLog.Author != expectedLog.Author {
				t.Errorf("Expected log with date %s to have author %s, but got %s",
					expectedLog.Date, expectedLog.Author, parsedLog.Author)
			}
		}
	}
}
