package kure

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/daniel-gut/kure/pkg/clients"
	"github.com/daniel-gut/kure/pkg/config"
	"github.com/daniel-gut/kure/pkg/graph"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type log struct {
	timestamp time.Time
	podName   string
	loglevel  string
	nodeName  string
}

type bucketData struct {
	logs       []log
	bucketName string // Field of the "log" type to be the bucket Key
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

	var bcData []string

	// Go through the loglist and calc stats by all fields

	s := reflect.ValueOf(&logList[0]).Elem()
	typeOfT := s.Type()

	for i := 1; i < s.NumField(); i++ {
		fieldName := typeOfT.Field(i).Name

		for _, l := range logList {

			key, err := l.normalizeDataForPrint(fieldName)
			bcData = append(bcData, key)
			if err != nil {
				return err
			}

		}
		graph.PrintBarChart(bcData)
		fmt.Println(strings.Repeat("-", 100))

		// empty log for next field name
		bcData = []string{}

	}

	return err
}

// func printLogHistogramTimestamp(logList []log) {

// 	var timestampList []float64

// 	for _, log := range logList {
// 		timestampList = append(timestampList, float64(log.timestamp.Unix()))
// 	}

// 	hist := histogram.Hist(10, timestampList)

// 	err := histogram.Fprintf(os.Stdout, hist, histogram.Linear(20), func(v float64) string {
// 		return time.Unix(int64(v), 0).Format("15:04:05")
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func printLogHistogramPodName(logList []log) {

// 	var timestampList []float64

// 	unique := make(map[string]bool)
// 	for _, l := range logList {
// 		if !unique[l.podName] {
// 			unique[l.podName] = true
// 		}
// 	}

// 	for _, log := range logList {
// 		timestampList = append(timestampList, float64(log.timestamp.Unix()))
// 	}

// 	hist := histogram.Hist(20, timestampList)

// 	err := histogram.Fprintf(os.Stdout, hist, histogram.Linear(10), func(v float64) string {
// 		return time.Unix(int64(v), 0).Format("15:04:05")
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

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
	nodeName := pod.Spec.NodeName

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

			logData, err := parseLog(line, podName, nodeName)

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

func parseLog(logRaw []byte, podName string, nodeName string) (log, error) {

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
			logData = log{timestamp: t, loglevel: "error", podName: podName, nodeName: nodeName}
		}
	}

	return logData, nil
}

func (log log) normalizeDataForPrint(keyName string) (string, error) {

	var key string

	switch keyName {
	case "podName":
		key = log.podName
	case "nodeName":
		key = log.nodeName
	case "loglevel":
		key = log.loglevel
	default:
		return "", fmt.Errorf("error, keyName fieldname unknown %w", keyName)
	}
	return key, nil
}
