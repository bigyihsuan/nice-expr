package lexer

import (
	"fmt"
	"unicode"

	"github.com/db47h/lex"
	"github.com/db47h/lex/state"
)

type NiceExprLexer struct {
	lex.Lexer
	Line int
}

func NewLexer(file *lex.File) *NiceExprLexer {
	l := &NiceExprLexer{}
	l.Line = 0
	l.Lexer = *lex.NewLexer(file, l.program)
	return l
}

func (nel *NiceExprLexer) LexAll() []Token {
	var tokens []Token
	for tok, pos, v := nel.Lex(); tok != toLt(EOF); tok, pos, v = nel.Lex() {
		s := fmt.Sprint(v)
		tokens = append(tokens, Token{toTt(tok), s, nel.Line, pos, pos + len(s)})
	}
	return tokens
}

func (nel *NiceExprLexer) program(s *lex.State) lex.StateFn {
	r := s.Next()
	pos := s.Pos()

	switch r { // single-character tokens
	case lex.EOF:
		// s.Emit(pos, Semicolon, ";")
		s.Emit(pos, toLt(EOF), "")
		return nil
	case ';': // newlines separate statements
		s.Emit(pos, toLt(Semicolon), ";")
		return nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return state.Number(toLt(Integer), toLt(Float), '.')
	case '"': // strings
		return state.QuotedString(toLt(String))
	case ',':
		s.Emit(pos, toLt(Comma), string(r))
		return nil
	case '+':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, toLt(PlusEqual), "+=")
			return nil
		}
		s.Emit(pos, toLt(Plus), string(r))
		return nil
	case '-': // either binary minus, unary minus, or -th operator
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, toLt(MinusEqual), "-=")
			return nil
		}
		s.Emit(pos, toLt(Minus), string(r))
		return nil
	case '*':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, toLt(StarEqual), "*=")
			return nil
		}
		s.Emit(pos, toLt(Star), string(r))
		return nil
	case '/':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, toLt(SlashEqual), "/=")
			return nil
		} else if s.Peek() == '/' {
			return nel.comment
		}
		s.Emit(pos, toLt(Slash), string(r))
		return nil
	case '%':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, toLt(PercentEqual), "%=")
			return nil
		}
		s.Emit(pos, toLt(Percent), string(r))
		return nil
	case '=':
		s.Emit(pos, toLt(Equal), string(r))
		return nil
	case '>':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, toLt(GreaterEqual), ">=")
		} else {
			s.Emit(pos, toLt(Greater), ">")
			s.Backup()
		}
		return nil
	case '<':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, toLt(LessEqual), "<=")
		} else {
			s.Emit(pos, toLt(Less), "<")
			s.Backup()
		}
		return nil
	case '(':
		s.Emit(pos, toLt(LeftParen), "(")
		return nil
	case ')':
		s.Emit(pos, toLt(RightParen), ")")
		return nil
	case '[':
		s.Emit(pos, toLt(LeftBracket), "[")
		return nil
	case ']':
		s.Emit(pos, toLt(RightBracket), "]")
		return nil
	case '{':
		s.Emit(pos, toLt(LeftBrace), "{")
		return nil
	case '}':
		s.Emit(pos, toLt(RightBrace), "}")
		return nil
	case '_':
		s.Emit(pos, toLt(Underscore), "_")
		return nil
	}

	switch {
	case unicode.IsSpace(r):
		for r = s.Next(); unicode.IsSpace(r); r = s.Next() {
			// nop
			if r == '\n' {
				nel.Line++
			}
		}
		s.Backup()
		return nil
	case unicode.IsLetter(r):
		return nel.ident_or_keyword
	}
	return nil
}

func (nel *NiceExprLexer) comment(s *lex.State) lex.StateFn {
	comment := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		comment = append(comment[:0], l.Current())
		for r := l.Next(); r != '\n'; r = l.Next() {
			comment = append(comment, r)
		}
		l.Backup()
		// completely ignore comments
		return nil
	}
}

func (nel *NiceExprLexer) ident_or_keyword(s *lex.State) lex.StateFn {
	name := make([]rune, 0, 64)
	return func(l *lex.State) lex.StateFn {
		pos := l.Pos()
		name = append(name[:0], l.Current())
		for r := l.Next(); unicode.IsLetter(r); r = l.Next() {
			name = append(name, r)
		}
		l.Backup()
		if tok, ok := Keywords[string(name)]; ok {
			l.Emit(pos, toLt(tok), string(name))
		} else {
			l.Emit(pos, toLt(Identifier), string(name))
		}
		return nil
	}
}
