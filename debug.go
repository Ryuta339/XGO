package main

import (
	"fmt"
	"os"
)

type Debuggable interface {
	debug()
}

var debugOutput = os.Stdout

func debugPrintf(format string, v ...interface{}) {
	debugPrintln(fmt.Sprintf(format, v...))
}

func debugPrintln(s string) {
	fmt.Fprintf(debugOutput, "# %s\n", s)
}

func debugPrint(s string) {
	fmt.Fprintf(debugOutput, "# %s", s)
}

func debugPrintlnWithVariable(name string, v interface{}) {
	debugPrintf("%s=%v", name, v)
}

func debugTokens(ts *TokenStream) {
	ts.debug()
}

func debugToken(tok *Token) {
	if tok == nil {
		debugPrintln("tok:nil")
		return
	}
	tok.debug()
}

func debugAst(ast Ast) {
	ast.debug()
}
