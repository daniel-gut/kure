package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kure",
	Long:  `All software has versions. This is kure's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kure - Keep K8S healthy v0.1 -- HEAD")
	},
}
