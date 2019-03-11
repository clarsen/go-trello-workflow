// expose trello API/cards as model to be exposed by graphql
package handle_graphql

import (
	"log"
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/trello"
)

// type Model_TaskInfo struct {
// 	ID    string
// 	Title string
// }

func TaskFor(card *trello.Card) (*Task, error) {
	local, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}
	createdDate := card.CreatedAt().In(local)
	url := card.ShortUrl
	return &Task{
		ID:          card.ID,
		Title:       card.Name,
		CreatedDate: &createdDate,
		URL:         &url,
	}, nil
}

func GetTasks(user, appkey, authtoken string,
	boardlist *BoardList,
) (tasks []Task, err error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}
	var list *trello.List
	if boardlist != nil {
		list, err = workflow.ListFor(cl, boardlist.Board, boardlist.List)
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
			t, err2 := TaskFor(card)
			if err2 != nil {
				err = err2
				return
			}
			tasks = append(tasks, *t)
		}
	} else {
		// get all tasks across all boards
		for _, bl := range workflow.AllLists {
			list, err = workflow.ListFor(cl, bl.Board, bl.List)
			if err != nil {
				// handle error
				return
			}
			cards, err2 := list.GetCards(trello.Defaults())
			if err2 != nil {
				// handle error
				err = err2
				return
			}
			for _, card := range cards {
				t, err2 := TaskFor(card)
				if err2 != nil {
					err = err2
					return
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
