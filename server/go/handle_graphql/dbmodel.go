// expose trello API/cards as model to be exposed by graphql
package handle_graphql

import (
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
	// Strip out created date from title
	title, created, period := workflow.GetTitleAndAttributes(card)
	if created != nil {
		maybeCreatedDate, err := time.Parse("2006-01-02", *created)
		if err == nil {
			createdDate = maybeCreatedDate
		}
	}

	return &Task{
		ID:          card.ID,
		Title:       title,
		CreatedDate: &createdDate,
		URL:         &url,
		Due:         card.Due,
		Period:      period,
	}, nil
}

func SetTaskDue(taskId string, due time.Time) (*Task, error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	card, err := cl.SetDue(taskId, due)
	if err != nil {
		return nil, err
	}
	return TaskFor(card)
}

func SetTaskDone(taskId string, done bool) (*Task, error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	var targetList *trello.List
	if done {
		targetList, err = workflow.ListFor(cl, "Kanban daily/weekly", "Done this week")
		if err != nil {
			return nil, err
		}
	} else {
		targetList, err = workflow.ListFor(cl, "Kanban daily/weekly", "Inbox")
		if err != nil {
			return nil, err
		}
	}
	card, err := cl.MoveToListOnBoard(taskId, targetList.ID, targetList.IDBoard)
	if err != nil {
		return nil, err
	}
	return TaskFor(card)
}

func GetTasks(user, appkey, authtoken string,
	boardlist *BoardListInput,
) (tasks []Task, err error) {
	cl, err2 := workflow.New(user, appkey, authtoken)
	if err2 != nil {
		err = err2
		return
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
			t.List = &BoardList{
				Board: boardlist.Board,
				List:  boardlist.List,
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
				t.List = &BoardList{
					Board: bl.Board,
					List:  bl.List,
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
