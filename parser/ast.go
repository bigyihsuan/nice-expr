package parser

import (
	"fmt"
	"nice-expr/lexer/token"
)

type Node interface{}

type TerminalNode struct {
	Node
	Token *token.Token
}
type UnaryNode struct {
	Node
	Child Node
}

type PrimitiveExpr struct {
	TerminalNode
}

type Identifier struct {
	Node
	Name *token.Token
}

// AssignmentExpr := Name is Value
type AssignmentExpr struct {
	Node
	Name  *Identifier
	Value Node
}

// DeclarationExpr := Name Type is Value
type DeclarationExpr struct {
	Node
	Name  *Identifier
	Type  TypeExpr
	Value Node
}

// VariableExpr := var Name Type is Value
type VariableDeclarationExpr struct {
	*DeclarationExpr
}

// ConstantExpr := const Name Type is Value
type ConstantDeclarationExpr struct {
	*DeclarationExpr
}

// --- TYPES --- //
// ------------- //
// TypeExpr := PrimitiveTypeExpr | ListTypeExpr | MapTypeExpr
type TypeExpr interface {
	Node
	fmt.Stringer
}

// PrimitiveTypeExpr := Name
type PrimitiveTypeExpr struct {
	Node
	Name *token.Token
}

func (t PrimitiveTypeExpr) String() string {
	return fmt.Sprint(t.Name.Tt)
}

// ListTypeExpr := list '[' TypeExpr ']'
type ListTypeExpr struct {
	Node
	ValueType TypeExpr
}

func (t ListTypeExpr) String() string {
	return fmt.Sprintf("list[%s]", t.ValueType)
}

// MapTypeExpr := map '[' TypeExpr ']' TypeExpr
type MapTypeExpr struct {
	Node
	KeyType   TypeExpr
	ValueType TypeExpr
}

func (t MapTypeExpr) String() string {
	return fmt.Sprintf("map[%s]%s", t.KeyType, t.ValueType)
}
