package evaluator_test

import (
	"bytes"
	"math/big"
	"nice-expr/evaluator"
	"nice-expr/lexer"
	"nice-expr/parser"
	"os"
	"testing"

	"github.com/db47h/lex"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateDeclaration(t *testing.T) {
	fileName := "./../test/declarations.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	nicerParser := parser.NewNiceExprParser(tokens)

	program, err := nicerParser.ParseProgram()
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}

	evaluator := evaluator.NewEvaluator()
	ee := evaluator.EvaluateProgram(program)
	if ee != nil {
		t.Fatal(ee)
	}
	k, kv := evaluator.GetConstant("k")
	assert.NotNil(t, k)
	assert.NotNil(t, kv)
	assert.Equal(t, "hello world", kv.V.(string))

	x, xv := evaluator.GetVariable("x")
	assert.NotNil(t, x)
	assert.NotNil(t, xv)
	xi := xv.V.(*big.Int).Int64()
	assert.Equal(t, int64(10), xi)

	n, nv := evaluator.GetVariable("n")
	assert.NotNil(t, n)
	assert.NotNil(t, nv)
	nf, _ := nv.V.(*big.Float).Float64()
	assert.Equal(t, float64(2.34), nf)
}
