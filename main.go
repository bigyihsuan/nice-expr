package main

import (
	"bytes"
	"fmt"
	"nice-expr/src/evaluator"
	"nice-expr/src/lexer"
	"nice-expr/src/parser"
	"os"

	"github.com/db47h/lex"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stdout, fmt.Errorf("not enough arguments; expected 1, got %d", len(os.Args)-1))
		return
	}
	fileName := os.Args[1]
	text, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
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
		return
	}

	fmt.Println(program)
	fmt.Println()

	nicerEvaluator := evaluator.NewEvaluator()
	ee := nicerEvaluator.EvaluateProgram(program)
	if ee != nil {
		fmt.Fprintln(os.Stderr, ee)
	}

	fmt.Println("Constants:", nicerEvaluator.Constants)
	fmt.Println("Variables:", nicerEvaluator.Variables)
	fmt.Println("ValueStack:", nicerEvaluator.ValueStack)
}
