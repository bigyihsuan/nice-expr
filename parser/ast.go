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

// Assignment := Name is Value
type Assignment struct {
	Node
	Name  *Identifier
	Op    *token.Token
	Value Node
}

func (ae Assignment) String() string {
	return fmt.Sprintf("set %v %v %v", ae.Name, ae.Op.Lexeme, ae.Value)
}

// Declaration := Name Type is Value
type Declaration struct {
	Node
	Name  *Identifier
	Type  Type
	Value Node
}

// Variable := var Name Type is Value
type VariableDeclaration struct {
	*Declaration
}

func (ae VariableDeclaration) String() string {
	return fmt.Sprintf("var %v %v is %v", ae.Name, ae.Type, ae.Value)
}

// Constant := const Name Type is Value
type ConstantDeclaration struct {
	*Declaration
}

func (ae ConstantDeclaration) String() string {
	return fmt.Sprintf("const %v %v is %v", ae.Name, ae.Type, ae.Value)
}

// --- TYPES --- //
// ------------- //
// Type := PrimitiveType | ListType | MapType | FuncType
type Type interface {
	Node
}

// PrimitiveType := Name
type PrimitiveType struct {
	Type
	Name *token.Token
}

func (t PrimitiveType) String() string {
	return fmt.Sprint(t.Name.Tt)
}

// ListType := list '[' Type ']'
type ListType struct {
	Type
	ValueType Type
}

func (t ListType) String() string {
	return fmt.Sprintf("list[%s]", t.ValueType)
}

// MapType := map '[' Type ']' Type
type MapType struct {
	Type
	KeyType   Type
	ValueType Type
}

func (t MapType) String() string {
	return fmt.Sprintf("map[%s]%s", t.KeyType, t.ValueType)
}

type FuncType struct {
	Type
	InputTypes []Type
	OutputType Type
}

func (t FuncType) String() string {
	return fmt.Sprintf("func(%s)%s", t.InputTypes, t.OutputType)
}
