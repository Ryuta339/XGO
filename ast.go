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


/*** Unary Expression ***/
type UnaryExpression struct {
	SuperAst
	operand *PrimaryExpression
}

func parseUnaryExpression () Ast {
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

// implements Ast
func (u *UnaryExpression) emit () {
	fmt.Printf ("\tmovl\t$%d, %%eax\n", u.operand.ival);
}

// implements Ast
func (u *UnaryExpression) debug () {
	debugPrint ("ast.unary_expression", u.operand);
}

/*** Primary Expression ***/
type PrimaryExpression struct {
	SuperAst
	ival    int
}


// implements Ast
func (p *PrimaryExpression) emit () {
}

// implements Ast
func (p *PrimaryExpression) debug () {
}
