package workflow

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	packr "github.com/gobuffalo/packr/v2"
	emoji "gopkg.in/kyokomi/emoji.v1"
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

var box *packr.Box

func init() {
	box = packr.New("workflow", "../templates")
}

func formatMonthAsString(month int) string {
	monthString := []string{
		"",
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}
	return monthString[month]
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

// CreateEmptyWeeklyRetrospective populates an empty template based on weekly summary
func CreateEmptyWeeklyRetrospective(summaryIn io.Reader, reviewOut io.Writer) error {
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

	review := WeeklyReview{
		GoingWell:        []string{"1", "2", "3"},
		NeedsImprovement: []string{"1", "2", "3"},
		Successes:        []string{"1", "2"},
		Challenges:       []string{"1", "2"},
		LearnAboutMyself: []string{"1", "2"},
		LearnAboutOthers: []string{"1", "2"},
	}

	// fill in pergoalreview from weekly goals
	for _, goal := range weekly.MonthlyGoals {
		if len(goal.WeeklyGoals) > 0 {
			var outcomes []string
			monthgoal := fmt.Sprintf("goal: %s", goal.Title)
			outcomes = append(outcomes, monthgoal)
			for _, weeklygoal := range goal.WeeklyGoals {
				note := ""
				if strings.Contains(weeklygoal.Status, "(done)") {
					note = ":green_heart:"
				}
				if strings.Contains(weeklygoal.Status, "(partial)") {
					note = ":yellow_heart:"
				}
				if strings.Contains(weeklygoal.Status, "(not done)") {
					note = ":broken_heart:"
				}
				outcome := emoji.Sprintf("%s %s %s", weeklygoal.Title, weeklygoal.Status, note)
				outcomes = append(outcomes, outcome)
			}
			outcomes = append(outcomes, "1")
			outcomes = append(outcomes, "2")
			review.PerGoalReviews = append(review.PerGoalReviews, PerGoalReview{
				DidToCreateOutcome: outcomes,
				KeepDoing:          []string{"1", "2"},
				DoDifferently:      []string{"1", "2"},
			})
		}
	}

	for _, goal := range weekly.MonthlySprints {
		if len(goal.WeeklyGoals) > 0 {
			var outcomes []string
			monthsprint := fmt.Sprintf("sprint: %s", goal.Title)
			outcomes = append(outcomes, monthsprint)
			for _, weeklygoal := range goal.WeeklyGoals {
				note := ""
				if strings.Contains(weeklygoal.Status, "(done)") {
					note = ":green_heart:"
				}
				if strings.Contains(weeklygoal.Status, "(partial)") {
					note = ":yellow_heart:"
				}
				if strings.Contains(weeklygoal.Status, "(not done)") {
					note = ":broken_heart:"
				}
				outcome := emoji.Sprintf("%s %s %s", weeklygoal.Title, weeklygoal.Status, note)
				outcomes = append(outcomes, outcome)
			}
			outcomes = append(outcomes, "1")
			outcomes = append(outcomes, "2")
			review.PerGoalReviews = append(review.PerGoalReviews, PerGoalReview{
				DidToCreateOutcome: outcomes,
				KeepDoing:          []string{"1", "2"},
				DoDifferently:      []string{"1", "2"},
			})
		}
	}

	buf2, err3 := yaml.Marshal(&review)
	if err3 != nil {
		log.Fatal(err3)
	}

	_, err = reviewOut.Write(buf2)
	return err
}

// CreateEmptyMonthlyRetrospective populates an empty template based on monthly summary
func CreateEmptyMonthlyRetrospective(summaryIn io.Reader, reviewOut io.Writer) error {
	buf, err := ioutil.ReadAll(summaryIn)
	if err != nil {
		return err
	}
	log.Printf("Got %+v", buf)
	retro := MonthlyRetrospective{
		MonthlyReview: MonthlyReview{
			Continue:         []string{"1", "2"},
			DoDifferently:    []string{"1", "2"},
			CandidateGoals:   []string{"1", "2"},
			CandidateSprints: []string{"1", "2"},
			Highlights:       []string{"1", "2"},
		},
	}

	err = yaml.Unmarshal(buf, &retro.MonthlySummary)
	if err != nil {
		return err
	}
	log.Printf("Got %+v", retro)

	for _, goal := range retro.MonthlyGoals {
		retro.MonthlyGoalReviews = append(retro.MonthlyGoalReviews, MonthlyGoalReview{
			Title:           goal.Title,
			Accomplishments: []string{"1", "2"},
			CreatedBy:       []string{"1", "2"},
		})
	}
	for _, goal := range retro.MonthlySprints {
		retro.MonthlySprintReviews = append(retro.MonthlySprintReviews, MonthlySprintReview{
			Title: goal.Title,
			LearningsAndResultsWhatContinueWhatChange: []string{"1", "2"},
		})
	}

	buf, err = yaml.Marshal(&retro.MonthlyReview)
	if err != nil {
		log.Fatal(err)
	}

	_, err = reviewOut.Write(buf)
	return err
}

// VisualizeWeeklyRetrospective writes out report of weekly tasks done, goals, sprints
func VisualizeWeeklyRetrospective(summaryIn, reviewIn io.Reader, visOut io.Writer) error {
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

	tmpl, err := box.FindString("weekly-retrospective.md")
	if err != nil {
		return err
	}
	t := template.Must(template.New("weekly-retrospective.md").Parse(tmpl))
	t.Execute(visOut, weekly)

	return nil
}

// VisualizeMonthlyInput writes out report of weekly review items collated for preparing monthly review
func VisualizeWeeklySummariesForMonthly(summaryIn io.Reader, visOut io.Writer) error {
	buf, err := ioutil.ReadAll(summaryIn)
	if err != nil {
		return err
	}
	var summary MonthlySummary
	err = yaml.Unmarshal(buf, &summary)
	if err != nil {
		return err
	}

	tmpl, err := box.FindString("monthly-input.md")
	if err != nil {
		return err
	}
	t := template.Must(template.New("monthly-input.md").Parse(tmpl))
	t.Execute(visOut, summary)

	return nil
}

// VisualizeMonthlyRetrospective writes out report of monthly goals, sprints
func VisualizeMonthlyRetrospective(summaryIn, reviewIn io.Reader, visOut io.Writer) error {
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

	tmpl, err := box.FindString("monthly-retrospective.md")
	if err != nil {
		return err
	}
	t := template.Must(template.New("monthly-retrospective.md").Parse(tmpl))
	t.Execute(visOut, retrospective)

	return nil
}

// VisualizeYearlyPlanSummary writes out report of monthly goals, sprints over the year
func VisualizeYearlyPlanSummary(summaryIn io.Reader, visOut io.Writer) error {
	buf, err := ioutil.ReadAll(summaryIn)
	if err != nil {
		return err
	}
	var summary YearlySummary
	err = yaml.Unmarshal(buf, &summary)
	if err != nil {
		return err
	}

	fmap := template.FuncMap{
		"formatMonthAsString": formatMonthAsString,
	}
	tmpl, err := box.FindString("yearly-plan-summary.md")
	if err != nil {
		return err
	}
	t := template.Must(template.New("yearly-plan-summary.md").Funcs(fmap).Parse(tmpl))
	t.Execute(visOut, summary)

	return nil
}
