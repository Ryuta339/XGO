package main

import (
	"strconv"
)

/*** interface definitioins ***/
type Constant interface {
	emitConstant()
	toStringValue() string
}

/* ===============================
 * Constants implementation
 * =============================== */
type RuneConstant struct {
	rval rune
}

// implements Constant
func (rc *RuneConstant) emitConstant() {
	emitCode("\tpushq\t$%d", rc.rval)
	frameHeight += 8
}

// implements Costant
func (rc *RuneConstant) toStringValue() string {
	return string(rc.rval)
}

type IntegerConstant struct {
	ival int
}

// implements Constant
func (ic *IntegerConstant) emitConstant() {
	emitCode("\tpushq\t$%d", ic.ival)
	frameHeight += 8
}

// implements Constant
func (ic *IntegerConstant) toStringValue() string {
	return strconv.Itoa(ic.ival)
}
