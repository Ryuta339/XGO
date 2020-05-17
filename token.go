package main

import (
	"fmt"
	"io/ioutil"
)

type TokenType string

const (
	T_EOF         TokenType = "EOF"
	T_INT         TokenType = "int"
	T_STRING      TokenType = "string"
	T_RUNE        TokenType = "rune"
	T_IDENTIFIER  TokenType = "identifier"
	T_PUNCTUATION TokenType = "punctuation"
	T_KEYWORD     TokenType = "keyword"
)

var keywordList = []string{
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
	typ  TokenType
	sval string
	SourceFile
}

// implements Debuggalbe
func (tok *Token) debug() {
	debugPrintf("tok:type=%s, sval=%s", tok.typ, tok.sval)
}

// implements fmt.Stringer
func (tok *Token) String() string {
	return fmt.Sprintf("(%s \"%s\" in%s: line %d: column %d)",
		tok.typ, tok.sval, tok.filename, tok.line, tok.column)
}

func (tok *Token) isEOF() bool {
	return tok != nil && tok.typ == T_EOF
}

func (tok *Token) isPunct(s string) bool {
	return tok != nil && tok.typ == T_PUNCTUATION && tok.sval == s
}

func (tok *Token) isKeyword(s string) bool {
	return tok != nil && tok.typ == T_KEYWORD && tok.sval == s
}

func (tok *Token) isString(s string) bool {
	return tok != nil && tok.typ == T_STRING && tok.sval == s
}

func (tok *Token) isIdentifier(s string) bool {
	return tok != nil && tok.typ == T_IDENTIFIER && tok.sval == s
}

func (tok *Token) isTypePunct() bool {
	return tok != nil && tok.typ == T_PUNCTUATION
}

func (tok *Token) isTypeKeyword() bool {
	return tok != nil && tok.typ == T_KEYWORD
}

func (tok *Token) isTypeString() bool {
	return tok != nil && tok.typ == T_STRING
}

func (tok *Token) isTypeIdentifier() bool {
	return tok != nil && tok.typ == T_IDENTIFIER
}

func (tok *Token) isTypeInt() bool {
	return tok != nil && tok.typ == T_INT
}

func (tok *Token) isTypeRune() bool {
	return tok != nil && tok.typ == T_RUNE
}

func (tok *Token) isSemicolon() bool {
	return tok.isPunct(";")
}

func newToken(typ TokenType, sval string) *Token {
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

var semicolon = &Token{
	typ:  T_PUNCTUATION,
	sval: ";",
}

func autoSemicolonInsert(last *Token) bool {
	return last.isTypeIdentifier() ||
		last.isTypeInt() || last.isTypeRune() || last.isTypeString() ||
		last.isKeyword("break") || last.isKeyword("continue") || last.isKeyword("fallthrough") || last.isKeyword("return") ||
		last.isPunct("++") || last.isPunct("--") || last.isPunct(")") || last.isPunct("]") || last.isPunct("}")
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
		typ:        T_EOF,
		sval:       "",
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
		if idx == ts.index {
			debugPrintln("\x1b[31m======== now parsing ========\x1b[39m")
		}
		debugToken(tok)
		if idx == ts.index {
			debugPrintln("\x1b[31m======== now parsing ========\x1b[39m")
		}
	}
}

func (ts *TokenStream) renderTokens() {
	debugPrintln("==== Start Dump Tokens ====")
	ts.debug()
	debugPrintln("==== End Dump Tokens ====")
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

func skip(isFunc func(byte) bool) {
	for {
		c, err := bStream.getc()
		if err != nil {
			return
		}
		if isFunc(c) {
			continue
		} else {
			bStream.ungetc()
			return
		}
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

func skipSpace() {
	skip(isSpace)
}

func isNewLine(b byte) bool {
	return b == '\n' || b == '\r'
}

func skipNewLine() {
	skip(isNewLine)
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

func isKeyword(word string) bool {
	for _, v := range keywordList {
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

func skipLine() {
	for {
		c, err := bStream.getc()
		if err != nil || isNewLine(c) {
			bStream.ungetc()
			return
		}
	}
}

func skipBlockComment() {
	prev, err := bStream.getc()
	if err != nil {
		bStream.ungetc()
		return
	}

	for {
		c, err := bStream.getc()
		if err != nil {
			putError("Premature end of block comment")
		}
		if prev == '*' && c == '/' {
			return
		}
		prev = c
	}
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
			tok = &Token{typ: T_INT, sval: sval}
		case c == '\'':
			sval := readChar()
			tok = &Token{typ: T_RUNE, sval: sval}
		case c == '"':
			sval := readString()
			tok = &Token{typ: T_STRING, sval: sval}
		case c == ' ' || c == '\t':
			skipSpace()
			continue
		case c == '\r' || c == '\n':
			// insert semicolon
			if len(r) > 0 {
				last := r[len(r)-1]
				if autoSemicolonInsert(last) {
					r = append(r, semicolon)
				}
			}
			skipNewLine()
			continue
		case c == '/':
			c, _ = bStream.getc()
			if c == '/' {
				skipLine()
				continue
			} else if c == '*' {
				skipBlockComment()
				continue
			} else if c == '=' {
				tok = &Token{typ: T_PUNCTUATION, sval: "/="}
			} else {
				bStream.ungetc()
				tok = &Token{typ: T_PUNCTUATION, sval: "/"}
			}
		case isPunctuation(c):
			tok = &Token{typ: T_PUNCTUATION, sval: fmt.Sprintf("%c", c)}
		case c == '=':
			tok = &Token{typ: T_PUNCTUATION, sval: fmt.Sprintf("%c", c)}
		default:
			sval := readName(c)
			if isKeyword(sval) {
				tok = &Token{typ: T_KEYWORD, sval: sval}
			} else {
				tok = &Token{typ: T_IDENTIFIER, sval: sval}
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
