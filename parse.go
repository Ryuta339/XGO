package main

import (
	"fmt"
	"strconv"
)

var stringIndex = 0
var stringList [] *AstString


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
		tok := lookahead (1)
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		switch tok.sval {
		case "+":
			consumeToken ("+")
			right := parseMultiplicativeExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &AdditiveOperator {},
				left:     ast,
				right:    right,
			}
		case "-":
			consumeToken ("-")
			right := parseMultiplicativeExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &SubtractionOperator {},
				left:     ast,
				right:    right,
			}
		default:
			return ast
		}
	}
	return ast
}

func parseMultiplicativeExpression () Ast {
	var ast Ast = parseUnaryExpression ()
	for {
		tok := lookahead (1)
		if tok == nil {
			return ast
		}
		if tok.typ != "punct" {
			return ast
		}
		switch tok.sval {
		case "*":
			consumeToken ("*")
			right := parseUnaryExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &MultiplicativeOperator {},
				left:     ast,
				right:    right,
			}
		case "/" :
			consumeToken ("/")
			right := parseUnaryExpression ()
			right.debug ()
			ast = &ArithmeticExpression {
				operator: &DivisionOperator {},
				left:     ast,
				right:    right,
			}
		case "+", "-":
			return ast
		default:
			return ast
		}
	}
	return ast
}


func parseUnaryExpression () *UnaryExpression {
	tok := lookahead (1)
	var ast Ast

	switch tok.typ {
	case "int", "rune", "string", "identifier":
		ast = parsePrimaryExpression ()
		ast.debug ()
		return &UnaryExpression {
			operand: ast,
		}
	default:
		fmt.Printf ("Unexpected token %v in parseUnaryExpression.\n", tok.sval)
		panic ("internal error")
	}
	debugPrint ("nil")

	return nil
}


func parsePrimaryExpression () Ast {
	tok := lookahead (1)
	if tok == nil {
		return nil
	}
	switch tok.typ {
	case "int", "rune", "string":
		ast := parseConstant ()
		ast.debug ()
		return &PrimaryExpression {
			child: ast,
		}
	case "identifier":
		ast := parseIdentifierOrFuncall ()
		ast.debug ()
		return &PrimaryExpression {
			child: ast,
		}
	default:
		fmt.Printf ("Unexpected token %v in parsePrimaryExpression.\n", tok.sval)
		panic ("internal error")
	}



	return nil
}

func parseConstant () Ast {
	tok := lookahead (1)
	if tok == nil {
		fmt.Printf ("tok is nil\n")
		panic ("internal error")
	}
	switch (tok.typ) {
	case "int":
		ival, _ := strconv.Atoi (tok.sval)
		nextToken ()
		return &AstConstant {
			constant: &IntegerConstant {
				ival: ival,
			},
		}
	case "rune":
		rarr := []rune (tok.sval)
		nextToken ()
		return &AstConstant {
			constant: &RuneConstant {
				rval: rarr[0],
			},
		}
	case "string":
		nextToken ()
		ast :=  &AstString {
			sval: tok.sval,
			slabel: fmt.Sprintf ("L%d", stringIndex),
		}
		stringIndex ++
		stringList = append (stringList, ast)
		return ast
	default:
		fmt.Printf ("unknown token %v in parseConstant\n", tok)
		debugToken (tok)
		panic ("internal error")
	}
	return nil
}

func parseSymbol () *Identifier {
	tok := lookahead (1)
	if tok == nil {
		fmt.Printf ("tok is nil\n")
		panic ("internal error")
	}
	if tok.typ == "symbol" {
		nextToken ()
		sym := findSymbol (tok.sval)
		if sym == nil {
			sym = makeSymbol (tok.sval, "int")
		}
		return &Identifier {
			symbol: sym,
		}
	} else {
		fmt.Printf ("Unexpected token %v in parseSymbol.\n", tok.sval)
		panic ("internal error")
	}
}


func parseIdentifierOrFuncall () Ast {
	tok := lookahead (1)
	name := tok.sval
	nextToken ()
	tok = lookahead (1)
	if tok != nil && tok.typ == "punct" && tok.sval == "(" {
		consumeToken ("(")
		args := parseArgumentList ()
		consumeToken (")")
		return &FunCall {
			fname : name,
			args  : args,
		}
	}

	fmt.Println ("TBD")
	return nil
}

func parseArgumentList () []Ast {
	var r []Ast
	for {
		tok := lookahead (1)
		if tok.sval == ")" {
			return r
		}
		arg := parseExpression ()
		r = append (r, arg)
		tok = lookahead (1)
		switch tok.sval {
		case ")":
			return r
		case ",":
			consumeToken (",")
			continue
		default:
			fmt.Printf ("Unexpected token %s, %s in parseArgumentList.\n", tok.typ, tok.sval)
			panic ("internal error")
		}
	}
}
