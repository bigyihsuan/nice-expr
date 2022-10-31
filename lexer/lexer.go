package lexer

import (
	"fmt"
	"nice-expr/token"
	TT "nice-expr/token/tokentype"
	"unicode"

	"github.com/db47h/lex"
	"github.com/db47h/lex/state"
)

type NiceExprLexer struct {
	lex.Lexer
	Line       int
	LineStart  int
	CurrentPos int
}

func NewLexer(file *lex.File) *NiceExprLexer {
	l := &NiceExprLexer{}
	l.Line = 1
	l.Lexer = *lex.NewLexer(file, l.program)
	return l
}

func (nel *NiceExprLexer) LexAll() []token.Token {
	var tokens []token.Token
	for tok, pos, v := nel.Lex(); tok != TT.ToLt(TT.EOF); tok, pos, v = nel.Lex() {
		s := fmt.Sprint(v)
		if nel.CurrentPos > pos {
			nel.CurrentPos += len(s) + 1
		} else {
			nel.CurrentPos = pos + 1
		}
		tok := token.Token{
			Tt:      TT.ToTt(tok),
			Lexeme:  s,
			Value:   v,
			CodePos: nel.CurrentPos,
			Line:    nel.Line,
			Start:   nel.CurrentPos - nel.LineStart,
			End:     nel.CurrentPos - nel.LineStart + len(s),
		}

		tokens = append(tokens, tok)
	}
	return tokens
}

func (nel *NiceExprLexer) program(s *lex.State) lex.StateFn {
	r := s.Next()
	pos := s.Pos()

	switch r { // single-character tokens
	case lex.EOF:
		// s.Emit(pos, Semicolon, ";")
		s.Emit(pos, TT.ToLt(TT.EOF), "")
		return nil
	case ';': // newlines separate statements
		s.Emit(pos, TT.ToLt(TT.Semicolon), ";")
		return nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return state.Number(TT.ToLt(TT.Integer), TT.ToLt(TT.Floating), '.')
	case '"': // strings
		return state.QuotedString(TT.ToLt(TT.String))
	case ',':
		s.Emit(pos, TT.ToLt(TT.Comma), string(r))
		return nil
	case '+':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.PlusEqual), "+=")
			return nil
		}
		s.Emit(pos, TT.ToLt(TT.Plus), string(r))
		return nil
	case '-': // either binary minus, unary minus, or -th operator
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.MinusEqual), "-=")
			return nil
		}
		s.Emit(pos, TT.ToLt(TT.Minus), string(r))
		return nil
	case '*':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.StarEqual), "*=")
			return nil
		}
		s.Emit(pos, TT.ToLt(TT.Star), string(r))
		return nil
	case '/':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.SlashEqual), "/=")
			return nil
		} else if s.Peek() == '/' {
			return nel.comment
		}
		s.Emit(pos, TT.ToLt(TT.Slash), string(r))
		return nil
	case '%':
		if s.Peek() == '=' {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.PercentEqual), "%=")
			return nil
		}
		s.Emit(pos, TT.ToLt(TT.Percent), string(r))
		return nil
	case '=':
		s.Emit(pos, TT.ToLt(TT.Equal), string(r))
		return nil
	case '>':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, TT.ToLt(TT.GreaterEqual), ">=")
		} else {
			s.Emit(pos, TT.ToLt(TT.Greater), ">")
			s.Backup()
		}
		return nil
	case '<':
		r := s.Next()
		if r == '=' {
			s.Emit(pos, TT.ToLt(TT.LessEqual), "<=")
		} else if r == '|' {
			s.Emit(pos, TT.ToLt(TT.LeftTriangle), "<|")
		} else {
			s.Emit(pos, TT.ToLt(TT.Less), "<")
			s.Backup()
		}
		return nil
	case '(':
		s.Emit(pos, TT.ToLt(TT.LeftParen), "(")
		return nil
	case ')':
		s.Emit(pos, TT.ToLt(TT.RightParen), ")")
		return nil
	case '[':
		s.Emit(pos, TT.ToLt(TT.LeftBracket), "[")
		return nil
	case ']':
		s.Emit(pos, TT.ToLt(TT.RightBracket), "]")
		return nil
	case '|':
		if s.Peek() != '>' {
			s.Emit(pos, TT.ToLt(TT.Invalid), string(r)+string(s.Peek()))
		} else {
			s.Next()
			s.Emit(pos, TT.ToLt(TT.RightTriangle), "|>")
		}
		return nil
	case '{':
		s.Emit(pos, TT.ToLt(TT.LeftBrace), "{")
		return nil
	case '}':
		s.Emit(pos, TT.ToLt(TT.RightBrace), "}")
		return nil
	case '_':
		s.Emit(pos, TT.ToLt(TT.Underscore), "_")
		return nil
	case ':':
		s.Emit(pos, TT.ToLt(TT.Colon), ":")
		return nil
	}

	switch {
	case unicode.IsSpace(r):
		for r = s.Next(); unicode.IsSpace(r); r = s.Next() {
			// nop
			if r == '\n' {
				nel.Line++
				nel.LineStart = s.Pos()
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
		if tok, ok := TT.Keywords[string(name)]; ok {
			switch tok {
			case TT.True:
				l.Emit(pos, TT.ToLt(tok), true)
			case TT.False:
				l.Emit(pos, TT.ToLt(tok), false)
			default:
				l.Emit(pos, TT.ToLt(tok), string(name))
			}
		} else {
			l.Emit(pos, TT.ToLt(TT.Identifier), string(name))
		}
		return nil
	}
}
