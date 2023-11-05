package types

import (
	"fmt"
	"time"
)

type Run struct {
	Name      string
	Starttime time.Time
	Elapsed   float64
}

func (r Run) RtnCSVrow() string {
	return fmt.Sprintf("%s,%s,%g", r.Name, r.Starttime.Format("2006-01-02 15:04:05"), r.Elapsed)
}
