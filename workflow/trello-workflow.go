package workflow

import (
	"fmt"
	"log"
	"strings"

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

// PrepareToday moves cards back to their respective boards at end of day
func (cl *Client) PrepareToday() error {
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

// DailyMaintenance moves cards in kanban board back to their homes at the end
// of the day
func DailyMaintenance(user, appkey, authtoken string) {
	log.Println("Running dailyMaintenance")
	wf, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}
	wf.PrepareToday()
	log.Println("Finished running dailyMaintenance")
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
