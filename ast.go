package main

import (
	"fmt"
)

/*** interface definitioins ***/

type Ast interface {
	emit()
	show(depth int)
	Debuggable
}

type LeftValue interface {
	Ast
	emitLeft()
}

/*** registers for function arguments ***/
var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

/*** default functions ***/
func printSpace(n int) {
	fmt.Printf("%*s", n, "")
}

func showAst(ast Ast, depth int) {
	ast.show(depth)
}

/* ================================================================ */

/* ================================
 * TranslationUnit
 *     implements Ast and Debuggale
 * ================================ */
type TranslationUnit struct {
	packname   string
	packages   []string
	childs     []Ast
	globalvars []*GlobalVariable
}

// implements Ast
func (tu *TranslationUnit) emit() {
	for _, sym := range tu.globalvars {
		emitCode(".global\t_%s", sym.name)
		emitCode("_%s:", sym.name)
		emitCode(".long\t%s", sym.initval.toStringValue())
	}
	for _, child := range tu.childs {
		child.emit()
	}
}

// implements Ast
func (tu *TranslationUnit) show(depth int) {
	printSpace(depth)
	fmt.Printf("TranslationUnit (%s) {\n", tu.packname)
	for _, pkg := range tu.packages {
		fmt.Printf("(import \"%s\")\n", pkg)
	}
	fmt.Printf("}\n")
	for _, child := range tu.childs {
		child.show(depth + 1)
	}
}

// implements Ast
func (tu *TranslationUnit) debug() {
	debugPrintln("ast.translation_unit")
	for _, child := range tu.childs {
		child.debug()
	}
}

/* ================================
 * Function Definition
 *     implements Ast
 * ================================ */
type FunctionDefinition struct {
	fname   string
	rettype string
	params  []*LocalVariable
	ast     Ast
	space   int
}

// implements Ast
func (fd *FunctionDefinition) emit() {
	emitFuncPrologue(fd.fname)
	var stacksize int = 0
	for _, v := range fd.params {
		stacksize += v.size
	}
	for i, _ := range fd.params {
		emitCode("\tpushq\t%%%s", regs[i])
		frameHeight += 8
	}
	if fd.space > 0 {
		emitCode("# allocate local variable area")
		emitCode("\tsubq\t$%d,\t%%rsp", fd.space)
		frameHeight += fd.space
		stacksize += fd.space
	}
	fd.ast.emit()
	if stacksize > 0 {
		emitCode("# free function argument and local variable area")
		emitCode("\taddq\t$%d,\t%%rsp", stacksize)
		frameHeight -= stacksize
	}
	emitCode("\tmovl\t$0, %%eax") // return 0
	emitFuncEpilogue()
}

// implements Ast
func (fd *FunctionDefinition) show(depth int) {
	printSpace(depth)
	fmt.Printf("FunctionDefinition(%s)\n", fd.fname)
	fd.ast.show(depth + 1)
}

// implements Ast
func (fd *FunctionDefinition) debug() {
	debugPrintf("funcdef: %s", fd.fname)
	fd.ast.debug()
}

/* ================================
 * Compound Statement
 *     implements Ast
 * ================================ */
type CompoundStatement struct {
	statements []Ast
	localvars  []*LocalVariable
}

// implements Ast
func (cs *CompoundStatement) emit() {
	for _, v := range cs.localvars {
		emitCode ("\tmovq\t$0,\t-%d(%%rbp)", v.offset)
	}
	for _, statement := range cs.statements {
		statement.emit()
	}
}

// implements Ast
func (cs *CompoundStatement) debug() {
	debugPrintln("ast.compound_statement")
	for _, statement := range cs.statements {
		statement.debug()
	}
}

// implements Ast
func (cs *CompoundStatement) show(depth int) {
	printSpace(depth)
	fmt.Printf("CompoundStatement\n")
	for _, ast := range cs.statements {
		ast.show(depth + 1)
	}
}

/* ================================
 * Statement
 *     implements Ast
 * ================================ */
type Statement struct {
	ast Ast
}

// implements Ast
func (s *Statement) emit() {
	s.ast.emit()
}

// implements Ast
func (s *Statement) debug() {
	debugPrintln("ast.statement")
	s.ast.debug()
}

// implements Ast
func (s *Statement) show(depth int) {
	printSpace(depth)
	fmt.Printf("Statement\n")
	s.ast.show(depth + 1)
}

/* ================================
 * DeclarationStatement
 *     implements Ast
 * ================================ */
type DeclarationStatement struct {
	sym    *LocalVariable
	assign Ast
}

// implements Ast
func (ds *DeclarationStatement) emit() {
	if ds.assign != nil {
		ds.assign.emit()
	}
}

// implements Ast
func (ds *DeclarationStatement) debug() {
	debugPrintln("ast.declaration_statement")
	if ds.assign != nil {
		ds.assign.debug()
	}
}

// implements Ast
func (ds *DeclarationStatement) show(depth int) {
	printSpace(depth)
	fmt.Printf("DeclarationStatement(%s)\n", ds.sym.name)
	if ds.assign != nil {
		ds.assign.show(depth + 1)
	}
}

/* ================================
 * Assignment Expression
 *     implements Ast
 * ================================ */
type AssignmentExpression struct {
	left  LeftValue
	right Ast
}

// implements Ast
func (ae *AssignmentExpression) emit() {
	ae.right.emit()
	ae.left.emitLeft()

	emitCode("\tpopq\t%%rax")
	emitCode("\tmovq\t0(%%rsp), %%rcx")
	emitCode("\tmovq\t%%rcx, 0(%%rax)")
	frameHeight -= 8
}

// implements Ast
func (ae *AssignmentExpression) debug() {
	debugPrintln("ast.assignment_expression")
	ae.left.debug()
	ae.right.debug()
}

//implements Ast
func (ae *AssignmentExpression) show(depth int) {
	printSpace(depth)
	fmt.Printf("AssignmentExpression\n")
	ae.left.show(depth + 1)
	ae.right.show(depth + 1)
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
func (ae *ArithmeticExpression) emit() {
	// emitCode (fmt.Sprintf ("\tmovl\t$%d, %%ebx\n", ae.right.operand.ival))
	ae.left.emit()
	ae.right.emit()
	emitCode("\tpopq\t%%rbx")
	emitCode("\tpopq\t%%rax")
	frameHeight -= 16
	ae.operator.emitOperator()
	emitCode("\tpushq\t%%rax")
	frameHeight += 8
}

// implements Ast
func (ae *ArithmeticExpression) debug() {
	debugPrintln("ast.arithmetic_expression")
	ae.left.debug()
	ae.right.debug()
}

// implements Ast
func (ae *ArithmeticExpression) show(depth int) {
	printSpace(depth)
	fmt.Printf("ArithemeticExpression\n")
	ae.left.show(depth + 1)
	ae.right.show(depth + 1)
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
func (ue *UnaryExpression) emit() {
	ue.operand.emit()
}

// implements Ast
func (ue *UnaryExpression) debug() {
	debugPrintln("ast.unary_expression")
	ue.operand.debug()
}

// implements Ast
func (ue *UnaryExpression) show(depth int) {
	printSpace(depth)
	fmt.Printf("UnaryExpression\n")
	ue.operand.show(depth + 1)
}

/* ================================
 * Primary Expression
 *     implements Ast
 * ================================ */
type PrimaryExpression struct {
	child Ast
}

// implements Ast
func (pe *PrimaryExpression) emit() {
	pe.child.emit()
}

// implements Ast
func (pe *PrimaryExpression) debug() {
	debugPrintln("ast.primary_expression")
	pe.child.debug()
}

// implements Ast
func (pe *PrimaryExpression) show(depth int) {
	printSpace(depth)
	fmt.Printf("PrimaryExpression\n")
	pe.child.show(depth + 1)
}

/* ================================
 * AstConstant
 *     implements Ast
 * ================================ */
type AstConstant struct {
	constant Constant
}

// implements Ast
func (ac *AstConstant) emit() {
	ac.constant.emitConstant()
}

// implements Ast
func (ac *AstConstant) debug() {
	debugPrintlnWithVariable("ast.constant", ac.constant)
}

// implements Ast
func (ac *AstConstant) show(depth int) {
	printSpace(depth)
	fmt.Printf("AstConstant (%s)\n", ac.constant.toStringValue())
}

/* ================================
 * Identifier
 *     implements Ast and LeftValue
 * ================================ */
type Identifier struct {
	symbol Symbol
}

// implements LeftValue
func (id *Identifier) emitLeft() {
	id.symbol.emitLeftValue()
}

// implements Ast
func (id *Identifier) emit() {
	id.symbol.emitRightValue()
}

// implements Ast
func (id *Identifier) debug() {
	debugPrintlnWithVariable("ast.identifier", id.symbol.getName())
}

// implemebts Ast
func (id *Identifier) show(depth int) {
	printSpace(depth)
	fmt.Printf("Identifier(%s)\n", id.symbol.getName())
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
func (fc *FunCall) emit() {
	for i, _ := range fc.args {
		emitCode("\tpushq\t%%%s", regs[i])
		frameHeight += 8
	}

	// stacking paddings
	var fh int
	// emitCode("# frame height %d before arguments", frameHeight)
	// fh = (frameHeight + 8*len(fc.args)) % 16   // for argument
	fh = frameHeight % 16
	if fh != 0 {
		padding := 16 - fh
		emitCode("\tsubq\t$%d, %%rsp  # stack padding", padding)
		frameHeight += padding
	}

	for _, arg := range fc.args {
		arg.emit()
	}

	for i, _ := range fc.args {
		j := len(fc.args) - 1 - i
		emitCode("\tpopq\t%%%s", regs[j])
		frameHeight -= 8
	}
	// emitCode("# frame height %d after arguments", frameHeight)
	emitCode("\tmovq\t$0, %%rax")
	emitCode("\tcallq\t_%s\t# frame height %d", fc.fname, frameHeight)

	if fh != 0 {
		padding := 16 - fh
		emitCode("\taddq\t$%d, %%rsp  # pop padding", padding)
		frameHeight -= padding
	}

	for i, _ := range fc.args {
		j := len(fc.args) - 1 - i
		emitCode("\tpopq\t%%%s", regs[j])
	}
}

// implements Ast
func (fc *FunCall) debug() {
	debugPrintln("ast.funcall")
	for _, v := range fc.args {
		v.debug()
	}
}

// implements Ast
func (fc *FunCall) show(depth int) {
	printSpace(depth)
	fmt.Printf("FunCall\n")
	for _, v := range fc.args {
		v.show(depth + 1)
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
func (as *AstString) emit() {
	emitCode("\tleaq\t.%s(%%rip), %%rax", as.slabel)
	emitCode("\tpushq\t%%rax")
	frameHeight += 8
}

// implement Ast
func (as *AstString) debug() {
	debugPrintln("ast.string")
}

// implements ast
func (as *AstString) show(depth int) {
	printSpace(depth)
	fmt.Printf("AstString (\"%s\")\n", as.sval)
}
