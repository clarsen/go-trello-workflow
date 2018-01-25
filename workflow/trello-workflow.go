package workflow

import (
	"fmt"
	"log"
	"regexp"
	"strings"
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

var boards []*trello.Board
var boardmap = make(map[string]*trello.Board)
var listmap = make(map[string]map[string]*trello.List)

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
		fmt.Println("would move", c.Name, "to", destlist.Name)
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

func listFor(m *trello.Member, b string, l string) (list *trello.List, err error) {
	board, err := boardFor(m, b)
	if err != nil {
		// handle error
		return
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		// handle error
		return
	}
	for _, li := range lists {
		// fmt.Println("examining board ", b)
		// fmt.Println(b.Name, "->", b)
		// if listmap[board.Name] == nil {
		//   listmap[board.Name] = map[string]*trello.List{}
		// }
		if listmap[board.Name] == nil {
			listmap[board.Name] = map[string]*trello.List{}
		}
		listmap[board.Name][li.Name] = li
	}
	list = listmap[b][l]
	return
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
	s := time.Now().Format("(2006-01-02)")
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

func (cl *Client) doMinutely() error {
	dateBoardsAndLists := []boardAndList{
		{"Kanban daily/weekly", "Inbox"},
		{"Kanban daily/weekly", "Today"},
		{"Backlog (Personal)", "Backlog"},
	}
	cherryPickBoardsAndLists := []boardAndList{
		{"Backlog (Personal)", "Backlog"},
		{"Backlog (work)", "Backlog"},
		{"Periodic board", "Often"},
		{"Periodic board", "Weekly"},
		{"Periodic board", "Bi-weekly to monthly"},
		{"Periodic board", "Quarterly to Yearly"},
	}
	reorderBoardsAndLists := []boardAndList{
		{"Backlog (Personal)", "Backlog"},
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

	cherryPickDestlist, err := listFor(cl.member, "Kanban daily/weekly", "Today")
	if err != nil {
		return err
	}

	for _, boardlist := range cherryPickBoardsAndLists {
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
				fmt.Println("moving to top", card.Name, card.Labels)
				removeLabel(card, toTopLabel)
				card.MoveToTopOfList()
			} else if hasLabel(card, toSomedayMaybeLabel) {
				fmt.Println("moving to", somedayDestlist.Name)
				removeLabel(card, toSomedayMaybeLabel)
				card.MoveToListOnBoard(somedayDestlist.ID,
					somedayDestlist.IDBoard, trello.Defaults())
			} else if hasLabel(card, toDoneLabel) {
				fmt.Println("moving to", doneDestList.Name)
				removeLabel(card, toDoneLabel)
				card.MoveToListOnBoard(doneDestList.ID,
					doneDestList.IDBoard, trello.Defaults())

			}
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
	for _, boardlist := range []string{"Inbox", "Today"} {
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
	// fmt.Println("got", client)
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
