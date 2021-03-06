package workflow

import (
	"log"
	"time"

	"github.com/clarsen/trello"
)

type WFBoardList struct {
	Board string
	List  string
}

// WFTask wraps underlying trello card representation
type WFTask struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	Desc             string       `json:"desc"`
	CreatedDate      *time.Time   `json:"createdDate"`
	URL              *string      `json:"url"`
	Due              *time.Time   `json:"due"`
	List             *WFBoardList `json:"list"`
	Period           *string      `json:"period"`
	DateLastActivity *time.Time   `json:"dateLastActivity"`
	ChecklistItems   []string     `json:"checklistItems"`
}

// MoveToListOnBoard moves card to board/list
func (cl *Client) MoveToListOnBoard(cardID string, listID, boardID string) (*trello.Card, error) {
	card, err := cl.Client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	log.Printf("Move card=%+v to list=%+v/board=%+v\n", cardID, listID, boardID)
	card.MoveToListOnBoard(listID, boardID, trello.Defaults())
	card, err = cl.Client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return card, nil
}

// AddComment adds comment to card
func (cl *Client) AddComment(cardID string, comment string) (*WFTask, error) {
	card, err := cl.Client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	// log.Printf("Would add %+v to card %+v, lastactivity %+v", comment, card.Name, card.DateLastActivity)
	_, err = card.AddComment(comment, trello.Arguments{"fields": "all"})
	if err != nil {
		return nil, err
	}
	card, err = cl.Client.GetCard(cardID, trello.Arguments{"fields": "all"})
	return cl.wfTaskFor(card)
}

// SetDue sets due date/time of card
func (cl *Client) SetDue(cardID string, due time.Time) (*WFTask, error) {
	card, err := cl.Client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	args := trello.Defaults()
	args["due"] = due.Format(time.RFC3339)
	// card.Due = &due
	err = card.Update(args)
	if err != nil {
		return nil, err
	}
	card, err = cl.Client.GetCard(cardID, trello.Arguments{"fields": "all"})
	if err != nil {
		return nil, err
	}
	return cl.wfTaskFor(card)
}

func (cl *Client) wfTaskFor(card *trello.Card) (*WFTask, error) {
	local, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}
	createdDate := card.CreatedAt().In(local)
	url := card.ShortUrl
	// Strip out created date from title
	title, created, period := GetTitleAndAttributes(card)
	if created != nil {
		maybeCreatedDate, err := time.Parse("2006-01-02", *created)
		if err == nil {
			createdDate = maybeCreatedDate
		}
	}
	if len(card.IDCheckLists) > 0 {
		card, err = cl.Client.GetCard(card.ID, trello.Arguments{"fields": "all"})
		if err != nil {
			return nil, err
		}
	}
	var checklistItems []string
	for _, cl := range card.Checklists {
		// log.Println("checklist:", cl)
		for _, item := range cl.CheckItems {
			// log.Printf("Card %+v, checklist item: %+v\n", card.Name, item.Name)
			checklistItems = append(checklistItems, item.Name)
		}
	}
	wfTask := &WFTask{
		ID:               card.ID,
		Title:            title,
		CreatedDate:      &createdDate,
		URL:              &url,
		Due:              card.Due,
		Period:           period,
		DateLastActivity: card.DateLastActivity,
		Desc:             card.Desc,
		ChecklistItems:   checklistItems,
	}
	bl, err := BoardListFor(cl, card.IDBoard, card.IDList)
	if err != nil {
		return nil, err
	}
	wfTask.List = &WFBoardList{
		bl.Board,
		bl.List,
	}
	return wfTask, nil
}

// Tasks returns all open tasks, possibly limited to a particular board/list
func (cl *Client) Tasks(board, boardList *string) (tasks []WFTask, err error) {
	var list *trello.List
	if boardList != nil {
		list, err = ListFor(cl, *board, *boardList)
		if err != nil {
			// handle error
			return
		}
		cards, err2 := list.GetCards(trello.Defaults())
		if err2 != nil {
			err = err2
			// handle error
			return
		}
		for _, card := range cards {
			t, err2 := cl.wfTaskFor(card)
			if err2 != nil {
				err = err2
				return
			}
			t.List = &WFBoardList{
				Board: *board,
				List:  *boardList,
			}
			tasks = append(tasks, *t)
		}
	} else {
		// XXX: more efficient query?
		// get all tasks across all boards
		var trelloBoard *trello.Board
		for _, boardInfo := range AllBoards {
			log.Printf("Getting tasks for %+v\n", boardInfo.Board)
			trelloBoard, err = BoardFor(cl.Member, boardInfo.Board)
			if err != nil {
				// handle error
				return
			}
			cards, err2 := trelloBoard.GetCards(trello.Arguments{"fields": "all"})
			if err2 != nil {
				// handle error
				err = err2
				return
			}
			for _, card := range cards {
				// log.Printf("Card %+v", card)
				t, err2 := cl.wfTaskFor(card)
				if err2 != nil {
					err = err2
					return
				}
				var listName string
				boardlist := ListForID(card.IDList)
				if boardlist != nil {
					listName = boardlist.List
				}
				t.List = &WFBoardList{
					Board: boardInfo.Board,
					List:  listName,
				}

				tasks = append(tasks, *t)
			}
		}

	}
	if err != nil {
		return
	}

	return tasks, nil
}
