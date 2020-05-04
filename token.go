package main

import (
	"fmt"
	"io/ioutil"
)

var reservedList = []string{
	"break",
	"default",
	"func",
	"interface",
	"select",
	"case",
	"defer",
	"go",
	"map",
	"struct",
	"chan",
	"else",
	"goto",
	"package",
	"switch",
	"const",
	"fallthrough",
	"if",
	"range",
	"type",
	"continue",
	"for",
	"import",
	"return",
	"var",
}

var tStream *TokenStream
var bStream *ByteStream

/* ================================
 * SorceFile
 * ================================ */
type SourceFile struct {
	filename string
	line     int
	column   int
}

/* ================================
 * Token
 *     implements Debuggable and fmt.Stringer
 * ================================ */
type Token struct {
	typ  string
	sval string
	SourceFile
}

// implements Debuggalbe
func (tok *Token) debug() {
	debugPrint(fmt.Sprintf("tok:type=%s, sval=%s", tok.typ, tok.sval))
}

// implements fmt.Stringer
func (tok *Token) String() string {
	return fmt.Sprintf("(%s \"%s\" in%s: line %d: column %d)",
		tok.typ, tok.sval, tok.filename, tok.line, tok.column)
}

func (tok *Token) isEOF() bool {
	return tok == nil || tok.typ == "eof"
}

func (tok *Token) isPunct(s string) bool {
	return tok != nil && tok.typ == "punct" && tok.sval == s
}

func (tok *Token) isReserved(s string) bool {
	return tok != nil && tok.typ == "reserved" && tok.sval == s
}

func (tok *Token) isString(s string) bool {
	return tok != nil && tok.typ == "string" && tok.sval == s
}

func (tok *Token) isIdentifier(s string) bool {
	return tok != nil && tok.typ == "identifier" && tok.sval == s
}

func (tok *Token) isTypePunct() bool {
	return tok != nil && tok.typ == "punct"
}

func (tok *Token) isTypeReserved() bool {
	return tok != nil && tok.typ == "keyword"
}

func (tok *Token) isTypeString() bool {
	return tok != nil && tok.typ == "string"
}

func (tok *Token) isTypeIdentifier() bool {
	return tok != nil && tok.typ == "identifier"
}

func (tok *Token) isTypeInt() bool {
	return tok != nil && tok.typ == "int"
}

func (tok *Token) isTypeRune() bool {
	return tok != nil && tok.typ == "rune"
}

func newToken(typ string, sval string) *Token {
	return &Token{
		typ:  typ,
		sval: sval,
		SourceFile: SourceFile{
			filename: bStream.filename,
			line:     bStream.line,
			column:   bStream.column,
		},
	}
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
	return &Token{
		typ : "eof",
		sval: "",
		SourceFile: bStream.SourceFile,
	}
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
	for idx, tok := range ts.tokens {
		if idx==ts.index {
			debugPrint("=== now parsing ===")
		}
		debugToken(tok)
		if idx==ts.index {
			debugPrint("=== now parsing ===")
		}
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
	tStream.nextToken()
}

func lookahead(num int) *Token {
	return tStream.lookahead(num)
}

func consumeToken(expected string) {
	tStream.consumeToken(expected)
}

func renderTokens() {
	tStream.renderTokens()
}

/* ================================ */

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
		c, err := bStream.getc()
		if err != nil {
			return string(chars)
		}
		if isNumber(c) {
			chars = append(chars, c)
			continue
		} else {
			bStream.ungetc()
			return string(chars)
		}
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func skipSpace() {
	for {
		c, err := bStream.getc()
		if err != nil {
			return
		}
		if isSpace(c) {
			continue
		} else {
			bStream.ungetc()
			return
		}
	}
}

func isAlphabet(b byte) bool {
	return ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z')
}

func isName(b byte) bool {
	return b == '_' || isAlphabet(b) || isNumber(b)
}

func readName(b byte) string {
	var bytes = []byte{b}
	for {
		c, err := bStream.getc()
		if err != nil {
			return string(bytes)
		}
		if isName(c) {
			bytes = append(bytes, c)
			continue
		} else {
			bStream.ungetc()
			return string(bytes)
		}
	}
}

func isReserved(word string) bool {
	for _, v := range reservedList {
		if word == v {
			return true
		}
	}
	return false
}

func readString() string {
	var bytes = []byte{}
	for {
		c, err := bStream.getc()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			c, err = bStream.getc()
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
	c, err := bStream.getc()
	if err != nil {
		panic("unexpected EOF")
	}
	if c != b {
		putError("char '%c' expected, but got '%c'\n", b, c)
	}
}

func readChar() string {
	c, err := bStream.getc()
	if err != nil {
		panic("invalid char literal")
	}
	if c == '\\' {
		c, err = bStream.getc()
	}
	expect('\'')
	return string([]byte{c})
}

func tokenize(filename string) {
	s := readFile(filename)
	var r []*Token
	bStream = &ByteStream{
		source: s,
		index:  0,
		SourceFile: SourceFile{
			filename: filename,
			line:     1,
			column:   0,
		},
	}
	for {
		c, err := bStream.getc()
		if err != nil {
			tStream = newTokenStream(r)
			return
		}
		var tok *Token
		switch {
		case c == 0:
			tStream = newTokenStream(r)
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
				// makeSymbol(sval, "int")
			}
		}
		r = append(r, tok)
	}
}

func readFile(filename string) string {
	bytes, ok := ioutil.ReadFile(filename)
	if ok != nil {
		panic(ok)
	}
	return string(bytes)
}
