package main

import (
	"fmt"
	"os"
)

var errorFlag = false
var astMode = false

func putError(errorMsg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "\x1b[31m")
	fmt.Fprintf(os.Stderr, errorMsg, v...)
	fmt.Fprintf(os.Stderr, "\x1b[39m\n")
	// 	errorFlag = true
	debugMode = true
	renderTokens()
	panic("internal error")
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

	tokenize(sourceFile)
	/*
		if debugMode {
			renderTokens()
		}
	*/
	ast := parse()
	if errorFlag {
		panic("internal error")
	}
	/*
		if debugMode {
			debugPrint("==== Start Dump Ast ====")
			debugAst(ast)
			debugPrint("==== End Dump Ast ====")
		}
	*/
	if astMode {
		showAst(ast, 0)
	} else {
		generate(ast)
	}
}
