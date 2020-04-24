package main

import (
	"fmt"
	"strconv"
)

var stringIndex = 0
var stringList [] *AstString

func parse () Ast {
	return parseTranslationUnit ()
}

func parseTranslationUnit () Ast {
	tok := lookahead (1)
	if tok == nil {
		return nil
	}
	packname := parsePackageDeclaration ()
	var childs []Ast

	for {
		tok = lookahead (1)
		if tok == nil {
			return &TranslationUnit {
				packname: packname,
				childs  : childs,
			}
		}
		
		switch tok.sval {
		case "func":
			ast := parseFunctionDefinition ()
			childs = append (childs, ast)
		default:
			putError ("func expected, but got %v.", tok.sval)
		}
	}

	return &TranslationUnit {
		packname: packname,
		childs  : childs,
	}
}

func parsePackageDeclaration () string {
	tok := lookahead (1)
	if tok == nil {
		putError ("No package declaration.")
	}
	if tok.typ != "reserved" || tok.sval != "package" {
		putError ("No package declaration.")
	}
	consumeToken ("package")

	tok = lookahead (1)
	if tok.typ != "string" {
		putError ("%s is not allowed to be package name.", tok.sval)
	}
	packname := tok.sval
	nextToken ()
	return packname
}


func parseFunctionDefinition () Ast {
	tok := lookahead (1)
	if tok.typ != "reserved" {
		putError ("func expected, but got %v", tok.typ)
	}
	consumeToken ("func")
	tok = lookahead (1)
	if tok.typ != "identifier" {
		putError ("Identifier expected, but got %v", tok.typ)
	}
	nextToken ()
	consumeToken ("(")
	consumeToken (")")
	// expect Type
	tok2 := lookahead (1)
	if tok2.sval != "{" {
		putError ("{ expected, but got %v", tok.sval)
	}
	ast := parseCompoundStatement ()
	return &FunctionDefinition {
		fname: tok.sval,
		ast  : ast,
	}
}


func parseCompoundStatement () Ast {
	var statements []Ast
	consumeToken ("{")
	for {
		tok := lookahead (1)
		switch tok.sval {
		case "}":
			consumeToken ("}")
			return &CompoundStatement {
				statements: statements,
			}
		default:
			var ast Ast = parseStatement ()
			statements = append (statements, ast)
		}
	}
	return parseStatement ()
}

func parseStatement () Ast {
	var ast Ast

	tok := lookahead (1)
	switch tok.sval {
	case "{":
		ast = parseCompoundStatement ()
	default:
		ast = parseExpression ()
	}

	return &Statement {
		ast: ast,
	}
}


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
			ast = &ArithmeticExpression {
				operator: &AdditiveOperator {},
				left:     ast,
				right:    right,
			}
		case "-":
			consumeToken ("-")
			right := parseMultiplicativeExpression ()
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
			ast = &ArithmeticExpression {
				operator: &MultiplicativeOperator {},
				left:     ast,
				right:    right,
			}
		case "/" :
			consumeToken ("/")
			right := parseUnaryExpression ()
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
		return &UnaryExpression {
			operand: ast,
		}
	default:
		putError ("Unexpected token %v in parseUnaryExpression.\n", tok.sval)
	}

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
		return &PrimaryExpression {
			child: ast,
		}
	case "identifier":
		ast := parseIdentifierOrFuncall ()
		return &PrimaryExpression {
			child: ast,
		}
	default:
		putError ("Unexpected token %v in parsePrimaryExpression.\n", tok.sval)
	}



	return nil
}

func parseConstant () Ast {
	tok := lookahead (1)
	if tok == nil {
		putError ("tok is nil\n")
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
		putError ("unknown token %v in parseConstant\n", tok)
		debugToken (tok)
	}
	return nil
}

func parseSymbol () *Identifier {
	tok := lookahead (1)
	if tok == nil {
		putError ("tok is nil\n")
		debugToken (tok)
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
		putError ("Unexpected token %v in parseSymbol.\n", tok.sval)
		debugToken (tok)
	}
	return nil
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
			putError ("Unexpected token %s in parseArgumentList.\n", tok.sval)
			debugToken (tok)
		}
	}
}


