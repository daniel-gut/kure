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

		logListPod, err := getLogs(p)
		if err != nil {
			return err
		}
		logList = append(logList, logListPod...)
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
		logList   []log
		logSince  int64
		namespace string
	)

	if logSince == 0 {
		logSince = config.LogSinceDefault
	}

	namespace, err := clients.GetNamespaceFromKubeconfig()
	if err != nil {
		namespace = ""
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err.Error())
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	containers := pod.Spec.Containers

	// Get logs for all containers
	logsAvailable := false
	for _, c := range containers {

		logOptions := &corev1.PodLogOptions{
			Container: c.Name,
			Follow:    false,
		}

		req := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)

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

		if len(logList) > 0 {
			logsAvailable = true
		}

	}

	if logsAvailable == false {
		fmt.Printf("No logs in specified pods since %ds\n", logSince)
		os.Exit(0)
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

	a := regexp.MustCompile(`\d{4}-\d{2}-\d{2}(\ |T)\d{2}:\d{2}:\d{2}`) // 2020-04-28 07:16:00 or 2020-04-14T07:04:19

	tsd := a.FindAll(logRaw, -1)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err.Error())
	}

	for _, ts := range tsd {

		t, err := dateparse.ParseIn(string(ts), loc)
		if err != nil {
			panic(err.Error())
		}

		if t.Unix() > (time.Now().Unix() - logSince) {
			logData = log{timestamp: t, loglevel: "error", podName: podName}
		}
	}

	return logData, nil
}
