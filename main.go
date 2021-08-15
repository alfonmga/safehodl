package main

import (
	"github.com/alfonmga/safehodl/cmd"
)

var PinCode string
var Secret32BytesKeyAES string

func main() {
	cmd.Execute(Secret32BytesKeyAES, PinCode)
}
