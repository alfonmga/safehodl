package cmd

import (
	safehodl "github.com/alfonmga/safehodl/lib"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure SafeHODL",
	Long:  `Configure SafeHODL`,
	Run: func(cmd *cobra.Command, args []string) {
		safehodl.StartInteractiveSafeHodlConfiguration()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
