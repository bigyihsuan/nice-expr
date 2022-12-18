package evaluator

import "nice-expr/src/ast"

var (
	BuiltinFunctionNames = []string{"print", "println", "len"}
)

type IdentifierEntry[T any] struct {
	Ident   *ast.Identifier
	Value   T
	VarType VariableType
}

//go:generate stringer -type=VariableType
type VariableType int

const (
	Invalid VariableType = iota
	Var
	Const
	Func
)
