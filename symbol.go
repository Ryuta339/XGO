package main

import (
	"fmt"
)

type NameSpace interface {
	emitLeftValue(sym *Symbol)
	emitRightValue(sym *Symbol)
}


/* ================================
 * LocalVariable
 *     implements NameSpace 
 * ================================ */
type LocalVariable struct {
	gtype  string
	offset int
}

// implements NameSpace
func (lv *LocalVariable) emitRightValue(sym *Symbol) {
	emitCode("\tpushq\t-%d(%%rbp)", sym.pos*lv.offset)
	frameHeight += 8
}
// implements NameSpace
func (lv *LocalVariable) emitLeftValue(sym *Symbol) {
	emitCode("\tleaq\t-%d(%%rbp), %%rax", sym.pos*lv.offset)
	emitCode("\tpushq\t%%rax")
	frameHeight += 8
}

/* ================================
 * GlobalVariable
 *     implements NameSpace 
 * ================================ */
type GlobalVariable struct {
	gtype string
}

// implements NameSpace
func (gv *GlobalVariable) emitRightValue(sym *Symbol) {
}
// implements NameSpace
func (gv *GlobalVariable) emitLeftValue(sym *Symbol) {
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

var symbolDepth int = 0

var globalsymlist []*Symbol          = make([]*Symbol,0)
var globalsymenv  map[string]*Symbol = make(map[string]*Symbol)

var localsymlist []*Symbol
var localsymenv  map[string]*Symbol
var localsymOffset int

func makeSymbol(name string, gtype string) *Symbol {
	var offset int = 8
	var sym *Symbol
	switch gtype {
	case "uint8", "int8", "byte", "bool":
		offset = 2
	case "uint16", "int16":
		offset = 4
	case "uint32", "int32", "uint", "int", "rune", "float":
		offset = 8
	case "uint64", "int64", "uintptr", "double":
		offset = 16
	}
	if symbolDepth == 0 {
		// global variable
		sym = &Symbol{
			pos:    len(globalsymlist) + 1,
			name:   name,
			nSpace:&GlobalVariable{
				gtype : gtype,
			},
		}
		globalsymlist = append(globalsymlist, sym)
		globalsymenv[name] = sym
	} else {
		// local variable
		sym = &Symbol{
			pos:    len(localsymlist) + 1,
			name:   name,
			nSpace: &LocalVariable{
				gtype : gtype,
				offset: localsymOffset,
			},
		}
		localsymOffset += offset
		localsymlist = append(localsymlist, sym)
		localsymenv[name] = sym
	}
	return sym
}

func findSymbol(name string) *Symbol {
	for _, sym := range localsymlist {
		if sym.name == name {
			return sym
		}
	}
	for _, sym := range globalsymlist {
		if sym.name == name {
			return sym
		}
	}
	fmt.Println("Undefined symbol %s.\n", name)
	panic("internal error")
}

func isDeclaredSymbol(name string) bool {
	_, okL := localsymenv[name]
	_, okG := globalsymenv[name]
	return okL || okG
}

func beginSymbolBlock() {
	localsymlist = make([]*Symbol, 0)
	localsymenv = make(map[string]*Symbol)
	localsymOffset = 8
	symbolDepth++
}

func endSymbolBlock() []*Symbol {
	if symbolDepth == 0 {
		putError ("global")
		return globalsymlist
	}
	symbolDepth--
	tmp := localsymlist
	localsymlist = nil
	return tmp
}
