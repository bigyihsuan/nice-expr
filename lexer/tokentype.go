//go:generate stringer -type=TokenType
package lexer

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
	Comma
	Identifier
	Semicolon
	Underscore
	// keywords
	Var
	Is
	For
	Break
	Return
	Func
	Not
	If
	Else
)

var Keywords = map[string]TokenType{
	"var":    Var,
	"is":     Is,
	"for":    For,
	"break":  Break,
	"return": Return,
	"func":   Func,
	"not":    Not,
	"if":     If,
	"else":   Else,
}

func toTt(lexTok lex.Token) TokenType {
	return TokenType(lexTok)
}
func toLt(tokType TokenType) lex.Token {
	return lex.Token(tokType)
}
