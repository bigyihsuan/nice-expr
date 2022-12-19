package main

import (
	"bytes"
	"fmt"
	"nice-expr/src/lexer"
	"nice-expr/src/parser"
	"nice-expr/src/visitor"
	"os"

	"github.com/db47h/lex"
	goflags "github.com/jessevdk/go-flags"
)

type Options struct {
	ShowTokens      bool   `short:"t" long:"showTokens"`
	ShowParseTree   bool   `short:"p" long:"showParseTree"`
	ShowStringVisit bool   `short:"s" long:"showStringVisit"`
	ShowTypeCheck   bool   `short:"T" long:"showTypeCheck"`
	ShowEvaluation  bool   `short:"e" long:"showEvaluation"`
	Code            string `short:"c" default:""`
}

func main() {
	var options Options
	var text []byte
	var fileName string
	remaining, err := goflags.Parse(&options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if options.Code != "" {
		text = []byte(options.Code)
		fileName = string(text)
	} else if len(remaining) < 1 {
		fmt.Fprintln(os.Stdout, fmt.Errorf("not enough arguments; expected 1 (filename), got %d", len(remaining)))
		return
	} else {
		fileName = remaining[0]
		text, err = os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
	text = append(text, '\n')
	byteReader := bytes.NewBuffer(text)
	file := lex.NewFile(fileName, byteReader)
	fmt.Println(file.Name())
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()

	if options.ShowTokens {
		fmt.Println("tokens:", tokens)
		fmt.Println()
	}

	nicerParser := parser.NewNiceExprParser(tokens)
	program, pe := nicerParser.Program()
	if pe != nil {
		fmt.Fprintln(os.Stderr, pe)
		fmt.Fprintf(os.Stderr, "last seen token: %v\n", nicerParser.LastSeen())
		return
	}

	if options.ShowParseTree {
		fmt.Println("program:", program.Statements)
		fmt.Println()
	}

	streval := visitor.NewStringVisitor()
	program.Accept(streval)

	if options.ShowStringVisit {
		fmt.Println("string visitor")
		fmt.Println("str:", streval.String())
		fmt.Println()
	}
	typevis := visitor.NewTypeChecker()
	program.Accept(typevis)

	errs := typevis.Errors()
	for errs.Len() > 0 {
		typeError, err := errs.Pop()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		} else if typeError != nil {
			fmt.Println("type error:", typeError)
		}
	}

	if options.ShowTypeCheck {
		fmt.Println("type checker")
		fmt.Println("type stack:", typevis.TypeStack())
		fmt.Println("identifiers:", typevis.Identifiers())
		fmt.Println()
	}
	if typevis.HasErrors() {
		return
	}
	nicerEvaluator := visitor.NewEvaluatingVisitor()
	program.Accept(nicerEvaluator)

	valErrs := nicerEvaluator.Errors()
	for valErrs.Len() > 0 {
		valErr, err := valErrs.Pop()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		} else if valErr != nil {
			fmt.Println("evaluation error:", valErr)
		}
	}
	if options.ShowEvaluation {
		fmt.Println("evaluator")
		fmt.Println("value stack:", nicerEvaluator.ValueStack())
		fmt.Println("identifiers:", nicerEvaluator.Identifiers())
	}
}
