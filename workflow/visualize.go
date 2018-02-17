package workflow

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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

// WeeklyRetrospective defines what goes into the visualization template.
type WeeklyRetrospective struct {
	WeeklySummary
	WeeklyReview
	NowHHMM        string
	ThisWeekSunday string
	DoneByDay      []DaySummary
}

// MonthlyRetrospective defines what goes into the visualization template.
type MonthlyRetrospective struct {
	MonthlySummary // mostly source
	MonthlyReview  // populated, joined with source
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

// VisualizeMonthlyRetrospective writes out report of monthly goals, sprints
func VisualizeMonthlyRetrospective(summaryIn, reviewIn io.Reader) error {
	buf, err := ioutil.ReadAll(summaryIn)
	if err != nil {
		return err
	}
	// log.Println("Read", buf)

	var retrospective MonthlyRetrospective
	err = yaml.Unmarshal(buf, &retrospective.MonthlySummary)
	if err != nil {
		return err
	}

	buf, err = ioutil.ReadAll(reviewIn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &retrospective.MonthlyReview)
	if err != nil {
		return err
	}

	shouldSee := map[string]MonthlyGoalInfo{}

	for _, goal := range retrospective.MonthlySummary.MonthlyGoals {
		shouldSee[goal.Title] = goal
	}
	// cross check and also fill in info
	for _, goal := range retrospective.MonthlyReview.MonthlyGoalReviews {
		if _, ok := shouldSee[goal.Title]; ok {
			// goal.Created = shouldSee[goal.Title].Created
			delete(shouldSee, goal.Title)
		} else {
			return fmt.Errorf("Goal %s not found as a goal", goal.Title)
		}
	}

	if len(shouldSee) > 0 {
		for k := range shouldSee {
			log.Printf("Need comment about goal %s", k)
		}
		return errors.New("Need comments")
	}

	shouldSee = map[string]MonthlyGoalInfo{}

	for _, goal := range retrospective.MonthlySummary.MonthlySprints {
		shouldSee[goal.Title] = goal
	}
	// cross check and also fill in info
	for _, goal := range retrospective.MonthlyReview.MonthlySprintReviews {
		if _, ok := shouldSee[goal.Title]; ok {
			// goal.Created = shouldSee[goal.Title].Created
			delete(shouldSee, goal.Title)
		} else {
			return fmt.Errorf("Sprint %s not found as a sprint", goal.Title)
		}
	}

	if len(shouldSee) > 0 {
		for k := range shouldSee {
			log.Printf("Need comment about sprint %s", k)
		}
		return errors.New("Need comments")
	}

	t, _ := template.ParseFiles("templates/monthly-retrospective.md")
	t.Execute(os.Stdout, retrospective)

	return nil
}
