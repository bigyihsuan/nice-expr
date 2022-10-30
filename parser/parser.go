package parser

import (
	"fmt"
	"nice-expr/lexer/token"
	TT "nice-expr/lexer/token/tokentype"
)

type NiceExprParser struct {
	Tokens    []token.Token
	lastToken *token.Token
}

func NewNiceExprParser(tokens []token.Token) NiceExprParser {
	return NiceExprParser{tokens, &token.Token{Tt: TT.EOF, Lexeme: "EOF", CodePos: -1, Line: -1, Start: -1, End: -1}}
}

// consume and return the next token in the token queue.
func (p *NiceExprParser) getNextToken(lastRule string) (*token.Token, *ParseError) {
	if len(p.Tokens) < 1 {
		return nil, NewParseError("out of tokens", nil, lastRule)
	}
	tok := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	p.lastToken = &tok
	return &tok, nil
}

// peek at the front of the token queue.
func (p *NiceExprParser) peekToken(lastRule string) (*token.Token, *ParseError) {
	if len(p.Tokens) < 1 {
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

func (p *NiceExprParser) ParseStatement() (Node, *ParseError) {
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

func (p *NiceExprParser) ParseExpr() (Node, *ParseError) {
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
	ok, err = p.checkAny(TT.PrimitiveLiterals, "Expr-Literals")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParsePrimitiveLiteral()
	}
	ok, err = p.optionalToken(TT.Identifier, "Expr-Identifier")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseIdentifier()
	}
	return nil, NewParseError("unknown expression", nil, "Expr")
}

func (p *NiceExprParser) ParseVariableDeclaration() (*VariableDeclaration, *ParseError) {
	expr := new(VariableDeclaration)
	_, err := p.expectToken(TT.Var, "VariableDeclaration-Var")
	if err != nil {
		return expr, err
	}
	decl, err := p.ParseDeclaration()
	if err != nil {
		return expr, err
	}
	expr.Declaration = decl
	return expr, nil
}
func (p *NiceExprParser) ParseConstantDeclaration() (*ConstantDeclaration, *ParseError) {
	expr := new(ConstantDeclaration)
	_, err := p.expectToken(TT.Const, "ConstantDeclaration-Const")
	if err != nil {
		return expr, err
	}
	decl, err := p.ParseDeclaration()
	if err != nil {
		return expr, err
	}
	expr.Declaration = decl
	return expr, nil
}

func (p *NiceExprParser) ParseDeclaration() (*Declaration, *ParseError) {
	ae := new(Declaration)
	name, err := p.ParseIdentifier()
	if err != nil {
		return ae, err
	}
	typeExpr, err := p.ParseType()
	if err != nil {
		return ae, err
	}
	_, err = p.expectToken(TT.Is, "Declaration-Is")
	if err != nil {
		return ae, err
	}
	value, err := p.ParseExpr()
	if err != nil {
		return ae, err
	}
	ae.Name = name
	ae.Type = typeExpr
	ae.Value = value
	return ae, nil
}

func (p *NiceExprParser) ParseAssignment() (*Assignment, *ParseError) {
	ae := new(Assignment)
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

func (p *NiceExprParser) ParseIdentifier() (*Identifier, *ParseError) {
	id := new(Identifier)
	token, err := p.expectToken(TT.Identifier, "Literal-CheckLiteral")
	if err != nil {
		return id, err
	} else if token == nil {
		return id, NewParseError("invalid token for literal", token, "Literal-CheckLiteral")
	}
	id.Name = token
	return id, nil
}

func (p *NiceExprParser) ParseLiteral() (Literal, *ParseError) {
	ok, err := p.checkAny(TT.CompositeLiteralStarts, "Literal-Start")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseCompoundLiteral()
	}
	return p.ParsePrimitiveLiteral()
}

func (p *NiceExprParser) ParsePrimitiveLiteral() (*PrimitiveLiteral, *ParseError) {
	pe := new(PrimitiveLiteral)
	token, err := p.expectAny(TT.PrimitiveLiterals, "Literal-CheckLiteral")
	if err != nil {
		return pe, err
	} else if token == nil {
		return pe, NewParseError("invalid token for literal", token, "Literal-CheckLiteral")
	}
	pe.Token = token
	return pe, nil
}

func (p *NiceExprParser) ParseCompoundLiteral() (CompoundLiteral, *ParseError) {
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

func (p *NiceExprParser) ParseListLiteral() (*ListLiteral, *ParseError) {
	l := new(ListLiteral)
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
		value, err := p.ParseLiteral()
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
func (p *NiceExprParser) ParseMapLiteral() (*MapLiteral, *ParseError) {
	m := new(MapLiteral)
	m.Values = make(map[Literal]Literal)
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
		key, err := p.ParseLiteral()
		if err != nil {
			return m, err
		}
		_, err = p.expectToken(TT.Colon, "MapLiteral-ItemColon")
		if err != nil {
			return m, err
		}
		value, err := p.ParseLiteral()
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

func (p *NiceExprParser) ParseType() (Type, *ParseError) {
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

func (p *NiceExprParser) ParsePrimitiveType() (*PrimitiveType, *ParseError) {
	pt := new(PrimitiveType)
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

func (p *NiceExprParser) ParseCompositeType() (Type, *ParseError) {
	ok, err := p.optionalToken(TT.ListType, "CompositeType-List")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseListType()
	}
	ok, err = p.optionalToken(TT.MapType, "CompositeType-Map")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseMapType()
	}
	return nil, NewParseError("expected `list` or `map`", nil, "CompositeType")
}

func (p *NiceExprParser) ParseListType() (*ListType, *ParseError) {
	lte := new(ListType)
	_, err := p.expectToken(TT.ListType, "ListType-List")
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

func (p *NiceExprParser) ParseMapType() (*MapType, *ParseError) {
	mte := new(MapType)
	_, err := p.expectToken(TT.MapType, "MapType-Map")
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

func (p *NiceExprParser) ParseFuncType() (*FuncType, *ParseError) {
	fte := new(FuncType)
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
		t, err := p.ParseType()
		if err != nil {
			return fte, err
		}
		_, err = p.expectToken(TT.Comma, "FuncType-CommaBetweenArgs")
		if err != nil {
			return fte, err
		}
		fte.InputTypes = append(fte.InputTypes, t)
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
	}
	// output type
	out, err := p.ParseType()
	if err != nil {
		return fte, err
	}
	fte.OutputType = out
	return fte, nil
}
