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
	Short: "acispanctl version",
	Long:  `Print the version number of acispanctl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("acispanctl v0.9 -- HEAD")
	},
}
