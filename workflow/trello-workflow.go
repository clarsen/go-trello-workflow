package workflow

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/clarsen/trello"
	"github.com/urfave/cli"
)

// Client wraps logged in member
type Client struct {
	member *trello.Member
}

// Test does nothing
func (c *Client) Test() {

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

func moveBackCard(m *trello.Member, c *trello.Card) (err error) {
	var destlist *trello.List
	for _, label := range c.Labels {
		fmt.Println("  ", label.Name, label.Color)
		switch label.Name {
		case "Personal":
			destlist, err = listFor(m, "Backlog (Personal)", "Backlog")
			// destlist, err = listFor(m, "Kanban daily/weekly", "Today")
			if err != nil {
				// handle error
				return err
			}
			fmt.Println("backlog is", destlist)
		}
	}
	if destlist != nil {
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

func (cl *Client) doToday(c *cli.Context) error {
	board, err := boardFor(cl.member, "Kanban daily/weekly")
	if err != nil {
		// handle error
		return err
	}
	fmt.Println("Kanban board is ", board.ID)
	fmt.Println("move items from Inbox to backlog based on label color")
	list, err := listFor(cl.member, "Kanban daily/weekly", "Inbox")
	if err != nil {
		// handle error
		return err
	}
	fmt.Println("kanban/inbox is ", list)
	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	for _, card := range cards {
		fmt.Println(card.Name, card.Labels)
		moveBackCard(cl.member, card)
	}
	return nil
}

// New create new client
func New(user string, appKey string, token string) (c *Client, err error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	client := trello.NewClient(appKey, token)
	client.Logger = logger
	fmt.Println("got", client)
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
