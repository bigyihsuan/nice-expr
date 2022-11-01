package ast

import (
	"fmt"
	"nice-expr/token"
	"nice-expr/value"
	"strings"
)

type Program struct {
	Statements []Expr
}

type Node interface{}

type Expr interface{ Node }

type BinaryExpr struct {
	Expr
	Left, Right Expr
	Op          *token.Token
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("(%v)%s(%v)", e.Left, e.Op.Lexeme, e.Right)
}

type UnaryExpr struct {
	Expr
	Op    *token.Token
	Right Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("%s(%v)", e.Op.Lexeme, e.Right)
}

type Literal interface{ Expr }

type PrimitiveLiteral struct {
	Literal
	Token *token.Token
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
	Values []Expr
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
	Values map[Expr]Expr
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
type Declaration interface {
	Expr
}

// Variable := var Name Type is Value
type VariableDeclaration struct {
	Declaration
	Name  *Identifier
	Type  Type
	Value Expr
}

func (ae VariableDeclaration) String() string {
	return fmt.Sprintf("var %v %v is (%v)", ae.Name, ae.Type, ae.Value)
}

// Constant := const Name Type is Value
type ConstantDeclaration struct {
	Declaration
	Name  *Identifier
	Type  Type
	Value Expr
}

func (ae ConstantDeclaration) String() string {
	return fmt.Sprintf("const %v %v is (%v)", ae.Name, ae.Type, ae.Value)
}

// --- TYPES --- //
// ------------- //
// Type := PrimitiveType | ListType | MapType | FuncType
type Type interface {
	Expr
	ToValueType() value.ValueType
}

// PrimitiveType := Name
type PrimitiveType struct {
	Type
	Name *token.Token
}

func (t PrimitiveType) String() string {
	return fmt.Sprint(t.Name.Tt)
}
func (t PrimitiveType) ToValueType() value.ValueType {
	var valType value.ValueType
	valType.Name = t.Name.Tt.String()
	return valType
}

// ListType := list '[' Type ']'
type ListType struct {
	Type
	ValueType Type
}

func (t ListType) String() string {
	return fmt.Sprintf("List[%s]", t.ValueType)
}
func (t ListType) ToValueType() value.ValueType {
	valType := value.NewValueType("List")
	valType.AddTypeArg(t.ValueType.ToValueType())
	return valType
}

// MapType := map '[' Type ']' Type
type MapType struct {
	Type
	KeyType   Type
	ValueType Type
}

func (t MapType) String() string {
	return fmt.Sprintf("Map[%s]%s", t.KeyType, t.ValueType)
}
func (t MapType) ToValueType() value.ValueType {
	valType := value.NewValueType("Map")
	valType.TypeArgs = append(valType.TypeArgs, t.KeyType.ToValueType(), t.ValueType.ToValueType())
	return valType
}

type FuncType struct {
	Type
	InputTypes []Type
	OutputType Type
}

func (t FuncType) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for _, e := range t.InputTypes {
		b.WriteString(fmt.Sprint(e))
		b.WriteRune(',')
	}
	b.WriteRune(']')
	return fmt.Sprintf("func(%s)%s", b.String(), t.OutputType)
}
