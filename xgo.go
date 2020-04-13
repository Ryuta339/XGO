package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"errors"
	"os"
)


type Token struct {
	typ  string
	sval string
}


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

func isPunctuation (b byte) bool {
	switch b {
	case '+', '-', '(', ')', '=', '{', '}', '*', '[', ']', ',', ':', '.', '!', '<', '>', '&', '|', '%', '/':
		return true
	default:
		return false
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
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
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

func isName (b byte) bool {
	return b == '_' || isAlphabet (b)
}

func isAlphabet (b byte) bool {
	return ('A'<=b && b<='Z') || ('a'<=b && b<='z')
}

func readName (b byte) string {
	var bytes = []byte {b}
	for {
		c, err := getc ()
		if err != nil {
			return string (bytes)
		}
		if isName (c) {
			bytes = append (bytes, c)
			continue
		} else {
			ungetc ()
			return string (bytes)
		}
	}
}

func readString () string {
	var bytes = []byte {}
	for {
		c, err := getc ()
		if err != nil {
			panic ("invalid string literal")
		}
		if c == '\\' {
			// この辺なんか気持ち悪い
			c, err = getc ()
			bytes = append (bytes, c)
			continue
		} else if c != '"' {
			// この辺なんか気持ち悪い
			bytes = append (bytes, c)
			continue
		} else {
			return string (bytes)
		}
	}
}

func expect (b byte) {
	c, err := getc ()
	if err != nil {
		panic ("unexpected EOF")
	}
	if c != b {
		fmt.Printf ("char '%c' expected, but got '%c'\n", b, c)
		panic ("unexpected char")
	}
}

func readChar () string {
	c, err := getc ()
	if err != nil {
		panic ("invalid char literal")
	}
	if c == '\\' {
		c, err = getc ()
	}
	debugPrint ("gotc:" + string(c))
	expect ('\'')
	return string ([]byte{c})
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
			sval := readNumber (c)
			tok = &Token {typ: "number", sval: sval}
		case c=='\'':
			sval := readChar ()
			tok = &Token {typ: "char", sval: sval}
		case c=='"':
			sval := readString ()
			tok = &Token {typ: "string", sval: sval}
		case c==' ' || c=='\t' || c=='\n' || c=='\r':
			skipSpace ()
			continue
			// tok = &Token {typ: "space", sval: " "}
		case isPunctuation (c):
			tok = &Token {typ: "punct", sval: fmt.Sprintf ("%c", c)}
		default:
			fmt.Printf ("c='%c'\n", c)
			panic ("unknown char")
		}
		debugToken (tok)
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
	fmt.Fprintf (os.Stderr, "# %s\n", s)
}

func debugPrintWithVariable (name string, v interface{}) {
	debugPrint (fmt.Sprintf ("%s=%v", name, v))
}

func debugToken (tok *Token) {
	if tok == nil {
		fmt.Fprintf (os.Stderr, "nil\n")
		return
	}
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
