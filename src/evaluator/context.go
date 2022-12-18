package evaluator

import (
	"nice-expr/src/ast"
)

type Context[T any] struct {
	Identifiers map[string]IdentifierEntry[T] // the variables that exist in this current context only
	// Functions map[string]IdentifierEntry[any]
	Parent *Context[T] // the context that contains this context, to allow for accessing outside names
}

// make a new context, with an optional parent. if no parent is given, the returned context's parent will be itself.
func NewContext[T any](parent ...*Context[T]) *Context[T] {
	context := new(Context[T])
	context.Identifiers = make(map[string]IdentifierEntry[T])
	if len(parent) > 0 {
		context.Parent = parent[0]
	} else {
		context.Parent = context
	}
	return context
}

// make a new copy of a context. if a parent is provided, the original is the parent of the returned context.
func CopyContext[T any](original *Context[T], parent ...*Context[T]) *Context[T] {
	context := NewContext(parent...)
	for k, v := range original.Identifiers {
		context.Identifiers[k] = v
	}
	return context
}

func (c *Context[T]) GetIdentifier(name string) (ident *ast.Identifier, value T, kind VariableType, sourceContext *Context[T]) {
	for n, entry := range c.Identifiers {
		if n == name {
			return entry.Ident, entry.Value, entry.VarType, c
		}
	}
	// not in this context, check the parent
	if c.Parent != c {
		return c.Parent.GetIdentifier(name)
	}
	return nil, value, Invalid, nil
}

func (c *Context[T]) AddIdentifier(name string, entry IdentifierEntry[T]) {
	c.Identifiers[name] = entry
}
