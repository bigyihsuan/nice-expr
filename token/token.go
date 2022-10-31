package token

import (
	"fmt"
	"nice-expr/token/tokentype"
	"strings"
)

// A lexed token.
type Token struct {
	Tt      tokentype.TokenType // the kind of token
	Lexeme  string              // the source string of the token
	Value   interface{}         // the value output by the lexer
	CodePos int                 // character index of this token
	Line    int                 // the 1-indexed line this token appears in
	Start   int                 // the 1-indexed starting index in the line the lexeme starts at, inclusive
	End     int                 // the 1-indexed ending index in the line the lexeme starts at, exclusive
}

func (t Token) String() string {
	return fmt.Sprintf(
		"{%v `%s` %d:%d:%d-%d}",
		t.Tt,
		strings.ReplaceAll(t.Lexeme, "\n", "\\n"),
		t.CodePos,
		t.Line,
		t.Start,
		t.End,
	)
}
