//go:generate stringer -type=TokenType
package tokentype

import (
	"nice-expr/value"

	"github.com/db47h/lex"
)

type TokenType lex.Token

const (
	Invalid TokenType = iota
	EOF
	// literals
	Integer
	Floating
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
	// logical ops
	And
	Or
	Not
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
	If
	Else
	// built-in
	None
	Int
	Dec
	Str
	Bool
	List
	Map

	True
	False
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
	"and":    And,
	"or":     Or,
	"not":    Not,
	"if":     If,
	"else":   Else,
	"none":   None,
	"int":    Int,
	"dec":    Dec,
	"str":    Str,
	"bool":   Bool,
	"list":   List,
	"map":    Map,
	"true":   True,
	"false":  False,
}

var (
	PrimitiveLiterals      = []TokenType{None, Integer, Floating, String, True, False}
	CompositeLiteralStarts = []TokenType{LeftBracket, LeftTriangle}

	PrimitiveTypes       = []TokenType{None, Int, Dec, Str, Bool}
	CompoundTypes        = []TokenType{List, Map}
	Types                = append(PrimitiveTypes, CompoundTypes...)
	AssignmentOperations = []TokenType{Is, PlusEqual, MinusEqual, StarEqual, SlashEqual, PercentEqual}
	LitToType            = func() map[TokenType]value.ValueType {
		m := make(map[TokenType]value.ValueType)
		for i := range PrimitiveTypes {
			m[PrimitiveLiterals[i]] = value.NewValueType(PrimitiveTypes[i].String())
		}
		// manually add true and false
		m[True] = value.NewValueType(Bool.String())
		m[False] = value.NewValueType(Bool.String())
		return m
	}()
	NoneType = value.NewValueType("None")
	IntType  = value.NewValueType("Int")
	DecType  = value.NewValueType("Dec")
	StrType  = value.NewValueType("Str")
	BoolType = value.NewValueType("Bool")
	ListType = value.NewValueType("List")
	MapType  = value.NewValueType("Map")

	BinMathOps        = []TokenType{Plus, Minus, Star, Slash, Percent}
	BinCompOps        = []TokenType{Equal, Greater, GreaterEqual, Less, LessEqual}
	BinLogOps         = []TokenType{And, Or}
	AssignmentMathOps = []TokenType{PlusEqual, MinusEqual, StarEqual, SlashEqual, PercentEqual}
	BinOps            = append(append(BinMathOps, BinCompOps...), BinLogOps...)
	UnaryOps          = []TokenType{Not, Minus}
)

func ToTt(lexTok lex.Token) TokenType {
	return TokenType(lexTok)
}
func ToLt(tokType TokenType) lex.Token {
	return lex.Token(tokType)
}
