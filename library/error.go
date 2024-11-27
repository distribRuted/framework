package Lib

import (
	"os"

	Log "github.com/distribRuted/framework/library/log"
	"github.com/fatih/color"
)

// Print a fatal error and exit
func ExitWithError(errMsg string) {
	errFormat := color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	Log.PrintMsg(errFormat(errMsg))
	os.Exit(1)
}

func Exit() {
	// TODO: Define the logging processes and the actions to be performed when the program terminates.
	os.Exit(0)
}
