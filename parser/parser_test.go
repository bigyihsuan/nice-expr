package parser_test

import (
	"nice-expr/lexer"
	"nice-expr/parser"
	"os"
	"strings"
	"testing"

	"github.com/db47h/lex"
)

func TestParsePrimitiveLiteral(t *testing.T) {
	fileName := "./../test/primitive-literal.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	cases := strings.Split(string(test), "\n")
	for _, c := range cases {
		if c == "\n" || len(c) <= 0 {
			continue
		}
		file := lex.NewFile(fileName, strings.NewReader(c))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		nicerParser := parser.NewNiceExprParser(tokens)
		expr, perr := nicerParser.ParsePrimitiveLiteral()
		if perr != nil {
			t.Fatal(perr)
		}
		if expr == nil {
			t.Fatal("parsed nil")
		}
		// t.Log(expr)
	}
}
func TestParseCompoundLiteral(t *testing.T) {
	fileName := "./../test/compound-literal.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	cases := strings.Split(string(test), "\n")
	for _, c := range cases {
		if c == "\n" || len(c) <= 0 {
			continue
		}
		file := lex.NewFile(fileName, strings.NewReader(c))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		nicerParser := parser.NewNiceExprParser(tokens)
		expr, perr := nicerParser.ParseCompoundLiteral()
		if perr != nil {
			t.Fatal(perr)
		}
		if expr == nil {
			t.Fatal("parsed nil")
		}
		// t.Log(expr)
	}
}

func TestParseType(t *testing.T) {
	fileName := "./../test/type.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	cases := strings.Split(string(test), "\n")
	for _, c := range cases {
		if c == "\n" || len(c) <= 0 {
			continue
		}
		file := lex.NewFile(fileName, strings.NewReader(c))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		nicerParser := parser.NewNiceExprParser(tokens)
		typeExpr, perr := nicerParser.ParseType()
		if perr != nil {
			t.Fatal(perr)
		}
		if typeExpr == nil {
			t.Fatal("parsed nil")
		}
		// t.Log(typeExpr)
	}
}

func TestParseAssignments(t *testing.T) {
	fileName := "./../test/assignments.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	cases := strings.Split(string(test), "\n")
	for _, c := range cases {
		if c == "\n" || len(c) <= 0 {
			continue
		}
		file := lex.NewFile(fileName, strings.NewReader(c))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		nicerParser := parser.NewNiceExprParser(tokens)
		expr, perr := nicerParser.ParseStatement()
		if perr != nil {
			t.Fatal(perr)
		}
		if expr == nil {
			t.Fatal("parsed nil")
		}
		// t.Log(expr)
	}
}

func TestParseUnary(t *testing.T) {
	fileName := "./../test/unary.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	cases := strings.Split(string(test), "\n")
	for _, c := range cases {
		if c == "\n" || len(c) <= 0 {
			continue
		}
		file := lex.NewFile(fileName, strings.NewReader(c))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()

		nicerParser := parser.NewNiceExprParser(tokens)
		expr, perr := nicerParser.ParseUnary()
		if perr != nil {
			t.Fatal(perr)
		}
		if expr == nil {
			t.Fatal("parsed nil")
		}
	}
}
