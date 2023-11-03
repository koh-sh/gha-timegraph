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

func getStartTime(run github.WorkflowRun) time.Time {
	// https://github.com/cli/cli/blob/trunk/pkg/cmd/run/shared/shared.go#L110
	if run.RunStartedAt.IsZero() {
		return run.RunStartedAt.Time
	}
	return run.CreatedAt.Time
}

func GetRuns(client *github.Client, count int, owner, repo, filename, branch, status string) ([]types.Run, error) {
	runs := make([]types.Run, 0, count)
	lopts := github.ListOptions{PerPage: 100}
	opts := github.ListWorkflowRunsOptions{ListOptions: lopts, Branch: branch, Status: status}
	for {
		wruns, resp, err := client.Actions.ListWorkflowRunsByFileName(context.Background(), owner, repo, filename, &opts)
		if err != nil {
			return nil, err
		}
		for _, v := range wruns.WorkflowRuns {
			endtime := v.UpdatedAt
			starttime := getStartTime(*v)
			elapsed := endtime.Sub(starttime).Round(time.Second).Seconds()
			runs = append(runs, types.Run{Name: *v.Name, Starttime: starttime, Elapsed: elapsed})
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
