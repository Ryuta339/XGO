package main

import (
	"fmt"
	"strings"
	"errors"
	"os"
)


type Token struct {
	typ  string
	sval string
}



func readToken () *Token {
	if tokenIndex <= len (tokens)-1 {
		r := tokens[tokenIndex]
		tokenIndex ++
		return r
	}
	return nil
}

func unreadToken () {
	if tokenIndex >= 0 {
		tokenIndex --
	}
}

func debugToken (tok *Token) {
	if tok == nil {
		fmt.Fprintf (os.Stderr, "nil\n")
	}
	debugPrint (fmt.Sprintf ("tok:type=%s, sval=%s", tok.typ, tok.sval))
}

func debugTokens (tokens []*Token) {
	for _, tok := range tokens {
		debugToken (tok)
	}
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
		} else if c == '"' {
			// この辺なんか気持ち悪い
			return string (bytes)
		} else {
			bytes = append (bytes, c)
			continue
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


func tokenize (s string) []*Token {
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
