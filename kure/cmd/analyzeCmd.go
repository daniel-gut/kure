package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/daniel-gut/kure/pkg/kure"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use: "analyze",
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("analyze COMMAND requires at least a resource type")
		}
		if kure.IsValidAnalyzeObj(args[0]) {
			return nil
		}

		return fmt.Errorf("%s is a unknown resource", args[0])
	},
	Short: "Analyze a component of the cluster. Usage: kure analyze [object]",
	Long:  `Analyze a component`,
	Run: func(cmd *cobra.Command, args []string) {
		err := kure.Analyze(args)
		if err != nil {
			log.Fatal(err)
		}

	},
}
