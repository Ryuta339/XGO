package main

import (
	"fmt"
	"io/ioutil"
	"os"
)


var debugMode = false
var errorFlag = false

func putError (errorMsg string, v ...interface{}) {
	fmt.Fprintf (os.Stderr, errorMsg, v)
	fmt.Fprintln (os.Stderr, "")
	errorFlag = true
}


func readFile (filename string) string {
	bytes, ok := ioutil.ReadFile (filename)
	if ok != nil {
		panic (ok)
	}
	return string (bytes);
}

func debugPrint (s string) {
	if debugMode {
		fmt.Fprintf (os.Stdout, "# %s\n", s)
	}
}

func debugPrintWithVariable (name string, v interface{}) {
	debugPrint (fmt.Sprintf ("%s=%v", name, v))
}

func main () {
	debugMode = true

	astMode := false

	var sourceFile string
	if len (os.Args) > 1 {
		sourceFile = os.Args[1] + ".go"
	} else {
		sourceFile = "/dev/stdin"
	}
	if len(os.Args)>2  && os.Args[2]=="-a" {
		astMode = true
	}

	s := readFile (sourceFile)

	tokens = tokenize (s)
	if len (tokens) == 0 {
		panic ("no tokens")
	}
	tokenIndex = 0
	if debugMode {
		debugPrint ("==== Start Dump Tokens ====")
		debugTokens (tokens)
		debugPrint ("==== End Dump Tokens ====")
	}
	ast := parse ()
	if errorFlag {
		panic ("internal error")
	}
	if debugMode {
		debugPrint ("==== Start Dump Ast ====")
		debugAst (ast)
		debugPrint ("==== End Dump Ast ====")
	}
	if astMode {
		showAst (ast, 0)
	} else {
		generate (ast)
	}
}
