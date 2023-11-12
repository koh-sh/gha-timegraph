package gha

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/koh-sh/gha-timegraph/internal/types"
	"github.com/schollz/progressbar/v3"
)

// return client of github api
func RtnClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}
	return github.NewClient(nil).WithAuthToken(token)
}

// return list of types.Run from GitHub Workflow runs
func GetRuns(client *github.Client, count int, owner, repo, filename, branch, status string, silent bool) ([]types.Run, error) {
	// setup
	runs := make([]types.Run, 0, count)
	lopts := github.ListOptions{PerPage: min(count, 100)}
	opts := github.ListWorkflowRunsOptions{Branch: branch, Status: status, ListOptions: lopts}

	// show progress bar
	var bar *progressbar.ProgressBar
	if !silent {
		bar = progressbar.Default(int64(count))
	}

	for {
		wfruns, resp, err := client.Actions.ListWorkflowRunsByFileName(context.Background(), owner, repo, filename, &opts)
		if err != nil {
			return nil, err
		}
		for _, v := range wfruns.WorkflowRuns {
			run := makeRun(*v)
			// If the run is older than a year, UpdatedAt is updated.
			// Too big Elapsed could be noise so should be omitted.
			if len(runs) == count || run.Elapsed > 60*60*24*365 {
				if !silent {
					bar.Finish()
				}
				return runs, nil
			}
			runs = append(runs, run)
			if !silent {
				bar.Add(1)
			}
		}
		if resp.NextPage == 0 {
			if !silent {
				bar.Finish()
			}
			return runs, nil
		}
		opts.Page = resp.NextPage
	}
}

// make types.Run from GitHub Workflow run
func makeRun(wfrun github.WorkflowRun) types.Run {
	endtime := wfrun.UpdatedAt
	starttime := getStartTime(wfrun)
	elapsed := endtime.Sub(starttime).Round(time.Second).Seconds()
	return types.Run{Name: *wfrun.Name, Starttime: starttime, Elapsed: elapsed}
}

// get start time
func getStartTime(wfrun github.WorkflowRun) time.Time {
	// https://github.com/cli/cli/blob/trunk/pkg/cmd/run/shared/shared.go#L110
	if wfrun.RunStartedAt.IsZero() {
		return wfrun.CreatedAt.Time
	}
	return wfrun.RunStartedAt.Time
}
