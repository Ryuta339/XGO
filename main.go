package main

import (
	"fmt"
	"os"
	"strings"
)

var errorFlag = false
var astMode = false
var tokenMode = false

func putError(errorMsg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "\x1b[31m")
	fmt.Fprintf(os.Stderr, errorMsg, v...)
	fmt.Fprintf(os.Stderr, "\x1b[39m\n")
	// 	errorFlag = true
	renderTokens()
	panic("internal error")
}

func parseOptions(args []string) string {
	var sourceFile string
	for _, opt := range args {
		if opt == "-t" {
			tokenMode = true
		}
		if opt == "-a" {
			astMode = true
		}

		if strings.HasSuffix(opt, ".go") {
			sourceFile = opt
		} else if opt == "-" {
			sourceFile = "/dev/stdin"
		}
	}
	if sourceFile == "" {
		putError("Unspecified source file.")
	}
	return sourceFile
}

func main() {
	var sourceFile string
	if len(os.Args) > 1 {
		sourceFile = parseOptions(os.Args[1:len(os.Args)])
	} else {
		putError("Usaga: xgo [-a][-t] sourceFile")
	}
	tokenize(sourceFile)
	if tokenMode {
		renderTokens()
	}
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
