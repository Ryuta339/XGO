package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var errorFlag = false
var astMode = false

func putError(errorMsg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, errorMsg, v)
	fmt.Fprintln(os.Stderr, "")
	errorFlag = true
}

func readFile(filename string) string {
	bytes, ok := ioutil.ReadFile(filename)
	if ok != nil {
		panic(ok)
	}
	return string(bytes)
}

func main() {
	debugMode = true

	var sourceFile string
	if len(os.Args) > 1 {
		sourceFile = os.Args[1] + ".go"
	} else {
		sourceFile = "/dev/stdin"
	}
	if len(os.Args) > 2 && os.Args[2] == "-a" {
		astMode = true
	}

	s := readFile(sourceFile)

	tokenize(s)
	/*
		if debugMode {
			renderTokens()
		}
	*/
	ast := parse()
	if errorFlag {
		panic("internal error")
	}
	if debugMode {
		debugPrint("==== Start Dump Ast ====")
		debugAst(ast)
		debugPrint("==== End Dump Ast ====")
	}
	if astMode {
		showAst(ast, 0)
	} else {
		generate(ast)
	}
}
