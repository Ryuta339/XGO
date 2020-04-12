package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
)

type Token struct {
	typ  string
	sval string
}


var tokens [] *Token
var tokenIndex int

func readFile (filename string) string {
	bytes, ok := ioutil.ReadFile (filename)
	if ok != nil {
		panic (ok)
	}
	return string (bytes);
}

func tokinize (s string) []*Token {
	var r [] *Token
	trimed := strings.Trim (s, "\n")
	chars := strings.Split (trimed, " ")
	var regexNumber = regexp.MustCompile (`^[0-9]+$`)
	for _, char := range chars {
		debugPrint ("char", char)
		var tok *Token
		if regexNumber.MatchString (char) {
			tok = &Token {typ: "number", sval: strings.Trim (char, " \n")}
		}

		r = append (r, tok)
	}

	return r
}

func readToken () *Token {
	if tokenIndex <= len (tokens)-1 {
		r := tokens[tokenIndex]
		tokenIndex ++
		return r
	}
	return nil
}


func parseExpr () Ast {
	ast := parseUnaryExpression ()
	return ast
}

func generate (ast Ast) {
	fmt.Println("\t.global _mymain")
	fmt.Println("_mymain:");
	emitAst (ast)
	fmt.Println("\tret");
}

func emitAst (ast Ast) {
	ast.emit ();
}

func debugPrint (name string, v interface{}) {
	fmt.Printf ("# %s=%v\n", name, v)
}

func debugTokens (tokens []*Token) {
	for _, tok := range tokens {
		debugPrint ("tok", tok)
	}
}

func debugAst (ast Ast) {
	ast.debug ();
}

func main () {
	s := readFile ("/dev/stdin")
	tokens = tokinize (s)
	tokenIndex = 0
	debugTokens (tokens)
	ast := parseExpr ()
	debugAst (ast)
	generate (ast)
}
