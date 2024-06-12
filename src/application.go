package importer

import (
	"fmt"
	"time"
)

type application struct {
	importer importer
}

func newApplication(
	args argParser,
	env enver,
	importer importer,
) (*application, error) {
	app := &application{}

	csvFileName, tableName, separator, err := args.pharse()
	if err != nil {
		return nil, err
	}

	err = env.loadEnv()
	if err != nil {
		return nil, err
	}

	err = importer.getCsvReader().init(csvFileName, separator)
	if err != nil {
		return nil, err
	}

	err = importer.getStorer().init(tableName)
	if err != nil {
		return nil, err
	}

	app.importer = importer

	return app, nil
}

func (a *application) displayTimeStat(startTime, analysisTime time.Time) {
	finisedTime := time.Now()
	fullAnalysisTime := a.durasionAsString(analysisTime.Sub(startTime).Seconds())
	fullDurationTime := a.durasionAsString(finisedTime.Sub(analysisTime).Seconds())
	totalTime := a.durasionAsString(finisedTime.Sub(startTime).Seconds())

	fmt.Printf("\nDone\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n", fullAnalysisTime, fullDurationTime, totalTime)
}

func (application) durasionAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", elapsed/60, int64(elapsed)%60)
}
