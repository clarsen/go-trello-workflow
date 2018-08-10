package workflow

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/clarsen/trello"
)

// Client wraps logged in member
type Client struct {
	member *trello.Member
}

// Test does nothing
func (cl *Client) Test() {

}

type CardAttrs struct {
	DoneDate string
	Title    string
}
type DayCards struct {
	Date string
	Done []CardAttrs
}

type MonthlyGoal struct {
	Title             string
	WeeklyGoals       []string
	WeeklyGoalsByWeek []string
	LatestWeekGoals   []string
}

type MonthlyGoals struct {
	Month string
	Goals []MonthlyGoal
}

type MonthlySprints struct {
	Month   string
	Sprints []string
}

type WeeklyReport struct {
	Week                 string
	DoneByDay            []DayCards
	MonthlyGoalsByMonth  []MonthlyGoals
	LatestMonthlySprints MonthlySprints
	LatestMonthGoals     MonthlyGoals
	NowHHMM              string
}

var boards []*trello.Board
var boardmap = make(map[string]*trello.Board)
var listmap = make(map[string]map[string]*trello.List)

func reportMonthlyGoal(card *trello.Card) (MonthlyGoal, error) {
	var g MonthlyGoal
	g.Title = card.Name
	goalsForWeek := map[int][]string{}
	re := regexp.MustCompile(`^week (\d+): (.*) (?:\(\d{4}-\d{2}-\d{2}\))(?: (\(.*\)))?$`)
	// re := regexp.MustCompile(`^week (\d+): (.*)$`)

	for _, cl := range card.Checklists {
		// log.Println("checklist:", cl)
		for _, item := range cl.CheckItems {
			expr := re.FindStringSubmatch(item.Name)
			if len(expr) > 0 {
				// log.Printf("for %s got match %+v\n", item.Name, expr)
				// log.Println("got week", expr[1])
				// log.Println("got text", expr[2])
				// log.Println("got status", expr[3])
				week, err := strconv.Atoi(expr[1])
				if err != nil {
					return g, err
				}
				text := expr[2]
				status := expr[3]
				if _, ok := goalsForWeek[week]; !ok {
					goalsForWeek[week] = []string{}
				}
				if len(status) > 0 {
					goalsForWeek[week] = append(goalsForWeek[week], text+" "+status)
				} else {
					goalsForWeek[week] = append(goalsForWeek[week], text)
				}
			}
			g.WeeklyGoals = append(g.WeeklyGoals, item.Name)
		}
	}

	var sortedWeeks []int
	for week := range goalsForWeek {
		sortedWeeks = append(sortedWeeks, week)
	}
	sort.IntSlice.Sort(sortedWeeks)
	for _, week := range sortedWeeks {
		joined := strings.Join(goalsForWeek[week], ", ")
		str := fmt.Sprintf("week %d: %s", week, joined)
		g.WeeklyGoalsByWeek = append(g.WeeklyGoalsByWeek, str)
	}

	latestWeek := sortedWeeks[len(sortedWeeks)-1]
	for _, goal := range goalsForWeek[latestWeek] {
		str := fmt.Sprintf("week %d: %s", latestWeek, goal)
		g.LatestWeekGoals = append(g.LatestWeekGoals, str)
	}

	return g, nil
}

// MonthlyCleanup moves cards to history board
func MonthlyCleanup(user, appkey, authtoken string) error {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	year, month, _ := ymwForTime(time.Now())

	monthlyGoalsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return err
	}
	cards, err := monthlyGoalsList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	destGoalsListName := fmt.Sprintf("%s goals", month)
	destGoalsList, err := listForCreate(cl.member, fmt.Sprintf("History %d", year), destGoalsListName)
	if err != nil {
		return err
	}

	for _, card := range cards {
		log.Println("copying", card.Name, "to", destGoalsListName)
		card.CopyToList(destGoalsList.ID,
			trello.Arguments{"idBoard": destGoalsList.IDBoard, "pos": "bottom", "keepFromSource": "all"})
	}

	monthlySprintsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Sprints")
	if err != nil {
		return err
	}
	cards, err = monthlySprintsList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}

	destSprintsListName := fmt.Sprintf("%s sprints", month)
	destSprintsList, err := listForCreate(cl.member, fmt.Sprintf("History %d", year), destSprintsListName)
	if err != nil {
		return err
	}

	for _, card := range cards {
		log.Println("copying", card.Name, "to", destSprintsListName)
		card.CopyToList(destSprintsList.ID,
			trello.Arguments{"idBoard": destSprintsList.IDBoard, "pos": "bottom", "keepFromSource": "all"})
	}

	return nil
}

func ymwForTime(t time.Time) (int, string, int) {
	ref := t.Add(-time.Hour * 72)
	year, week := ref.ISOWeek()
	month := ref.Month().String()
	return year, month, week
}

// WeeklyCleanup moves cards to history board, copies periodic cards to history,
// moves periodic cards back
func WeeklyCleanup(user, appkey, authtoken string) error {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}
	doneList, err := listFor(cl.member, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return err
	}
	// log.Println(time.Now().Add(-time.Hour * 72))
	year, month, week := ymwForTime(time.Now())
	destListName := fmt.Sprintf("%02d %s", week, month)

	destList, err := listForCreate(cl.member, fmt.Sprintf("History %d", year), destListName)
	if err != nil {
		return err
	}
	destGoalsListName := fmt.Sprintf("%02d %s goals", week, month)
	destGoalsList, err := listForCreate(cl.member, fmt.Sprintf("History %d", year), destGoalsListName)
	if err != nil {
		return err
	}
	var _ = destList      // TODO: for real
	var _ = destGoalsList // TODO: for real

	cards, err := doneList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	for _, card := range cards {
		if isPeriodic(card) {
			log.Println("copying", card.Name, "to", destListName)
			card.CopyToList(destList.ID,
				trello.Arguments{"idBoard": destList.IDBoard, "pos": "bottom", "keepFromSource": "all"})
		} else {
			log.Println("moving", card.Name, "to", destListName)
			card.MoveToListOnBoard(destList.ID, destList.IDBoard, trello.Arguments{"pos": "bottom"})
		}
	}
	for _, card := range cards {
		if isPeriodic(card) {
			log.Println("moving", card.Name, "back to periodic")
			err2 := moveBackPeriodic(cl.member, card)
			if err2 != nil {
				return err2
			}
		}
	}
	monthlyGoalsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return err
	}
	cards, err = monthlyGoalsList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	for _, card := range cards {
		log.Println("copying", card.Name, "to", destGoalsListName)
		card.CopyToList(destGoalsList.ID,
			trello.Arguments{"idBoard": destGoalsList.IDBoard, "pos": "bottom", "keepFromSource": "all"})
	}

	return nil
}

// Weekly writes out report of weekly tasks done, goals, sprints
func Weekly(user, appkey, authtoken string) error {

	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	doneList, err := listFor(cl.member, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return err
	}

	cards, err := doneList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	report := WeeklyReport{}
	report.Week = "2018-01-28" // TODO: make correct week
	doneByDay := map[string][]CardAttrs{}

	for _, card := range cards {
		dateForSort := card.DateLastActivity.Format("2006-01-02")
		date := card.DateLastActivity.Format("2006-Jan-2 (Mon)")
		done := CardAttrs{
			DoneDate: date,
			Title:    card.Name,
		}
		if _, ok := doneByDay[dateForSort]; !ok {
			doneByDay[dateForSort] = []CardAttrs{}
		}
		doneByDay[dateForSort] = append(doneByDay[dateForSort], done)
	}
	var sortedDates []string
	for date := range doneByDay {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates)
	for _, d := range sortedDates {
		report.DoneByDay = append(report.DoneByDay,
			DayCards{Date: doneByDay[d][0].DoneDate, Done: doneByDay[d]})
	}

	// historyBoard, err := boardFor(cl.member, "History 2018")
	// if err != nil {
	// 	// handle error
	// 	return err
	// }

	months := []string{"January", "February", "March", "April",
		"May", "June", "July", "August", "September", "October",
		"November", "December",
	}

	var currentMonth string

	for _, month := range months {
		olderMonthlyGoals, err := listFor(cl.member, "History 2018", month+" goals")
		if err != nil {
			return err
		}
		if olderMonthlyGoals == nil {
			// assume this month is current
			currentMonth = month
			break
		}
		goalCards, err := olderMonthlyGoals.GetCards(trello.Defaults())
		if err != nil {
			return err
		}
		monthlyGoals := MonthlyGoals{}
		monthlyGoals.Month = month

		for _, goalCard := range goalCards {
			rmg, err := reportMonthlyGoal(goalCard)
			if err != nil {
				return err
			}
			monthlyGoals.Goals = append(monthlyGoals.Goals, rmg)
		}
		report.MonthlyGoalsByMonth = append(report.MonthlyGoalsByMonth, monthlyGoals)
	}

	monthlyGoalsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return err
	}
	goalCards, err := monthlyGoalsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return err
	}
	monthlyGoals := MonthlyGoals{}
	monthlyGoals.Month = currentMonth

	for _, goalCard := range goalCards {
		rmg, err := reportMonthlyGoal(goalCard)
		if err != nil {
			return err
		}
		monthlyGoals.Goals = append(monthlyGoals.Goals, rmg)
	}
	report.MonthlyGoalsByMonth = append(report.MonthlyGoalsByMonth, monthlyGoals)
	report.LatestMonthGoals = monthlyGoals

	monthlySprintsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Sprints")
	if err != nil {
		return err
	}
	sprintCards, err := monthlySprintsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return err
	}
	monthlySprints := MonthlySprints{}
	monthlySprints.Month = currentMonth

	for _, sprintCard := range sprintCards {
		monthlySprints.Sprints = append(monthlySprints.Sprints, sprintCard.Name)
	}
	report.LatestMonthlySprints = monthlySprints

	// log.Printf("report is %+v\n", report)
	report.NowHHMM = time.Now().Format("15:04")

	t, _ := template.ParseFiles("templates/weekly.md")
	t.Execute(os.Stdout, report)

	return nil
}

func sortList(m *trello.Member, list *trello.List) error {
	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	sort.Sort(byDue(cards))
	for idx, card := range cards {
		fmt.Printf("%f: %v\n", card.Pos, card.Name)
		card.SetPos(float64(idx)*100.0 + 1.0)
	}
	return nil
}

func boardFor(m *trello.Member, s string) (board *trello.Board, err error) {
	if boards == nil {
		boards, err = m.GetBoards(trello.Defaults())
		if err != nil {
			fmt.Println("error")
			return
			// Handle error
		}
		// fmt.Println("got", boards)
		for _, b := range boards {
			// fmt.Println("examining board ", b)
			// fmt.Println(b.Name, "->", b)
			boardmap[b.Name] = b
		}
		// fmt.Println(boardmap)
	}
	board = boardmap[s]
	return
}

func moveBackPeriodic(m *trello.Member, c *trello.Card) (err error) {
	periodicToList := map[string]string{
		"(po)":   "Often",
		"(p1w)":  "Weekly",
		"(p2w)":  "Bi-weekly to monthly",
		"(p4w)":  "Bi-weekly to monthly",
		"(p2m)":  "Quarterly to Yearly",
		"(p3m)":  "Quarterly to Yearly",
		"(p12m)": "Quarterly to Yearly",
	}

	var destlist *trello.List
	name := c.Name

	for substr, listname := range periodicToList {
		if strings.Contains(name, substr) {
			destlist, err = listFor(m, "Periodic board", listname)
			if err != nil {
				return
			}
			break
		}
	}
	if destlist != nil {
		fmt.Println("moving", c.Name, "to", destlist.Name)
		c.MoveToListOnBoard(destlist.ID, destlist.IDBoard, trello.Defaults())
		c.MoveToTopOfList()
	}
	return
}

func moveBackCard(m *trello.Member, c *trello.Card) (err error) {
	var destlist *trello.List

	// first do periodics which also have personal/work labels

	for _, label := range c.Labels {
		switch label.Name {
		case "Periodic":
			fmt.Println("  ", label.Name, label.Color)
			return moveBackPeriodic(m, c)
		}
	}

	for _, label := range c.Labels {
		fmt.Println("  ", label.Name, label.Color)
		switch label.Name {
		case "Personal":
			destlist, err = listFor(m, "Backlog (Personal)", "Backlog")
			if err != nil {
				// handle error
				return err
			}
		case "Process":
			destlist, err = listFor(m, "Backlog (Personal)", "Backlog")
			if err != nil {
				// handle error
				return err
			}
		case "Work":
			destlist, err = listFor(m, "Backlog (work)", "Backlog")
			if err != nil {
				// handle error
				return err
			}
		}
	}
	if destlist != nil {
		fmt.Println("backlog is", destlist)
		fmt.Println("would move", c.Name, "to", destlist.Name)
		c.MoveToListOnBoard(destlist.ID, destlist.IDBoard, trello.Defaults())
		c.MoveToTopOfList()
	}
	return
}

func listForCreate(m *trello.Member, b string, l string) (*trello.List, error) {
	list, err := listFor(m, b, l)
	if list != nil {
		return list, err
	}
	board, err := boardFor(m, b)
	if err != nil {
		return nil, err
	}
	list, err = board.CreateList(l, trello.Arguments{"pos": "bottom"})
	if err != nil {
		return nil, err
	}
	listmap[b][l] = list
	return list, nil
}

func listFor(m *trello.Member, b string, l string) (*trello.List, error) {
	if list, ok := listmap[b][l]; ok {
		return list, nil
	}

	board, err := boardFor(m, b)
	if err != nil {
		// handle error
		return nil, err
	}
	if board == nil {
		return nil, fmt.Errorf("Board %s not found", b)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		// handle error
		return nil, err
	}
	for _, li := range lists {
		// fmt.Println("examining board ", b)
		// fmt.Println(b.Name, "->", b)
		// if listmap[board.Name] == nil {
		//   listmap[board.Name] = map[string]*trello.List{}
		// }
		// fmt.Println("list ", li.Name)
		if listmap[board.Name] == nil {
			listmap[board.Name] = map[string]*trello.List{}
		}
		listmap[board.Name][li.Name] = li
	}
	list := listmap[b][l]
	return list, nil
}

func hasDate(card *trello.Card) bool {
	re := regexp.MustCompile("\\(\\d{4}-\\d{2}-\\d{2}\\)")
	date := re.FindString(card.Name)
	return date != ""
}

func isPeriodic(card *trello.Card) bool {
	re := regexp.MustCompile("\\((po|p1w|p2w|p4w|p2m|p3m|p12m)\\)")
	period := re.FindString(card.Name)
	return period != ""
}

func addDateToName(card *trello.Card) {
	log.Println("Add date to ", card.Name)
	local, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Println("Can't find timezone")
	}
	s := time.Now().In(local).Format("(2006-01-02)")
	card.Update(trello.Arguments{"name": card.Name + " " + s})
}

type boardAndList struct {
	Board string
	List  string
}

const (
	cherryPickLabel     = "orange"
	toTopLabel          = "sky"
	toSomedayMaybeLabel = "lime"
	toDoneLabel         = "pink"
)

func removeLabel(card *trello.Card, color string) {
	for _, label := range card.Labels {
		if label.Color == color {
			card.RemoveLabel(label.ID)
			return
		}
	}
	return
}

func hasLabel(card *trello.Card, color string) bool {
	for _, label := range card.Labels {
		if label.Color == color {
			return true
		}
	}
	return false
}

func titleStartsWith(card *trello.Card, prefix string) bool {
	return strings.HasPrefix(card.Name, prefix)
}

func (cl *Client) doMinutely() error {
	dateBoardsAndLists := []boardAndList{
		{"Kanban daily/weekly", "Inbox"},
		{"Kanban daily/weekly", "Today"},
		{"Backlog (Personal)", "Backlog"},
	}
	cherryPickBoardsAndLists := []boardAndList{
		{"Backlog (Personal)", "Backlog"},
		{"Backlog (Personal)", "Projects"},
		{"Backlog (work)", "Backlog"},
		{"Periodic board", "Often"},
		{"Periodic board", "Weekly"},
		{"Periodic board", "Bi-weekly to monthly"},
		{"Periodic board", "Quarterly to Yearly"},
	}
	reorderBoardsAndLists := []boardAndList{
		{"Kanban daily/weekly", "Waiting on"},
		{"Backlog (Personal)", "Backlog"},
		{"Backlog (Personal)", "Projects"},
		{"Backlog (Personal)", "Projects: delegated"},
		{"Backlog (Personal)", "Projects - Soon"},
		{"Backlog (Personal)", "Projects (not yet)"},
		{"Backlog (Personal)", "Area: Finance"},
		{"Backlog (Personal)", "Area: Friends"},
		{"Backlog (work)", "Backlog"},
		{"Periodic board", "Often"},
		{"Periodic board", "Weekly"},
		{"Periodic board", "Bi-weekly to monthly"},
		{"Periodic board", "Quarterly to Yearly"},
	}

	for _, boardlist := range dateBoardsAndLists {
		list, err := listFor(cl.member, boardlist.Board, boardlist.List)
		if err != nil {
			// handle error
			return err
		}

		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			// handle error
			return err
		}
		for _, card := range cards {

			if !hasDate(card) && !isPeriodic(card) {
				addDateToName(card)
			}
		}
	}

	// Move cards from inbox to backlog
	for _, boardlist := range []string{"Inbox"} {
		// fmt.Printf("move items from %s to backlog based on label color\n", boardlist)

		list, err := listFor(cl.member, "Kanban daily/weekly", boardlist)
		if err != nil {
			// handle error
			return err
		}
		// fmt.Printf("kanban/%s is %v", boardlist, list)
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			// handle error
			return err
		}

		somedaylist, err := listFor(cl.member, "Someday/Maybe", "Maybe")
		if err != nil {
			return err
		}

		for _, card := range cards {
			if titleStartsWith(card, "? ") {
				fmt.Printf("move %s %+v to someday/maybe\n", card.Name, card.Labels)
				card.Update(trello.Arguments{"name": strings.TrimPrefix(card.Name, "? ")})
				card.MoveToListOnBoard(somedaylist.ID, somedaylist.IDBoard, trello.Defaults())
				card.MoveToTopOfList()

			} else {
				fmt.Printf("move %s %+v to backlog\n", card.Name, card.Labels)
				moveBackCard(cl.member, card)
			}
		}
	}

	cherryPickDestlist, err := listFor(cl.member, "Kanban daily/weekly", "Today")
	if err != nil {
		return err
	}

	for _, boardlist := range cherryPickBoardsAndLists {
		list, err2 := listFor(cl.member, boardlist.Board, boardlist.List)
		if err2 != nil {
			// handle error
			return err2
		}

		cards, err2 := list.GetCards(trello.Defaults())
		if err2 != nil {
			// handle error
			return err2
		}
		for _, card := range cards {

			if hasLabel(card, cherryPickLabel) {
				fmt.Println("cherry picking", card.Name, card.Labels)
				removeLabel(card, cherryPickLabel)
				card.MoveToListOnBoard(cherryPickDestlist.ID,
					cherryPickDestlist.IDBoard, trello.Defaults())
			}
		}
	}

	somedayDestlist, err := listFor(cl.member, "Someday/Maybe", "Maybe")
	if err != nil {
		return err
	}
	doneDestList, err := listFor(cl.member, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return err
	}

	for _, boardlist := range reorderBoardsAndLists {
		list, err := listFor(cl.member, boardlist.Board, boardlist.List)
		if err != nil {
			// handle error
			return err
		}

		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			// handle error
			return err
		}
		for _, card := range cards {

			if hasLabel(card, toTopLabel) {
				fmt.Println("moving", card.Name, card.Labels, "to top")
				removeLabel(card, toTopLabel)
				card.MoveToTopOfList()
			} else if hasLabel(card, toSomedayMaybeLabel) {
				fmt.Println("moving", card.Name, card.Labels, "to", somedayDestlist.Name)
				removeLabel(card, toSomedayMaybeLabel)
				card.MoveToListOnBoard(somedayDestlist.ID,
					somedayDestlist.IDBoard, trello.Defaults())
			} else if hasLabel(card, toDoneLabel) {
				fmt.Println("moving", card.Name, card.Labels, "to", doneDestList.Name)
				removeLabel(card, toDoneLabel)
				card.MoveToListOnBoard(doneDestList.ID,
					doneDestList.IDBoard, trello.Defaults())

			}
		}
		err = sortList(cl.member, list)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareToday moves cards back to their respective boards at end of day
func (cl *Client) prepareToday() error {
	board, err := boardFor(cl.member, "Kanban daily/weekly")
	if err != nil {
		// handle error
		return err
	}
	fmt.Println("Kanban board is ", board.ID)
	for _, boardlist := range []string{"Today"} {
		fmt.Printf("move items from %s to backlog based on label color\n", boardlist)

		list, err := listFor(cl.member, "Kanban daily/weekly", boardlist)
		if err != nil {
			// handle error
			return err
		}
		fmt.Printf("kanban/%s is %v", boardlist, list)
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			// handle error
			return err
		}
		for _, card := range cards {
			fmt.Println(card.Name, card.Labels)
			moveBackCard(cl.member, card)
		}
	}
	return nil
}

// MinutelyMaintenance does things like cherry picking, moving, adding dates to
// titles
func MinutelyMaintenance(user, appkey, authtoken string) {
	// log.Println("Running MinutelyMaintenance")
	wf, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}
	wf.doMinutely()
}

// DailyMaintenance moves cards in kanban board back to their homes at the end
// of the day
func DailyMaintenance(user, appkey, authtoken string) {
	log.Println("Running DailyMaintenance")
	wf, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}
	wf.prepareToday()
	log.Println("Finished running DailyMaintenance")
}

// New create new client
func New(user string, appKey string, token string) (c *Client, err error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	client := trello.NewClient(appKey, token)
	client.Logger = logger

	member, err := client.GetMember(user, trello.Defaults())
	if err != nil {
		// Handle error
		return nil, err
	}
	c = &Client{
		member,
	}
	return
}
