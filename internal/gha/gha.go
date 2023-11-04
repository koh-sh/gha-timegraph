package gha

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/koh-sh/gha-timegraph/internal/types"
)

func RtnClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}
	return github.NewClient(nil).WithAuthToken(token)
}

func GetRuns(client *github.Client, count int, owner, repo, filename, branch, status string) ([]types.Run, error) {
	runs := make([]types.Run, 0, count)
	lopts := github.ListOptions{PerPage: min(count, 100)}
	opts := github.ListWorkflowRunsOptions{Branch: branch, Status: status, ListOptions: lopts}
	for {
		wfruns, resp, err := client.Actions.ListWorkflowRunsByFileName(context.Background(), owner, repo, filename, &opts)
		if err != nil {
			return nil, err
		}
		for _, v := range wfruns.WorkflowRuns {
			runs = append(runs, makeRun(*v))
			if len(runs) == count {
				return runs, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return runs, nil
}

func makeRun(wfrun github.WorkflowRun) types.Run {
	endtime := wfrun.UpdatedAt
	starttime := getStartTime(wfrun)
	elapsed := endtime.Sub(starttime).Round(time.Second).Seconds()
	return types.Run{Name: *wfrun.Name, Starttime: starttime, Elapsed: elapsed}
}

func getStartTime(wfrun github.WorkflowRun) time.Time {
	// https://github.com/cli/cli/blob/trunk/pkg/cmd/run/shared/shared.go#L110
	if wfrun.RunStartedAt.IsZero() {
		return wfrun.RunStartedAt.Time
	}
	return wfrun.CreatedAt.Time
}
