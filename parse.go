package main

import (
	"fmt"
	"strconv"
)

var stringIndex = 0
var stringList []*AstString

func parse() Ast {
	return parseTranslationUnit()
}

func parseTranslationUnit() Ast {
	tok := lookahead(1)
	if tok == nil {
		return nil
	}
	packname := parsePackageDeclaration()
	packages := parseImport()
	var childs []Ast

	for {
		tok = lookahead(1)

		switch {
		case tok == nil:
			return &TranslationUnit{
				packname: packname,
				packages: packages,
				childs:   childs,
			}
		case tok.isReserved("func"):
			ast := parseFunctionDefinition()
			childs = append(childs, ast)
		default:
			putError("func expected, but got %v.", tok.sval)
		}
	}

	return &TranslationUnit{
		packname: packname,
		packages: packages,
		childs:   childs,
	}
}

func parsePackageDeclaration() string {
	tok := lookahead(1)
	packname := ""

	switch {
	case tok == nil:
		putError("No package declaration.")
	case tok.isReserved("package"):
		consumeToken("package")
		tok = lookahead(1)
		if tok.isTypeString() {
			packname = tok.sval
			nextToken()
		} else {
			putError("%s is not allowed to be package name.", tok.sval)
		}
	default:
		putError("No package declaration.")
	}

	return packname
}

func parseImport() []string {
	var packages []string
	for {
		tok := lookahead(1)

		switch {
		case tok == nil:
			return packages
		case tok.isReserved("import"):
			consumeToken("import")
			ps := parseImportPackageNames()
			packages = append(packages, ps...)
		default:
			return packages
		}
	}
	return packages
}

func parseImportPackageNames() []string {
	tok := lookahead(1)
	switch {
	case tok == nil:
		putError("Unexpected termination")
		return nil
	case tok.isTypeString():
		nextToken()
		return []string{tok.sval}
	case tok.isPunct("("):
		return parseImportParenthesis()
	default:
		return nil
	}
}

func parseImportParenthesis() []string {
	consumeToken("(")
	var packages []string
	for {
		tok := lookahead(1)
		switch {
		case tok == nil:
			putError("Expected ), but got EOF.")
			return packages
		case tok.isTypeString():
			packages = append(packages, tok.sval)
			nextToken()
		case tok.isPunct(")"):
			consumeToken(")")
			return packages
		default:
			putError("Expected \" or ), but got %s.", tok.sval)
			return packages
		}
	}
}

func parseFunctionDefinition() Ast {
	tok := lookahead(1)
	if !tok.isReserved("func") {
		putError("Expected func, but got %s", tok.typ)
		return nil
	}
	consumeToken("func")
	tok = lookahead(1)
	if !tok.isTypeIdentifier() {
		putError("Expected identifier, but got %s", tok.typ)
		return nil
	}
	nextToken()
	consumeToken("(")
	consumeToken(")")
	// expect Type
	tok2 := lookahead(1)
	if !tok2.isPunct("{") {
		putError("Expected {, but got %s", tok.sval)
		return &FunctionDefinition{
			fname: tok.sval,
			ast:   nil,
		}
	}
	ast := parseCompoundStatement()
	return &FunctionDefinition{
		fname: tok.sval,
		ast:   ast,
	}
}

func parseCompoundStatement() Ast {
	var statements []Ast
	consumeToken("{")
	for {
		tok := lookahead(1)
		switch {
		case tok.isPunct("}"):
			consumeToken("}")
			return &CompoundStatement{
				statements: statements,
			}
		default:
			var ast Ast = parseStatement()
			statements = append(statements, ast)
		}
	}
	return parseStatement()
}

func parseStatement() Ast {
	var ast Ast

	tok := lookahead(1)
	switch {
	case tok.isPunct("{"):
		ast = parseCompoundStatement()
	default:
		ast = parseExpression()
	}

	return &Statement{
		ast: ast,
	}
}

func parseExpression() Ast {
	ast := parseAdditiveExpression()
	return ast
}

func parseAssignmentExpression() Ast {
	return nil
}

func parseAdditiveExpression() Ast {
	var ast Ast = parseMultiplicativeExpression()
	for {
		tok := lookahead(1)
		switch {
		case tok == nil:
			return ast
		case tok.isPunct("+"):
			consumeToken("+")
			right := parseMultiplicativeExpression()
			ast = &ArithmeticExpression{
				operator: &AdditiveOperator{},
				left:     ast,
				right:    right,
			}
		case tok.isPunct("-"):
			consumeToken("-")
			right := parseMultiplicativeExpression()
			ast = &ArithmeticExpression{
				operator: &SubtractionOperator{},
				left:     ast,
				right:    right,
			}
		default:
			return ast
		}
	}
	return ast
}

func parseMultiplicativeExpression() Ast {
	var ast Ast = parseUnaryExpression()
	for {
		tok := lookahead(1)
		switch {
		case tok == nil:
			return ast
		case tok.isPunct("*"):
			consumeToken("*")
			right := parseUnaryExpression()
			ast = &ArithmeticExpression{
				operator: &MultiplicativeOperator{},
				left:     ast,
				right:    right,
			}
		case tok.isPunct("/"):
			consumeToken("/")
			right := parseUnaryExpression()
			ast = &ArithmeticExpression{
				operator: &DivisionOperator{},
				left:     ast,
				right:    right,
			}
		case tok.isPunct("+"), tok.isPunct("-"):
			return ast
		default:
			return ast
		}
	}
	return ast
}

func parseUnaryExpression() *UnaryExpression {
	tok := lookahead(1)
	var ast Ast

	switch {
	case tok == nil:
		return nil
	case tok.isTypeString(), tok.isTypeIdentifier(), tok.isTypeInt(), tok.isTypeRune():
		ast = parsePrimaryExpression()
		return &UnaryExpression{
			operand: ast,
		}
	default:
		putError("Unexpected token %v in parseUnaryExpression.\n", tok.sval)
	}

	return nil
}

func parsePrimaryExpression() Ast {
	tok := lookahead(1)
	switch {
	case tok == nil:
		return nil
	case tok.isTypeInt(), tok.isTypeRune(), tok.isTypeString():
		ast := parseConstant()
		return &PrimaryExpression{
			child: ast,
		}
	case tok.isTypeIdentifier():
		ast := parseIdentifierOrFuncall()
		return &PrimaryExpression{
			child: ast,
		}
	default:
		putError("Unexpected token %v in parsePrimaryExpression.\n", tok.sval)
	}

	return nil
}

func parseConstant() Ast {
	tok := lookahead(1)
	switch {
	case tok == nil:
		putError("tok is nil\n")
	case tok.isTypeInt():
		ival, _ := strconv.Atoi(tok.sval)
		nextToken()
		return &AstConstant{
			constant: &IntegerConstant{
				ival: ival,
			},
		}
	case tok.isTypeRune():
		rarr := []rune(tok.sval)
		nextToken()
		return &AstConstant{
			constant: &RuneConstant{
				rval: rarr[0],
			},
		}
	case tok.isTypeString():
		nextToken()
		ast := &AstString{
			sval:   tok.sval,
			slabel: fmt.Sprintf("L%d", stringIndex),
		}
		stringIndex++
		stringList = append(stringList, ast)
		return ast
	default:
		putError("unknown token %v in parseConstant\n", tok)
		debugToken(tok)
	}
	return nil
}

func parseIdentifier() *Identifier {
	tok := lookahead(1)
	switch {
	case tok == nil:
		putError("tok is nil\n")
		debugToken(tok)
	case tok.isTypeIdentifier():
		nextToken()
		sym := findSymbol(tok.sval)
		if sym == nil {
			sym = makeSymbol(tok.sval, "int")
		}
		return &Identifier{
			symbol: sym,
		}
	default:
		putError("Unexpected token %v in parseSymbol.\n", tok.sval)
		debugToken(tok)
	}
	return nil
}

func parseIdentifierOrFuncall() Ast {
	tok := lookahead(1)
	name := tok.sval
	nextToken()
	tok = lookahead(1)
	switch {
	case tok == nil:
		return nil
	case tok.isPunct("("):
		consumeToken("(")
		args := parseArgumentList()
		consumeToken(")")
		return &FunCall{
			fname: name,
			args:  args,
		}
	}

	fmt.Println("TBD")
	return nil
}

func parseArgumentList() []Ast {
	var r []Ast
	for {
		tok := lookahead(1)
		if tok == nil {
			putError("Expected ), but got EOF")
			return r
		}
		if tok.isPunct(")") {
			return r
		}
		arg := parseExpression()
		r = append(r, arg)
		tok = lookahead(1)
		switch {
		case tok == nil:
			putError("Expected ), but got EOF")
		case tok.isPunct(")"):
			return r
		case tok.isPunct(","):
			consumeToken(",")
			continue
		default:
			putError("Unexpected token %s in parseArgumentList.\n", tok.sval)
			debugToken(tok)
		}
	}
}
