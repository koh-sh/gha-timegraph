package types

import (
	"testing"
	"time"
)

func TestRun_RtnCSVrow(t *testing.T) {
	type fields struct {
		Name      string
		Starttime time.Time
		Elapsed   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				Name:      "basic",
				Starttime: time.Date(2022, 4, 1, 9, 5, 0, 0, time.UTC),
				Elapsed:   300,
			},
			want: "basic,2022-04-01 09:05:00,300",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Run{
				Name:      tt.fields.Name,
				Starttime: tt.fields.Starttime,
				Elapsed:   tt.fields.Elapsed,
			}
			if got := r.RtnCSVrow(); got != tt.want {
				t.Errorf("Run.RtnCSVrow() = %v, want %v", got, tt.want)
			}
		})
	}
}
