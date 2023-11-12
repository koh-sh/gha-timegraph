package stdout

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/koh-sh/gha-timegraph/internal/types"
)

func TestPrintRuns(t *testing.T) {
	runs := []types.Run{
		{
			Name:      "test",
			Starttime: time.Date(2022, 4, 1, 9, 0, 0, 0, time.UTC),
			Elapsed:   300,
		},
		{
			Name:      "test",
			Starttime: time.Date(2022, 4, 1, 8, 0, 0, 0, time.UTC),
			Elapsed:   200,
		},
	}
	csvwant := `Name,StartTime(UTC),Elapsed(Sec)
test,2022-04-01 09:00:00,300
test,2022-04-01 08:00:00,200`
	tblwant := `|StartTime(UTC)|Elapsed(Sec)|
|---|---|
|2022-04-01 09:00:00|300|
|2022-04-01 08:00:00|200|`
	jsnwant := `[
     {
          "Name": "test",
          "Starttime": "2022-04-01T09:00:00Z",
          "Elapsed": 300
     },
     {
          "Name": "test",
          "Starttime": "2022-04-01T08:00:00Z",
          "Elapsed": 200
     }
]`

	type args struct {
		runs []types.Run
		out  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "csv",
			args:    args{runs: runs, out: "csv"},
			want:    csvwant,
			wantErr: false,
		},
		{
			name:    "table",
			args:    args{runs: runs, out: "table"},
			want:    tblwant,
			wantErr: false,
		},
		{
			name:    "json",
			args:    args{runs: runs, out: "json"},
			want:    jsnwant,
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{runs: []types.Run{}, out: "csv"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "not supported format",
			args:    args{runs: runs, out: "notSupportedFormat"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := out2pipe(t, PrintRuns, tt.args.runs, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrintRuns() = %v, want %v", got, tt.want)
			}
		})
	}
}

// get stdout of function and return as string
func out2pipe(t *testing.T, fnc func([]types.Run, string) error, runs []types.Run, out string) (string, error) {
	// https://zenn.dev/glassonion1/articles/8ac939208bd455
	t.Helper()

	// keep original Stdout
	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()
	// make pipe
	r, w, _ := os.Pipe()
	os.Stdout = w
	// call fnc
	if e := fnc(runs, out); e != nil {
		return "", e
	}
	// close writer
	w.Close()
	// get string from Buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}
