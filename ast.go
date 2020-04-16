package main

import (
	"fmt"
	"strconv"
);



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


/* ================================ */

func parseExpression () Ast {
	ast := parseAdditiveExpression ()
	return ast
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

func parseAdditiveExpression () Ast {
	var ast Ast = parseMultiplicativeExpression ()
	for {
		tok := readToken ()
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		switch tok.sval {
		case "+":
			right := parseMultiplicativeExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &AdditiveOperator {},
				left:     ast,
				right:    right,
			}
		case "-":
			right := parseMultiplicativeExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &SubtractionOperator {},
				left:     ast,
				right:    right,
			}
		default:
			// fmt.Printf ("unknown token %v in parseAdditiveExpression\n", tok)
			// debugToken (tok)
			// panic ("internal error")
			unreadToken ()
			return ast
		}
	}
	return ast
}

func parseMultiplicativeExpression () Ast {
	var ast Ast = parseUnaryExpression ()
	for {
		tok := readToken ()
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		switch tok.sval {
		case "*":
			right := parseUnaryExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &MultiplicativeOperator {},
				left:     ast,
				right:    right,
			}
		case "/" :
			right := parseUnaryExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &DivisionOperator {},
				left:     ast,
				right:    right,
			}
		case "+", "-":
			unreadToken ()
			return ast
		default:
			// fmt.Printf ("unknown token %v.\n", tok)
			// debugToken (tok)
			// panic ("internal error")
			unreadToken ()
			return ast
		}
	}
	return ast
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

func parseUnaryExpression () *UnaryExpression {
	ast := parsePrimaryExpression ()
	ast.debug ()
	return &UnaryExpression {
		operand: ast,
	}
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

func parsePrimaryExpression () Ast {
	ast := parseConstant ()
	ast.debug ()
	return &PrimaryExpression {
		child: ast,
	}
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

func parseConstant () Ast {
	tok := readToken ()
	if tok == nil {
		fmt.Printf ("tok is nil\n")
		panic ("internal error")
	}
	switch (tok.typ) {
	case "int":
		ival, _ := strconv.Atoi (tok.sval)
		return &AstConstant {
			constant: &IntegerConstant {
				ival: ival,
			},
		}
	case "rune":
		rarr := []rune (tok.sval)
		return &AstConstant {
			constant: &RuneConstant {
				rval: rarr[0],
			},
		}
	case "string":
		return nil
	default:
		fmt.Printf ("unknown token %v in parseConstant\n", tok)
		debugToken (tok)
		panic ("internal error")
	}
	return nil
}
