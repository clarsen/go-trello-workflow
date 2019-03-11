// expose trello API/cards as model to be exposed by graphql
package handle_graphql

import (
	"log"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/trello"
)

// type Model_TaskInfo struct {
// 	ID    string
// 	Title string
// }

func GetTasks(user, appkey, authtoken string,
) (tasks []Task, err error) {
	cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	list, err := workflow.ListFor(cl, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return
	}

	cards, err := list.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return
	}
	for _, card := range cards {
		tasks = append(tasks,
			Task{
				ID:    card.ID,
				Title: card.Name,
			})
	}
	return tasks, nil
}
