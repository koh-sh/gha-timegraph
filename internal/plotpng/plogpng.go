package plotpng

import (
	"github.com/koh-sh/gha-timegraph/internal/types"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func SavePng(runs []types.Run) error {
	p := plot.New()

	p.Title.Text = runs[0].Name
	p.X.AutoRescale = true
	p.X.Label.Text = "StartTime(UTC)"
	p.Y.Label.Text = "Elapsed (Sec)"
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "2006-01-02 15:04:05",
	}

	err := plotutil.AddLinePoints(p, "Elapsed", rtnPlots(runs))
	if err != nil {
		return err
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		return err
	}
	return nil
}

func rtnPlots(runs []types.Run) plotter.XYs {
	pts := make(plotter.XYs, len(runs))
	for i := range pts {
		pts[i].X = float64(runs[i].Starttime.Unix())
		pts[i].Y = runs[i].Elapsed
	}
	return pts
}
