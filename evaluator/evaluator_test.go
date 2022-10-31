package evaluator_test

import (
	"bytes"
	"math/big"
	"nice-expr/ast"
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

func TestEvaluateNestedDeclaration(t *testing.T) {
	fileName := "./../test/nested-declarations.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	// t.Log(tokens)

	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.ParseProgram()
	if pe != nil {
		t.Fatal(pe)
	}
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}
	// t.Log(program)

	evaluator := evaluator.NewEvaluator()
	ee := evaluator.EvaluateProgram(program)
	if ee != nil {
		t.Fatal(ee)
	}

	var x, y, z *ast.Identifier
	var xv, yv, zv *value.Value
	var xi, yi, zi int64

	x, xv = evaluator.GetVariable("x")
	assert.NotNil(t, x)
	xNotNil := assert.NotNil(t, xv)
	if xNotNil {
		xi = xv.V.(*big.Int).Int64()
		assert.Equal(t, int64(123), xi)
	}

	y, yv = evaluator.GetVariable("y")
	assert.NotNil(t, y)
	yNotNil := assert.NotNil(t, yv)
	if yNotNil {
		yi = xv.V.(*big.Int).Int64()
		assert.Equal(t, int64(123), yi)
	}

	z, zv = evaluator.GetVariable("z")
	assert.NotNil(t, z)
	zNotNil := assert.NotNil(t, zv)
	if zNotNil {
		zi = xv.V.(*big.Int).Int64()
		assert.Equal(t, int64(123), zi)
	}

	if xNotNil && yNotNil && zNotNil {
		assert.True(t, assert.ObjectsAreEqual(xi, yi) && assert.ObjectsAreEqual(yi, zi) && assert.ObjectsAreEqual(xi, zi))
	}

	var a, b, c *ast.Identifier
	var av, bv, cv *value.Value
	var ai, bi, ci string

	a, av = evaluator.GetVariable("a")
	assert.NotNil(t, a)
	aNotNil := assert.NotNil(t, av)
	if aNotNil {
		ai = av.V.(string)
		assert.Equal(t, "nesting is fun", ai)
	}

	b, bv = evaluator.GetConstant("b")
	assert.NotNil(t, b)
	bNotNil := assert.NotNil(t, bv)
	if bNotNil {
		bi = av.V.(string)
		assert.Equal(t, "nesting is fun", bi)
	}

	c, cv = evaluator.GetVariable("c")
	assert.NotNil(t, c)
	cNotNil := assert.NotNil(t, cv)
	if cNotNil {
		ci = av.V.(string)
		assert.Equal(t, "nesting is fun", ci)
	}

	if aNotNil && bNotNil && cNotNil {
		assert.True(t, assert.ObjectsAreEqual(ai, bi) && assert.ObjectsAreEqual(bi, ci) && assert.ObjectsAreEqual(ai, ci))
	}

	t.Log("Constants:", evaluator.Constants)
	t.Log("Variables:", evaluator.Variables)
	t.Log("ValueStack:", evaluator.ValueStack)
}
