package kure

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aybabtme/uniplot/histogram"
	"github.com/daniel-gut/kure/pkg/clients"
	"github.com/daniel-gut/kure/pkg/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type log struct {
	timestamp time.Time
	podName   string
	loglevel  string
}

func analyzeLog(podList []string) error {

	var (
		logList []log
		err     error
	)

	for _, p := range podList {
		fmt.Printf("Gettings logs for %s\n", p)

		logList, err = getLogs(p)
		if err != nil {
			return err
		}

		// fmt.Println(logData)
		// // logData = []byte(logMockData)

		// logList, err = parseLog(logData, p)
		// if err != nil {
		// 	return err
		// }
	}

	printLogHistogram(logList)

	return err
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

func getLogs(podName string) ([]log, error) {

	var (
		// logs     []byte
		logList []log
		// logData  log
		logSince int64
	)

	if logSince == 0 {
		logSince = config.LogSinceDefault
	}

	k8sconfig := clients.GetConfig()

	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		panic(err.Error())
	}

	pod, err := clientset.CoreV1().Pods("api-services").Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting pod for %w", err)
	}
	containers := pod.Spec.Containers

	// Get logs for all containers
	for _, c := range containers {

		logOptions := &corev1.PodLogOptions{
			Container: c.Name,
			Follow:    false,
		}

		req := clientset.CoreV1().Pods("api-services").GetLogs(podName, logOptions)

		stream, err := req.Stream()
		if err != nil {
			return nil, fmt.Errorf("error opening stream to Pod: %w", err)
		}
		defer stream.Close()

		reader := bufio.NewReader(stream)

		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error reading stream: %w", err)
			}

			logData, err := parseLog(line, podName)

			if logData != (log{}) {
				logList = append(logList, logData)
			}

		}

		if len(logList) == 0 {
			err := fmt.Errorf("No logs since %d seconds", logSince)
			return nil, err
		}

	}

	return logList, nil

}

func parseLog(logRaw []byte, podName string) (log, error) {

	var (
		logSince int64
		logData  log
	)

	if logSince == 0 {
		logSince = config.LogSinceDefault
	}

	a := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`) // 2020-04-14T07:04:19
	tsd := a.FindAll(logRaw, -1)

	loc, err := time.LoadLocation("UCT")
	if err != nil {
		panic(err.Error())
	}

	for _, ts := range tsd {

		t, err := dateparse.ParseIn(string(ts), loc)
		if err != nil {
			panic(err.Error())
		}

		// fmt.Println(t)

		if t.Unix() > (time.Now().Unix() - logSince) {
			logData = log{timestamp: t, loglevel: "error", podName: podName}
		}

	}

	return logData, nil
}
