package workflow

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/clarsen/trello"
	"gopkg.in/yaml.v2"
)

type TaskInfo struct {
	CreatedDate string `yaml:"createdDate"`
	DoneDate    string `yaml:"doneDate"`
	Title       string `yaml:"title"`
}

type WeeklyGoalInfo struct {
	Title   string `yaml:"title"`
	Created string `yaml:"createdDate"`
	Status  string `yaml:"status"`
	Week    int    `yaml:"weekNumber"`
	Year    int    `yaml:"year"`
}

type MonthlyGoalInfo struct {
	Title       string           `yaml:"title"`
	Created     string           `yaml:"createdDate"`
	WeeklyGoals []WeeklyGoalInfo `yaml:"weeklyGoals"`
}

type WeeklySummary struct {
	Year           int               `yaml:"year"`
	Week           int               `yaml:"weekNumber"`
	Done           []TaskInfo        `yaml:"doneThisWeek"`
	MonthlyGoals   []MonthlyGoalInfo `yaml:"monthlyGoals"`
	MonthlySprints []MonthlyGoalInfo `yaml:"monthlySprints"`
}

func dumpMonthlyGoalInfo(card *trello.Card, weekNumber *int) (MonthlyGoalInfo, error) {
	var g MonthlyGoalInfo
	g.Title = card.Name

	reMDY := regexp.MustCompile(`\((\d{4}-\d{2}-\d{2})\)`)
	exprGoalMDY := reMDY.FindStringSubmatch(card.Name)
	if len(exprGoalMDY) <= 0 {
		return g, errors.New("Couldn't parse date out of title")
	}
	g.Created = exprGoalMDY[1]

	re := regexp.MustCompile(`^week (\d+): (.*) (?:\((\d{4}-\d{2}-\d{2})\))(?: (\(.*\)))?$`)
	reYear := regexp.MustCompile(`^(\d{4})-\d{2}-\d{2}$`)

	for _, cl := range card.Checklists {
		// log.Println("checklist:", cl)
		for _, item := range cl.CheckItems {
			expr := re.FindStringSubmatch(item.Name)
			if len(expr) > 0 {
				// log.Printf("for %s got match %+v\n", item.Name, expr)
				// log.Println("got week", expr[1])
				// log.Println("got text", expr[2])
				// log.Println("got created", expr[3])
				// log.Println("got status", expr[4])
				week, err := strconv.Atoi(expr[1])

				if weekNumber != nil && week < *weekNumber {
					continue
				}
				if err != nil {
					return g, err
				}
				text := expr[2]
				created := expr[3]
				exprYear := reYear.FindStringSubmatch(created)
				if len(exprYear) <= 0 {
					continue
				}
				year, err := strconv.Atoi(exprYear[1])
				if err != nil {
					return g, err
				}
				status := expr[4]
				wgi := WeeklyGoalInfo{
					Year:    year,
					Week:    week,
					Created: created,
					Title:   text,
					Status:  status,
				}
				g.WeeklyGoals = append(g.WeeklyGoals, wgi)
			}
		}
	}

	return g, nil
}

// DumpSummaryForWeek dumps current content of Trello board to summary file for week
func DumpSummaryForWeek(user, appkey, authtoken string, out io.Writer) error {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		log.Fatal(err)
	}

	doneList, err := listFor(cl.member, "Kanban daily/weekly", "Done this week")
	if err != nil {
		return err
	}

	cards, err := doneList.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}
	summary := WeeklySummary{}
	summary.Year, summary.Week = time.Now().ISOWeek()

	for _, card := range cards {
		summary.Done = append(summary.Done, TaskInfo{
			DoneDate:    card.DateLastActivity.Format("2006-01-02"),
			Title:       card.Name,
			CreatedDate: card.CreatedAt().Format("2006-01-02"),
		})
	}

	monthlyGoalsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Goals")
	if err != nil {
		return err
	}
	goalCards, err := monthlyGoalsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return err
	}

	for _, goalCard := range goalCards {
		dmgi, err := dumpMonthlyGoalInfo(goalCard, &summary.Week)
		if err != nil {
			return err
		}
		summary.MonthlyGoals = append(summary.MonthlyGoals, dmgi)

	}

	monthlySprintsList, err := listFor(cl.member, "Kanban daily/weekly", "Monthly Sprints")
	if err != nil {
		return err
	}
	sprintCards, err := monthlySprintsList.GetCards(trello.Arguments{"fields": "all"})
	if err != nil {
		// handle error
		return err
	}

	for _, sprintCard := range sprintCards {
		dmsi, err := dumpMonthlyGoalInfo(sprintCard, &summary.Week)
		if err != nil {
			return err
		}
		summary.MonthlySprints = append(summary.MonthlySprints, dmsi)

	}

	buf, err := yaml.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}

	_, err = out.Write(buf)
	return err
}
