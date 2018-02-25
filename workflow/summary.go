package workflow

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/clarsen/trello"
	"gopkg.in/yaml.v2"
)

// TaskInfo defines what a task looks like -- just reflects what's currenty in trello.
type TaskInfo struct {
	CreatedDate string `yaml:"createdDate"`
	DoneDate    string `yaml:"doneDate"`
	Title       string `yaml:"title"`
}

// WeeklyGoalInfo defines what a weekly goal looks like
type WeeklyGoalInfo struct {
	Title   string `yaml:"title"`
	Created string `yaml:"createdDate"`
	Status  string `yaml:"status"`
	Week    int    `yaml:"weekNumber"`
	Year    int    `yaml:"year"`
}

// MonthlyGoalInfo defines what a monthly goal looks like - a title, when it was
// created and a list of weekly goals.
type MonthlyGoalInfo struct {
	Title       string           `yaml:"title"`
	Created     string           `yaml:"createdDate"`
	WeeklyGoals []WeeklyGoalInfo `yaml:"weeklyGoals"`
}

// WeeklySummary defines the summarization data that is dumped for downstream
// consumption independent of task management tool (in this case Trello).
type WeeklySummary struct {
	Year           int               `yaml:"year"`
	Month          int               `yaml:"month"`
	Week           int               `yaml:"weekNumber"`
	Done           []TaskInfo        `yaml:"doneThisWeek"`
	MonthlyGoals   []MonthlyGoalInfo `yaml:"monthlyGoals"`
	MonthlySprints []MonthlyGoalInfo `yaml:"monthlySprints"`
}

// MonthlySummary defines the summarization data that is dumped for downstream
// consumption independent of task management tool (in this case Trello).
//
// for monthly summary, monthly goal and monthly sprint info will merge together
// weekly goals
type MonthlySummary struct {
	Year           int               `yaml:"year"`
	Month          int               `yaml:"month"`
	MonthlyGoals   []MonthlyGoalInfo `yaml:"monthlyGoals"`
	MonthlySprints []MonthlyGoalInfo `yaml:"monthlySprints"`
}

// YearlySummary defines what's dumped for rolling up into the yearly plan summary
type YearlySummary struct {
	Year             int              `yaml:"year"`
	MonthlySummaries []MonthlySummary `yaml:"monthlySummaries"`
}

func dumpMonthlyGoalInfo(card *trello.Card, weekNumber *int) (MonthlyGoalInfo, error) {
	var g MonthlyGoalInfo
	g.Title = card.Name

	reMDY := regexp.MustCompile(`\((\d{4}-\d{2}-\d{2})\)`)
	exprGoalMDY := reMDY.FindStringSubmatch(card.Name)
	if len(exprGoalMDY) <= 0 {
		return g, errors.New("Couldn't parse date out of title")
	}
	g.Created = exprGoalMDY[1]

	re := regexp.MustCompile(`^week (\d+): (.*) (?:\((\d{4}-\d{2}-\d{2})\))(?: (\(.*\)))?$`)
	reYear := regexp.MustCompile(`^(\d{4})-\d{2}-\d{2}$`)

	for _, cl := range card.Checklists {
		// log.Println("checklist:", cl)
		for _, item := range cl.CheckItems {
			expr := re.FindStringSubmatch(item.Name)
			if len(expr) > 0 {
				// log.Printf("for %s got match %+v\n", item.Name, expr)
				// log.Println("got week", expr[1])
				// log.Println("got text", expr[2])
				// log.Println("got created", expr[3])
				// log.Println("got status", expr[4])
				week, err := strconv.Atoi(expr[1])

				if weekNumber != nil && week != *weekNumber {
					continue
				}
				if err != nil {
					return g, err
				}
				text := expr[2]
				created := expr[3]
				exprYear := reYear.FindStringSubmatch(created)
				if len(exprYear) <= 0 {
					continue
				}
				year, err := strconv.Atoi(exprYear[1])
				if err != nil {
					return g, err
				}
				status := expr[4]
				wgi := WeeklyGoalInfo{
					Year:    year,
					Week:    week,
					Created: created,
					Title:   text,
					Status:  status,
				}
				g.WeeklyGoals = append(g.WeeklyGoals, wgi)
			}
		}
	}

	return g, nil
}

func dumpSummaryForWeek(
	summary *WeeklySummary,
	out io.Writer,
) error {
	buf, err := yaml.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}

	_, err = out.Write(buf)
	return err
}

func generateWeeklySummary(
	year, week, month int,
	doneCards, goalCards, sprintCards []*trello.Card,
) (
	*WeeklySummary, error,
) {
	summary := WeeklySummary{}
	summary.Year = year
	summary.Week = week
	summary.Month = month
	local, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}

	for _, card := range doneCards {

		actions, err2 := card.GetActions(trello.Arguments{"filter": "copyCard,updateCard:idList,moveCardToBoard"})
		if err2 != nil {
			return nil, err2
		}
		var latest *time.Time
		for _, a := range actions {
			if latest == nil || a.Date.After(*latest) {
				latest = &a.Date
			}
		}

		// well, unfortunately, Trello loses action data after a while when moving
		// card to another board, so the last action may only be the time when it
		// was moved to the history board and no actions when in a previous board.
		var doneDate string
		if latest == nil {
			doneDate = "unknown"
		} else {
			doneDate = latest.In(local).Format("2006-01-02 (Mon)")
		}
		summary.Done = append(summary.Done, TaskInfo{
			DoneDate:    doneDate,
			Title:       card.Name,
			CreatedDate: card.CreatedAt().In(local).Format("2006-01-02"),
		})
	}

	for _, goalCard := range goalCards {
		dmgi, err2 := dumpMonthlyGoalInfo(goalCard, &summary.Week)
		if err2 != nil {
			return nil, err2
		}
		summary.MonthlyGoals = append(summary.MonthlyGoals, dmgi)
	}

	for _, sprintCard := range sprintCards {
		dmsi, err2 := dumpMonthlyGoalInfo(sprintCard, &summary.Week)
		if err2 != nil {
			return nil, err2
		}
		summary.MonthlySprints = append(summary.MonthlySprints, dmsi)
	}
	return &summary, nil
}

func prepareSummaryForWeek(
	user, appkey, authtoken string,
) (
	year, week, month int,
	doneCards, goalCards, sprintCards []*trello.Card,
	err error) {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	doneList, err := listFor(cl.member, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return
	}

	doneCards, err = doneList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return
	}

	monthlyGoalsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return
	}
	goalCards, err = monthlyGoalsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return
	}

	monthlySprintsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Sprints")
	if err != nil {
		return
	}
	sprintCards, err = monthlySprintsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return
	}

	year, week = time.Now().ISOWeek()
	month = int(time.Now().Month())

	return year, week, month, doneCards, goalCards, sprintCards, nil
}

// GetSummaryForWeek returns a summary structure usable by other downstream
// in-memory pipelines like daily reminder.
func GetSummaryForWeek(user, appkey, authtoken string) (*WeeklySummary, error) {
	year, week, month, doneCards, goalCards, sprintCards, err := prepareSummaryForWeek(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	return generateWeeklySummary(year, week, month, doneCards, goalCards, sprintCards)
}

// DumpSummaryForWeek dumps current content of Trello board to summary file for week
func DumpSummaryForWeek(user, appkey, authtoken string, out io.Writer) error {
	year, week, month, doneCards, goalCards, sprintCards, err := prepareSummaryForWeek(user, appkey, authtoken)
	if err != nil {
		return err
	}
	summary, err := generateWeeklySummary(year, week, month, doneCards, goalCards, sprintCards)
	if err != nil {
		return err
	}
	return dumpSummaryForWeek(summary, out)

}

func mergeMonthlyGoalInfo(goalsAcrossWeeks []MonthlyGoalInfo) []MonthlyGoalInfo {
	allWeeklyGoals := map[string][]WeeklyGoalInfo{}
	created := map[string]string{}
	for _, mg := range goalsAcrossWeeks {
		allWeeklyGoals[mg.Title] = append(allWeeklyGoals[mg.Title], mg.WeeklyGoals...)
		created[mg.Title] = mg.Created
	}

	var monthlyGoals []MonthlyGoalInfo
	for title, goals := range allWeeklyGoals {
		mg := MonthlyGoalInfo{Title: title,
			Created:     created[title],
			WeeklyGoals: goals,
		}
		monthlyGoals = append(monthlyGoals, mg)
	}
	return monthlyGoals
}

// GenerateSummaryForMonth rolls up weekly summaries to month level
func GenerateSummaryForMonth(year, month int, summaryIn [][]byte, out io.Writer) error {
	var weeklies []WeeklySummary
	for _, buf := range summaryIn {
		var weekly WeeklySummary
		err := yaml.Unmarshal(buf, &weekly)
		if err != nil {
			return err
		}
		if weekly.Year == year && weekly.Month == month {
			weeklies = append(weeklies, weekly)
		}
	}
	if len(weeklies) <= 0 {
		return fmt.Errorf("No summaries for %d-%02d", year, month)
	}

	summary := MonthlySummary{}
	summary.Year = year
	summary.Month = month

	var allMonthlyGoals []MonthlyGoalInfo
	for _, ws := range weeklies {
		allMonthlyGoals = append(allMonthlyGoals, ws.MonthlyGoals...)
	}
	summary.MonthlyGoals = mergeMonthlyGoalInfo(allMonthlyGoals)

	allMonthlyGoals = []MonthlyGoalInfo{}
	for _, ws := range weeklies {
		allMonthlyGoals = append(allMonthlyGoals, ws.MonthlySprints...)
	}
	summary.MonthlySprints = mergeMonthlyGoalInfo(allMonthlyGoals)
	// log.Printf("filtered to %+v\n", weeklies)

	buf, err := yaml.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}

	_, err = out.Write(buf)
	return err
}

// GenerateSummaryForYear rolls up monthly summaries to year level
func GenerateSummaryForYear(year int, summaryIn [][]byte, out io.Writer) error {
	var monthlies []MonthlySummary
	for _, buf := range summaryIn {
		var monthly MonthlySummary
		err := yaml.Unmarshal(buf, &monthly)
		if err != nil {
			return err
		}
		monthlies = append(monthlies, monthly)
	}
	if len(monthlies) <= 0 {
		return fmt.Errorf("No summaries for %d", year)
	}

	summary := YearlySummary{Year: year}
	for _, m := range monthlies {
		summary.MonthlySummaries = append(summary.MonthlySummaries, m)
	}

	buf, err := yaml.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}

	_, err = out.Write(buf)

	return err
}

// DumpSummaryForWeekFromHistory dumps current content of Trello history board to summary file for week
func DumpSummaryForWeekFromHistory(user, appkey, authtoken string, week int, out io.Writer) error {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	months := []string{"January", "February", "March", "April",
		"May", "June", "July", "August", "September", "October",
		"November", "December",
	}

	var doneCards, goalCards, sprintCards []*trello.Card

	for idx, month := range months {
		olderMonthlyDoneList, err := listFor(cl.member, "History 2018", fmt.Sprintf("%02d %s", week, month))
		if err != nil {
			return err
		}
		if olderMonthlyDoneList == nil {
			continue
		}
		doneCards, err = olderMonthlyDoneList.GetCards(trello.Defaults())
		if err != nil {
			return err
		}

		monthlyGoalsList, err := listFor(cl.member, "History 2018", fmt.Sprintf("%02d %s goals", week, month))
		if err != nil {
			return err
		}
		goalCards, err = monthlyGoalsList.GetCards(trello.Arguments{"fields": "all"})
		if err != nil {
			// handle error
			return err
		}

		var monthlySprintsList *trello.List
		if week >= 5 {
			// for the one-off capture, we're mid-february where there is no sprints
			// list (only created at monthly review)
			monthlySprintsList, err = listFor(cl.member, "Kanban daily/weekly", "Monthly Sprints")
		} else {
			monthlySprintsList, err = listFor(cl.member, "History 2018", fmt.Sprintf("%s sprints", month))
		}
		if err != nil {
			return err
		}

		sprintCards, err = monthlySprintsList.GetCards(trello.Arguments{"fields": "all"})
		if err != nil {
			// handle error
			return err
		}
		summary, err := generateWeeklySummary(2018, week, idx+1, doneCards, goalCards, sprintCards)
		if err != nil {
			return err
		}
		return dumpSummaryForWeek(summary, out)
	}

	// doesn't reach here
	return nil
}
