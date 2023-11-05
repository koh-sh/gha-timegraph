package gha

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/koh-sh/gha-timegraph/internal/types"
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
