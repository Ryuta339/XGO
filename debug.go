package main

import (
	"fmt"
	"os"
)


var debugMode = false

func debugPrint (s string) {
	if debugMode {
		fmt.Fprintf (os.Stdout, "# %s\n", s)
	}
}
func debugPrintWithVariable (name string, v interface{}) {
	debugPrint (fmt.Sprintf ("%s=%v", name, v))
}
