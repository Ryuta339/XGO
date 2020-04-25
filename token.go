package main

import (
	"errors"
	"fmt"
	"strings"
)

var ts *TokenStream
var source string
var sourceIndex int

/* ================================
 * Token
 *     implements Debuggable
 * ================================ */
type Token struct {
	typ  string
	sval string
}

// implements Debuggalbe
func (tok *Token) debug() {
	debugPrint(fmt.Sprintf("tok:type=%s, sval=%s", tok.typ, tok.sval))
}

/* ================================
 * TokenStream
 *     implements Debuggable
 * ================================ */
type TokenStream struct {
	index  int
	tokens []*Token
}

func (ts *TokenStream) nextToken() {
	if ts.index <= len(ts.tokens)-1 {
		ts.index++
	}
}

func (ts *TokenStream) lookahead(num int) *Token {
	idx := ts.index + num - 1
	if idx <= len(ts.tokens)-1 {
		return ts.tokens[idx]
	}
	return nil
}

func (ts *TokenStream) consumeToken(expected string) {
	tok := ts.lookahead(1)
	if tok == nil {
		putError("Unexpected termination.\n")
	}
	if expected == tok.sval {
		ts.nextToken()
	} else {
		putError("Expected token %s, but got %s.\n", expected, tok.sval)
	}
}

// implements Debuggable
func (ts *TokenStream) debug() {
	for _, tok := range ts.tokens {
		debugToken(tok)
	}
}

func (ts *TokenStream) renderTokens() {
	debugPrint("==== Start Dump Tokens ====")
	ts.debug()
	debugPrint("==== End Dump Tokens ====")
}

/* ================================ */
func newTokenStream(tokens []*Token) *TokenStream {
	return &TokenStream{
		index:  0,
		tokens: tokens,
	}
}

// wrapper
func nextToken() {
	ts.nextToken()
}

func lookahead(num int) *Token {
	return ts.lookahead(num)
}

func consumeToken(expected string) {
	ts.consumeToken(expected)
}

/* ================================ */

func getc() (byte, error) {
	if sourceIndex >= len(source) {
		return 0, errors.New("EOF")
	}
	r := source[sourceIndex]
	sourceIndex++
	return r, nil
}

func ungetc() {
	if sourceIndex > 0 {
		sourceIndex--
	}
}

func isPunctuation(b byte) bool {
	switch b {
	case '+', '-', '(', ')', '=', '{', '}', '*', '[', ']', ',', ':', '.', '!', '<', '>', '&', '|', '%', '/':
		return true
	default:
		return false
	}
}

func isNumber(b byte) bool {
	ret := '0' <= b && b <= '9'
	return ret
}

func readNumber(b byte) string {
	var chars = []byte{b}
	for {
		c, err := getc()
		if err != nil {
			return string(chars)
		}
		if isNumber(c) {
			chars = append(chars, c)
			continue
		} else {
			ungetc()
			return string(chars)
		}
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func skipSpace() {
	for {
		c, err := getc()
		if err != nil {
			return
		}
		if isSpace(c) {
			continue
		} else {
			ungetc()
			return
		}
	}
}

func isAlphabet(b byte) bool {
	return ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z')
}

func isName(b byte) bool {
	return b == '_' || isAlphabet(b)
}

func readName(b byte) string {
	var bytes = []byte{b}
	for {
		c, err := getc()
		if err != nil {
			return string(bytes)
		}
		if isName(c) {
			bytes = append(bytes, c)
			continue
		} else {
			ungetc()
			return string(bytes)
		}
	}
}

func isReserved(word string) bool {
	return word == "func" || word == "package"
}

func readString() string {
	var bytes = []byte{}
	for {
		c, err := getc()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			c, err = getc()
			if err != nil {
				panic("invalid string literal")
			}
			switch c {
			case 'n':
				bytes = append(bytes, '\\', 'n')
			case 'r':
				bytes = append(bytes, '\\', 'r')
			case 't':
				bytes = append(bytes, '\\', 't')
			default:
				bytes = append(bytes, c)
			}
			continue
		} else if c == '"' {
			return string(bytes)
		} else {
			bytes = append(bytes, c)
			continue
		}
	}
}

func expect(b byte) {
	c, err := getc()
	if err != nil {
		panic("unexpected EOF")
	}
	if c != b {
		fmt.Printf("char '%c' expected, but got '%c'\n", b, c)
		panic("unexpected char")
	}
}

func readChar() string {
	c, err := getc()
	if err != nil {
		panic("invalid char literal")
	}
	if c == '\\' {
		c, err = getc()
	}
	expect('\'')
	return string([]byte{c})
}

func tokenize(s string) {
	var r []*Token
	s = strings.Trim(s, "\n")
	source = s
	for {
		c, err := getc()
		if err != nil {
			ts = newTokenStream(r)
			return
		}
		var tok *Token
		switch {
		case c == 0:
			ts = newTokenStream(r)
			return
		case isNumber(c):
			sval := readNumber(c)
			tok = &Token{typ: "int", sval: sval}
		case c == '\'':
			sval := readChar()
			tok = &Token{typ: "rune", sval: sval}
		case c == '"':
			sval := readString()
			tok = &Token{typ: "string", sval: sval}
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			skipSpace()
			continue
			// tok = &Token {typ: "space", sval: " "}
		case isPunctuation(c):
			tok = &Token{typ: "punct", sval: fmt.Sprintf("%c", c)}
		case c == '=':
			tok = &Token{typ: "assignment", sval: fmt.Sprintf("%c", c)}
		default:
			sval := readName(c)
			if isReserved(sval) {
				tok = &Token{typ: "reserved", sval: sval}
			} else {
				tok = &Token{typ: "identifier", sval: sval}
				makeSymbol(sval, "int")
			}
		}
		r = append(r, tok)
	}
}
