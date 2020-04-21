package main

import (
	"fmt"
)


/*** interface definitioins ***/

type Ast interface {
	emit ()
	show (depth int)
	debug ()
}

type ArithmeticOperator interface {
	emitOperator ()
}

type Constant interface {
	emitConstant ()
}


/*** default functions ***/
func printSpace (n int) {
	fmt.Printf ("%*s", n, "")
}


func showAst (ast Ast, depth int) {
	ast.show (depth)
}

func debugAst (ast Ast) {
	ast.debug ()
}



/* ===============================
 * Arithmetic operators implementation
 * =============================== */
type AdditiveOperator struct {
}
// implements ArithmeticOperator
func (ao *AdditiveOperator) emitOperator () {
	emitCode ("\taddl\t%%ebx, %%eax")
}

type SubtractionOperator struct {
}
// implements ArithmeticOperator
func (so *SubtractionOperator) emitOperator () {
	emitCode ("\tsubl\t%%ebx, %%eax")
}

type MultiplicativeOperator struct {
}
// implements ArithmeticOperator
func (mo *MultiplicativeOperator) emitOperator () {
	emitCode ("\tpushq\t%%rdx")
	emitCode ("\timul\t%%ebx, %%eax")
	emitCode ("\tpopq\t%%rdx")
}

type DivisionOperator struct {
}
// implements AritheticOperator
func (do *DivisionOperator) emitOperator () {
	emitCode ("\tidivl\t%%ebx, %%eax")
}


/* ===============================
 * Constants implementation
 * =============================== */
type RuneConstant struct {
	rval rune
}
// implements Constant
func (rc *RuneConstant) emitConstant () {
	emitCode ("\tpushq\t$%d", rc.rval)
	frameHeight += 8
}

type IntegerConstant struct {
	ival int
}
// implements Constant
func (ic *IntegerConstant) emitConstant () {
	emitCode ("\tpushq\t$%d", ic.ival)
	frameHeight += 8
}


/* ================================================================ */




/* ================================
 * Assignment Expression
 *     implements Ast
 * ================================ */
type AssignmentExpression struct {
	left  Identifier
	right Ast
}

// implements Ast
func (ae *AssignmentExpression) emit () {
	ae.right.emit ()
	ae.left.emitLeft ()

	emitCode ("\tpopq\t%%rax")
	emitCode ("\tmovq\t0(%%rsp), %%rcx")
	emitCode ("\tmovq\t%%rcx, 0(%%rax)")
	frameHeight -= 8
}

// implements Ast
func (ae *AssignmentExpression) debug () {
	debugPrint ("assignment_expression")
	ae.left.debug ()
	ae.right.debug ()
}

//implements Ast
func (ae *AssignmentExpression) show (depth int) {
	printSpace (depth)
	fmt.Printf ("AssignmentExpression\n")
	ae.left.show (depth+1)
	ae.right.show (depth+1)
}



/* ================================
 * Arithmetic Expression
 *     implements Ast
 * ================================ */
type ArithmeticExpression struct {
	operator ArithmeticOperator
	left     Ast
	right    Ast
}

// implements Ast
func (ae *ArithmeticExpression) emit () {
	// emitCode (fmt.Sprintf ("\tmovl\t$%%d, %%%eax\n", ae.left.operand.ival))
	// emitCode (fmt.Sprintf ("\tmovl\t$%d, %%ebx\n", ae.right.operand.ival))
	ae.left.emit ()
	ae.right.emit ()
	emitCode ("\tpopq\t%%rbx")
	emitCode ("\tpopq\t%%rax")
	frameHeight -= 16
	ae.operator.emitOperator ()
	emitCode ("\tpushq\t%%rax")
	frameHeight += 8
}

// implements Ast
func (ae *ArithmeticExpression) debug () {
	debugPrint ("ast.arithmetic_expression")
	ae.left.debug ()
	ae.right.debug ()
}

// implements Ast
func (ae *ArithmeticExpression) show (depth int) {
	printSpace (depth)
	fmt.Printf ("ArithemeticExpression\n")
	ae.left.show (depth+1)
	ae.right.show (depth+1)
}


/* ================================
 * Unary Expression 
 *     implements Ast
 * ================================ */
type UnaryExpression struct {
//	operand *PrimaryExpression
	operand Ast
}

// implements Ast
func (ue *UnaryExpression) emit () {
	// emitCode ("\tpushq\t$%d", u.operand.ival)
	ue.operand.emit ()
}

// implements Ast
func (ue *UnaryExpression) debug () {
	debugPrint ("ast.unary_expression");
	ue.operand.debug ()
}

// implements Ast
func (ue *UnaryExpression) show (depth int) {
	printSpace (depth)
	fmt.Printf ("UnaryExpression\n")
	ue.operand.show (depth+1)
}


/* ================================
 * Primary Expression 
 *     implements Ast
 * ================================ */
type PrimaryExpression struct {
	child Ast
}


// implements Ast
func (pe *PrimaryExpression) emit () {
	pe.child.emit ()
}

// implements Ast
func (pe *PrimaryExpression) debug () {
	debugPrint ("ast.primary_expression")
}

// implements Ast
func (pe *PrimaryExpression) show (depth int) {
	printSpace (depth)
	fmt.Printf ("PrimaryExpression\n")
	pe.child.show (depth+1)
}


/* ================================
 * AstConstant
 *     implements Ast
 * ================================ */
type AstConstant struct {
	constant Constant
}

// implements Ast
func (ac *AstConstant) emit () {
	ac.constant.emitConstant ()
}

// implements Ast
func (ac *AstConstant) debug () {
	debugPrintWithVariable ("ast.constant", ac.constant)
}

// implements Ast
func (ac *AstConstant) show (depth int) {
	printSpace (depth)
	fmt.Printf ("AstConstant\n")
}


/* ================================
 * Identifier
 *     implements Ast
 * ================================ */
type Identifier struct {
	symbol *Symbol
}

func (id *Identifier) emitLeft () {
	id.symbol.emitSymbol (LEFT)
}

// implements Ast
func (id *Identifier) emit () {
	id.symbol.emitSymbol (RIGHT)
}

// implements Ast 
func (id *Identifier) debug () {
	debugPrintWithVariable ("ast.identifier", id.symbol.name)
}

// implemebts Ast
func (id *Identifier) show (depth int) {
	printSpace (depth)
	fmt.Printf ("Identifier")
}


/* ================================
 * FunCall
 *     implements Ast
 * ================================ */
type FunCall struct {
	fname string
	args  []Ast
}

// implements Ast
func (fc *FunCall) emit () {
	var regs = [] string {"rdi", "rsi"}
	for i, _ := range fc.args {
		emitCode ("\tpushq\t%%%s", regs[i])
		frameHeight += 8
	}

	// stacking paddings
	var fh int
	fh = (frameHeight + 8*len(fc.args)) % 16   // for argument
	if fh != 0 {
		padding := 16 - fh
		emitCode ("\tsubq\t$%d, %%rsp  # stack padding", padding)
		frameHeight += padding
	}

	for _, arg := range fc.args {
		arg.emit ()
		// emitCode ("\tpushq\t%%rax")
		// frameHeight += 8
	}

	for i, _ := range fc.args {
		j := len (fc.args) - 1 - i
		emitCode ("\tpopq\t%%%s", regs[j])
		frameHeight -= 8
	}
	emitCode ("# frame height %d", frameHeight)
	emitCode ("\tmovq\t$0, %%rax")
	emitCode ("\tcallq\t_%s\t# frame height %d", fc.fname, frameHeight)

	if fh != 0 {
		padding := 16 - fh
		emitCode ("\taddq\t$%d, %%rsp  # pop padding", padding)
		frameHeight -= padding
	}

	for i, _ := range fc.args {
		j := len (fc.args) - 1 - i
		emitCode ("\tpopq\t%%%s", regs[j])
	}
}

// implements Ast
func (fc *FunCall) debug () {
}

// implements Ast
func (fc *FunCall) show (depth int) {
	printSpace (depth)
	fmt.Printf ("FunCall\n")
	for _, v := range fc.args {
		v.show (depth+1)
	}
}



/* ================================
 * AstString
 *     implements Ast
 * ================================ */
type AstString struct {
	sval   string
	slabel string
}

// implement Ast
func (as *AstString) emit () {
	emitCode ("\tleaq\t.%s(%%rip), %%rax", as.slabel)
	emitCode ("\tpushq\t%%rax")
	frameHeight += 8
}

// implement Ast
func (as *AstString) debug () {
}

// implements ast
func (as *AstString) show (depth int) {
	printSpace (depth)
	fmt.Printf ("AstString (%s)\n", as.sval)
}
