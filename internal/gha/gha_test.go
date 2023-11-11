package gha

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/koh-sh/gha-timegraph/internal/types"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func Test_makeRun(t *testing.T) {
	wftitle := "basic"
	wftitle2 := "RunStarted is zero"
	type args struct {
		wfrun github.WorkflowRun
	}
	tests := []struct {
		name string
		args args
		want types.Run
	}{
		{
			name: "basic",
			args: args{github.WorkflowRun{
				Name: &wftitle,
				CreatedAt: &github.Timestamp{
					Time: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				},
				UpdatedAt: &github.Timestamp{
					Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC),
				},
				RunStartedAt: &github.Timestamp{
					Time: time.Date(2022, 4, 1, 9, 2, 0, 0, time.UTC),
				},
			}},
			want: types.Run{
				Name:      wftitle,
				Starttime: time.Date(2022, 4, 1, 9, 2, 0, 0, time.UTC),
				Elapsed:   180,
			},
		},
		{
			name: "RunStarted is zero",
			args: args{github.WorkflowRun{
				Name: &wftitle2,
				CreatedAt: &github.Timestamp{
					Time: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				},
				UpdatedAt: &github.Timestamp{
					Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC),
				},
				RunStartedAt: &github.Timestamp{},
			}},
			want: types.Run{
				Name:      wftitle2,
				Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				Elapsed:   300,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeRun(tt.args.wfrun); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRtnClient(t *testing.T) {
	tests := []struct {
		name string
		want *github.Client
	}{
		{
			name: "basic",
			want: github.NewClient(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RtnClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RtnClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRuns(t *testing.T) {
	type args struct {
		client   *github.Client
		count    int
		owner    string
		repo     string
		filename string
		branch   string
		status   string
		silent   bool
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Run
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				client:   mockClient("basic"),
				count:    30,
				owner:    "owner",
				repo:     "repo",
				filename: "test.yml",
				branch:   "",
				status:   "success",
				silent:   false,
			},
			want: []types.Run{{
				Name:      "test",
				Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				Elapsed:   300,
			}},
			wantErr: false,
		},
		{
			name: "old",
			args: args{
				client:   mockClient("old"),
				count:    30,
				owner:    "owner",
				repo:     "repo",
				filename: "test.yml",
				branch:   "",
				status:   "success",
				silent:   false,
			},
			want: []types.Run{{
				Name:      "test",
				Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				Elapsed:   300,
			}},
			wantErr: false,
		},
		{
			name: "pages",
			args: args{
				client:   mockClient("pages"),
				count:    4,
				owner:    "owner",
				repo:     "repo",
				filename: "test.yml",
				branch:   "",
				status:   "success",
				silent:   true,
			},
			want: []types.Run{
				{
					Name:      "test",
					Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
					Elapsed:   300,
				},
				{
					Name:      "test",
					Starttime: time.Date(2022, 4, 1, 8, 0, 0, 0, time.UTC),
					Elapsed:   300,
				},
				{
					Name:      "test",
					Starttime: time.Date(2022, 4, 1, 7, 0, 0, 0, time.UTC),
					Elapsed:   300,
				},
				{
					Name:      "test",
					Starttime: time.Date(2022, 4, 1, 6, 0, 0, 0, time.UTC),
					Elapsed:   300,
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				client:   mockClient("empty"),
				count:    30,
				owner:    "owner",
				repo:     "repo",
				filename: "test.yml",
				branch:   "",
				status:   "success",
				silent:   true,
			},
			want:    []types.Run{},
			wantErr: false,
		},
		{
			name: "ratelimit",
			args: args{
				client:   mockClient("ratelimit"),
				count:    1,
				owner:    "owner",
				repo:     "repo",
				filename: "test.yml",
				branch:   "",
				status:   "success",
				silent:   true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRuns(tt.args.client, tt.args.count, tt.args.owner, tt.args.repo, tt.args.filename, tt.args.branch, tt.args.status, tt.args.silent)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRuns() = %v, want %v", got, tt.want)
			}
		})
	}
}

// return mock GitHub Client
func mockClient(ptn string) *github.Client {
	switch ptn {
	case "ratelimit":
		return github.NewClient(mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{},
				},
			),
			mock.WithRateLimit(0, 0),
		),
		)

	case "old":
		return github.NewClient(mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC)},
						},
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2021, 4, 1, 9, 0, 0, 0, time.UTC)},
						},
					},
				},
			),
		))

	case "empty":
		return github.NewClient(mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{},
				},
			),
		))
	case "pages":
		return github.NewClient(mock.NewMockedHTTPClient(
			mock.WithRequestMatchPages(
				mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC)},
						},
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 8, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 8, 0, 0, 0, time.UTC)},
						},
					},
				},
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 7, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 7, 0, 0, 0, time.UTC)},
						},
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 6, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 6, 0, 0, 0, time.UTC)},
						},
					},
				},
			),
		))

	default:
		return github.NewClient(mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{
					WorkflowRuns: []*github.WorkflowRun{
						{
							Name:         github.String("test"),
							UpdatedAt:    &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC)},
							RunStartedAt: &github.Timestamp{Time: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC)},
						},
					},
				},
			),
		))
	}
}
