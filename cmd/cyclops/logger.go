package cyclops

import (
	"os"

	"github.com/fatih/color"
)

var SuccessPrint = color.New(color.FgHiGreen).PrintlnFunc()
var WarningPrint = color.New(color.FgHiYellow).PrintlnFunc()
var FatalPrint = func(a ...interface{}) {
	color.New(color.FgHiRed).PrintlnFunc()(a)
	os.Exit(1)
}
