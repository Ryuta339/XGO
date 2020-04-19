package main

import (
	"fmt"
	"io/ioutil"
	"os"
)


var debugMode = false


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

	var sourceFile string
	if len (os.Args) > 1 {
		sourceFile = os.Args[1] + ".go"
	} else {
		sourceFile = "/dev/stdin"
	}

	s := readFile (sourceFile)

	tokens = tokenize (s)
	if len (tokens) == 0 {
		panic ("no tokens")
	}
	tokenIndex = 0
	debugTokens (tokens)
	ast := parseExpression ()
	debugAst (ast)
	generate (ast)
}
