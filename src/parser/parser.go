package parser

import (
	"fmt"
	"nice-expr/src/ast"
	"nice-expr/src/token"
	TT "nice-expr/src/token/tokentype"
	"nice-expr/src/util"

	"golang.org/x/exp/slices"
)

type NiceExprParser struct {
	Tokens     []token.Token
	lastTokens *util.Stack[token.Token]
}

func NewNiceExprParser(tokens []token.Token) NiceExprParser {
	stack := new(util.Stack[token.Token])
	stack.Push(token.Token{Tt: TT.EOF, Lexeme: "EOF", CodePos: -1, Line: -1, Start: -1, End: -1})
	return NiceExprParser{tokens, stack}
}

func (p NiceExprParser) LastSeen() *token.Token {
	return p.lastTokens.Peek()
}

func (p NiceExprParser) hasMore() bool {
	return len(p.Tokens) > 0
}

// consume and return the next token in the token queue.
func (p *NiceExprParser) getNextToken() (*token.Token, *ParseError) {
	if !p.hasMore() {
		return nil, NewParseError("out of tokens", nil, "getNextToken")
	}
	tok := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	p.lastTokens.Push(tok)
	return &tok, nil
}

// peek at the front of the token queue.
func (p *NiceExprParser) peekToken() (*token.Token, *ParseError) {
	if !p.hasMore() {
		return nil, NewParseError("out of tokens", nil, "peekToken")
	}
	return &(p.Tokens[0]), nil
}

// put the last-consumed token back onto the front of the token queue.
func (p *NiceExprParser) putBackToken() {
	tok, err := p.lastTokens.Pop()
	if err != nil {
		return
	}
	p.Tokens = append([]token.Token{tok}, p.Tokens...)
}

// consume a token and determine if it is of a desired token type.
func (p *NiceExprParser) expectToken(tokType TT.TokenType) (*token.Token, *ParseError) {
	token, err := p.getNextToken()
	if err != nil {
		return token, err.addRule("expectToken")
	}
	if token.Tt != tokType {
		return nil, NewParseError(fmt.Sprintf("expected token `%v`", tokType), token, "expectToken")
	}
	return token, nil
}

// consume a token and determine if it is one of a desired token type.
func (p *NiceExprParser) expectAny(tokTypes []TT.TokenType) (*token.Token, *ParseError) {
	token, err := p.getNextToken()
	if err != nil {
		return nil, err.addRule("expectAny")
	}
	for _, tokType := range tokTypes {
		if token.Is(tokType) {
			return token, nil
		}
	}
	return nil, NewParseError(fmt.Sprintf("expected one of %v", tokTypes), token, "expectAny")
}

// peek at a token and determine if it is of a desired token type.
func (p *NiceExprParser) checkToken(tokType TT.TokenType) (bool, *ParseError) {
	token, err := p.peekToken()
	if err != nil {
		return false, err.addRule("checkToken")
	}
	if token.Tt != tokType {
		return false, nil
	}
	return true, nil
}

// peek at a token and determine if it is of a desired token type.
func (p *NiceExprParser) checkAny(tokTypes []TT.TokenType) (bool, *ParseError) {
	token, err := p.peekToken()
	if err != nil {
		return false, err.addRule("checkToken")
	}
	for _, tokType := range tokTypes {
		if token.Is(tokType) {
			return true, nil
		}
	}
	return false, nil
}

// peek at a token and determine if it is of a desired token type.
func (p *NiceExprParser) optionalToken(tokType TT.TokenType) (bool, *ParseError) {
	token, err := p.peekToken()
	if err != nil {
		return false, err.addRule("optionalToken")
	}
	if token.Tt != tokType {
		return false, nil
	}
	return true, nil
}

// peek at a token and determine if it is of one of a desired token type.
func (p *NiceExprParser) optionalAny(tokTypes []TT.TokenType) (bool, *ParseError) {
	token, err := p.peekToken()
	if err != nil {
		return false, err.addRule("optionalAny")
	}
	for _, tokType := range tokTypes {
		if token.Is(tokType) {
			return true, nil
		}
	}
	return false, nil
}

func (p *NiceExprParser) Program() (*ast.Program, *ParseError) {
	var program = new(ast.Program)
	for p.hasMore() {
		if isSc, err := p.optionalToken(TT.Semicolon); err != nil {
			return program, err.addRule("Program.EarlySemicolon")
		} else if isSc {
			// skip empty statements
			p.getNextToken()
			continue
		}
		stmt, err := p.Statement()
		if err != nil {
			return program, err.addRule("Program")
		} else if stmt == nil {
			continue
		}
		program.Statements = append(program.Statements, stmt)
	}
	return program, nil
}

func (p *NiceExprParser) Statement() (ast.Expr, *ParseError) {
	var expr ast.Expr
	tok, err := p.peekToken()
	if err != nil {
		return nil, err.addRule("Statement.Return?")
	} else if tok.Is(TT.Return) {
		p.getNextToken()
		ret, err := p.Return()
		if err != nil {
			return ret, err
		}
		expr = ret
	} else {
		expr, err = p.Expr()
		if err != nil {
			return expr, err.addRule("Statement")
		} else if expr == nil {
			return nil, nil
		}
	}
	_, err = p.expectToken(TT.Semicolon)
	if err != nil {
		return expr, err.addRule("Statement.Semicolon")
	}
	return expr, nil
}

func (p *NiceExprParser) Expr() (ast.Expr, *ParseError) {
	var expr ast.Expr

	switch tok, err := p.peekToken(); {
	case err != nil:
		return expr, err.addRule("Expr")
	case tok.Is(TT.LeftParen): // nested
		p.getNextToken()
		expr, err = p.Expr()
		if err != nil {
			return expr, err.addRule("Expr.ParenExpr")
		}
		if _, err := p.expectToken(TT.RightParen); err != nil {
			return expr, err.addRule("Expr.ParenExprEnd")
		}
		return expr, nil
	case tok.Is(TT.LeftBrace): // block
		return p.Block()
	case tok.Is(TT.If):
		return p.If()
	case slices.Contains(TT.VarConstSet, tok.Tt): // decl or assignment
		if expr, err = p.AssOrDecl(); err != nil {
			return expr, err.addRule("Expr.AssOrDecl")
		}
		return expr, nil
	default:
		expr, err = p.Test()
		if err != nil {
			return expr, err.addRule("Expr.Test")
		}
		if hasUnderscore, err := p.optionalToken(TT.Underscore); err != nil {
			return expr, err.addRule("Expr.Indexing?")
		} else if hasUnderscore {
			return p.Indexing(expr)
		} else {
			return expr, nil
		}
	}
}

func (p *NiceExprParser) Return() (*ast.Return, *ParseError) {
	ret := new(ast.Return)
	unary := new(ast.UnaryExpr)
	if tok, err := p.peekToken(); err != nil {
		return ret, err.addRule("Return.End?")
	} else if tok.Is(TT.Semicolon) {
		ret.UnaryExpr = *unary
		ret.UnaryExpr.Right = nil // no expr after return
		return ret, nil
	}
	// optional expr
	expr, err := p.Expr()
	if err != nil {
		return ret, err.addRule("Return.Expr")
	}
	unary.Right = expr
	ret.UnaryExpr = *unary
	return ret, nil
}

func (p *NiceExprParser) Block() (*ast.Block, *ParseError) {
	var block = new(ast.Block)
	p.getNextToken()
	for {
		if tok, err := p.peekToken(); err != nil {
			return block, err.addRule("Block")
		} else if tok.Is(TT.RightBrace) { // leaving block
			p.getNextToken()
			break
		}
		if isSc, err := p.optionalToken(TT.Semicolon); err != nil {
			return block, err.addRule("Block.EarlySemicolon")
		} else if isSc {
			// skip empty statements
			p.getNextToken()
			continue
		}
		stmt, err := p.Statement()
		if err != nil {
			return block, err.addRule("Block")
		} else if stmt == nil {
			continue
		}
		block.Statements = append(block.Statements, stmt)
		if tok, err := p.peekToken(); err != nil {
			return block, err.addRule("Block")
		} else if tok.Is(TT.RightBrace) { // leaving block
			p.getNextToken()
			break
		}
	}
	return block, nil
}
func (p *NiceExprParser) If() (*ast.If, *ParseError) {
	ifExpr := new(ast.If)
	if _, err := p.expectToken(TT.If); err != nil {
		return ifExpr, err.addRule("If.If")
	}
	condition, err := p.Expr()
	if err != nil {
		return ifExpr, err.addRule("If.Condition")
	}
	ifExpr.Condition = condition
	if _, err := p.expectToken(TT.Then); err != nil {
		return ifExpr, err.addRule("If.Then")
	}
	then, err := p.Block()
	if err != nil {
		return ifExpr, err.addRule("If.ThenBlock")
	}
	ifExpr.Then = then
	// optional else
	if hasElse, err := p.optionalToken(TT.Else); err != nil {
		return ifExpr, err.addRule("If.Else?")
	} else if !hasElse {
		// no else, just if-then
		return ifExpr, nil
	}
	// has an else, check if nested if
	p.getNextToken()
	if hasNestedIf, err := p.optionalToken(TT.If); err != nil {
		return ifExpr, err.addRule("If.NestedIf?")
	} else if hasNestedIf {
		nestedIf, err := p.If()
		if err != nil {
			return ifExpr, err.addRule("If.NestedIf")
		}
		ifExpr.ElseIf = nestedIf
	} else {
		// get else block
		elseBlock, err := p.Block()
		if err != nil {
			return ifExpr, err.addRule("If.Else")
		}
		ifExpr.Else = elseBlock
	}
	return ifExpr, nil
}

func (p *NiceExprParser) Indexing(left ast.Expr) (ast.Expr, *ParseError) {
	indexing := new(ast.Indexing)
	indexing.Left = left
	if _, err := p.expectToken(TT.Underscore); err != nil {
		return indexing, err.addRule("Expr.Indexing.Underscore")
	}
	if right, err := p.Expr(); err != nil {
		return indexing, err.addRule("Expr.Indexing.Right")
	} else {
		indexing.Right = right
		return indexing, nil
	}
}

// test ::= notTest | test ("and"|"or") test ;
func (p *NiceExprParser) Test() (ast.Expr, *ParseError) {
	var (
		test = new(ast.BinaryExpr)
		and  = new(ast.AndTest)
		or   = new(ast.OrTest)
		err  *ParseError
	)
	if test.Left, err = p.NotTest(); err != nil {
		return test, err.addRule("Test.NotTest")
	}
	if ok, err := p.optionalAny(TT.BinLogOps); err != nil {
		return test, err.addRule("Test.LogOp?")
	} else if !ok {
		// notTest only
		return test.Left, nil
	}
	// test "and" test
	tok, err := p.expectAny(TT.BinLogOps)
	if err != nil {
		return test, err.addRule("Test.TestOp")
	}
	if test.Right, err = p.Test(); err != nil {
		return test, err.addRule("Test.Test")
	}
	switch tok.Tt {
	case TT.And:
		and.BinaryExpr = *test
		return and, nil
	case TT.Or:
		or.BinaryExpr = *test
		return or, nil
	default:
		return nil, NewParseError("unknown test operator %v", tok, "Test.TestOp")
	}
}

// notTest ::= comparison | "not" notTest ;
func (p *NiceExprParser) NotTest() (ast.Expr, *ParseError) {
	notTest := new(ast.NotTest)
	var unary ast.UnaryExpr
	var ok bool
	var err *ParseError
	if ok, err = p.optionalToken(TT.Not); err != nil {
		return notTest, err.addRule("NotTest.Not?")
	} else if !ok {
		// just a comparison
		return p.Comparison()
	} else if _, err = p.expectToken(TT.Not); err != nil {
		return notTest, err.addRule("NotTest.Not")
	}
	// "not" notTest
	if unary.Right, err = p.NotTest(); err != nil {
		return notTest, err.addRule("NotTest.NotTest")
	}
	notTest.UnaryExpr = unary
	return notTest, nil
}

// comparison ::= addExpr (("<"|">"|"<="|">="|"=") comparison)* ;
func (p *NiceExprParser) Comparison() (ast.Expr, *ParseError) {
	var (
		equal        = new(ast.Equal)
		greater      = new(ast.Greater)
		less         = new(ast.Less)
		greaterEqual = new(ast.GreaterEqual)
		lessEqual    = new(ast.LessEqual)
		comparison   = new(ast.BinaryExpr)
		err          *ParseError
	)
	// required left addExpr
	if comparison.Left, err = p.AddExpr(); err != nil {
		return nil, err.addRule("Comparison.AddExprRightLeft")
	}
	if ok, err := p.optionalAny(TT.ComparisonOps); err != nil {
		return nil, err.addRule("Comparison.ComparisonOp?")
	} else if !ok {
		// only "addExpr"
		return comparison.Left, nil
	}
	tok, err := p.expectAny(TT.ComparisonOps)
	if err != nil {
		return nil, err.addRule("Comparison.ComparisonOp")
	}
	// optional right comparison
	if comparison.Right, err = p.Comparison(); err != nil {
		return nil, err.addRule("Comparison.AddExprRight")
	}
	switch tok.Tt {
	case TT.Equal:
		equal.BinaryExpr = *comparison
		return equal, nil
	case TT.Greater:
		greater.BinaryExpr = *comparison
		return greater, nil
	case TT.Less:
		less.BinaryExpr = *comparison
		return less, nil
	case TT.GreaterEqual:
		greaterEqual.BinaryExpr = *comparison
		return greaterEqual, nil
	case TT.LessEqual:
		lessEqual.BinaryExpr = *comparison
		return lessEqual, nil
	default:
		return nil, NewParseError("unknown comparison operator: %v", tok, "Comparison.ComparisonOp")
	}
}

// addExpr ::= mulExpr | addExpr ("+"|"-") addExpr ;
func (p *NiceExprParser) AddExpr() (ast.Expr, *ParseError) {
	var (
		add     = new(ast.Add)
		sub     = new(ast.Sub)
		addExpr = new(ast.BinaryExpr)
		err     *ParseError
	)
	if addExpr.Left, err = p.MulExpr(); err != nil {
		return addExpr, err.addRule("AddExpr.MulExpr")
	}
	if ok, err := p.optionalAny(TT.AddOps); err != nil {
		return addExpr, err.addRule("AddExpr.AddOp?")
	} else if !ok {
		// only "mulExpr"
		return addExpr.Left, nil
	}
	tok, err := p.expectAny(TT.AddOps)
	if err != nil {
		return addExpr, err.addRule("AddExpr.AddOp")
	}
	if addExpr.Right, err = p.AddExpr(); err != nil {
		return addExpr, err.addRule("AddExpr.AddExprRight")
	}
	switch tok.Tt {
	case TT.Plus:
		add.BinaryExpr = *addExpr
		return add, nil
	case TT.Minus:
		sub.BinaryExpr = *addExpr
		return sub, nil
	default:
		return nil, NewParseError("unknown addition op: %v", tok, "AddExpr.AddOp")
	}
}

// mulExpr ::= unaryMinusExpr | mulExpr ("*"|"/"|"%") mulExpr ;
func (p *NiceExprParser) MulExpr() (ast.Expr, *ParseError) {
	var (
		mul     = new(ast.Mul)
		div     = new(ast.Div)
		mod     = new(ast.Mod)
		mulExpr = new(ast.BinaryExpr)
		err     *ParseError
	)
	if mulExpr.Left, err = p.UnaryMinusExpr(); err != nil {
		return nil, err.addRule("MulExpr.UnaryMinusExpr")
	}
	if ok, err := p.optionalAny(TT.MulOps); err != nil {
		return nil, err.addRule("MulExpr.MulOp?")
	} else if !ok {
		// only "unaryMinusExpr"
		return mulExpr.Left, nil
	}
	tok, err := p.expectAny(TT.MulOps)
	if err != nil {
		return nil, err.addRule("MulExpr.MulOp")
	}
	if mulExpr.Right, err = p.MulExpr(); err != nil {
		return mulExpr, err.addRule("MulExpr.MulExprRight")
	}
	switch tok.Tt {
	case TT.Star:
		mul.BinaryExpr = *mulExpr
		return mul, nil
	case TT.Slash:
		div.BinaryExpr = *mulExpr
		return div, nil
	case TT.Percent:
		mod.BinaryExpr = *mulExpr
		return mod, nil
	default:
		return nil, NewParseError("unknown multiplication op: %v", tok, "MulExpr.Op")
	}
}

// unaryMinusExpr ::= "-" unaryMinusExpr | expr ;
func (p *NiceExprParser) UnaryMinusExpr() (ast.Expr, *ParseError) {
	unaryMinusExpr := new(ast.UnaryMinus)
	if ok, err := p.optionalToken(TT.Minus); err != nil {
		return unaryMinusExpr, err.addRule("UnaryMinusExpr.Minus?")
	} else if ok {
		// consume minus
		p.getNextToken()
		// have a minus, get another unary minus expr
		if unaryMinusExpr.Right, err = p.UnaryMinusExpr(); err != nil {
			// "-" unaryMinusExp
			return unaryMinusExpr, err.addRule("UnaryMinusExpr.UnaryMinusExpr")
		}
		return unaryMinusExpr, nil
	}
	// no minus, just an expr
	if tok, err := p.peekToken(); err != nil {
		return unaryMinusExpr, err.addRule("UnaryMinusExpr.Primary?")
	} else if slices.Contains(TT.Primaries, tok.Tt) { // primary
		return p.Primary()
	} else {
		return nil, NewParseError("something in unary minus", tok, "UnaryMinusExpr")
	}
}

func (p *NiceExprParser) Primary() (ast.Expr, *ParseError) {
	var expr ast.Expr
	switch tok, err := p.peekToken(); {
	case err != nil:
		return nil, err.addRule("Primary")
	case slices.Contains(TT.Literals, tok.Tt): // primary
		expr, err = p.Literal()
		if err != nil {
			return expr, err.addRule("Primary.Literal")
		}
		return expr, err
	case tok.Is(TT.Identifier): // primary
		expr, err = p.IdentifierOrFuncCall()
		if err != nil {
			return expr, err.addRule("Primary.IdentOrFuncCall")
		}
		return expr, err
	default:
		return nil, NewParseError("unkown primary", tok, "Primary")
	}
}

func (p *NiceExprParser) AssOrDecl() (ast.Expr, *ParseError) {
	switch tok, err := p.peekToken(); {
	case tok.Is(TT.Var):
		return p.VariableDeclaration()
	case tok.Is(TT.Const):
		return p.ConstantDeclaration()
	case tok.Is(TT.Set):
		return p.Assignment()
	case err != nil:
		return nil, err.addRule("AssOrDecl.Start")
	default:
		return nil, NewParseError("unknown leading token", tok, "AssOrDecl.Start")
	}
}

func (p *NiceExprParser) IdentifierOrFuncCall() (ast.Expr, *ParseError) {
	if _, err := p.expectToken(TT.Identifier); err != nil {
		return nil, err.addRule("IdentifierOrFuncCall.Ident")
	}
	if ok, err := p.optionalToken(TT.LeftParen); err != nil {
		return nil, err.addRule("IdentifierOrFuncCall.FuncCall?")
	} else if ok {
		p.putBackToken() // put back the ident for func call
		return p.FunctionCall()
	} else {
		p.putBackToken() // put back the ident for ident
		return p.Identifier()
	}
}

func (p *NiceExprParser) VariableDeclaration() (*ast.VariableDeclaration, *ParseError) {
	expr := new(ast.VariableDeclaration)
	if _, err := p.expectToken(TT.Var); err != nil {
		return expr, err.addRule("VariableDeclaration.Var")
	}
	name, err := p.Identifier()
	if err != nil {
		return expr, err.addRule("VariableDeclaration")
	}
	expr.Name = name
	if _, err = p.expectToken(TT.Is); err != nil {
		return expr, err.addRule("VariableDeclaration.Is")
	}
	typeExpr, err := p.Type()
	if err != nil {
		return expr, err.addRule("VariableDeclaration.Type")
	}
	expr.Type = typeExpr
	value, err := p.Expr()
	if err != nil {
		return expr, err.addRule("VariableDeclaration.Value")
	}
	expr.Value = value
	return expr, nil
}

func (p *NiceExprParser) ConstantDeclaration() (*ast.ConstantDeclaration, *ParseError) {
	expr := new(ast.ConstantDeclaration)
	if _, err := p.expectToken(TT.Const); err != nil {
		return expr, err.addRule("ConstantDeclaration.Const")
	}
	name, err := p.Identifier()
	if err != nil {
		return expr, err.addRule("ConstantDeclaration")
	}
	expr.Name = name
	if _, err = p.expectToken(TT.Is); err != nil {
		return expr, err.addRule("ConstantDeclaration.Is")
	}
	typeExpr, err := p.Type()
	if err != nil {
		return expr, err.addRule("ConstantDeclaration.Type")
	}
	expr.Type = typeExpr
	value, err := p.Expr()
	if err != nil {
		return expr, err.addRule("ConstantDeclaration.Value")
	}
	expr.Value = value
	return expr, nil
}

func (p *NiceExprParser) Assignment() (*ast.Assignment, *ParseError) {
	ae := new(ast.Assignment)

	if _, err := p.expectToken(TT.Set); err != nil {
		return ae, err.addRule("Assignment.Set")
	}
	name, err := p.Identifier()
	if err != nil {
		return ae, err.addRule("Assignment.Identifier")
	}
	ae.Name = name
	op, err := p.expectAny(TT.AssignmentOperations)
	if err != nil {
		return ae, err.addRule("Assignment.AssignmentOp")
	}
	ae.Op = op
	value, err := p.Expr()
	if err != nil {
		return ae, err.addRule("Assignment.Value")
	}
	ae.Value = value

	return p.DesugarAssignment(ae), nil
}

// convert an assignment from:
//
//	set name += val
//
// to:
//
//	set name is name + val
func (p *NiceExprParser) DesugarAssignment(a *ast.Assignment) *ast.Assignment {
	var (
		ass     = new(ast.Assignment)
		binExpr = ast.BinaryExpr{Left: a.Name, Right: a.Value}
		add     ast.Add
		sub     ast.Sub
		mul     ast.Mul
		div     ast.Div
		mod     ast.Mod
	)

	ass.Name = a.Name
	ass.Op = a.Op
	ass.Op.Tt, ass.Op.Lexeme = TT.Is, "is"

	switch a.Op.Tt {
	case TT.PlusEqual:
		add = ast.Add{BinaryExpr: binExpr}
		ass.Value = add
	case TT.MinusEqual:
		sub = ast.Sub{BinaryExpr: binExpr}
		ass.Value = sub
	case TT.StarEqual:
		mul = ast.Mul{BinaryExpr: binExpr}
		ass.Value = mul
	case TT.SlashEqual:
		div = ast.Div{BinaryExpr: binExpr}
		ass.Value = div
	case TT.PercentEqual:
		mod = ast.Mod{BinaryExpr: binExpr}
		ass.Value = mod
	default:
		return a
	}
	return ass
}

func (p *NiceExprParser) Identifier() (*ast.Identifier, *ParseError) {
	id := new(ast.Identifier)
	token, err := p.expectToken(TT.Identifier)
	if err != nil {
		return id, err.addRule("Identifier.CheckIdentifier")
	} else if token == nil {
		return id, NewParseError("invalid token for literal", token, "Identifier.CheckIdentifier")
	}
	id.Tok = token
	return id, nil
}

func (p *NiceExprParser) FunctionCall() (*ast.FunctionCall, *ParseError) {
	funcCall := new(ast.FunctionCall)
	ident, err := p.Identifier()
	if err != nil {
		return funcCall, err.addRule("FunctionCall")
	}
	if _, err := p.expectToken(TT.LeftParen); err != nil {
		return funcCall, err.addRule("FunctionCall.ArgsStart")
	}
	funcCall.Ident = ident
	funcCall.Arguments, err = p.ExprList(TT.RightParen)
	if err != nil {
		return funcCall, err.addRule("FuncCall.Arguments")
	}
	if _, err := p.expectToken(TT.RightParen); err != nil {
		return funcCall, err.addRule("FunctionCall.ArgsEnd")
	}
	return funcCall, nil
}

func (p *NiceExprParser) ExprList(endingToken TT.TokenType) ([]ast.Expr, *ParseError) {
	l := []ast.Expr{}
	for {
		if ok, err := p.optionalToken(endingToken); err != nil {
			return l, err.addRule("ExprList.Ending")
		} else if ok {
			break
		}
		expr, err := p.Expr()
		if err != nil {
			return l, err.addRule("ExprList.Expr")
		}

		// optional trailing comma
		if ok, err := p.optionalToken(endingToken); err != nil {
			return l, err.addRule("ExprList.Ending")
		} else if ok {
			l = append(l, expr)
			break
		}

		if _, err = p.expectToken(TT.Comma); err != nil {
			return l, err.addRule("ExprList.Comma")
		}
		l = append(l, expr)
	}
	return l, nil
}

func (p *NiceExprParser) Literal() (ast.Literal, *ParseError) {
	ok, err := p.checkAny(TT.CompositeLiteralStarts)
	if err != nil {
		return nil, err.addRule("Literal.Start")
	} else if ok {
		return p.CompoundLiteral()
	}
	return p.PrimitiveLiteral()
}

func (p *NiceExprParser) PrimitiveLiteral() (*ast.PrimitiveLiteral, *ParseError) {
	pe := new(ast.PrimitiveLiteral)
	token, err := p.expectAny(TT.PrimitiveLiterals)
	if err != nil {
		return pe, err.addRule("PrimitiveLiteral.CheckLiteral")
	} else if token == nil {
		return pe, NewParseError("invalid token for literal", token, "Literal.CheckLiteral")
	}
	pe.Token = token
	return pe, nil
}

func (p *NiceExprParser) CompoundLiteral() (ast.CompoundLiteral, *ParseError) {
	ok, err := p.optionalToken(TT.LeftBracket)
	if err != nil {
		return nil, err.addRule("CompoundLiteral.List?")
	} else if ok {
		return p.ListLiteral()
	}
	ok, err = p.optionalToken(TT.LeftTriangle)
	if err != nil {
		return nil, err.addRule("CompoundLiteral.Map?")
	} else if ok {
		return p.MapLiteral()
	}
	return nil, NewParseError("unknown compound literal start", nil, "CompoundLiteral")
}

func (p *NiceExprParser) ListLiteral() (*ast.ListLiteral, *ParseError) {
	l := new(ast.ListLiteral)
	var err *ParseError
	_, err = p.expectToken(TT.LeftBracket)
	if err != nil {
		return nil, err.addRule("ListLiteral.Start")
	}
	l.Values, err = p.ListElements()
	if err != nil {
		return l, err.addRule("ListLiteral.ListElements")
	}
	if _, err = p.expectToken(TT.RightBracket); err != nil {
		return l, err.addRule("ListLiteral.End")
	}
	return l, nil
}

func (p *NiceExprParser) ListElements() ([]ast.Expr, *ParseError) {
	l := []ast.Expr{}
	for {
		// list items are comma.separated, and trailing comma is optional
		if ok, err := p.checkToken(TT.RightBracket); err != nil {
			return l, err.addRule("ListElements.End")
		} else if ok {
			break
		}
		value, err := p.Expr()
		if err != nil {
			return l, err.addRule("ListElements.Elements")
		}
		if ok, err := p.checkToken(TT.RightBracket); err != nil {
			return l, err.addRule("ListElements.End")
		} else if ok {
			l = append(l, value)
			break
		}
		if _, err = p.expectToken(TT.Comma); err != nil {
			return l, err.addRule("ListElements.Elements")
		}
		l = append(l, value)
	}
	return l, nil
}

func (p *NiceExprParser) MapLiteral() (*ast.MapLiteral, *ParseError) {
	m := new(ast.MapLiteral)
	var err *ParseError
	if _, err := p.expectToken(TT.LeftTriangle); err != nil {
		return nil, err.addRule("MapLiteral")
	}
	m.Values, err = p.MapEntries()
	if err != nil {
		return m, err.addRule("MapLiteral.MapEntries")
	}
	if _, err = p.expectToken(TT.RightTriangle); err != nil {
		return m, err
	}
	return m, nil
}

func (p *NiceExprParser) MapEntries() (map[ast.Expr]ast.Expr, *ParseError) {
	m := make(map[ast.Expr]ast.Expr)
	for {
		// list items are comma.separated, and trailing comma is optional
		if ok, err := p.checkToken(TT.RightTriangle); err != nil {
			return m, err.addRule("MapEntry.End")
		} else if ok {
			break
		}
		key, err := p.Expr()
		if err != nil {
			return m, err.addRule("MapEntry.Key")
		}
		_, err = p.expectToken(TT.Colon)
		if err != nil {
			return m, err.addRule("MapEntry.Colon")
		}
		value, err := p.Expr()
		if err != nil {
			return m, err.addRule("MapEntry.Value")
		}
		if ok, err := p.checkToken(TT.RightTriangle); err != nil {
			return m, err.addRule("MapEntry.End")
		} else if ok {
			m[key] = value
			break
		}
		_, err = p.expectToken(TT.Comma)
		if err != nil {
			return m, err.addRule("MapEntry.Comma")
		}
		m[key] = value
	}
	return m, nil
}

func (p *NiceExprParser) Type() (ast.Type, *ParseError) {
	ok, err := p.checkAny(TT.PrimitiveTypes)
	if err != nil {
		return nil, err.addRule("Type")
	} else if ok {
		return p.PrimitiveType()
	}
	ok, err = p.checkAny(TT.CompoundTypes)
	if err != nil {
		return nil, err.addRule("Type")
	} else if ok && err == nil {
		return p.CompoundType()
	}
	ok, err = p.checkToken(TT.Func)
	if err != nil {
		return nil, err.addRule("Type")
	} else if ok {
		return p.FuncType()
	}
	return nil, NewParseError("type not found", nil, "Type")
}

func (p *NiceExprParser) PrimitiveType() (*ast.PrimitiveType, *ParseError) {
	pt := new(ast.PrimitiveType)
	name, err := p.expectAny(TT.PrimitiveTypes)
	if err != nil {
		return nil, err.addRule("PrimitiveType")
	}
	if name == nil {
		return nil, NewParseError("nil typename", name, "PrimitiveType")
	}
	pt.Name = name
	return pt, nil
}

func (p *NiceExprParser) CompoundType() (ast.Type, *ParseError) {
	ok, err := p.optionalToken(TT.List)
	if err != nil {
		return nil, err.addRule("CompoundType.List")
	} else if ok && err == nil {
		return p.ListType()
	}
	ok, err = p.optionalToken(TT.Map)
	if err != nil {
		return nil, err.addRule("CompoundType.Map")
	} else if ok && err == nil {
		return p.MapType()
	}
	return nil, NewParseError("expected `list` or `map`", nil, "CompoundType")
}

func (p *NiceExprParser) ListType() (*ast.ListType, *ParseError) {
	lte := new(ast.ListType)
	_, err := p.expectToken(TT.List)
	if err != nil {
		return lte, err.addRule("ListType.List")
	}
	_, err = p.expectToken(TT.LeftBracket)
	if err != nil {
		return lte, err.addRule("ListType.Start")
	}
	valueType, err := p.Type()
	if err != nil {
		return lte, err.addRule("ListType.Value")
	}
	_, err = p.expectToken(TT.RightBracket)
	if err != nil {
		return lte, err.addRule("ListType.End")
	}
	lte.ValueType = valueType
	return lte, nil
}

func (p *NiceExprParser) MapType() (*ast.MapType, *ParseError) {
	mte := new(ast.MapType)
	_, err := p.expectToken(TT.Map)
	if err != nil {
		return mte, err.addRule("MapType.Map")
	}
	_, err = p.expectToken(TT.LeftBracket)
	if err != nil {
		return mte, err.addRule("MapType.KeyStart")
	}
	keyType, err := p.Type()
	if err != nil {
		return mte, err.addRule("MapType.Key")
	}
	_, err = p.expectToken(TT.RightBracket)
	if err != nil {
		return mte, err.addRule("MapType.KeyEnd")
	}
	valueType, err := p.Type()
	if err != nil {
		return mte, err.addRule("MapType.Value")
	}
	mte.KeyType = keyType
	mte.ValueType = valueType
	return mte, nil
}

func (p *NiceExprParser) FuncType() (*ast.FuncType, *ParseError) {
	fte := new(ast.FuncType)
	_, err := p.expectToken(TT.Func)
	if err != nil {
		return fte, nil
	}
	if _, err = p.expectToken(TT.LeftParen); err != nil {
		return fte, err.addRule("FuncType.InputTypesStart")
	}
	if ok, err := p.checkToken(TT.RightParen); err != nil {
		return fte, err.addRule("FuncType.InputTypesEnd?")
	} else if ok {
		if _, err := p.expectToken(TT.RightParen); err != nil {
			return fte, err.addRule("FuncType.InputTypesEnd")
		}
	}
	fte.InputTypes, err = p.TypeList()
	if err != nil {
		return fte, err.addRule("FuncType.InputTypes")
	}
	// output type, optional
	if ok, err := p.checkAny(TT.Types); err != nil {
		return fte, err.addRule("FuncType.CheckOutputType")
	} else if ok {
		out, err := p.Type()
		if err != nil {
			return fte, err.addRule("FuncType.OutputType")
		}
		fte.OutputType = out
	}
	return fte, nil
}

func (p *NiceExprParser) TypeList() ([]ast.Type, *ParseError) {
	types := []ast.Type{}

	for {
		if ok, err := p.checkAny(TT.Types); err != nil {
			return types, err.addRule("TypeList.MoreTypes?")
		} else if !ok {
			break
		}
		t, err := p.Type()
		if err != nil {
			return types, err.addRule("TypeList.Type")
		}
		if _, err = p.expectToken(TT.Comma); err != nil {
			return types, err.addRule("TypeList.Comma")
		}
		types = append(types, t)
	}
	return types, nil
}
