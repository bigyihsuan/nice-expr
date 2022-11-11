package evaluator_test

import (
	"bytes"
	"io"
	"math/big"
	"nice-expr/src/ast"
	"nice-expr/src/evaluator"
	"nice-expr/src/lexer"
	"nice-expr/src/parser"
	"nice-expr/src/util"
	"nice-expr/src/value"
	"os"
	"strings"
	"testing"

	"github.com/db47h/lex"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	rescueOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueOutput
	return string(out)
}

func TestEvaluateDeclarations(t *testing.T) {
	fileName := "./../../test/declarations.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.Program()
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}
	if pe != nil {
		t.Fatal(pe)
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
	var intStrMapActual map[int64]string
	assert.NotNil(t, intStrMap)
	if assert.NotNil(t, intStrMapVal) {
		intStrMapActual = func() map[int64]string {
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

	varOk, varOkVal := evaluator.GetVariable("ok")
	assert.NotNil(t, varOk)
	if assert.NotNil(t, varOkVal) {
		assert.Equal(t, true, varOkVal.V.(bool))
	}

	notOk, notOkVal := evaluator.GetConstant("notOk")
	assert.NotNil(t, notOk)
	if assert.NotNil(t, notOkVal) {
		assert.Equal(t, false, notOkVal.V.(bool))
	}

	// variable usage
	reuse, reuseVal := evaluator.GetVariable("reuse")
	assert.NotNil(t, reuse)
	if assert.NotNil(t, reuseVal) {
		assert.Equal(t, int64(10), reuseVal.V.(*big.Int).Int64())
		assert.Equal(t, xv.V.(*big.Int).Int64(), reuseVal.V.(*big.Int).Int64())
	}
	copied, copiedVal := evaluator.GetConstant("copied")
	assert.NotNil(t, copied)
	if assert.NotNil(t, copiedVal) {
		copiedActual := func() map[int64]string {
			m := make(map[int64]string)
			for k, v := range copiedVal.V.(map[*value.Value]*value.Value) {
				kv := k.V.(*big.Int).Int64()
				vv := v.V.(string)
				m[kv] = vv
			}
			return m
		}()
		assert.Equal(t, map[int64]string{1: "a", 2: "b", 3: "c"}, copiedActual)
		assert.Equal(t, intStrMapActual, copiedActual)
	}

	t.Log("Constants:", evaluator.Constants)
	t.Log("Variables:", evaluator.Variables)
	t.Log("ValueStack:", evaluator.ValueStack)
}

func TestEvaluateNestedDeclarations(t *testing.T) {
	fileName := "./../../test/nested-declarations.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	t.Log(tokens)

	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.Program()
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

func TestEvaluateUnaryMinus(t *testing.T) {
	cases := []util.TestCase{
		{Code: "89", Expected: int64(89), ExpectedType: value.IntType},
		{Code: "-89", Expected: int64(-89), ExpectedType: value.IntType},
		{Code: "12.34", Expected: float64(12.34), ExpectedType: value.DecType},
		{Code: "-12.34", Expected: float64(-12.34), ExpectedType: value.DecType},
		{Code: "-0", Expected: int64(0), ExpectedType: value.IntType},
	}
	for _, tc := range cases {
		file := lex.NewFile(tc.Code, strings.NewReader(tc.Code))
		nicerLexer := lexer.NewLexer(file)
		tokens := nicerLexer.LexAll()
		t.Log(tokens)
		nicerParser := parser.NewNiceExprParser(tokens)
		expr, perr := nicerParser.UnaryMinusExpr()
		if perr != nil {
			t.Fatal(perr)
		}
		if expr == nil {
			t.Fatal("parsed nil")
		}
		evaluator := evaluator.NewEvaluator()
		val, ee := evaluator.EvaluateUnaryMinusExpr(expr)
		if ee != nil {
			t.Fatal(ee)
		}
		switch {
		case tc.ExpectedType.Is(value.IntType):
			i, err := val.Int()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.Expected, i)
		case tc.ExpectedType.Is(value.DecType):
			f, err := val.Dec()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.Expected, f)
		}
	}
}

func TestEvaluateBuiltinFunctionsPrint(t *testing.T) {
	fileName := "./../../test/func-builtin-print.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	// t.Log(tokens)

	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.Program()
	if pe != nil {
		t.Fatal(pe)
	}
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}
	// t.Log(program)

	evaluator := evaluator.NewEvaluator()

	output := captureOutput(func() {
		ee := evaluator.EvaluateProgram(program)
		if ee != nil {
			t.Fatal(ee)
		}
	})

	expected := `Hello world!
this is a line
1234
56.78
true
false
[1,2,3,4,5,]
multiplestringsinasinglecall
each
on
a
new
line
`
	assert.Equal(t, expected, output)
}
func TestEvaluateBuiltinFunctionsLen(t *testing.T) {
	fileName := "./../../test/func-builtin-len.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	t.Log(tokens)

	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.Program()
	if pe != nil {
		t.Fatal(pe)
	}
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}
	t.Log(program)

	evaluator := evaluator.NewEvaluator()

	output := captureOutput(func() {
		ee := evaluator.EvaluateProgram(program)
		if ee != nil {
			t.Fatal(ee)
		}
	})

	expected := `3
2
5
`
	assert.Equal(t, expected, output)
}

func TestEvaluateBinaryOperators(t *testing.T) {
	fileName := "./../../test/binary.test.ne"
	test, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	file := lex.NewFile(fileName, bytes.NewReader(test))
	nicerLexer := lexer.NewLexer(file)
	tokens := nicerLexer.LexAll()
	// t.Log(tokens)

	nicerParser := parser.NewNiceExprParser(tokens)

	program, pe := nicerParser.Program()
	if pe != nil {
		t.Fatal(pe)
	}
	if len(program.Statements) <= 0 {
		t.Fatal("parsed nil")
	}
	t.Log(program)

	evaluator := evaluator.NewEvaluator()

	output := captureOutput(func() {
		ee := evaluator.EvaluateProgram(program)
		if ee != nil {
			t.Fatal(ee)
		}
	})

	expected := `2
	0
	6.6
	1
	1
	hello world
	heo
	hllo
	hlo
	[1,2,3,4,5,6,]
	[1,2,4,5,]
	true
	true
	true
	true
	true
	true
	true
	2
	5
	1
	0
	`
	assert.Equal(t, expected, output)
}
