//go:generate go run github.com/99designs/gqlgen
package handle_graphql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/clarsen/go-trello-workflow/workflow"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var appkey, authtoken, user string

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

}

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
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

func (r *mutationResolver) GenerateWeeklySummary(ctx context.Context, year *int, week *int) (*bool, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
	summarydir := "task-summary"
	fname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, _year, _week)

	out, err := wd.fs.Create(fname)
	if err != nil {
		log.Printf("DumpSummaryForWeek %+v", err)
		return nil, err
	}

	err = workflow.DumpSummaryForWeek(user, appkey, authtoken, _year, _week, out)
	if err != nil {
		log.Printf("DumpSummaryForWeek %+v", err)
		return nil, err
	}

	// add file
	wd.worktree.Add(fname)
	status, err := wd.worktree.Status()
	if err != nil {
		return nil, err
	}
	if status.IsClean() {
		log.Printf("no change, not commiting")
		res := true
		return &res, nil
	}

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("dump summary for %04d-%02d", _year, _week))
	if err != nil {
		return nil, err
	}
	res := true
	return &res, nil
}

func (r *mutationResolver) GenerateWeeklyReviewTemplate(ctx context.Context, year *int, week *int) (*bool, error) {
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
	reviewdir := "reviews"

	summaryFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", summarydir, _year, _week)
	if _, err = wd.fs.Stat(summaryFname); os.IsNotExist(err) {
		return nil, err
	}
	inSummary, err2 := wd.fs.Open(summaryFname)
	if err2 != nil {
		return nil, err2
	}

	templateFname := fmt.Sprintf("%s/weekly-%d-%02d.yaml", reviewdir, _year, _week)
	if _, err = wd.fs.Stat(templateFname); err == nil {
		return nil, errors.New(fmt.Sprintf("%s exists already", templateFname))
	}

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

	// commit and push
	err = wd.commitAndPushData(fmt.Sprintf("dump summary for %04d-%02d", _year, _week))
	if err != nil {
		return nil, err
	}
	res := true
	return &res, nil
}

func (r *mutationResolver) SetDueDate(ctx context.Context, taskID string, due time.Time) (*Task, error) {
	t, err := SetTaskDue(taskID, due)
	return t, err
}

func (r *queryResolver) WeeklyVisualization(ctx context.Context) (*string, error) {
	// setup working directory
	wd, err := getData()
	if err != nil {
		return nil, err
	}

	summarydir := "task-summary"
	reviewdir := "reviews"

	_year, _week := time.Now().Add(-time.Hour * 72).ISOWeek()
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

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tasks(ctx context.Context, dueBefore *int, inBoardList *BoardList) ([]Task, error) {

	tasks, err := GetTasks(user, appkey, authtoken, inBoardList)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
