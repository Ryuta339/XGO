package main

import ()

type Symbol interface {
	emitLeftValue()
	emitRightValue()
	getName() string
}

type SymbolBase struct {
	pos   int
	name  string
	gtype string
}

/* ================================
 * LocalVariable
 *     implements Symbol
 * ================================ */
type LocalVariable struct {
	offset int
	SymbolBase
}

// implements Symbol
func (lv *LocalVariable) emitRightValue() {
	emitCode("\tpushq\t-%d(%%rbp)", lv.offset)
	frameHeight += 8
}

// implements Symbol
func (lv *LocalVariable) emitLeftValue() {
	emitCode("\tleaq\t-%d(%%rbp), %%rax", lv.offset)
	emitCode("\tpushq\t%%rax")
	frameHeight += 8
}

// implements Symbol
func (lv *LocalVariable) getName() string {
	return lv.name
}

/* ================================
 * GlobalVariable
 *     implements Symbol
 * ================================ */
type GlobalVariable struct {
	initval Constant
	SymbolBase
}

// implements Symbol
func (gv *GlobalVariable) emitRightValue() {
	emitCode("\tpushq\t_%s(%%rip)", gv.name)
}

// implements Symbol
func (gv *GlobalVariable) emitLeftValue() {
	// Global Offset Table
	emitCode("\tpushq\t_%s@GOTPCREL(%%rip)", gv.name)
}

// implements Symbol
func (gv *GlobalVariable) getName() string {
	return gv.name
}

/* ================================ */

var symbolDepth int = 0

var globalsymlist []*GlobalVariable = make([]*GlobalVariable, 0)
var globalsymenv map[string]*GlobalVariable = make(map[string]*GlobalVariable)

var localsymlist []*LocalVariable
var localsymenv map[string]*LocalVariable
var localsymOffset int

func makeSymbol(name string, gtype string) Symbol {
	var offset int = 8
	var sym Symbol
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
		gv := &GlobalVariable{
			SymbolBase: SymbolBase{
				pos:   len(globalsymlist) + 1,
				name:  name,
				gtype: gtype,
			},
		}
		globalsymlist = append(globalsymlist, gv)
		globalsymenv[name] = gv
		sym = gv
	} else {
		// local variable
		lv := &LocalVariable{
			SymbolBase: SymbolBase{
				pos:   len(localsymlist) + 1,
				name:  name,
				gtype: gtype,
			},
			offset: localsymOffset,
		}
		localsymOffset += offset
		localsymlist = append(localsymlist, lv)
		localsymenv[name] = lv
		sym = lv
	}
	return sym
}

func findSymbol(name string) Symbol {
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
	putError("Undefined symbol %s.\n", name)
	return nil
}

func isDeclaredSymbol(name string) bool {
	_, okL := localsymenv[name]
	_, okG := globalsymenv[name]
	return okL || okG
}

func beginSymbolBlock() {
	localsymlist = make([]*LocalVariable, 0)
	localsymenv = make(map[string]*LocalVariable)
	localsymOffset = 8
	symbolDepth++
}

func endSymbolBlock() []*LocalVariable {
	if symbolDepth == 0 {
		putError("Out of block")
		return nil
	}
	symbolDepth--
	tmp := localsymlist
	localsymlist = nil
	return tmp
}
