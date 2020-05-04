package main

import (
	"fmt"
)

type NameSpace interface {
	emitLeftValue(sym *Symbol)
	emitRightValue(sym *Symbol)
}

type LocalVariable struct {
	gtype  string
	offset int
}

// implements NameSpace
func (lv *LocalVariable) emitRightValue(sym *Symbol) {
	emitCode("\tpushq\t-%d(%%rbp)", sym.pos*lv.offset)
	frameHeight += lv.offset
}
func (lv *LocalVariable) emitLeftValue(sym *Symbol) {
	emitCode("\tleaq\t-%d(%%rbp), %%rax", sym.pos*lv.offset)
	emitCode("\tpushq\t%%rax")
	frameHeight += lv.offset
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

var symlist []*Symbol
var symenv  map[string]*Symbol

func makeSymbol(name string, gtype string) *Symbol {
	sym := &Symbol{
		pos:  len(symlist) + 1,
		name: name,
		nSpace: &LocalVariable{
			gtype : gtype,
			offset: 8,
		},
	}
	symlist = append(symlist, sym)
	symenv[name] = sym
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

func isDeclaredSymbol(name string) bool {
	_, ok := symenv[name]
	return ok
}

func beginSymbolBlock() {
	symlist = make([]*Symbol, 0)
	symenv = make(map[string]*Symbol)
}

func endSymbolBlock() []*Symbol {
	tmp := symlist
	symlist = nil
	return tmp
}
