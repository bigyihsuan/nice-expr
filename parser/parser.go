package parser

import (
	"fmt"
	"nice-expr/ast"
	"nice-expr/token"
	TT "nice-expr/token/tokentype"
)

type NiceExprParser struct {
	Tokens    []token.Token
	lastToken *token.Token
}

func NewNiceExprParser(tokens []token.Token) NiceExprParser {
	return NiceExprParser{tokens, &token.Token{Tt: TT.EOF, Lexeme: "EOF", CodePos: -1, Line: -1, Start: -1, End: -1}}
}

func (p NiceExprParser) hasMore() bool {
	return len(p.Tokens) > 0
}

// consume and return the next token in the token queue.
func (p *NiceExprParser) getNextToken(lastRule string) (*token.Token, *ParseError) {
	if !p.hasMore() {
		return nil, NewParseError("out of tokens", nil, lastRule)
	}
	tok := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	p.lastToken = &tok
	return &tok, nil
}

// peek at the front of the token queue.
func (p *NiceExprParser) peekToken(lastRule string) (*token.Token, *ParseError) {
	if !p.hasMore() {
		return nil, NewParseError("out of tokens", nil, lastRule)
	}
	return &(p.Tokens[0]), nil
}

// put the last-consumed token back onto the front of the token queue.
func (p *NiceExprParser) putBackToken() {
	p.Tokens = append([]token.Token{*p.lastToken}, p.Tokens...)
}

// consume a token and determine if it is of a desired token type.
func (p *NiceExprParser) expectToken(tokType TT.TokenType, lastRule string) (*token.Token, *ParseError) {
	token, err := p.getNextToken(lastRule)
	if err != nil {
		return token, err
	}
	if token.Tt != tokType {
		return nil, NewParseError(fmt.Sprintf("expected token `%v`", tokType), token, lastRule)
	}
	return token, nil
}

// consume a token and determine if it is one of a desired token type.
func (p *NiceExprParser) expectAny(tokTypes []TT.TokenType, lastRule string) (*token.Token, *ParseError) {
	token, err := p.getNextToken(lastRule)
	if err != nil {
		return nil, err
	}
	for _, tokType := range tokTypes {
		if token.Tt == tokType {
			return token, nil
		}
	}
	return nil, NewParseError(fmt.Sprintf("expected one of %v", tokTypes), token, lastRule)
}

// peek at a token and determine if it is of a desired token type.
func (p *NiceExprParser) checkToken(tokType TT.TokenType, lastRule string) (bool, *ParseError) {
	token, err := p.peekToken(lastRule)
	if err != nil {
		return false, err
	}
	if token.Tt != tokType {
		return false, NewParseError(fmt.Sprintf("expected token `%v`", tokType), token, lastRule)
	}
	return true, nil
}

// peek at a token and determine if it is of a desired token type.
func (p *NiceExprParser) optionalToken(tokType TT.TokenType, lastRule string) (bool, *ParseError) {
	token, err := p.peekToken(lastRule)
	if err != nil {
		return false, err
	}
	if token.Tt != tokType {
		return false, nil
	}
	return true, nil
}

// peek at a token and determine if it is of one of a desired token type.
func (p *NiceExprParser) checkAny(tokTypes []TT.TokenType, lastRule string) (bool, *ParseError) {
	token, err := p.peekToken(lastRule)
	if err != nil {
		return false, err
	}
	for _, tokType := range tokTypes {
		if token.Tt == tokType {
			return true, nil
		}
	}
	return false, nil
}

func (p *NiceExprParser) ParseProgram() (ast.Program, *ParseError) {
	var program ast.Program
	for p.hasMore() {
		stmt, err := p.ParseStatement()
		if err != nil {
			return program, err
		}
		program.Statements = append(program.Statements, stmt)
	}
	return program, nil
}

func (p *NiceExprParser) ParseStatement() (ast.Node, *ParseError) {
	node, err := p.ParseExpr()
	if err != nil {
		return node, err
	}
	_, err = p.expectToken(TT.Semicolon, "Statement-Semicolon")
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *NiceExprParser) ParseExpr() (ast.Node, *ParseError) {
	ok, err := p.optionalToken(TT.Var, "Expr-VarDecl")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseVariableDeclaration()
	}
	ok, err = p.optionalToken(TT.Const, "Expr-ConstDecl")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseConstantDeclaration()
	}
	ok, err = p.optionalToken(TT.Set, "Expr-Assignment")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseAssignment()
	}
	ok, err = p.checkAny(TT.PrimitiveLiterals, "Expr-PrimitiveLiterals")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParsePrimitiveLiteral()
	}
	ok, err = p.checkAny(TT.CompositeLiteralStarts, "Expr-CompositeLiteralStarts")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseCompoundLiteral()
	}
	ok, err = p.optionalToken(TT.Identifier, "Expr-Identifier")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseIdentifier()
	}
	return nil, NewParseError("unknown expression", nil, "Expr")
}

func (p *NiceExprParser) ParseVariableDeclaration() (*ast.VariableDeclaration, *ParseError) {
	expr := new(ast.VariableDeclaration)
	_, err := p.expectToken(TT.Var, "VariableDeclaration-Var")
	if err != nil {
		return expr, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return expr, err
	}
	typeExpr, err := p.ParseType()
	if err != nil {
		return expr, err
	}
	_, err = p.expectToken(TT.Is, "VariableDeclaration-Is")
	if err != nil {
		return expr, err
	}
	value, err := p.ParseExpr()
	if err != nil {
		return expr, err
	}
	expr.Name = name
	expr.Type = typeExpr
	expr.Value = value
	return expr, nil
}
func (p *NiceExprParser) ParseConstantDeclaration() (*ast.ConstantDeclaration, *ParseError) {
	expr := new(ast.ConstantDeclaration)
	_, err := p.expectToken(TT.Const, "ConstantDeclaration-Const")
	if err != nil {
		return expr, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return expr, err
	}
	typeExpr, err := p.ParseType()
	if err != nil {
		return expr, err
	}
	_, err = p.expectToken(TT.Is, "ConstantDeclaration-Is")
	if err != nil {
		return expr, err
	}
	value, err := p.ParseExpr()
	if err != nil {
		return expr, err
	}
	expr.Name = name
	expr.Type = typeExpr
	expr.Value = value
	return expr, nil
}

func (p *NiceExprParser) ParseAssignment() (*ast.Assignment, *ParseError) {
	ae := new(ast.Assignment)
	_, err := p.expectToken(TT.Set, "Assignment-Set")
	if err != nil {
		return ae, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return ae, err
	}
	op, err := p.expectAny(TT.AssignmentOperations, "Assignment-Op")
	if err != nil {
		return ae, err
	}
	value, err := p.ParseExpr()
	if err != nil {
		return ae, err
	}
	ae.Name = name
	ae.Op = op
	ae.Value = value
	return ae, nil
}

func (p *NiceExprParser) ParseIdentifier() (*ast.Identifier, *ParseError) {
	id := new(ast.Identifier)
	token, err := p.expectToken(TT.Identifier, "Literal-CheckLiteral")
	if err != nil {
		return id, err
	} else if token == nil {
		return id, NewParseError("invalid token for literal", token, "Literal-CheckLiteral")
	}
	id.Name = token
	return id, nil
}

func (p *NiceExprParser) ParseLiteral() (ast.Literal, *ParseError) {
	ok, err := p.checkAny(TT.CompositeLiteralStarts, "Literal-Start")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseCompoundLiteral()
	}
	return p.ParsePrimitiveLiteral()
}

func (p *NiceExprParser) ParsePrimitiveLiteral() (*ast.PrimitiveLiteral, *ParseError) {
	pe := new(ast.PrimitiveLiteral)
	token, err := p.expectAny(TT.PrimitiveLiterals, "Literal-CheckLiteral")
	if err != nil {
		return pe, err
	} else if token == nil {
		return pe, NewParseError("invalid token for literal", token, "Literal-CheckLiteral")
	}
	pe.Token = token
	return pe, nil
}

func (p *NiceExprParser) ParseCompoundLiteral() (ast.CompoundLiteral, *ParseError) {
	ok, err := p.optionalToken(TT.LeftBracket, "CompoundLiteral-List?")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseListLiteral()
	}
	ok, err = p.optionalToken(TT.LeftTriangle, "CompoundLiteral-Map?")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseMapLiteral()
	}
	return nil, NewParseError("unknown compound literal start", nil, "CompoundLiteral")
}

func (p *NiceExprParser) ParseListLiteral() (*ast.ListLiteral, *ParseError) {
	l := new(ast.ListLiteral)
	_, err := p.expectToken(TT.LeftBracket, "ListLiteral-Start")
	if err != nil {
		return nil, err
	}
	for {
		// list items are comma-separated, and trailing comma is needed
		ok, err := p.optionalToken(TT.RightBracket, "ListLiteral-End")
		if err != nil {
			return l, err
		} else if ok {
			break
		}
		value, err := p.ParseExpr()
		if err != nil {
			return l, err
		}
		_, err = p.expectToken(TT.Comma, "ListLiteral-ItemComma")
		if err != nil {
			return l, err
		}
		l.Values = append(l.Values, value)
	}
	_, err = p.expectToken(TT.RightBracket, "ListLiteral-End")
	if err != nil {
		return l, err
	}
	return l, nil
}
func (p *NiceExprParser) ParseMapLiteral() (*ast.MapLiteral, *ParseError) {
	m := new(ast.MapLiteral)
	m.Values = make(map[ast.Expr]ast.Expr)
	_, err := p.expectToken(TT.LeftTriangle, "MapLiteral-Start")
	if err != nil {
		return nil, err
	}
	for {
		// list items are comma-separated, and trailing comma is needed
		ok, err := p.optionalToken(TT.RightTriangle, "MapLiteral-End")
		if err != nil {
			return m, err
		} else if ok {
			break
		}
		key, err := p.ParseExpr()
		if err != nil {
			return m, err
		}
		_, err = p.expectToken(TT.Colon, "MapLiteral-ItemColon")
		if err != nil {
			return m, err
		}
		value, err := p.ParseExpr()
		if err != nil {
			return m, err
		}
		_, err = p.expectToken(TT.Comma, "MapLiteral-ItemComma")
		if err != nil {
			return m, err
		}
		m.Values[key] = value
	}
	_, err = p.expectToken(TT.RightTriangle, "MapLiteral-End")
	if err != nil {
		return m, err
	}
	return m, nil
}

func (p *NiceExprParser) ParseType() (ast.Type, *ParseError) {
	ok, err := p.checkAny(TT.PrimitiveTypes, "Type-CheckPrimitive")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParsePrimitiveType()
	}
	ok, err = p.checkAny(TT.CompositeTypes, "Type-CheckComposite")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseCompositeType()
	}
	ok, err = p.checkToken(TT.Func, "Type-CheckFunc")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseFuncType()
	}
	return nil, NewParseError("type not found", nil, "Type")
}

func (p *NiceExprParser) ParsePrimitiveType() (*ast.PrimitiveType, *ParseError) {
	pt := new(ast.PrimitiveType)
	name, err := p.expectAny(TT.PrimitiveTypes, "PrimitiveType")
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, NewParseError("nil typename", name, "PrimitiveType")
	}
	pt.Name = name
	return pt, nil
}

func (p *NiceExprParser) ParseCompositeType() (ast.Type, *ParseError) {
	ok, err := p.optionalToken(TT.List, "CompositeType-List")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseListType()
	}
	ok, err = p.optionalToken(TT.Map, "CompositeType-Map")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseMapType()
	}
	return nil, NewParseError("expected `list` or `map`", nil, "CompositeType")
}

func (p *NiceExprParser) ParseListType() (*ast.ListType, *ParseError) {
	lte := new(ast.ListType)
	_, err := p.expectToken(TT.List, "ListType-List")
	if err != nil {
		return lte, err
	}
	_, err = p.expectToken(TT.LeftBracket, "ListType-LeftBracket")
	if err != nil {
		return lte, err
	}
	valueType, err := p.ParseType()
	if err != nil {
		return lte, err
	}
	_, err = p.expectToken(TT.RightBracket, "ListType-RightBracket")
	if err != nil {
		return lte, err
	}
	lte.ValueType = valueType
	return lte, nil
}

func (p *NiceExprParser) ParseMapType() (*ast.MapType, *ParseError) {
	mte := new(ast.MapType)
	_, err := p.expectToken(TT.Map, "MapType-Map")
	if err != nil {
		return mte, err
	}
	_, err = p.expectToken(TT.LeftBracket, "MapType-LeftBracket")
	if err != nil {
		return mte, err
	}
	keyType, err := p.ParseType()
	if err != nil {
		return mte, err
	}
	_, err = p.expectToken(TT.RightBracket, "MapType-RightBracket")
	if err != nil {
		return mte, err
	}
	valueType, err := p.ParseType()
	if err != nil {
		return mte, err
	}
	mte.KeyType = keyType
	mte.ValueType = valueType
	return mte, nil
}

func (p *NiceExprParser) ParseFuncType() (*ast.FuncType, *ParseError) {
	fte := new(ast.FuncType)
	_, err := p.expectToken(TT.Func, "FuncType-Func")
	if err != nil {
		return fte, nil
	}
	_, err = p.expectToken(TT.LeftParen, "FuncType-ArgsStart")
	if err != nil {
		return fte, nil
	}
	// arguments
	for {
		ok, err := p.optionalToken(TT.RightParen, "FuncType-ArgsEnd")
		if err != nil {
			return fte, err
		} else if ok {
			_, err := p.getNextToken("FuncType-ArgsEnd")
			if err != nil {
				return fte, err
			}
			break
		}
		t, err := p.ParseType()
		if err != nil {
			return fte, err
		}
		_, err = p.expectToken(TT.Comma, "FuncType-CommaBetweenArgs")
		if err != nil {
			return fte, err
		}
		fte.InputTypes = append(fte.InputTypes, t)
	}
	// output type, optional
	out, err := p.ParseType()
	if err != nil {
		return fte, err
	}
	fte.OutputType = out
	return fte, nil
}