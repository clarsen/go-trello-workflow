//go:generate go run github.com/99designs/gqlgen
package handle_graphql

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	yaml "gopkg.in/yaml.v2"

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/gtoggl-api/gthttp"
	"github.com/clarsen/gtoggl-api/gttimeentry"
	"github.com/clarsen/gtoggl-api/gtworkspace"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var appkey, authtoken, user string
var cl *workflow.Client
var teCL *gttimeentry.TimeEntryClient
var togglWorkspace *gtworkspace.Workspace

const (
	// SummaryDir is where dumped data about tasks is written
	SummaryDir = "task-summary"
	// ReviewDir is where templates are written and updated by humans.
	ReviewDir = "reviews"
	// ReviewVisDir is where templates from ReviewDir are processed into human friendly documents
	ReviewVisDir = "visualized_reviews"
	debug        = false
)

func init() {
	appkey = os.Getenv("appkey")
	if appkey == "" {
		log.Fatal("$appkey must be set")
	}
	authtoken = os.Getenv("authtoken")
	if authtoken == "" {
		log.Fatal("$authtoken must be set")
	}
	user = os.Getenv("user")
	if user == "" {
		log.Fatal("$user must be set")
	}
	_cl, err := workflow.New(user, appkey, authtoken)
	if err != nil {
		log.Fatal("Couldn't create new client")
	}
	cl = _cl
	togglAPIKey := os.Getenv("togglAPIKey")
	if togglAPIKey == "" {
		log.Fatal("$togglAPIKey must be set")
	}

	var _thc *gthttp.TogglHttpClient
	if debug {
		logger := log.New(os.Stderr, "gthttp.trace: ", log.LstdFlags|log.Llongfile)
		_thc, err = gthttp.NewClient(togglAPIKey, gthttp.SetTraceLogger(logger))
		if err != nil {
			log.Fatal("Couldn't create new gthttp client")
		}
	} else {
		_thc, err = gthttp.NewClient(togglAPIKey)
		if err != nil {
			log.Fatal("Couldn't create new gthttp client")
		}
	}

	_tec := gttimeentry.NewClient(_thc)
	teCL = _tec

	_tew := gtworkspace.NewClient(_thc)

	wlist, err := _tew.List()
	if err != nil {
		log.Fatal("Couldn't get workspaces")
	}
	for _, w := range wlist {
		log.Printf("%d: %s\n", w.Id, w.Name)
		togglWorkspace = &w
	}
}

type Resolver struct{}

func (r *Resolver) MonthlyGoal() MonthlyGoalResolver {
	return &monthlyGoalResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type monthlyGoalResolver struct{ *Resolver }

func (r *monthlyGoalResolver) WeeklyGoals(ctx context.Context, obj *MonthlyGoal) ([]WeeklyGoal, error) {
	goals := MonthlyGoalToWeeklyGoals(obj)
	return goals, nil
}

type mutationResolver struct{ *Resolver }

type workdir struct {
	fs       billy.Filesystem
	repo     *git.Repository
	worktree *git.Worktree
}

func getData() (*workdir, error) {
	// Filesystem abstraction based on memory
	fs := memfs.New()

	// Git objects storer based on memory
	storer := memory.NewStorage()
	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: "https://github.com/clarsen/data-and-reviews.git",
		Auth: &http.BasicAuth{
			Username: "abc123", // anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	if err != nil {
		return nil, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	return &workdir{
		fs:       fs,
		repo:     repo,
		worktree: w,
	}, nil
}

func (wd *workdir) commitAndPushData(message string) error {
	// commit
	commit, err := wd.worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Workflow",
			Email: "clarsen@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	_, err = wd.repo.CommitObject(commit)
	if err != nil {
		return err
	}
	// push
	err = wd.repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "abc123", // anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	return err
}

func (r *mutationResolver) PrepareWeeklyReview(ctx context.Context, year *int, week *int) (*GenerateResult, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
	if year != nil {
		_year = *year
	}
	if week != nil {
		_week = *week
	}

	summaryFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", SummaryDir, _year, _year, _week)

	out, err := wd.fs.Create(summaryFname)
	if err != nil {
		log.Printf("DumpSummaryForWeek %+v", err)
		return nil, err
	}

	err = workflow.DumpSummaryForWeek(user, appkey, authtoken, _year, _week, out)
	if err != nil {
		log.Printf("DumpSummaryForWeek %+v", err)
		return nil, err
	}

	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	templateFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", ReviewDir, _year, _year, _week)
	// overwrite if exists
	outReview, err2 := wd.fs.Create(templateFname)
	if err2 != nil {
		return nil, err2
	}
	err2 = workflow.CreateEmptyWeeklyRetrospective(inSummary, outReview)
	if err2 != nil {
		return nil, err2
	}
	outReview.Close()

	// add (possibly changed) file
	wd.worktree.Add(templateFname)

	// add (possibly changed) file
	wd.worktree.Add(summaryFname)
	status, err := wd.worktree.Status()
	if err != nil {
		return nil, err
	}

	if status.IsClean() {
		log.Printf("no change, not commiting")
		msg := fmt.Sprintf("No change in %s, not commiting", summaryFname)
		result := GenerateResult{
			Message: &msg,
			Ok:      true,
		}
		return &result, nil
	}

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("dump summary for %04d-%02d", _year, _week))
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Updated %s, template at %s", summaryFname, templateFname)
	result := GenerateResult{
		Message: &msg,
		Ok:      true,
	}
	return &result, nil
}

// FinishWeeklyReview writes/commit visual to visualized-reviews
func (r *mutationResolver) FinishWeeklyReview(ctx context.Context, year *int, week *int) (*FinishResult, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
	if year != nil {
		_year = *year
	}
	if week != nil {
		_week = *week
	}

	summaryFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", SummaryDir, _year, _year, _week)
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	reviewFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", ReviewDir, _year, _year, _week)
	inReview, err2 := wd.fs.Open(reviewFname)
	if err2 != nil {
		return nil, err2
	}

	visualFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.md", ReviewVisDir, _year, _year, _week)
	out, err2 := wd.fs.Create(visualFname)
	if err2 != nil {
		return nil, err2
	}

	log.Printf("generating %s\n", visualFname)
	err2 = workflow.VisualizeWeeklyRetrospective(inSummary, inReview, out)
	if err2 != nil {
		return nil, err2
	}

	// add (possibly changed) file
	wd.worktree.Add(visualFname)
	status, err := wd.worktree.Status()
	if err != nil {
		return nil, err
	}

	if status.IsClean() {
		log.Printf("no change, not commiting")
		err = workflow.WeeklyCleanup(user, appkey, authtoken)
		if err != nil {
			return nil, err
		}
		msg := fmt.Sprintf("No change in %s, not commiting, reset cards", visualFname)
		result := FinishResult{
			Message: &msg,
			Ok:      true,
		}
		return &result, nil
	}

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("visualized review for %04d-%02d", _year, _week))
	if err != nil {
		return nil, err
	}

	err = workflow.WeeklyCleanup(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Updated %s and reset cards", visualFname)
	result := FinishResult{
		Message: &msg,
		Ok:      true,
	}
	return &result, nil
}

func (r *mutationResolver) SetDueDate(ctx context.Context, taskID string, due time.Time) (*Task, error) {
	t, err := SetTaskDue(taskID, due)
	return t, err
}

func (r *mutationResolver) SetDone(ctx context.Context, taskID string, done bool, status *string, titleComment *string, nextDue *time.Time) (*Task, error) {
	t, err := SetTaskDone(taskID, done, titleComment)
	if err != nil {
		return nil, err
	}
	if nextDue != nil {
		t, err = SetTaskDue(taskID, *nextDue)
	}
	t.List = &BoardList{
		Board: "Kanban daily/weekly",
		List:  "Done this week",
	}
	return t, err
}

func (r *mutationResolver) SetGoalDone(ctx context.Context, taskID string, checkitemID string, done bool, status *string) ([]MonthlyGoal, error) {
	// update checkitem title and done state
	card, err := cl.GetCard(taskID)
	if err != nil {
		return nil, err
	}
	for _, cl := range card.Checklists {
		for _, item := range cl.CheckItems {
			if item.ID == checkitemID {
				title, week, created, estDuration, curStatus := workflow.GetAttributesFromChecklistTitle(item.Name)
				if status == nil {
					status = curStatus
				}
				s := workflow.ChecklistTitleFromAttributes(title, week, created, estDuration, status)
				state := "incomplete"
				if done {
					state = "complete"
				}
				item.SetNameAndState(s, state)
				break
			}
		}
	}
	return GetMonthlyGoals(cl)
}

func TimeEntryToTimer(te *gttimeentry.TimeEntry) (*Timer, error) {
	return &Timer{
		ID:    strconv.FormatUint(te.Id, 16),
		Title: te.Description,
	}, nil
}

func (r *mutationResolver) AddWeeklyGoal(ctx context.Context, taskID string, title string, week int) ([]MonthlyGoal, error) {
	mg, err := cl.GetMonthlyGoal(taskID)
	if err != nil {
		return nil, err
	}

	err = cl.AddWeeklyGoal(mg, title, week)
	if err != nil {
		return nil, err
	}

	goals, err := GetMonthlyGoals(cl)
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *mutationResolver) AddMonthlyGoal(ctx context.Context, title string) ([]MonthlyGoal, error) {
	log.Printf("add goal %s\n", title)
	err := cl.AddMonthlyGoal(title)
	if err != nil {
		return nil, err
	}
	goals, err := GetMonthlyGoals(cl)
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *mutationResolver) StartTimer(ctx context.Context, taskID string, checkitemID *string) (*Timer, error) {
	if checkitemID != nil {
		log.Printf("StartTimer %s, %s\n", taskID, *checkitemID)
	} else {
		log.Printf("StartTimer %s\n", taskID)
	}
	card, err := cl.GetCard(taskID)
	if err != nil {
		return nil, err
	}
	title, _, _ := workflow.GetTitleAndAttributes(card)
	if checkitemID != nil {
		log.Printf("Find %s in checklist\n", *checkitemID)
		log.Printf("card: %+v\n", card)
		log.Printf("card checklists: %+v\n", card.Checklists)
		for _, cl := range card.Checklists {
			log.Printf("  checklist: %+v\n", cl)
			for _, item := range cl.CheckItems {
				log.Printf("    checkitem: %+v\n", item)
				if item.ID == *checkitemID {
					title, _, _, _, _ = workflow.GetAttributesFromChecklistTitle(item.Name)
					break
				}
			}
		}
	}
	tEntry := gttimeentry.TimeEntry{
		Description: title,
		Start:       time.Now(),
		// Workspace:   togglWorkspace,
		Wid: togglWorkspace.Id,
	}
	updatedTEntry, err := teCL.CreateAndStart(&tEntry)
	if err != nil {
		return nil, err
	}
	return TimeEntryToTimer(updatedTEntry)
}

func (r *mutationResolver) StopTimer(ctx context.Context, timerID string) (*bool, error) {
	id, err := strconv.ParseUint(timerID, 16, 64)
	if err != nil {
		return nil, err
	}
	_, err = teCL.Stop(id)
	if err != nil {
		ret := false
		return &ret, err
	}
	ret := true
	return &ret, nil
}

func (r *queryResolver) ActiveTimer(ctx context.Context) (*Timer, error) {
	current, err := teCL.GetCurrent()
	if err != nil {
		return nil, err
	}
	if current != nil {
		return TimeEntryToTimer(current)
	}
	return nil, nil
}

func (r *queryResolver) WeeklyVisualization(ctx context.Context, year *int, week *int) (*string, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
	if year != nil {
		_year = *year
	}
	if week != nil {
		_week = *week
	}
	summaryFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", SummaryDir, _year, _year, _week)
	if _, err = wd.fs.Stat(summaryFname); os.IsNotExist(err) {
		return nil, err
	}
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}
	reviewFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", ReviewDir, _year, _year, _week)
	if _, err = wd.fs.Stat(reviewFname); os.IsNotExist(err) {
		return nil, err
	}
	inReview, err2 := wd.fs.Open(reviewFname)
	if err2 != nil {
		return nil, err2
	}
	log.Printf("got inSummary and inReview: inSummary: %+v inReview: %+v\n", inSummary, inReview)
	var out strings.Builder

	err2 = workflow.VisualizeWeeklyRetrospective(inSummary, inReview, &out)
	if err2 != nil {
		return nil, err2
	}
	res := out.String()
	return &res, nil
}

func (r *mutationResolver) MoveTaskToList(ctx context.Context, taskID string, list BoardListInput) (*Task, error) {
	trelloList, err := workflow.ListFor(cl, list.Board, list.List)
	if err != nil {
		return nil, err
	}

	card, err := cl.MoveToListOnBoard(taskID, trelloList.ID, trelloList.IDBoard)
	if err != nil {
		return nil, err
	}
	t, err := TaskFor(card)
	if err != nil {
		return nil, err
	}
	t.List = &BoardList{
		Board: list.Board,
		List:  list.List,
	}
	return t, nil
}

func (r *mutationResolver) AddTask(ctx context.Context, title string, board *string, list *string) (*Task, error) {
	_board := "Kanban daily/weekly"
	_list := "Today"
	if board == nil {
		board = &_board
	}
	if list == nil {
		list = &_list
	}
	log.Printf("AddTask %s, %s, %s", title, *board, *list)
	_, err := workflow.ListFor(cl, *board, *list)
	if err != nil {
		return nil, err
	}
	card, err := cl.CreateCard(title, *board, *list)
	if err != nil {
		return nil, err
	}

	t, err := TaskFor(card)
	if err != nil {
		return nil, err
	}
	t.List = &BoardList{
		Board: *board,
		List:  *list,
	}
	return t, nil

}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tasks(ctx context.Context, dueBefore *int, inBoardList *BoardListInput) ([]Task, error) {
	tasks, err := GetTasks(cl, inBoardList)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *queryResolver) MonthlyGoals(ctx context.Context) ([]MonthlyGoal, error) {
	goals, err := GetMonthlyGoals(cl)
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *mutationResolver) PrepareMonthlyReview(ctx context.Context, year *int, month *int) (*GenerateResult, error) {
	// Summarize month in progress for given year/month
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year := time.Now().Year()
	_month := int(time.Now().Month())
	if year != nil {
		_year = *year
	}
	if month != nil {
		_month = *month
	}

	var inSummaries [][]byte
	var inReviews []workflow.WeeklyReviewData

	for week := 1; week <= 53; week++ {
		weekSummary := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", SummaryDir, _year, _year, week)
		if _, err1 := wd.fs.Stat(weekSummary); os.IsNotExist(err1) {
			// log.Printf("%+v doesn't exist, skipping", weekSummary)
			continue
		}
		inSummary, err2 := wd.fs.Open(weekSummary)
		if err2 != nil {
			log.Fatal(err2)
		}
		buf, err2 := ioutil.ReadAll(inSummary)
		if err2 != nil {
			log.Fatal(err2)
		}
		inSummaries = append(inSummaries, buf)

		// get month
		var weekly workflow.WeeklySummary
		err3 := yaml.Unmarshal(buf, &weekly)
		if err3 != nil {
			log.Fatal(err3)
		}

		// get weekly review for month
		reviewFname := fmt.Sprintf("%s/%d/weekly/weekly-%d-%02d.yaml", ReviewDir, _year, _year, week)
		if _, err = wd.fs.Stat(reviewFname); os.IsNotExist(err) {
			log.Printf("%+v doesn't exist", reviewFname)
			continue
		}

		inReview, err2 := wd.fs.Open(reviewFname)
		if err2 != nil {
			return nil, err2
		}
		buf2, err2 := ioutil.ReadAll(inReview)
		if err2 != nil {
			return nil, err2
		}
		if _month == weekly.Month {
			inReviews = append(inReviews, workflow.WeeklyReviewData{
				Week:    week,
				Month:   weekly.Month,
				Year:    weekly.Year,
				Content: buf2,
			})
		}
	}

	monthlySummary := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", SummaryDir, _year, _year, _month)
	out, err := wd.fs.Create(monthlySummary)
	if err != nil {
		log.Printf("PrepareMonthlyReview %+v", err)
		return nil, err
	}

	err = workflow.GenerateSummaryForMonth(user, appkey, authtoken, _year, _month, inSummaries, inReviews, out)
	if err != nil {
		log.Printf("GenerateSummaryForMonth %+v, %+v", monthlySummary, err)
		return nil, err
	}
	out.Close()

	// Generate monthly review template
	inMonthlySummary, err := wd.fs.Open(monthlySummary)
	if err != nil {
		log.Printf("Open %+v, %+v", monthlySummary, err)
		return nil, err
	}

	templateFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", ReviewDir, _year, _year, _month)
	outReview, err2 := wd.fs.Create(templateFname)
	if err2 != nil {
		log.Printf("Create %+v, %+v", templateFname, err2)
		return nil, err2
	}

	err2 = workflow.CreateEmptyMonthlyRetrospective(inMonthlySummary, outReview)
	if err2 != nil {
		log.Printf("CreateEmptyMonthlyRetrospective %+v, %+v", templateFname, err)
		return nil, err2
	}

	// Generate monthly summary of weeks for visualization input
	inMonthlySummary, err = wd.fs.Open(monthlySummary)
	if err != nil {
		log.Printf("Open %+v, %+v", monthlySummary, err)
		return nil, err
	}
	monthlyInputFname := fmt.Sprintf("%s/%d/monthlyinput/monthlyinput-%d-%02d.md", ReviewVisDir, _year, _year, _month)
	outMonthlyVisInput, err := wd.fs.Create(monthlyInputFname)
	if err != nil {
		log.Printf("Create %+v, %+v", monthlyInputFname, err)
		return nil, err
	}

	err = workflow.VisualizeWeeklySummariesForMonthly(inMonthlySummary, outMonthlyVisInput)
	if err != nil {
		log.Printf("VisualizeWeeklySummariesForMonthly %+v, %+v", monthlyInputFname, err)
		return nil, err
	}

	// add (possibly changed) files
	wd.worktree.Add(monthlySummary)
	wd.worktree.Add(templateFname)
	wd.worktree.Add(monthlyInputFname)

	status, err := wd.worktree.Status()
	if err != nil {
		return nil, err
	}
	if status.IsClean() {
		log.Printf("no change, not commiting")
		msg := fmt.Sprintf("No change in %s, not commiting", monthlySummary)
		result := GenerateResult{
			Message: &msg,
			Ok:      true,
		}
		return &result, nil
	}

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("dump summary for %04d-%02d", _year, _month))
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Updated %s", monthlySummary)
	result := GenerateResult{
		Message: &msg,
		Ok:      true,
	}
	return &result, nil
}

func (r *queryResolver) MonthlyVisualization(ctx context.Context, year *int, month *int) (*string, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}
	_year := time.Now().Year()
	_month := int(time.Now().Month())
	if year != nil {
		_year = *year
	}
	if month != nil {
		_month = *month
	}

	summaryFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", SummaryDir, _year, _year, _month)
	inSummary, err := wd.fs.Open(summaryFname)
	if err != nil {
		log.Printf("Open %+v, %+v", summaryFname, err)
		return nil, err
	}

	reviewFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", ReviewDir, _year, _year, _month)
	inReview, err := wd.fs.Open(reviewFname)
	if err != nil {
		log.Printf("Open %+v, %+v", reviewFname, err)
		return nil, err
	}
	var out strings.Builder

	err = workflow.VisualizeMonthlyRetrospective(inSummary, inReview, &out)
	if err != nil {
		return nil, err
	}
	res := out.String()
	return &res, nil
}

func (r *mutationResolver) FinishMonthlyReview(ctx context.Context, year *int, month *int) (*FinishResult, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year := time.Now().Year()
	_month := int(time.Now().Month())
	if year != nil {
		_year = *year
	}
	if month != nil {
		_month = *month
	}

	summaryFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", SummaryDir, _year, _year, _month)
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	reviewFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.yaml", ReviewDir, _year, _year, _month)
	inReview, err2 := wd.fs.Open(reviewFname)
	if err2 != nil {
		return nil, err2
	}

	visualFname := fmt.Sprintf("%s/%d/monthly/monthly-%d-%02d.md", ReviewVisDir, _year, _year, _month)
	out, err2 := wd.fs.Create(visualFname)
	if err2 != nil {
		return nil, err2
	}

	log.Printf("generating %s\n", visualFname)
	err2 = workflow.VisualizeMonthlyRetrospective(inSummary, inReview, out)
	if err2 != nil {
		return nil, err2
	}

	// add (possibly changed) file
	wd.worktree.Add(visualFname)
	status, err := wd.worktree.Status()
	if err != nil {
		return nil, err
	}

	if status.IsClean() {
		log.Printf("no change, not commiting")
		err = workflow.WeeklyCleanup(user, appkey, authtoken)
		if err != nil {
			return nil, err
		}
		msg := fmt.Sprintf("No change in %s, not commiting, reset cards", visualFname)
		result := FinishResult{
			Message: &msg,
			Ok:      true,
		}
		return &result, nil
	}

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("visualized review for %04d-%02d", _year, _month))
	if err != nil {
		return nil, err
	}

	err = workflow.MonthlyCleanup(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Updated %s and reset cards", visualFname)
	result := FinishResult{
		Message: &msg,
		Ok:      true,
	}
	return &result, nil
}
