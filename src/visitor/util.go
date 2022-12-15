package visitor

import "nice-expr/src/ast"

type IdentifierEntry[T any] struct {
	Ident *ast.Identifier
	Value T
}
