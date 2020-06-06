// expose trello API/cards as model to be exposed by graphql
package handle_graphql

import (
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/trello"
)

type MonthlyGoal struct {
	IDCard string
	Title  string // `json:"title"`
	// WeeklyGoals []WeeklyGoal `json:"weeklyGoals"`
	wfGoal *workflow.WFMonthlyGoal
}

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

// AddComment adds comment to card
func AddComment(taskID string, comment string) (*Task, error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	wfTask, err := cl.AddComment(taskID, comment)
	if err != nil {
		return nil, err
	}
	// XXX: refactor
	task := Task{
		ID:               wfTask.ID,
		Title:            wfTask.Title,
		CreatedDate:      wfTask.CreatedDate,
		URL:              wfTask.URL,
		Due:              wfTask.Due,
		Period:           wfTask.Period,
		DateLastActivity: wfTask.DateLastActivity,
		Desc:             wfTask.Desc,
		ChecklistItems:   wfTask.ChecklistItems,
		List: &BoardList{
			wfTask.List.Board,
			wfTask.List.List,
		},
	}
	return &task, nil
}

func SetTaskDue(taskID string, due time.Time) (*Task, error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	wfTask, err := cl.SetDue(taskID, due)
	if err != nil {
		return nil, err
	}
	// XXX: refactor
	task := Task{
		ID:               wfTask.ID,
		Title:            wfTask.Title,
		CreatedDate:      wfTask.CreatedDate,
		URL:              wfTask.URL,
		Due:              wfTask.Due,
		Period:           wfTask.Period,
		DateLastActivity: wfTask.DateLastActivity,
		Desc:             wfTask.Desc,
		ChecklistItems:   wfTask.ChecklistItems,
		List: &BoardList{
			wfTask.List.Board,
			wfTask.List.List,
		},
	}
	return &task, nil
}

// SetTaskDone moves task to done this week list or back to inbox (depending on done) and optionally prepends a title comment
func SetTaskDone(taskID string, done bool, titleComment *string) (*Task, error) {
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
	card, err := cl.MoveToListOnBoard(taskID, targetList.ID, targetList.IDBoard)
	if err != nil {
		return nil, err
	}

	if titleComment != nil {
		args := trello.Defaults()
		args["name"] = (*titleComment) + card.Name
		err = card.Update(args)
		if err != nil {
			return nil, err
		}
	}

	_, _, period := workflow.GetTitleAndAttributes(card)
	if period == nil {
		// mark Due done if present
		err = card.Update(trello.Arguments{"dueComplete": "true"})
		if err != nil {
			return nil, err
		}
	}
	return TaskFor(card)
}

func GetTasks(cl *workflow.Client,
	boardlist *BoardListInput,
) (tasks []Task, err error) {
	var wfTasks []workflow.WFTask
	if boardlist != nil {
		wfTasks, err = cl.Tasks(&boardlist.Board, &boardlist.List)
		if err != nil {
			return
		}
	} else {
		wfTasks, err = cl.Tasks(nil, nil)
		if err != nil {
			return
		}
	}
	for _, t := range wfTasks {
		tasks = append(tasks, Task{
			ID:               t.ID,
			Title:            t.Title,
			CreatedDate:      t.CreatedDate,
			URL:              t.URL,
			Due:              t.Due,
			Period:           t.Period,
			DateLastActivity: t.DateLastActivity,
			Desc:             t.Desc,
			ChecklistItems:   t.ChecklistItems,
			List: &BoardList{
				t.List.Board,
				t.List.List,
			},
		})
	}
	return tasks, nil
}

func GetMonthlyGoals(cl *workflow.Client) (goals []MonthlyGoal, err error) {
	wfGoals, err := cl.MonthlyGoals()
	if err != nil {
		return goals, err
	}
	for idx, _ := range wfGoals {
		g := wfGoals[idx]
		goals = append(goals, MonthlyGoal{
			IDCard: g.IDCard,
			Title:  g.Title,
			wfGoal: &g,
		})
	}
	return goals, nil
}

func MonthlyGoalToWeeklyGoals(mg *MonthlyGoal) []WeeklyGoal {
	wfGoals := workflow.WeeklyGoals(mg.wfGoal)
	goals := []WeeklyGoal{}
	for _, wg := range wfGoals {
		goals = append(goals, WeeklyGoal{
			IDCard:      wg.IDCard,
			IDCheckitem: wg.IDCheckitem,
			Title:       wg.Title,
			Year:        wg.Year,
			Month:       wg.Month,
			Week:        wg.Week,
			Tasks:       []*Task{},
			Done:        wg.Done,
			Status:      wg.Status,
		})
	}
	return goals

}
