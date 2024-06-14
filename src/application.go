package importer

import (
	"fmt"
	"math"
	"time"
)

type application struct {
	importer importer
}

func newApplication(
	args argParser,
	env enver,
	connector dBConnector,
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

	err = connector.init()
	if err != nil {
		return nil, err
	}

	err = importer.getCsvReader().init(csvFileName, separator)
	if err != nil {
		return nil, err
	}

	err = importer.getStorer().init(connector.getDBConfig(), tableName)
	if err != nil {
		return nil, err
	}

	app.importer = importer

	return app, nil
}

func (a *application) displayTimeStat(startTime, analysisTime time.Time) {
	// Record the finished time
	finishedTime := time.Now()

	// Calculate the durations
	fullAnalysisDuration := analysisTime.Sub(startTime).Seconds()
	fullDurationDuration := finishedTime.Sub(analysisTime).Seconds()
	totalDuration := finishedTime.Sub(startTime).Seconds()

	// Convert durations to string representations
	fullAnalysisTime := a.durationAsString(fullAnalysisDuration)
	fullDurationTime := a.durationAsString(fullDurationDuration)
	totalTime := a.durationAsString(totalDuration)

	// Print the results
	fmt.Printf("\nDone\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n", fullAnalysisTime, fullDurationTime, totalTime)
}

func (application) durationAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", math.Floor(elapsed/60), int64(elapsed)%60)
}
