//go:generate go run github.com/99designs/gqlgen
package handle_graphql

import (
	"context"
	"fmt"
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

	"github.com/clarsen/go-trello-workflow/workflow"
	"github.com/clarsen/gtoggl-api/gthttp"
	"github.com/clarsen/gtoggl-api/gttimeentry"
	"github.com/clarsen/gtoggl-api/gtworkspace"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var appkey, authtoken, user string
var cl *workflow.Client
var teCL *gttimeentry.TimeEntryClient
var togglWorkspace *gtworkspace.Workspace

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

	// logger := log.New(os.Stderr, "gthttp.trace: ", log.LstdFlags|log.Llongfile)
	// _thc, err := gthttp.NewClient(togglAPIKey, gthttp.SetTraceLogger(logger))
	_thc, err := gthttp.NewClient(togglAPIKey)
	if err != nil {
		log.Fatal("Couldn't create new gthttp client")
	}

	//	_tc := gtclient.NewClient(_thc)
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
	summarydir := "task-summary"
	summaryFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, _year, _week)

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

	reviewdir := "reviews"

	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	templateFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, _year, _week)
	if _, err = wd.fs.Stat(templateFname); err != nil {
		// doesn't exist
		outReview, err2 := wd.fs.Create(templateFname)
		if err2 != nil {
			return nil, err2
		}

		err2 = workflow.CreateEmptyWeeklyRetrospective(inSummary, outReview)
		if err2 != nil {
			return nil, err2
		}
		outReview.Close()

		// add file
		wd.worktree.Add(templateFname)
	}

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

	summarydir := "task-summary"
	summaryFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, _year, _week)
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	reviewdir := "reviews"
	reviewFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, _year, _week)
	inReview, err2 := wd.fs.Open(reviewFname)
	if err2 != nil {
		return nil, err2
	}

	reviewvisdir := "visualized_reviews"
	visualFname := fmt.Sprintf("%s/weekly-%d-%02d.md", reviewvisdir, _year, _week)
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
		msg := fmt.Sprintf("No change in %s, not commiting", visualFname)
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
	msg := fmt.Sprintf("Updated %s", visualFname)
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

func (r *mutationResolver) SetDone(ctx context.Context, taskID string, done bool, status *string, nextDue *time.Time) (*Task, error) {
	t, err := SetTaskDone(taskID, done)
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

func TimeEntryToTimer(te *gttimeentry.TimeEntry) (*Timer, error) {
	return &Timer{
		ID:    strconv.FormatUint(te.Id, 16),
		Title: te.Description,
	}, nil
}

func (r *mutationResolver) StartTimer(ctx context.Context, taskID string, checkitemID *string) (*Timer, error) {
	log.Printf("StartTimer %s, %s\n", taskID, checkitemID)
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

	summarydir := "task-summary"
	reviewdir := "reviews"

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
	if year != nil {
		_year = *year
	}
	if week != nil {
		_week = *week
	}
	summaryFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, _year, _week)
	if _, err = wd.fs.Stat(summaryFname); os.IsNotExist(err) {
		return nil, err
	}
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}
	reviewFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, _year, _week)
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

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tasks(ctx context.Context, dueBefore *int, inBoardList *BoardListInput) ([]Task, error) {
	tasks, err := GetTasks(user, appkey, authtoken, inBoardList)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *queryResolver) MonthlyGoals(ctx context.Context) ([]MonthlyGoal, error) {
	goals, err := GetMonthlyGoals(user, appkey, authtoken)
	if err != nil {
		return nil, err
	}
	return goals, nil
}
