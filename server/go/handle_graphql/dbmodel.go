// expose trello API/cards as model to be exposed by graphql
package handle_graphql

import (
	"time"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/trello"
)

type MonthlyGoal struct {
	Title string // `json:"title"`
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

func SetTaskDue(taskId string, due time.Time) (*Task, error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	wfTask, err := cl.SetDue(taskId, due)
	if err != nil {
		return nil, err
	}
	task := Task{
		ID:          wfTask.ID,
		Title:       wfTask.Title,
		CreatedDate: wfTask.CreatedDate,
		URL:         wfTask.URL,
		Due:         wfTask.Due,
		Period:      wfTask.Period,
		List: &BoardList{
			wfTask.List.Board,
			wfTask.List.List,
		},
	}
	return &task, nil
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
			ID:          t.ID,
			Title:       t.Title,
			CreatedDate: t.CreatedDate,
			URL:         t.URL,
			Due:         t.Due,
			Period:      t.Period,
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
