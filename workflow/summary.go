package workflow

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

// WeeklySummaryItems is array of feedback from a particular week
type WeeklySummaryItems struct {
	Week    int      `yaml:"week"`
	Content []string `yaml:"content"`
}

// PerGoalWeeklySummaryItems is array of feedback from a particular week for a particular goal
type PerGoalWeeklySummaryItems struct {
	Goal               string               `yaml:"goal"`
	DidToCreateOutcome []WeeklySummaryItems `yaml:"didToCreateOutcome"`
	KeepDoing          []WeeklySummaryItems `yaml:"keepDoing"`
	DoDifferently      []WeeklySummaryItems `yaml:"doDifferently"`
}

// MonthlySummary defines the summarization data that is dumped for downstream
// consumption independent of task management tool (in this case Trello).
//
// for monthly summary, monthly goal and monthly sprint info will merge together
// weekly goals
type MonthlySummary struct {
	Year                 int                         `yaml:"year"`
	Month                int                         `yaml:"month"`
	WeeksOfYear          string                      `yaml:"weeksOfYear"`
	MonthlyGoals         []MonthlyGoalInfo           `yaml:"monthlyGoals"`
	MonthlySprints       []MonthlyGoalInfo           `yaml:"monthlySprints"`
	Events               []string                    `yaml:"events"`
	GoingWell            []WeeklySummaryItems        `yaml:"goingWell"`
	NeedsImprovement     []WeeklySummaryItems        `yaml:"needsImprovement"`
	Successes            []WeeklySummaryItems        `yaml:"successes"`
	Challenges           []WeeklySummaryItems        `yaml:"challenges"`
	MonthlyGoalSummaries []PerGoalWeeklySummaryItems `yaml:"monthlyGoalSummaries"`
}

// YearlySummary defines what's dumped for rolling up into the yearly plan summary
type YearlySummary struct {
	Year             int              `yaml:"year"`
	MonthlySummaries []MonthlySummary `yaml:"monthlySummaries"`
}

// WeeklyReviewData holds raw review content and extra info about the year, month, week it is from
type WeeklyReviewData struct {
	Year    int
	Month   int
	Week    int
	Content []byte
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
	year, week int,
) (
	month int,
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
	type MonthForWeek struct {
		year      int
		month     int
		weekBegin int
		weekEnd   int
	}

	// Monthly reviews don't always fall strictly after the end of the month.
	monthForWeekYearRange := []MonthForWeek{
		{2018, 1, 1, 4},
		{2018, 2, 5, 9},
		{2018, 3, 10, 13},
		{2018, 4, 14, 17},
		{2018, 5, 18, 22},
		{2018, 6, 23, 26},
		{2018, 7, 27, 30},
		{2018, 8, 31, 35},
		{2018, 9, 36, 39},
		{2018, 10, 40, 44},
		{2018, 11, 45, 48},
		{2018, 12, 49, 52},
		{2019, 1, 1, 5},
		{2019, 2, 6, 9},
		{2019, 3, 10, 13},
		{2019, 4, 14, 17},
	}

	// Although this would allow us to do review on monday or tuesday instead of
	// just sunday, shifting back -3 days results in daily email not reporting on
	// upcoming week until wednesday.
	// XXX: should make week an external parameter
	// year, week = time.Now().ISOWeek()
	month = 0
	for _, ymw := range monthForWeekYearRange {
		if year == ymw.year && week >= ymw.weekBegin && week <= ymw.weekEnd {
			month = ymw.month
			break
		}
	}
	// month = int(time.Now().Month())
	if month == 0 {
		log.Fatalf("no month mapping for year %d week %d", year, week)
	}
	return month, doneCards, goalCards, sprintCards, nil
}

// GetSummaryForWeek returns a summary structure usable by other downstream
// in-memory pipelines like daily reminder.
func GetSummaryForWeek(user, appkey, authtoken string, year, week int) (*WeeklySummary, error) {
	// XXX: should make week an external parameter so email reminder can obey calendar
	month, doneCards, goalCards, sprintCards, err := prepareSummaryForWeek(user, appkey, authtoken, year, week)
	if err != nil {
		return nil, err
	}
	return generateWeeklySummary(year, week, month, doneCards, goalCards, sprintCards)
}

// DumpSummaryForWeek dumps current content of Trello board to summary file for week
func DumpSummaryForWeek(user, appkey, authtoken string, year, week int, out io.Writer) error {
	// XXX: should make week an external parameter so review can lag a bit into the week
	month, doneCards, goalCards, sprintCards, err := prepareSummaryForWeek(user, appkey, authtoken, year, week)
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
	var titles []string
	for title := range allWeeklyGoals {
		titles = append(titles, title)
	}
	sort.Strings(titles)

	for _, title := range titles {
		goals := allWeeklyGoals[title]
		mg := MonthlyGoalInfo{Title: title,
			Created:     created[title],
			WeeklyGoals: goals,
		}
		monthlyGoals = append(monthlyGoals, mg)
	}
	return monthlyGoals
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

// GenerateSummaryForMonth rolls up weekly summaries to month level
func GenerateSummaryForMonth(user, appkey, authtoken string, year, month int, summaryIn [][]byte, reviewIn []WeeklyReviewData, out io.Writer) error {

	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	eventsList, err := listFor(cl.member, "Kanban daily/weekly", "Upcoming events")
	if err != nil {
		return err
	}

	eventCards, err := eventsList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	_ = eventCards

	summary := MonthlySummary{}
	summary.Year = year
	summary.Month = month

	// events
	// parse out event info, limit to particular month
	re := regexp.MustCompile(`^((?:\?|X)?\s*)(\d+)/(\d+)(?:/(\d+))?(?:-(\d+/\d+)(?:/(\d+))?)? (.*)$`)
	for _, event := range eventCards {
		expr := re.FindStringSubmatch(event.Name)
		if len(expr) > 0 {
			log.Printf("for %s got match %+v\n", event.Name, expr)
			// log.Println("got (opt) maybe/didn't do", expr[1])
			// log.Println("got month", expr[2])
			// log.Println("got day", expr[3])
			// log.Println("got (opt) year", expr[4])
			// log.Println("got (opt) end date", expr[5])
			// log.Println("got (opt) end date year", expr[6])
			// log.Println("got details", expr[7])
			mon, err2 := strconv.Atoi(expr[2])
			if err2 != nil {
				log.Println("error parsing month in", event.Name)
				summary.Events = append(summary.Events, event.Name)
				continue
			}
			if mon == month {
				summary.Events = append(summary.Events, event.Name)
			}
		} else {
			log.Printf("%s didn't parse\n", event.Name)
			summary.Events = append(summary.Events, event.Name)
		}
	}

	var weeklies []WeeklySummary
	var weeknums = []int{}
	for _, buf := range summaryIn {
		var weekly WeeklySummary
		err2 := yaml.Unmarshal(buf, &weekly)
		if err2 != nil {
			return err2
		}
		if weekly.Year == year && weekly.Month == month {
			weeklies = append(weeklies, weekly)
			weeknums = append(weeknums, weekly.Week)
		}
	}
	summary.WeeksOfYear = arrayToString(weeknums, ",")

	var goingWell []WeeklySummaryItems
	var needsImprovement []WeeklySummaryItems
	var successes []WeeklySummaryItems
	var challenges []WeeklySummaryItems
	goalReviews := make(map[string]PerGoalWeeklySummaryItems)

	for _, wrd := range reviewIn {
		var weekly WeeklyReview
		err2 := yaml.Unmarshal(wrd.Content, &weekly)
		if err2 != nil {
			return err2
		}
		if wrd.Year == year && wrd.Month == month {
			goingWell = append(goingWell, WeeklySummaryItems{
				Week:    wrd.Week,
				Content: weekly.GoingWell,
			})
			needsImprovement = append(needsImprovement, WeeklySummaryItems{
				Week:    wrd.Week,
				Content: weekly.NeedsImprovement,
			})
			successes = append(successes, WeeklySummaryItems{
				Week:    wrd.Week,
				Content: weekly.Successes,
			})
			challenges = append(challenges, WeeklySummaryItems{
				Week:    wrd.Week,
				Content: weekly.Challenges,
			})

			for _, goal := range weekly.PerGoalReviews {
				var title string
				var ginfo PerGoalWeeklySummaryItems
				if len(goal.DidToCreateOutcome) == 0 {
					title = "no goal specified"
					ginfo = goalReviews[title]
					ginfo.Goal = title
				} else {
					title = goal.DidToCreateOutcome[0]
					ginfo = goalReviews[title]
					ginfo.Goal = title
					ginfo.DidToCreateOutcome = append(ginfo.DidToCreateOutcome, WeeklySummaryItems{
						Week:    wrd.Week,
						Content: goal.DidToCreateOutcome[1:],
					})
				}
				if len(goal.KeepDoing) > 0 {
					ginfo.KeepDoing = append(ginfo.KeepDoing, WeeklySummaryItems{
						Week:    wrd.Week,
						Content: goal.KeepDoing,
					})
				}
				if len(goal.DoDifferently) > 0 {
					ginfo.DoDifferently = append(ginfo.DoDifferently, WeeklySummaryItems{
						Week:    wrd.Week,
						Content: goal.DoDifferently,
					})
				}
				goalReviews[title] = ginfo
			}
		}
	}

	if len(weeklies) > 0 {
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
		summary.GoingWell = goingWell
		summary.NeedsImprovement = needsImprovement
		summary.Successes = successes
		summary.Challenges = challenges
		for _, gs := range goalReviews {
			summary.MonthlyGoalSummaries = append(summary.MonthlyGoalSummaries, gs)
		}
	}

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
