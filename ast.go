package main

import (
	"fmt"
	"strconv"
);



type Ast interface {
	emit ()
	debug ()
}

type SuperAst struct {
	typ     string
}




func parseExpression () Ast {
	ast := parseBinaryExpression ()
	return ast
}


/* ================================
 * Binary Expression
 * ================================ */
type BinaryExpression struct {
	SuperAst
	operator string
	left     *UnaryExpression
	right    *UnaryExpression
}

// implements Ast
func (b *BinaryExpression) emit () {
	fmt.Printf ("\tmovl\t$%d, %%ebx\n", b.left.operand.ival);
	fmt.Printf ("\tmovl\t$%d, %%eax\n", b.right.operand.ival);
	var s string
	if b.operator == "+" {
		s = "addl"
	} else if b.operator == "-" {
		s = "subl"
	}
	fmt.Printf ("\t%s\t%%ebx, %%eax\n", s)
}

// implements Ast
func (b *BinaryExpression) debug () {
	debugPrintWithVariable ("ast.binary_exression", b.typ)
	b.left.debug ()
	b.right.debug ()
}

func parseBinaryExpression () Ast {
	ast := parseUnaryExpression ()
	for {
		tok := readToken ()
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		if tok.sval == "+" || tok.sval == "-" {
			right := parseUnaryExpression ()
			right.debug ()
			return &BinaryExpression {
				SuperAst: SuperAst {
					typ: "binary_expression",
				},
				operator: tok.sval,
				left:     ast,
				right:    right,
			}
		} else {
			fmt.Printf ("unknown token%v\n", tok)
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
	SuperAst
	operand *PrimaryExpression
}


// implements Ast
func (u *UnaryExpression) emit () {
	fmt.Printf ("\tmovl\t$%d, %%eax\n", u.operand.ival);
}

// implements Ast
func (u *UnaryExpression) debug () {
	debugPrintWithVariable ("ast.unary_expression", u.operand);
}

func parseUnaryExpression () *UnaryExpression {
	tok := readToken ()
	ival, _ := strconv.Atoi (tok.sval)
	return &UnaryExpression {
		SuperAst: SuperAst{
			typ : "unary_expression",
		},
		operand: &PrimaryExpression {
			SuperAst: SuperAst{
				typ:  "int",
			},
			ival: ival,
		},
	}
}


/* ================================
 * Primary Expression 
 * ================================ */
type PrimaryExpression struct {
	SuperAst
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
