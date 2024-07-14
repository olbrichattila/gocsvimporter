package importer

import (
	"fmt"
	"math"
)

type application struct {
	importer importer
}

func newApplication(
	importer importer,
) (*application, error) {
	app := &application{}
	app.importer = importer

	return app, nil
}

func (a *application) displayTimeStat(analysisTime, importTime, totalTime float64) {
	fmt.Printf(
		"\n\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n",
		a.durationAsString(analysisTime),
		a.durationAsString(importTime),
		a.durationAsString(totalTime),
	)
}

func (application) durationAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", math.Floor(elapsed/60), int64(elapsed)%60)
}
