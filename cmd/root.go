package cmd

import (
	"fmt"
	"os"

	safehodl "github.com/alfonmga/safehodl/lib"
	"github.com/spf13/cobra"
)

var PinCode string

var rootCmd = &cobra.Command{
	Use:   "safehodl",
	Short: "SafeHodl",
	Long: `Track your Bitcoin holdings value in a safe way
https://github.com/alfonmga/safehodl`,
	Run: func(cmd *cobra.Command, args []string) {
		safehodl.AssertPinCodeForUsage(PinCode)

		hasStoredAmount, _ := safehodl.GetHodlAmount()
		if !hasStoredAmount {
			fmt.Println(`Error: Execute first "safehodl config" to configure SafeHODL.`)
			os.Exit(0)
		}

		safehodl.DisplayHodlInfo()
	},
}

func Execute(secret32BytesKeyAES string, pinCode string) {
	PinCode = pinCode
	safehodl.Secret32BytesKeyAES = secret32BytesKeyAES

	cobra.CheckErr(rootCmd.Execute())
}
