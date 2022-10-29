package lexer

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("{%v `%s` %d:%d-%d}", t.Tt, strings.ReplaceAll(t.Lexeme, "\n", "\\n"), t.Line, t.Start, t.End)
}
