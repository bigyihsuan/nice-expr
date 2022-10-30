package parser

import (
	"fmt"
	"nice-expr/lexer/token"
	"strings"
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

type Literal interface{}
type PrimitiveLiteral struct {
	Literal
	TerminalNode
}

func (pl PrimitiveLiteral) String() string {
	return fmt.Sprint(pl.Token.Lexeme)
}

type CompoundLiteral interface {
	Literal
}

type ListLiteral struct {
	Node
	CompoundLiteral
	Values []Literal // TODO: Expand this to allow expressions
}

func (ll ListLiteral) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for _, e := range ll.Values {
		b.WriteString(fmt.Sprint(e))
		b.WriteRune(',')
	}
	b.WriteRune(']')
	return b.String()
}

type MapLiteral struct {
	Node
	CompoundLiteral
	Values map[Literal]Literal // TODO: Expand this to allow expressions
}

func (ll MapLiteral) String() string {
	var b strings.Builder
	b.WriteString("<|")
	for k, v := range ll.Values {
		b.WriteString(fmt.Sprint(k))
		b.WriteRune(':')
		b.WriteString(fmt.Sprint(v))
		b.WriteRune(',')
	}
	b.WriteString("|>")
	return b.String()
}

type Identifier struct {
	Node
	Name *token.Token
}

func (id Identifier) String() string {
	return id.Name.Lexeme
}

// AssignmentExpr := Name is Value
type AssignmentExpr struct {
	Node
	Name  *Identifier
	Op    *token.Token
	Value Node
}

func (ae AssignmentExpr) String() string {
	return fmt.Sprintf("set %v %v %v", ae.Name, ae.Op.Lexeme, ae.Value)
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

func (ae VariableDeclarationExpr) String() string {
	return fmt.Sprintf("var %v %v is %v", ae.Name, ae.Type, ae.Value)
}

// ConstantExpr := const Name Type is Value
type ConstantDeclarationExpr struct {
	*DeclarationExpr
}

func (ae ConstantDeclarationExpr) String() string {
	return fmt.Sprintf("const %v %v is %v", ae.Name, ae.Type, ae.Value)
}

// --- TYPES --- //
// ------------- //
// TypeExpr := PrimitiveTypeExpr | ListTypeExpr | MapTypeExpr
type TypeExpr interface {
	Node
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
