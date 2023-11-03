package types

import "time"

type Run struct {
	Name      string
	Starttime time.Time
	Elapsed   float64
}
