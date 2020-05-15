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
	// renderTokens()
	panic(s)
}

func parseOptions(args []string) {
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
}

func main() {
	if len(os.Args) > 1 {
		parseOptions(os.Args[1:len(os.Args)])
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
	if astMode {
		showAst(ast, 0)
	} else {
		generate(ast)
	}
}
