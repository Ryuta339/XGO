package main




/*** interface definitioins ***/

type Ast interface {
	emit ()
	debug ()
}

type ArithmeticOperator interface {
	emitOperator ()
}

type Constant interface {
	emitConstant ()
}


/*** default functions ***/

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
}

type IntegerConstant struct {
	ival int
}
// implements Constant
func (ic *IntegerConstant) emitConstant () {
	emitCode ("\tpushq\t$%d", ic.ival)
}


/* ================================================================ */




/* ================================
 * Assignment Expression
 *     implements Ast
 * ================================ */
type AssignmentExpression struct {
	left  AstSymbol
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
	frameHeight -= 8
	ae.operator.emitOperator ()
	emitCode ("\tpushq\t%%rax")
	frameHeight += 4
}

// implements Ast
func (ae *ArithmeticExpression) debug () {
	debugPrint ("ast.arithmetic_expression")
	ae.left.debug ()
	ae.right.debug ()
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
func (u *UnaryExpression) emit () {
	// emitCode ("\tpushq\t$%d", u.operand.ival)
	u.operand.emit ()
	frameHeight += 4
}

// implements Ast
func (u *UnaryExpression) debug () {
	debugPrint ("ast.unary_expression");
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

func (ac *AstConstant) debug () {
	debugPrintWithVariable ("ast.constant", ac.constant)
}



/* ================================
 * AstSymbol
 *     implements Ast
 * ================================ */
type AstSymbol struct {
	symbol *Symbol
}

func (as *AstSymbol) emitLeft () {
	as.symbol.emitSymbol (LEFT)
}

// implements Ast
func (as *AstSymbol) emit () {
	as.symbol.emitSymbol (RIGHT)
}

// implements Ast 
func (as *AstSymbol) debug () {
	debugPrintWithVariable ("ast.symbol", as.symbol.name)
}

