package ast

import (
	"fmt"
	"nice-expr/src/token"
	"nice-expr/src/value"
	"strings"
)

type Visitable interface {
	Accept(Visitor)
}
type Expr interface{ Visitable }

type Program struct {
	Statements []Expr
}

func (p *Program) Accept(v Visitor) {
	v.Program(v, p)
}

type UnaryExpr struct {
	Expr
	Right Expr
}

func (e *UnaryExpr) Accept(v Visitor) {
	v.UnaryExpr(v, e)
}

type BinaryExpr struct {
	Expr
	Left, Right Expr
}

func (e *BinaryExpr) Accept(v Visitor) {
	v.BinaryExpr(v, e)
}

// "statements"

type Assignment struct {
	Name  *Identifier
	Op    *token.Token
	Value Expr
}

func (s *Assignment) Accept(v Visitor) {
	v.Assignment(v, s)
}

func (ae Assignment) String() string {
	return fmt.Sprintf("(set (%v) (%v) (%v))", ae.Name, ae.Op.Lexeme, ae.Value)
}

type Declaration interface{ Expr }

type VariableDeclaration struct {
	Declaration
	Name  *Identifier
	Type  Type
	Value Expr
}

func (s *VariableDeclaration) Accept(v Visitor) {
	v.VariableDeclaration(v, s)
}

func (ae VariableDeclaration) String() string {
	return fmt.Sprintf("(var (%v) (%v) is (%v))", ae.Name, ae.Type, ae.Value)
}

type ConstantDeclaration struct {
	Declaration
	Name  *Identifier
	Type  Type
	Value Expr
}

func (s *ConstantDeclaration) Accept(v Visitor) {
	v.ConstantDeclaration(v, s)
}

func (ae ConstantDeclaration) String() string {
	return fmt.Sprintf("(const (%v) (%v) is (%v))", ae.Name, ae.Type, ae.Value)
}

// tests

type Test interface {
	AcceptTest(v Visitor)
}

type AndTest struct {
	BinaryExpr
}

func (t *AndTest) Accept(v Visitor) {
	v.AndTest(v, t)
}
func (t *AndTest) AcceptTest(v Visitor) {
	v.AndTest(v, t)
}

type OrTest struct {
	BinaryExpr
}

func (t *OrTest) Accept(v Visitor) {
	v.OrTest(v, t)
}
func (t *OrTest) AcceptTest(v Visitor) {
	v.OrTest(v, t)
}

type NotTest struct {
	UnaryExpr
}

func (t *NotTest) Accept(v Visitor) {
	v.NotTest(v, t)
}

func (t NotTest) String() string {
	return fmt.Sprintf("-%s", t.Right)
}

// comparisons

type Comparison interface {
	AcceptCompare(v Visitor)
}

type Equal struct {
	BinaryExpr
}

func (c *Equal) AcceptCompare(v Visitor) {
	v.Equal(v, c)
}
func (c *Equal) Accept(v Visitor) {
	v.Equal(v, c)
}

type Greater struct {
	BinaryExpr
}

func (c *Greater) AcceptCompare(v Visitor) {
	v.Greater(v, c)
}
func (c *Greater) Accept(v Visitor) {
	v.Greater(v, c)
}

type Less struct {
	BinaryExpr
}

func (c *Less) AcceptCompare(v Visitor) {
	v.Less(v, c)
}
func (c *Less) Accept(v Visitor) {
	v.Less(v, c)
}

type GreaterEqual struct {
	BinaryExpr
}

func (c *GreaterEqual) AcceptCompare(v Visitor) {
	v.GreaterEqual(v, c)
}
func (c *GreaterEqual) Accept(v Visitor) {
	v.GreaterEqual(v, c)
}

type LessEqual struct {
	BinaryExpr
}

func (c *LessEqual) AcceptCompare(v Visitor) {
	v.LessEqual(v, c)
}
func (c *LessEqual) Accept(v Visitor) {
	v.LessEqual(v, c)
}

// arithmetic

type AddExpr interface {
	AcceptAddExpr(v Visitor)
}

type Add struct {
	AddExpr
	BinaryExpr
}

func (a *Add) Accept(v Visitor) {
	v.Add(v, a)
}
func (a *Add) AcceptAddExpr(v Visitor) {
	v.Add(v, a)
}

type Sub struct {
	AddExpr
	BinaryExpr
}

func (a *Sub) Accept(v Visitor) {
	v.Sub(v, a)
}
func (a *Sub) AcceptAddExpr(v Visitor) {
	v.Sub(v, a)
}

type MulExpr interface {
	AcceptMulExpr(v Visitor)
}

type Mul struct {
	MulExpr
	BinaryExpr
}

func (m *Mul) Accept(v Visitor) {
	v.Mul(v, m)
}
func (m *Mul) AcceptMulExpr(v Visitor) {
	v.Mul(v, m)
}

type Div struct {
	MulExpr
	BinaryExpr
}

func (m *Div) Accept(v Visitor) {
	v.Div(v, m)
}
func (m *Div) AcceptMulExpr(v Visitor) {
	v.Div(v, m)
}

type Mod struct {
	MulExpr
	BinaryExpr
}

func (m *Mod) Accept(v Visitor) {
	v.Mod(v, m)
}
func (m *Mod) AcceptMulExpr(v Visitor) {
	v.Mod(v, m)
}

type UnaryMinus struct {
	UnaryExpr
}

func (u *UnaryMinus) Accept(v Visitor) {
	v.UnaryMinus(v, u)
}

// primaries
type Primary interface{ Expr }

type Literal interface{ Primary }

type PrimitiveLiteral struct {
	Literal
	Token *token.Token
}

func (l *PrimitiveLiteral) Accept(v Visitor) {
	v.PrimitiveLiteral(v, l)
}

func (pl PrimitiveLiteral) String() string {
	return fmt.Sprint(pl.Token.Lexeme)
}

type CompoundLiteral interface {
	Literal
	AcceptCompoundLiteral(v Visitor)
}

type ListLiteral struct {
	CompoundLiteral
	Values []Expr
}

func (l *ListLiteral) Accept(v Visitor) {
	v.ListLiteral(v, l)
}
func (l *ListLiteral) AcceptCompoundLiteral(v Visitor) {
	v.ListLiteral(v, l)
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
	CompoundLiteral
	Values map[Expr]Expr
}

func (l *MapLiteral) Accept(v Visitor) {
	v.MapLiteral(v, l)
}
func (l *MapLiteral) AcceptCompoundLiteral(v Visitor) {
	v.MapLiteral(v, l)
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
	Primary
	Tok *token.Token
}

func (i *Identifier) Accept(v Visitor) {
	v.Identifier(v, i)
}

func (id Identifier) String() string {
	return fmt.Sprintf("(Identifier %s)", id.Tok.Lexeme)
}

func (id Identifier) Name() string {
	return id.Tok.Lexeme
}

type FunctionCall struct {
	Primary
	Ident     *Identifier
	Arguments []Expr
}

func (f *FunctionCall) Accept(v Visitor) {
	v.FunctionCall(v, f)
}

func (fn FunctionCall) String() string {
	return fmt.Sprintf("(%s (%s))", fn.Ident, fn.Arguments)
}

// other exprs

type Indexing struct {
	BinaryExpr
}

func (i *Indexing) Accept(v Visitor) {
	v.Indexing(v, i)
}

// types
type Type interface {
	Expr
	ToValueType() value.ValueType
}

// PrimitiveType := Name
type PrimitiveType struct {
	Name *token.Token
}

func (t *PrimitiveType) Accept(v Visitor) {
	v.PrimitiveType(v, t)
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

func (t *ListType) Accept(v Visitor) {
	v.ListType(v, t)
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

func (t *MapType) Accept(v Visitor) {
	v.MapType(v, t)
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
