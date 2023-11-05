package plotpng

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/koh-sh/gha-timegraph/internal/types"
	"gonum.org/v1/plot/plotter"
)

func Test_rtnPlots(t *testing.T) {
	type args struct {
		runs []types.Run
	}
	tests := []struct {
		name string
		args args
		want plotter.XYs
	}{
		{
			name: "basic",
			args: args{runs: []types.Run{{
				Name:      "basic",
				Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
				Elapsed:   300,
			}}},
			want: plotter.XYs{plotter.XY{
				X: float64(time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC).Unix()),
				Y: 300,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rtnPlots(tt.args.runs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rtnPlots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSavePng(t *testing.T) {
	type args struct {
		runs    []types.Run
		outfile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{runs: []types.Run{}, outfile: "tmp/out.png"},
			wantErr: false,
		},
		{
			name: "basic",
			args: args{runs: []types.Run{{
				Name:      "basic",
				Starttime: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC),
				Elapsed:   300,
			}}, outfile: "tmp/out.png"},
			wantErr: false,
		},
		{
			name: "err",
			args: args{runs: []types.Run{{
				Name:      "err",
				Starttime: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC),
				Elapsed:   300,
			}}, outfile: "tmp/out.out"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SavePng(tt.args.runs, tt.args.outfile); (err != nil) != tt.wantErr {
				t.Errorf("SavePng() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// check the output from basic case
	if getFileMD5("tmp/out.png") != getFileMD5("testdata/out.png_") {
		t.Errorf("output file is not same")
	}
}

// get MD5Sum from png file
func getFileMD5(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum(b)
	return fmt.Sprintf("%x", hash)
}
