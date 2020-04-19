package main

import (
	"fmt"
	"strconv"
)


func parseExpression () Ast {
	ast := parseAdditiveExpression ()
	return ast
}

func parseAssignmentExpression () Ast {
	return nil
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


func parseUnaryExpression () *UnaryExpression {
	ast := parsePrimaryExpression ()
	ast.debug ()
	return &UnaryExpression {
		operand: ast,
	}
}



func parsePrimaryExpression () Ast {
	tok := readToken ()
	if tok == nil {
		return nil
	}
	unreadToken ()
	switch tok.typ {
	case "int":
		ast := parseConstant ()
		ast.debug ()
		return &PrimaryExpression {
			child: ast,
		}
	case "symbol":
		ast := parseSymbol ()
		ast.debug ()
		return &PrimaryExpression {
			child: ast,
		}
	default:
		fmt.Printf ("Unexpected token %v.\n", tok.sval)
		panic ("internal error")
	}
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

func parseSymbol () *AstSymbol {
	tok := readToken ()
	if tok == nil {
		fmt.Printf ("tok is nil\n")
		panic ("internal error")
	}
	if tok.typ == "symbol" {
		sym := findSymbol (tok.sval)
		if sym == nil {
			sym = makeSymbol (tok.sval, "int")
		}
		return &AstSymbol {
			symbol: sym,
		}
	} else {
		fmt.Println ("Unexpected token %v.\n", tok.sval)
		panic ("internal error")
	}
}

