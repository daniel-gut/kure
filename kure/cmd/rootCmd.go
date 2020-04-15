package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kure",
	Short: "Cure your Kubernetes Cluster",
	Long:  "Helper tool to find and cure unhealthy parts within a Kubernetes Cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute is the main function for the CLI, called by the main.go
func Execute() {
	if len(os.Args) < 2 {
		rootCmd.Usage()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	rootCmd.PersistentFlags().StringP("context", "c", "K8S Context", "the name of the cluster context you want to communicate with")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "the Namespace you want to interact with, default is 'all'")
}
