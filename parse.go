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
	if tok.isEOF() {
		return nil
	}
	packname := parsePackageDeclaration()
	packages := parseImport()
	var childs []Ast

	for {
		tok = lookahead(1)

		switch {
		case tok.isEOF():
			return &TranslationUnit{
				packname: packname,
				packages: packages,
				childs:   childs,
			}
		case tok.isReserved("func"):
			ast := parseFunctionDefinition()
			childs = append(childs, ast)
		case tok.isReserved("var"):
			ast := parseDeclarationStatement()
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
	case tok.isEOF():
		putError("No package declaration.")
	case tok.isReserved("package"):
		consumeToken("package")
		tok = lookahead(1)
		if tok.isTypeIdentifier() {
			packname = tok.sval
			nextToken()
		} else {
			putError("%s is not allowed to be package name.", tok.sval)
		}
	default:
		putError("No package declaration.")
	}
	consumeToken(";")

	return packname
}

func parseImport() []string {
	var packages []string
	for {
		tok := lookahead(1)

		switch {
		case tok.isEOF():
			return packages
		case tok.isReserved("import"):
			consumeToken("import")
			ps := parseImportPackageNames()
			packages = append(packages, ps...)
			consumeToken(";")
		default:
			return packages
		}
	}
	return packages
}

func parseImportPackageNames() []string {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
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
		case tok.isEOF():
			putError("Expected ), but got EOF.")
			return packages
		case tok.isTypeString():
			packages = append(packages, tok.sval)
			nextToken()
			consumeToken(";")
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
	// argument ?
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
	beginSymbolBlock()
	for {
		tok := lookahead(1)
		switch {
		case tok.isPunct("}"):
			consumeToken("}")
			consumeToken(";")
			localvars := endSymbolBlock()
			return &CompoundStatement{
				statements: statements,
				localvars : localvars,
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
	case tok.isReserved("var"):
		ast = parseDeclarationStatement()
	default:
		ast = parseExpression()
		consumeToken(";")
	}

	return &Statement{
		ast: ast,
	}
}

func parseDeclarationStatement() Ast {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
		return nil
	case tok.isReserved("var"):
		consumeToken("var")
		tok2 := lookahead(1)
		if !tok2.isTypeIdentifier() {
			putError("Expected identifier, but got %s.", tok2.sval)
			return nil
		}
		nextToken()

		tok3 := lookahead(1)
		if !tok3.isTypeIdentifier() {
			putError("Expected type, but got %s.", tok3.sval)
			return nil
		}
		sym := makeSymbol(tok2.sval, tok3.sval)
		nextToken()

		tok4 := lookahead(1)
		if tok4.isPunct("=") {
			id := &Identifier {
				symbol: sym,
			}
			ast := parseAssignmentExpressionRightHand(id)
			consumeToken(";")
			return &DeclarationStatement{
				sym   : sym,
				assign: ast,
			}
		}
		consumeToken(";")
		return &DeclarationStatement{
			sym  : sym,
			assign: nil,
		}
	default:
		putError("Expected var, but got %s.", tok.sval)
		return nil
	}
	return nil
}

func parseExpression() Ast {
	ast := parseAssignmentExpression()
	return ast
}

func parseAssignmentExpression() Ast {
	var ast Ast = parseAdditiveExpression()
	return parseAssignmentExpressionRightHand(ast)
}

func parseAssignmentExpressionRightHand(ast Ast) Ast {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
		return ast
	case tok.isPunct("="):
		consumeToken("=")
		var right Ast = parseAdditiveExpression()
		left, ok := ast.(LeftValue)
		if !ok {
			putError("fatal: cannot cast %T.", ast)
			panic("internal error")
		}
		return &AssignmentExpression {
			left : left,
			right: right,
		}
	case tok.isSemicolon():
		return ast
	case tok.isPunct(")"), tok.isPunct("}"):
		return ast
	default:
		return ast
	}
	return ast
}

func parseAdditiveExpression() Ast {
	var ast Ast = parseMultiplicativeExpression()
	for {
		tok := lookahead(1)
		switch {
		case tok.isEOF():
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
		case tok.isSemicolon():
			return ast
		case tok.isPunct("="):
			return ast
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
		case tok.isEOF():
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
		case tok.isSemicolon():
			return ast
		case tok.isPunct("+"), tok.isPunct("-"), tok.isPunct("="):
			return ast
		default:
			return ast
		}
	}
	return ast
}

func parseUnaryExpression() Ast {
	tok := lookahead(1)
	var ast Ast

	switch {
	case tok.isEOF():
		return nil
	case tok.isTypeString(), tok.isTypeIdentifier(), tok.isTypeInt(), tok.isTypeRune():
		ast = parsePrimaryExpression()
		/*
		return &UnaryExpression{
			operand: ast,
		}
		*/
		return ast
	default:
		putError("Unexpected token %v in parseUnaryExpression.\n", tok.sval)
	}

	return nil
}

func parsePrimaryExpression() Ast {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
		return nil
	case tok.isTypeInt(), tok.isTypeRune(), tok.isTypeString():
		ast := parseConstant()
		return &PrimaryExpression{
			child: ast,
		}
	case tok.isTypeIdentifier(), tok.isTypeReserved():
		ast := parseIdentifierOrFuncall()
		return ast
	default:
		putError("Unexpected token %v in parsePrimaryExpression.\n", tok.sval)
	}

	return nil
}

func parseConstant() Ast {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
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
		return getAstString(tok.sval)
	default:
		putError("unknown token %v in parseConstant\n", tok)
		debugToken(tok)
	}
	return nil
}

func parseIdentifier() *Identifier {
	tok := lookahead(1)
	switch {
	case tok.isEOF():
		putError("tok is nil\n")
		debugToken(tok)
	case tok.isTypeIdentifier():
		nextToken()
		sym := findSymbol(tok.sval)
		
		if sym == nil {
			// sym = makeSymbol(tok.sval)
			putError("Undefined variable %s.\n", tok.sval)
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
	if isDeclaredSymbol(name) {
		return parseIdentifier()
	}
	nextToken()
	tok = lookahead(1)
	switch {
	case tok.isEOF():
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
		if tok.isEOF() {
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
		case tok.isEOF():
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


/* ================================================================ */

func getAstString(sval string) *AstString {
	for i:=0; i<stringIndex; i++ {
		if sval==stringList[i].sval {
			// This is probably preferable 
			// because ast does not become a tree structure
			return stringList[i]
		}
	}
	ast := &AstString{
		sval   : sval,
		slabel : fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	stringList = append(stringList, ast)
	return ast
}
