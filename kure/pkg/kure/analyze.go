package kure

import (
	"fmt"
)

func IsValidAnalyzeArg(arg string) bool {
	switch arg {
	case "deployment":
		return true

	default:
		return false
	}
}

//Analyze starts a pod in the cluster
func Analyze(args []string) {
	fmt.Printf("Analyze with args %v \n", args)
}
