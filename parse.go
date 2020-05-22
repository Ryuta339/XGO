package main

import (
	"fmt"
	"strconv"
)

var stringIndex = 0
var stringList []*AstString

func parse() Ast {
	currentScope = globalScope
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
				packname:   packname,
				packages:   packages,
				childs:     childs,
				globalvars: getGlobalSymList(),
			}
		case tok.isKeyword("func"):
			ast := parseFunctionDefinition()
			childs = append(childs, ast)
		case tok.isKeyword("var"):
			ast := parseGlobalDeclaration()
			childs = append(childs, ast)
		default:
			putError("func expected, but got %v.", tok.sval)
		}
	}
	return &TranslationUnit{
		packname:   packname,
		packages:   packages,
		childs:     childs,
		globalvars: getGlobalSymList(),
	}
}

func parsePackageDeclaration() string {
	tok := lookahead(1)
	packname := ""

	switch {
	case tok.isEOF():
		putError("No package declaration.")
	case tok.isKeyword("package"):
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
		case tok.isKeyword("import"):
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

func parseGlobalDeclaration() Ast {
	sym := parseDeclarationStatementCommon().(*GlobalVariable)
	if sym==nil {
		return &GlobalDeclaration{
			sym: nil,
		}
	}
	tok := lookahead(1)

	initstr := "0"
	if tok.isPunct("=") {
		consumeToken("=")
		tok2 := lookahead(1)
		initstr = tok2.sval
		nextToken()
	}
	consumeToken(";")

	var initval Constant
	switch sym.gtype {
	case "int":
		ival, _ := strconv.Atoi(initstr)
		initval = &IntegerConstant{
			ival: ival,
		}
	default:
		putError("Acceptable global variable is int, but got %s", sym.gtype)
	}
	sym.initval = initval
	return &GlobalDeclaration{
		sym: sym,
	}
}


func parseFunctionDefinition() Ast {
	tok := lookahead(1)
	if !tok.isKeyword("func") {
		putError("Expected func, but got %s", tok.typ)
		return nil
	}
	consumeToken("func")
	beginFunction()
	tok = lookahead(1)
	if !tok.isTypeIdentifier() {
		putError("Expected identifier, but got %s", tok.typ)
		return nil
	}
	nextToken()
	consumeToken("(")

	beginSymbolBlock()

PARSE_ARGUMENT_LIST_LOOP:
	for {
		tok := lookahead(1)
		switch {
		case tok.isPunct(")"):
			break PARSE_ARGUMENT_LIST_LOOP
		case tok.isTypeIdentifier():
			nextToken()
			tok2 := lookahead(1)
			makeSymbol(tok.sval, tok2.sval)
			nextToken()
			tok3 := lookahead(1)
			if tok3.isPunct(",") {
				consumeToken(",")
				continue
			} else if tok3.isPunct(")") {
				break PARSE_ARGUMENT_LIST_LOOP
			} else {
				putError("Expected ) or \",\", but got %s.", tok3.sval)
			}
		default:
			putError("Expected ) or identifier, but got %s.", tok.sval)
		}
	}

	consumeToken(")")
	// expect Type
	tok2 := lookahead(1)
	var rettype string
	if tok2.isTypeIdentifier() {
		rettype = tok2.sval
		nextToken()
	} else {
		rettype = "void"
	}
	tok3 := lookahead(1)
	if !tok3.isPunct("{") {
		putError("Expected {, but got %s", tok3.sval)
	}
	ast := parseCompoundStatement()
	params := endSymbolBlock()
	space := endFunction()
	return &FunctionDefinition{
		fname:   tok.sval,
		rettype: rettype,
		params:  params,
		ast:     ast,
		space:   space,
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
				localvars:  localvars,
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
	case tok.isKeyword("var"):
		ast = parseDeclarationStatement()
	default:
		ast = parseExpression()
		tok2 := lookahead(1)
		switch {
		case tok2.isPunct("="):
			ast = parseAssignmentExpressionRightHand(ast)
			consumeToken(";")
		default:
			consumeToken(";")
		}
	}

	return &Statement{
		ast: ast,
	}
}

func parseDeclarationStatementCommon() Symbol {

	tok := lookahead(1)
	switch {
	case tok.isEOF():
		return nil
	case tok.isKeyword("var"):
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
		return sym
	default:
		putError("Expected var, but got %s.", tok.sval)
		return nil
	}
	return nil
}

func parseDeclarationStatement() Ast {
	sym := parseDeclarationStatementCommon().(*LocalVariable)
	tok := lookahead(1)
	if tok.isPunct("=") {
		id := &Identifier{
			symbol: sym,
		}
		ast := parseAssignmentExpressionRightHand(id)
		consumeToken(";")
		return &DeclarationStatement{
			sym:    sym,
			assign: ast,
		}
	}
	consumeToken(";")
	return &DeclarationStatement{
		sym:    sym,
		assign: nil,
	}
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
		return &AssignmentExpression{
			left:  left,
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

func parseExpression() Ast {
	ast := parseAdditiveExpression()
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
		return ast
	default:
		putError("Unexpected token %v in parseUnaryExpression.", tok.sval)
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
	case tok.isTypeIdentifier(), tok.isTypeKeyword():
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
		sym := currentScope.findSymbol(tok.sval)

		if sym == nil {
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
	if currentScope.isDeclaredSymbol(name) {
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
	default:
		putError("Undeclared identifier %s.", name)
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
	for i := 0; i < stringIndex; i++ {
		if sval == stringList[i].sval {
			// This is probably preferable
			// because ast does not become a tree structure
			return stringList[i]
		}
	}
	ast := &AstString{
		sval:   sval,
		slabel: fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	stringList = append(stringList, ast)
	return ast
}
