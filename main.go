package main

import (
	"github.com/alfonmga/safehodl/cmd"
)

// PinCode access pin code for secure access to SafeHODL.
var PinCode string

// Secret32BytesKeyAES 32 bytes AES secret key for .safehodl data file encryption.
var Secret32BytesKeyAES string

func main() {
	cmd.Execute(Secret32BytesKeyAES, PinCode)
}
