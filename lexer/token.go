//go:generate stringer -type=TokenType
package lexer

import "fmt"

type TokenType int

const (
	Invalid TokenType = iota - 1
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
	Equal
	PlusEqual
	MinusEqual
	StarEqual
	SlashEqual
	// misc
	Comma
	Identifier
	// keywords
	Var
	Is
	For
	Break
	Return
	Func
)

// A lexed token.
type Token struct {
	Tt     TokenType // the kind of token
	Lexeme string    // the source string of the token
	Line   int       // the 1-indexed line this token appears in
	Start  int       // the 1-indexed starting index in the line the lexeme starts at, inclusive
	End    int       // the 1-indexed ending index in the line the lexeme starts at, exclusive
}

func (t Token) String() string {
	return fmt.Sprintf("{%s `%s` %d:%d,%d}", t.Tt, t.Lexeme, t.Line, t.Start, t.End)
}
