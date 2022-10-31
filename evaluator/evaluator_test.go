package evaluator_test

import (
	"bytes"
	"math/big"
	"nice-expr/evaluator"
	"nice-expr/lexer"
	"nice-expr/parser"
	"nice-expr/value"
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
	if assert.NotNil(t, kv) {
		assert.Equal(t, "hello world", kv.V.(string))
	}

	x, xv := evaluator.GetVariable("x")
	assert.NotNil(t, x)
	if assert.NotNil(t, xv) {
		xi := xv.V.(*big.Int).Int64()
		assert.Equal(t, int64(10), xi)
	}

	n, nv := evaluator.GetVariable("n")
	assert.NotNil(t, n)
	if assert.NotNil(t, nv) {
		nf, _ := nv.V.(*big.Float).Float64()
		assert.Equal(t, float64(2.34), nf)
	}

	intList, intListVal := evaluator.GetConstant("intList")
	assert.NotNil(t, intList)
	if assert.NotNil(t, intListVal) {
		intListActual := func() (i []int64) {
			for _, e := range intListVal.V.([]*value.Value) {
				i = append(i, e.V.(*big.Int).Int64())
			}
			return
		}()
		assert.Equal(t, []int64{1, 2, 3, 4, 5}, intListActual)
	}

	floatList, floatListVal := evaluator.GetVariable("decList")
	assert.NotNil(t, floatList)
	if assert.NotNil(t, floatListVal) {
		floatListActual := func() (i []float64) {
			for _, e := range floatListVal.V.([]*value.Value) {
				f, _ := e.V.(*big.Float).Float64()
				i = append(i, f)
			}
			return
		}()
		assert.Equal(t, []float64{1.1, 2.2, 3.3, 4.4, 5.5}, floatListActual)
	}

	intStrMap, intStrMapVal := evaluator.GetVariable("intStrMap")
	assert.NotNil(t, intStrMap)
	if assert.NotNil(t, intStrMapVal) {
		intStrMapActual := func() map[int64]string {
			m := make(map[int64]string)
			for k, v := range intStrMapVal.V.(map[*value.Value]*value.Value) {
				kv := k.V.(*big.Int).Int64()
				vv := v.V.(string)
				m[kv] = vv
			}
			return m
		}()
		assert.Equal(t, map[int64]string{1: "a", 2: "b", 3: "c"}, intStrMapActual)
	}

	t.Log("Constants:", evaluator.Constants)
	t.Log("Variables:", evaluator.Variables)
	t.Log("ValueStack:", evaluator.ValueStack)
}
