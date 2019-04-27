package workflow

import (
	"fmt"
	"log"
	"time"

	"github.com/clarsen/trello"
)

// WFMonthlyGoal is a wrapper for underlying trello card representation
type WFMonthlyGoal struct {
	IDCard string
	Title  string // `json:"title"`
	// WeeklyGoals []WeeklyGoal `json:"weeklyGoals"`
	card *trello.Card
}

// WFWeeklyGoal is a wrapper for underlying trello card checklist item
type WFWeeklyGoal struct {
	IDCard      string  `json:"idCard"`
	IDCheckitem string  `json:"idCheckitem"`
	Title       string  `json:"title"`
	Year        *int    `json:"year"`
	Month       *int    `json:"month"`
	Week        *int    `json:"week"`
	Done        *bool   `json:"done"`
	Status      *string `json:"status"`
}

func wfMonthlyGoalFor(card *trello.Card) (WFMonthlyGoal, error) {
	title, _, _ := GetTitleAndAttributes(card)
	return WFMonthlyGoal{
		Title:  title,
		IDCard: card.ID,
	}, nil
}

func (cl *Client) AddWeeklyGoal(mg *WFMonthlyGoal, title string, week int) error {
	card, err := cl.Client.GetCard(mg.IDCard, trello.Defaults())
	if err != nil {
		return err
	}

	created := time.Now()
	cstr := created.Format(" (2006-01-02)")
	wgTitle := fmt.Sprintf("week %d: %s%s", week, title, cstr)

	err = card.Checklists[0].AddCheckItem(wgTitle)
	return err
}

func (cl *Client) AddMonthlyGoal(title string) error {
	card, err := cl.CreateCard(title, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return err
	}

	err = card.AddChecklist("Weekly goals")
	return err
}

func (cl *Client) GetMonthlyGoal(id string) (*WFMonthlyGoal, error) {
	card, err := cl.Client.GetCard(id, trello.Defaults())
	if err != nil {
		return nil, err
	}
	g, err := wfMonthlyGoalFor(card)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// MonthlyGoals returns all goals for this month
func (cl *Client) MonthlyGoals() (goals []WFMonthlyGoal, err error) {
	list, err2 := ListFor(cl, "Kanban daily/weekly", "Monthly Goals")
	if err2 != nil {
		err = err2
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
		g, err3 := wfMonthlyGoalFor(card)
		if err3 != nil {
			err = err3
			return
		}
		g.card = card
		goals = append(goals, g)
	}
	list, err2 = ListFor(cl, "Kanban daily/weekly", "Monthly Sprints")
	if err2 != nil {
		err = err2
		// handle error
		return
	}
	cards, err2 = list.GetCards(trello.Defaults())
	if err2 != nil {
		err = err2
		// handle error
		return
	}
	for _, card := range cards {
		g, err2 := wfMonthlyGoalFor(card)
		if err2 != nil {
			err = err2
			return
		}
		g.card = card
		goals = append(goals, g)
	}

	return goals, nil
}

// WeeklyGoals returns WFWeeklyGoals for a given monthly goal
func WeeklyGoals(mg *WFMonthlyGoal) []WFWeeklyGoal {
	goals := []WFWeeklyGoal{}
	for _, cl := range mg.card.Checklists {
		log.Println("checklist:", cl)
		for _, item := range cl.CheckItems {
			title, week, created, _, status := GetAttributesFromChecklistTitle(item.Name)
			month := int(created.Month())
			year, _ := created.ISOWeek()
			done := item.State == "complete"
			wg := WFWeeklyGoal{
				IDCard:      item.IDCard,
				IDCheckitem: item.ID,
				Title:       title,
				Year:        &year,
				Month:       &month,
				Week:        &week,
				Done:        &done,
				Status:      status,
			}
			goals = append(goals, wg)
		}
	}

	return goals
}
