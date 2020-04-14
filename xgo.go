package main

import (
	"fmt"
	"io/ioutil"
	"os"
)


var tokens [] *Token
var tokenIndex int
var source string
var sourceIndex int
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

func renderTokens (tokens []*Token) {
	debugPrint ("==== Start Dump Tokens ====")
	for _, tok := range tokens {
		if tok.typ == "string" {
			fmt.Fprintf (os.Stderr, "\"%s\"\n", tok.sval)
		} else {
			fmt.Fprintf (os.Stderr, "%s\n", tok.sval)
		}
	}
	debugPrint ("==== End Dump Tokens ====")
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
