package cmd

import (
	"fmt"
	"os"

	safehodl "github.com/alfonmga/safehodl/lib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "safehodl",
	Short: "SafeHodl",
	Long: `Track your Bitcoin holdings value in a safe way
https://github.com/alfonmga/safehodl`,
	Run: func(cmd *cobra.Command, args []string) {
		safehodl.AssertPassphrase()

		hasStoredAmount, _ := safehodl.GetBtcAmount()
		if !hasStoredAmount {
			fmt.Println(`Error: Execute first "safehodl config" to configure SafeHODL.`)
			os.Exit(0)
		}

		safehodl.DisplayHodlInfo()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
