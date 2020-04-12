package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"errors"
)

type Token struct {
	typ  string
	sval string
}


var tokens [] *Token
var tokenIndex int
var source string
var sourceIndex int



func readFile (filename string) string {
	bytes, ok := ioutil.ReadFile (filename)
	if ok != nil {
		panic (ok)
	}
	return string (bytes);
}

func getc () (byte, error) {
	if sourceIndex >= len (source) {
		return 0, errors.New ("EOF")
	}
	r := source[sourceIndex]
	sourceIndex ++
	return r, nil
}

func ungetc () {
	if sourceIndex > 0 {
		sourceIndex --
	}
}

func isNumber (b byte) bool {
	debugPrint (fmt.Sprintf ("is_numeric %c", b))
	return '0' <= b && b <= '9'
}

func readNumber (b byte) string {
	var chars = []byte{b}
	for {
		debugPrint ("read number");
		c, err := getc ()
		if err != nil {
			return string (chars)
		}
		if isNumber (c) {
			chars = append (chars, c)
			continue
		} else {
			ungetc ()
			return string (chars)
		}
	}
}

func isSpace (b byte) bool {
	return b == ' ' || b == '\t'
}

func skipSpace () {
	for {
		c, err := getc ()
		if err != nil {
			return
		}
		if isSpace (c) {
			continue
		} else {
			ungetc ()
			return
		}
	}
}

func tokinize (s string) []*Token {
	var r [] *Token
	s = strings.Trim (s, "\n")
	source = s
	for {
		c, err := getc ()
		if err != nil {
			return r
		}
		var tok *Token
		switch {
		case c == 0:
			return r
		case isNumber (c):
			val := readNumber (c)
			tok = &Token {typ: "number", sval: val}
		case c==' ' || c=='\t':
			skipSpace ()
			tok = &Token {typ: "space", sval: " "}
		case c=='+':
			tok = &Token {typ: "punct", sval: fmt.Sprintf ("%c", c)}
		default:
			fmt.Printf ("c='%c'\n", c)
			panic ("unknown char")
		}

		r = append (r, tok)
	}
}

func readToken () *Token {
	if tokenIndex <= len (tokens)-1 {
		r := tokens[tokenIndex]
		tokenIndex ++
		return r
	}
	return nil
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

func debugPrint (s string) {
	fmt.Printf ("# %s\n", s)
}

func debugPrintWithVariable (name string, v interface{}) {
	debugPrint (fmt.Sprintf ("%s=%v\n", name, v))
}

func debugToken (tok *Token) {
	debugPrint (fmt.Sprintf ("tok:type=%s, sval=%s", tok.typ, tok.sval))
}

func debugTokens (tokens []*Token) {
	for _, tok := range tokens {
		debugToken (tok)
	}
}

func debugAst (ast Ast) {
	ast.debug ();
}

func main () {
	s := readFile ("/dev/stdin")
	tokens = tokinize (s)
	if len (tokens) == 0 {
		panic ("no tokens")
	}
	tokenIndex = 0
	debugTokens (tokens)
	ast := parseExpression ()
	debugAst (ast)
	generate (ast)
}
