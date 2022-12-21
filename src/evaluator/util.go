package evaluator

import (
	"fmt"
	"nice-expr/src/ast"
)

var (
	BuiltinFunctionNames = []string{"print", "println", "len"}
)

type IdentifierEntry[T any] struct {
	Ident   *ast.Identifier
	Value   T
	VarType VariableType
}

func (i IdentifierEntry[T]) String() string {
	return fmt.Sprintf("{%s %v %s}", i.Ident, i.Value, i.VarType)
}

//go:generate stringer -type=VariableType
type VariableType int

const (
	Invalid VariableType = iota
	Var
	Const
	Func
)
