package kure

import (
	"os"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aybabtme/uniplot/histogram"
)

type log struct {
	timestamp time.Time
	podName   string
	loglevel  string
}

func analyzeLog(resourceName string) (string, error) {

	logList, err := parseLog(logMockData)
	if err != nil {
		return "", err
	}

	printLogHistogram(logList)

	return "Successfull", err
}

func parseLog(logData string) ([]log, error) {
	a := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`) // 2020-04-14T07:04:19
	tsd := a.FindAll([]byte(logData), -1)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err.Error())
	}

	var logList []log

	for _, ts := range tsd {
		t, err := dateparse.ParseIn(string(ts), loc)
		if err != nil {
			panic(err.Error())
		}

		if t.Unix() > (time.Now().Unix() - 3600) {
			logList = append(logList, log{timestamp: t, loglevel: "error"})
		}

	}

	return logList, nil
}

func printLogHistogram(logList []log) {

	var timestampList []float64

	for _, log := range logList {
		timestampList = append(timestampList, float64(log.timestamp.Unix()))
	}

	hist := histogram.Hist(10, timestampList)

	err := histogram.Fprintf(os.Stdout, hist, histogram.Linear(10), func(v float64) string {
		return time.Unix(int64(v), 0).Format("15:04:05")
	})
	if err != nil {
		panic(err)
	}

}
