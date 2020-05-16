package main

import ()

type Symbol interface {
	emitLeftValue()
	emitRightValue()
	getName() string
}

type SymbolBase struct {
	name  string
	gtype string
	size  int
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

type Scope struct {
	symenv map[string]Symbol
	outer  *Scope
	offset int
}

func (sc *Scope) findSymbol(name string) Symbol {
	for s := sc; s != nil; s = s.outer {
		sym := s.symenv[name]
		if sym != nil {
			return sym
		}
	}
	putError("Undefined symbol %s.\n", name)
	return nil
}

func (sc *Scope) setSymbol(name string, sym Symbol) {
	sc.symenv[name] = sym
}

func (sc *Scope) isDeclaredSymbol(name string) bool {
	for s := sc; s != nil; s = s.outer {
		sym := s.symenv[name]
		if sym != nil {
			return true
		}
	}
	return false
}

func newLocalScope(outer *Scope) *Scope {
	var offset int = 8
	if outer.offset != -1 {
		offset = outer.offset
	}
	return &Scope{
		outer:  outer,
		symenv: make(map[string]Symbol),
		offset: offset,
	}
}

var globalScope *Scope = &Scope{
	outer:  nil,
	symenv: make(map[string]Symbol),
	offset: -1,
}
var currentScope *Scope

/* ================================ */

func makeSymbol(name string, gtype string) Symbol {
	var size int = 8
	var sym Symbol
	switch gtype {
	case "uint8", "int8", "byte", "bool":
		size = 2
	case "uint16", "int16":
		size = 4
	case "uint32", "int32", "uint", "int", "rune", "float":
		size = 8
	case "uint64", "int64", "uintptr", "double":
		size = 16
	}

	if currentScope.outer == nil {
		// global variable
		sym = &GlobalVariable{
			SymbolBase: SymbolBase{
				name:  name,
				gtype: gtype,
				size:  size,
			},
		}
	} else {
		// local variable
		sym = &LocalVariable{
			SymbolBase: SymbolBase{
				name:  name,
				gtype: gtype,
				size:  size,
			},
			offset: currentScope.offset,
		}
		currentScope.offset += size
	}
	currentScope.setSymbol(name, sym)
	return sym
}

func beginSymbolBlock() {
	currentScope = newLocalScope(currentScope)
}

func endSymbolBlock() []*LocalVariable {
	if currentScope.outer == nil {
		putError("Out of block")
		return nil
	}
	var lvs []*LocalVariable
	for _, sym := range currentScope.symenv {
		lvs = append(lvs, sym.(*LocalVariable))
	}
	currentScope = currentScope.outer
	return lvs
}

func getGlobalSymList() []*GlobalVariable {
	var gvs []*GlobalVariable
	for _, sym := range globalScope.symenv {
		gvs = append(gvs, sym.(*GlobalVariable))
	}
	return gvs
}
