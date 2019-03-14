package workflow

import (
	"testing"
	"time"
)

func TestChecklistItemParsing(t *testing.T) {
	var tests = []struct {
		inputTitle string
		title      string
		week       int
		created    time.Time
		status     string
	}{
		{"week 11: (2h) implement weekly review repository manipulation (2019-03-10) (partial)",
			"(2h) implement weekly review repository manipulation",
			11,
			time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC),
			"(partial)",
		},
	}
	for _, tvec := range tests {
		title, week, created, status := GetAttributesFromChecklistTitle(tvec.inputTitle)
		if title != tvec.title || week != tvec.week || (status != nil && *status != tvec.status) || (created != nil && tvec.created != *created) {
			t.Errorf("GetAttributesFromChecklistTitle(%+v): expected (%s,%d,%s,%+v), actual (%s,%d,%s,%+v)",
				tvec.inputTitle, tvec.title, tvec.week, tvec.status, tvec.created, title, week, *status, *created)
		}
	}
}

func TestYMWForTime(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fail()
	}
	var mdyTests = []struct {
		// input
		time time.Time
		// expected results
		year  int
		month string
		week  int
	}{
		{time.Date(2017, 12, 31, 12, 0, 0, 0, loc), 2017, "December", 52},
		{time.Date(2018, 1, 7, 12, 0, 0, 0, loc), 2018, "January", 1},
		{time.Date(2018, 4, 1, 12, 0, 0, 0, loc), 2018, "March", 13},
	}
	for _, tvec := range mdyTests {
		year, month, week := ymwForTime(tvec.time)
		if year != tvec.year || month != tvec.month || week != tvec.week {
			t.Errorf("YMWForTime(%+v): expected (%d,%s,%d), actual (%d,%s,%d)",
				tvec.time, tvec.year, tvec.month, tvec.week, year, month, week)
		}
	}
}