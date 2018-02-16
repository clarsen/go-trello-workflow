package workflow

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// DaySummary is generated from WeeklySummary.Done
type DaySummary struct {
	Date string
	Done []TaskInfo
}

type WeeklyRetrospective struct {
	WeeklySummary
	WeeklyReview
	NowHHMM        string
	ThisWeekSunday string
	DoneByDay      []DaySummary
}

func summarizeByDay(summary WeeklySummary) ([]DaySummary, error) {
	doneByDay := map[string][]TaskInfo{}

	for _, done := range summary.Done {
		doneByDay[done.DoneDate] = append(doneByDay[done.DoneDate], done)
	}

	var sortedDates []string
	for date := range doneByDay {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates)
	var ds []DaySummary

	for _, d := range sortedDates {
		ds = append(ds, DaySummary{
			Date: d,
			Done: doneByDay[d],
		})
	}
	return ds, nil
}

// VisualizeWeeklyRetrospective writes out report of weekly tasks done, goals, sprints
func VisualizeWeeklyRetrospective(summaryIn, reviewIn io.Reader) error {
	buf, err := ioutil.ReadAll(summaryIn)
	if err != nil {
		return err
	}
	// log.Println("Read", buf)

	var weekly WeeklyRetrospective
	err = yaml.Unmarshal(buf, &weekly.WeeklySummary)
	if err != nil {
		return err
	}

	buf, err = ioutil.ReadAll(reviewIn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &weekly.WeeklyReview)
	if err != nil {
		return err
	}

	weekly.NowHHMM = time.Now().Format("15:04")
	sunday := time.Now()
	for sunday.Weekday() != time.Sunday {
		sunday = sunday.AddDate(0, 0, -1)
	}
	weekly.ThisWeekSunday = sunday.Format("2006-01-02")

	weekly.DoneByDay, err = summarizeByDay(weekly.WeeklySummary)
	if err != nil {
		return err
	}
	// log.Printf("Got %+v\n", weekly)

	t, _ := template.ParseFiles("templates/weekly-retrospective.md")
	t.Execute(os.Stdout, weekly)

	return nil
}
