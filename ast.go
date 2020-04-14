package main

import (
	"fmt"
	"strconv"
);

var frameHeight int = 0;


/*** Interface Definitioin ***/

type Ast interface {
	emit ()
	debug ()
}

type ArithmeticOperator interface {
	emit ()
}


func emitCode (code string) {
	fmt.Println (code)
}

func generate (ast Ast) {
	emitCode ("\t.global _mymain")
	emitCode ("_mymain:");
	ast.emit ()
	emitCode ("\tpopq\t%rax")
	frameHeight -= 4
	emitCode ("\tret");
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
func (ao *AdditiveOperator) emit () {
	emitCode ("\taddl\t%ebx, %eax")
}

type SubtractionOperator struct {
}
// implements ArithmeticOperator
func (so *SubtractionOperator) emit () {
	emitCode ("\tsubl\t%ebx, %eax")
}

type MultiplicativeOperator struct {
}
// implements ArithmeticOperator
func (mo *MultiplicativeOperator) emit () {
	emitCode ("\timul\t%ebx, %eax")
}

type DivisionOperator struct {
}
// implements AritheticOperator
func (do *DivisionOperator) emit () {
	emitCode ("\tidivl\t%ebx, %eax")
}



/* ================================ */

func parseExpression () Ast {
	ast := parseAdditiveExpression ()
	return ast
}


/* ================================
 * Arithmetic Expression
 * ================================ */
type ArithmeticExpression struct {
	operator ArithmeticOperator
	left     Ast
	right    Ast
}

// implements Ast
func (ae *ArithmeticExpression) emit () {
	// emitCode (fmt.Sprintf ("\tmovl\t$%d, %%eax\n", ae.left.operand.ival))
	// emitCode (fmt.Sprintf ("\tmovl\t$%d, %%ebx\n", ae.right.operand.ival))
	ae.left.emit ()
	ae.right.emit ()
	emitCode ("\tpopq\t%rbx")
	emitCode ("\tpopq\t%rax")
	frameHeight -= 8
	ae.operator.emit ()
	emitCode ("\tpushq\t%rax")
	frameHeight += 4
}

// implements Ast
func (ae *ArithmeticExpression) debug () {
	debugPrint ("ast.arithmetic_exression")
	ae.left.debug ()
	ae.right.debug ()
}

func parseAdditiveExpression () Ast {
	ast := parseMultiplicativeExpression ()
	for {
		tok := readToken ()
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		if tok.sval == "+" {
			right := parseMultiplicativeExpression ()
			right.debug ()
			return &ArithmeticExpression {
				operator: &AdditiveOperator {},
				left:     ast,
				right:    right,
			}
		} else if tok.sval == "-" {
			right := parseMultiplicativeExpression ()
			right.debug ()
			return &ArithmeticExpression {
				operator: &SubtractionOperator {},
				left:     ast,
				right:    right,
			}
		} else {
			fmt.Printf ("unknown token %v in parseAdditiveExpression\n", tok)
			debugToken (tok)
			panic ("internal error")
		}
	}
	return ast
}

func parseMultiplicativeExpression () Ast {
	ast := parseUnaryExpression ()
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
			return &ArithmeticExpression {
				operator: &MultiplicativeOperator {},
				left:     ast,
				right:    right,
			}
		case "/" :
			right := parseUnaryExpression ()
			right.debug ()
			return &ArithmeticExpression {
				operator: &DivisionOperator {},
				left:     ast,
				right:    right,
			}
		case "+", "-":
			unreadToken ()
			return ast
		default:
			fmt.Printf ("unknown token %v.\n", tok)
			debugToken (tok)
			panic ("internal error")
		}
	}
	return ast
}

/* ================================
 * Unary Expression 
 * ================================ */
type UnaryExpression struct {
	operand *PrimaryExpression
}


// implements Ast
func (u *UnaryExpression) emit () {
	emitCode (fmt.Sprintf ("\tpushq\t$%d", u.operand.ival))
	frameHeight += 4
}

// implements Ast
func (u *UnaryExpression) debug () {
	debugPrintWithVariable ("ast.unary_expression", u.operand);
}

func parseUnaryExpression () *UnaryExpression {
	tok := readToken ()
	ival, _ := strconv.Atoi (tok.sval)
	return &UnaryExpression {
		operand: &PrimaryExpression {
			typ:  "int",
			ival: ival,
		},
	}
}


/* ================================
 * Primary Expression 
 * ================================ */
type PrimaryExpression struct {
	typ     string
	ival    int
}


// implements Ast
func (p *PrimaryExpression) emit () {
}

// implements Ast
func (p *PrimaryExpression) debug () {
	debugPrintWithVariable ("ast.primary_expression", fmt.Sprintf ("%d", p.ival))
}

func parsePrimaryExpression () Ast {
	return nil
}
