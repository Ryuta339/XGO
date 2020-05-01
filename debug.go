package main

import (
	"fmt"
	"os"
)

type Debuggable interface {
	debug()
}

var debugMode = false

func debugPrint(s string) {
	if debugMode {
		fmt.Fprintf(os.Stdout, "# %s\n", s)
	}
}
func debugPrintWithVariable(name string, v interface{}) {
	debugPrint(fmt.Sprintf("%s=%v", name, v))
}

func debugTokens(ts *TokenStream) {
	ts.debug()
}

func debugToken(tok *Token) {
	if tok == nil {
		debugPrint("tok:nil")
		return
	}
	tok.debug()
}

func debugAst(ast Ast) {
	ast.debug()
}
