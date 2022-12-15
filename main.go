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
	Code string `short:"c" default:""`
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

	fmt.Println(tokens)
	fmt.Println()

	nicerParser := parser.NewNiceExprParser(tokens)
	program, pe := nicerParser.Program()
	if pe != nil {
		fmt.Fprintln(os.Stderr, pe)
		fmt.Fprintf(os.Stderr, "last seen token: %v\n", nicerParser.LastSeen())
		return
	}

	fmt.Println("program:", program.Statements)
	fmt.Println()

	fmt.Println("string visitor")
	streval := visitor.NewStringVisitor()
	program.Accept(streval)

	fmt.Println("str:", streval.String())
	fmt.Println()

	typevis := visitor.NewTypeChecker()
	program.Accept(typevis)

	fmt.Println(typevis.TypeStack())
	errs := typevis.Errors()
	typeError, err := errs.Pop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	} else if typeError != nil {
		fmt.Println("got type error: ", typeError)
	}
	fmt.Println(typevis.Identifiers())

	// nicerEvaluator := evaluator.NewEvaluator()
	// ee := nicerEvaluator.EvaluateProgram(program)
	// if ee != nil {
	// 	fmt.Fprintln(os.Stderr, ee)
	// }

	// fmt.Println("Constants:", nicerEvaluator.Constants)
	// fmt.Println("Variables:", nicerEvaluator.Variables)
	// fmt.Println("ValueStack:", nicerEvaluator.ValueStack)
}
