//go:generate stringer -type=TokenType
package tokentype

import (
	"github.com/db47h/lex"
)

type TokenType lex.Token

const (
	Invalid TokenType = iota
	EOF
	// literals
	Integer
	Float
	String
	// braces
	LeftBracket
	RightBracket
	LeftBrace
	RightBrace
	LeftParen
	RightParen
	LeftTriangle
	RightTriangle
	// math operators
	Plus
	Minus
	Star
	Slash
	Percent
	PlusEqual
	MinusEqual
	StarEqual
	SlashEqual
	PercentEqual
	// comparison ops
	Equal
	Greater
	GreaterEqual
	Less
	LessEqual
	// misc
	Identifier
	Comma
	Colon
	Semicolon
	Underscore
	// keywords
	Var
	Const
	Set
	Is
	For
	Break
	Return
	Func
	Not
	If
	Else
	// built-in
	NoneType
	IntType
	FloatType
	StringType
	ListType
	MapType
)

var Keywords = map[string]TokenType{
	"var":    Var,
	"const":  Const,
	"set":    Set,
	"is":     Is,
	"for":    For,
	"break":  Break,
	"return": Return,
	"func":   Func,
	"not":    Not,
	"if":     If,
	"else":   Else,
	"none":   NoneType,
	"int":    IntType,
	"float":  FloatType,
	"string": StringType,
	"list":   ListType,
	"map":    MapType,
}

var (
	PrimitiveLiterals      = []TokenType{Integer, Float, String}
	CompositeLiteralStarts = []TokenType{LeftBracket, LeftTriangle}
	PrimitiveTypes         = []TokenType{NoneType, IntType, FloatType, StringType}
	CompositeTypes         = []TokenType{ListType, MapType}
	AssignmentOperations   = []TokenType{Is, PlusEqual, MinusEqual, StarEqual, SlashEqual, PercentEqual}
)

func ToTt(lexTok lex.Token) TokenType {
	return TokenType(lexTok)
}
func ToLt(tokType TokenType) lex.Token {
	return lex.Token(tokType)
}
