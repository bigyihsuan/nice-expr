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
		return p.ParseVariableDeclarationExpr()
	}
	ok, err = p.optionalToken(TT.Const, "Expr-ConstDecl")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseConstantDeclarationExpr()
	}
	ok, err = p.optionalToken(TT.Set, "Expr-Assignment")
	if err != nil {
		return nil, err
	} else if ok {
		return p.ParseAssignmentExpr()
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

func (p *NiceExprParser) ParseVariableDeclarationExpr() (*VariableDeclarationExpr, *ParseError) {
	expr := new(VariableDeclarationExpr)
	_, err := p.expectToken(TT.Var, "VariableDeclarationExpr-Var")
	if err != nil {
		return expr, err
	}
	decl, err := p.ParseDeclarationExpr()
	if err != nil {
		return expr, err
	}
	expr.DeclarationExpr = decl
	return expr, nil
}
func (p *NiceExprParser) ParseConstantDeclarationExpr() (*ConstantDeclarationExpr, *ParseError) {
	expr := new(ConstantDeclarationExpr)
	_, err := p.expectToken(TT.Const, "ConstantDeclarationExpr-Const")
	if err != nil {
		return expr, err
	}
	decl, err := p.ParseDeclarationExpr()
	if err != nil {
		return expr, err
	}
	expr.DeclarationExpr = decl
	return expr, nil
}

func (p *NiceExprParser) ParseDeclarationExpr() (*DeclarationExpr, *ParseError) {
	ae := new(DeclarationExpr)
	name, err := p.ParseIdentifier()
	if err != nil {
		return ae, err
	}
	typeExpr, err := p.ParseTypeExpr()
	if err != nil {
		return ae, err
	}
	_, err = p.expectToken(TT.Is, "DeclarationExpr-Is")
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

func (p *NiceExprParser) ParseAssignmentExpr() (*AssignmentExpr, *ParseError) {
	ae := new(AssignmentExpr)
	_, err := p.expectToken(TT.Set, "AssignmentExpr-Set")
	if err != nil {
		return ae, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return ae, err
	}
	op, err := p.expectAny(TT.AssignmentOperations, "AssignmentExpr-Op")
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

func (p *NiceExprParser) ParseTypeExpr() (TypeExpr, *ParseError) {
	ok, err := p.checkAny(TT.PrimitiveTypes, "TypeExpr-CheckPrimitive")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParsePrimitiveTypeExpr()
	}
	ok, err = p.checkAny(TT.CompositeTypes, "TypeExpr-CheckComposite")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseCompositeTypeExpr()
	}
	return nil, NewParseError("type not found", nil, "TypeExpr")
}

func (p *NiceExprParser) ParseCompositeTypeExpr() (TypeExpr, *ParseError) {
	ok, err := p.optionalToken(TT.ListType, "CompositeTypeExpr-List")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseListTypeExpr()
	}
	ok, err = p.optionalToken(TT.MapType, "CompositeTypeExpr-Map")
	if err != nil {
		return nil, err
	} else if ok && err == nil {
		return p.ParseMapTypeExpr()
	}
	return nil, NewParseError("expected `list` or `map`", nil, "CompositeTypeExpr")
}

func (p *NiceExprParser) ParseListTypeExpr() (*ListTypeExpr, *ParseError) {
	lte := new(ListTypeExpr)
	_, err := p.expectToken(TT.ListType, "ListTypeExpr-List")
	if err != nil {
		return lte, err
	}
	_, err = p.expectToken(TT.LeftBracket, "ListTypeExpr-LeftBracket")
	if err != nil {
		return lte, err
	}
	valueType, err := p.ParseTypeExpr()
	if err != nil {
		return lte, err
	}
	_, err = p.expectToken(TT.RightBracket, "ListTypeExpr-RightBracket")
	if err != nil {
		return lte, err
	}
	lte.ValueType = valueType
	return lte, nil
}
func (p *NiceExprParser) ParseMapTypeExpr() (*MapTypeExpr, *ParseError) {
	mte := new(MapTypeExpr)
	_, err := p.expectToken(TT.MapType, "MapTypeExpr-Map")
	if err != nil {
		return mte, err
	}
	_, err = p.expectToken(TT.LeftBracket, "MapTypeExpr-LeftBracket")
	if err != nil {
		return mte, err
	}
	keyType, err := p.ParseTypeExpr()
	if err != nil {
		return mte, err
	}
	_, err = p.expectToken(TT.RightBracket, "MapTypeExpr-RightBracket")
	if err != nil {
		return mte, err
	}
	valueType, err := p.ParseTypeExpr()
	if err != nil {
		return mte, err
	}
	mte.KeyType = keyType
	mte.ValueType = valueType
	return mte, nil
}

func (p *NiceExprParser) ParsePrimitiveTypeExpr() (*PrimitiveTypeExpr, *ParseError) {
	pt := new(PrimitiveTypeExpr)
	name, err := p.expectAny(TT.PrimitiveTypes, "PrimitiveTypeExpr")
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, NewParseError("nil typename", name, "PrimitiveTypeExpr")
	}
	pt.Name = name
	return pt, nil
}
