package cmd

import (
	"github.com/daniel-gut/kure/pkg/kure"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helperPodCmd)
}

var helperPodCmd = &cobra.Command{
	Use:   "helperPod",
	Short: "Start a helper pod",
	Long:  `Start a helper pod with ready tooling and helper scripts`,
	Run: func(cmd *cobra.Command, args []string) {
		kure.HelperPod()
	},
}
