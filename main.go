package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v56/github"
)

func main() {
	// Construct a new GitHub client
	client := rtnClient()

	// parameters
	owner := "koh-sh"
	repo := "codebuild-multirunner"
	filename := "go-test.yml"
	branch := ""
	status := "success"
	count := 300

	runs, err := getRuns(client, count, owner, repo, filename, branch, status)
	if err != nil {
		log.Fatal(err)
	}
	printAsCSV(runs)
}

func rtnClient() *github.Client {
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

func getRuns(client *github.Client, count int, owner, repo, filename, branch, status string) ([]Run, error) {
	runs := make([]Run, 0, count)
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
			runs = append(runs, Run{Name: *v.Name, Starttime: starttime, Elapsed: elapsed})
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

func printAsCSV(runs []Run) {
	// print Columns
	fmt.Printf("%s,%s,%s\n", "Name", "StartTime(UTC)", "Elapsed")
	for _, v := range runs {
		fmt.Printf("%s,%s,%g\n", v.Name, v.Starttime.Format("2006-01-02 15:04:05"), v.Elapsed)
	}
}

type Run struct {
	Name      string
	Starttime time.Time
	Elapsed   float64
}
