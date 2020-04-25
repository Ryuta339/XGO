package main

import (
	"fmt"
)

type NameSpace interface {
	emitLeftValue(sym *Symbol)
	emitRightValue(sym *Symbol)
}

type LocalVariable struct {
	typ string
}

// implements NameSpace
func (lv *LocalVariable) emitRightValue(sym *Symbol) {
	emitCode("\tpushq\t-%d(%%rbp)", sym.pos*8)
	frameHeight += 8
}
func (lv *LocalVariable) emitLeftValue(sym *Symbol) {
	emitCode("\tleaq\t-%d(%%rbp), %%rax", sym.pos*8)
	emitCode("\tpushq\t%%rax")
	frameHeight += 8
}

/* ================================ */

type ValueType int

const (
	RIGHT = iota
	LEFT
)

type Symbol struct {
	pos    int
	name   string
	nSpace NameSpace
}

func (s *Symbol) emitSymbol(vType ValueType) {
	if vType == RIGHT {
		s.nSpace.emitRightValue(s)
	} else if vType == LEFT {
		s.nSpace.emitLeftValue(s)
	} else {
		fmt.Printf("undefined ValueType %v.\n", vType)
		panic("internal error")
	}
}

var symlist []*Symbol = make([]*Symbol, 0, 4)

func makeSymbol(name string, typ string) *Symbol {
	sym := &Symbol{
		pos:  len(symlist) + 1,
		name: name,
		nSpace: &LocalVariable{
			typ: typ,
		},
	}
	symlist = append(symlist, sym)
	return sym
}

func findSymbol(name string) *Symbol {
	for _, sym := range symlist {
		if sym.name == name {
			return sym
		}
	}
	fmt.Println("Undefined symbol %s.\n", name)
	panic("internal error")
}
