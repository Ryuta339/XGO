package main

import (
	"errors"
)

/* ================================
 * ByteStream
 * ================================ */

type ByteStream struct {
	source string
	index  int
	SourceFile
}

func (bs *ByteStream) getc() (byte, error) {
	if bs.index >= len(bs.source) {
		return 0, errors.New("EOF")
	}
	r := bs.source[bs.index]
	if r == '\r' || r == '\n' {
		bs.line++
		bs.column = 0
	}
	bs.index++
	bs.column++
	return r, nil
}

func (bs *ByteStream) ungetc() {
	if bs.index > 0 {
		bs.index--
		r := bs.source[bs.index]
		if r == '\r' || r == '\n' {
			bs.line--
			bs.column = -1
		}
	}
}
