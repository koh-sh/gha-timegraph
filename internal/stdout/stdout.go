package stdout

import (
	"encoding/json"
	"fmt"

	"github.com/koh-sh/gha-timegraph/internal/types"
)

// print runs for specified output format
func PrintRuns(runs []types.Run, out string) error {
	if len(runs) == 0 {
		return nil
	}

	titleN := "Name"
	titleX := "StartTime(UTC)"
	titleY := "Elapsed(Sec)"
	timefmt := "2006-01-02 15:04:05"

	switch out {
	case "csv":
		fmt.Printf("%s,%s,%s\n", titleN, titleX, titleY)
		for _, v := range runs {
			fmt.Printf("%s,%s,%g\n", v.Name, v.Starttime.Format(timefmt), v.Elapsed)
		}
	case "table":
		fmt.Printf("|%s|%s|\n", titleX, titleY)
		fmt.Println("|---|---|")
		for _, v := range runs {
			fmt.Printf("|%s|%g|\n", v.Starttime.Format(timefmt), v.Elapsed)
		}
	case "json":
		j, err := json.MarshalIndent(runs, "", "     ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
	default:
		return fmt.Errorf("Not Supported Format: %s", out)
	}
	return nil
}
