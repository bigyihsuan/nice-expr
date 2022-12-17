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
	Debug bool   `short:"d" long:"debug"`
	Code  string `short:"c" default:""`
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
	byteReader := bytes.NewBuffer(text)
	file := lex.NewFile(fileName, byteReader)
	fmt.Println(file.Name())
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()

	if options.Debug {
		fmt.Println(tokens)
		fmt.Println()
	}

	nicerParser := parser.NewNiceExprParser(tokens)
	program, pe := nicerParser.Program()
	if pe != nil {
		fmt.Fprintln(os.Stderr, pe)
		fmt.Fprintf(os.Stderr, "last seen token: %v\n", nicerParser.LastSeen())
		return
	}

	if options.Debug {
		fmt.Println("program:", program.Statements)
		fmt.Println()

		fmt.Println("string visitor")
	}

	streval := visitor.NewStringVisitor()
	program.Accept(streval)

	if options.Debug {

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
			fmt.Println("got type error: ", typeError)
		}
	}

	if options.Debug {
		fmt.Println(typevis.TypeStack())
		fmt.Println(typevis.Identifiers())
		fmt.Println()
		fmt.Println("evaluator")
	}

	nicerEvaluator := visitor.NewEvaluatingVisitor()
	program.Accept(nicerEvaluator)

	if options.Debug {
		fmt.Println(nicerEvaluator.ValueStack())
		valErrs := nicerEvaluator.Errors()
		for valErrs.Len() > 0 {
			valErr, err := valErrs.Pop()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			} else if valErr != nil {
				fmt.Println("got evaluation error:", valErr)
			}
		}
		fmt.Println(nicerEvaluator.Identifiers())
	}
}
