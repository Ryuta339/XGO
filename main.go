package main

import (
	"fmt"
	"os"
	"strings"
)

var errorFlag = false
var astMode = false
var tokenMode = false
var sourceFile string

func putError(errorMsg string, v ...interface{}) {
	s := fmt.Sprintf("\x1b[31m"+errorMsg, v...)
	s += " [" + sourceFile + "]\x1b[39m\n"
	// 	errorFlag = true
	// renderTokens()
	panic(s)
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
