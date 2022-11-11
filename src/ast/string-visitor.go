package ast

import (
	"fmt"
	"nice-expr/src/util"
	"strings"
)

type StringVisitor struct {
	DefaultVisitor
	stringStack *util.Stack[string]
}

func NewStringVisitor() *StringVisitor {
	sv := new(StringVisitor)
	sv.stringStack = new(util.Stack[string])
	return sv
}

func (v *StringVisitor) String() string {
	s, _ := v.stringStack.Pop()
	return s
}

// TODO: func (v *StringVisitor) Assignment(v Visitor, s *Assignment)

func (v *StringVisitor) VariableDeclaration(_ Visitor, s *VariableDeclaration) {
	s.Name.Accept(v)
	name, _ := v.stringStack.Pop()
	s.Type.Accept(v)
	typ, _ := v.stringStack.Pop()
	s.Value.Accept(v)
	expr, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(var %s is %s %s)", name, typ, expr))
}
func (v *StringVisitor) ConstantDeclaration(_ Visitor, s *ConstantDeclaration) {
	s.Name.Accept(v)
	name, _ := v.stringStack.Pop()
	s.Type.Accept(v)
	typ, _ := v.stringStack.Pop()
	s.Value.Accept(v)
	expr, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(const %s is %s %s)", name, typ, expr))
}

func (v *StringVisitor) AndTest(_ Visitor, t *AndTest) {
	t.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(and %s %s)", left, right))
}
func (v *StringVisitor) OrTest(_ Visitor, t *OrTest) {
	t.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(or %s %s)", left, right))
}
func (v *StringVisitor) NotTest(_ Visitor, t *NotTest) {
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(not %s)", right))
}
func (v *StringVisitor) Equal(_ Visitor, c *Equal) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(= %s %s)", left, right))
}
func (v *StringVisitor) Greater(_ Visitor, c *Greater) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(> %s %s)", left, right))
}
func (v *StringVisitor) Less(_ Visitor, c *Less) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(< %s %s)", left, right))
}
func (v *StringVisitor) GreaterEqual(_ Visitor, c *GreaterEqual) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(>= %s %s)", left, right))
}
func (v *StringVisitor) LessEqual(_ Visitor, c *LessEqual) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(<= %s %s)", left, right))
}
func (v *StringVisitor) Add(_ Visitor, a *Add) {
	a.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	a.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(+ %s %s)", left, right))
}
func (v *StringVisitor) Sub(_ Visitor, a *Sub) {
	a.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	a.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(- %s %s)", left, right))
}
func (v *StringVisitor) Mul(_ Visitor, m *Mul) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(* %s %s)", left, right))
}
func (v *StringVisitor) Div(_ Visitor, m *Div) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(/ %s %s)", left, right))
}
func (v *StringVisitor) Mod(_ Visitor, m *Mod) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(%% %s %s)", left, right))
}
func (v *StringVisitor) UnaryMinus(_ Visitor, t *UnaryMinus) {
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(- %s)", right))
}
func (v *StringVisitor) Identifier(_ Visitor, i *Identifier) {
	v.stringStack.Push(i.Name.Lexeme)
}
func (v *StringVisitor) FunctionCall(_ Visitor, f *FunctionCall) {
	f.Ident.Accept(v)
	name, _ := v.stringStack.Pop()
	args := []string{}
	for _, a := range f.Arguments {
		a.Accept(v)
		arg, _ := v.stringStack.Pop()
		args = append(args, arg)
	}
	v.stringStack.Push(fmt.Sprintf("(%s %s)", name, strings.Join(args, " ")))
}
func (v *StringVisitor) PrimitiveLiteral(_ Visitor, l *PrimitiveLiteral) {
	v.stringStack.Push(l.String())
}
func (v *StringVisitor) ListLiteral(_ Visitor, l *ListLiteral) {
	eles := []string{}
	for _, expr := range l.Values {
		expr.Accept(v)
		e, _ := v.stringStack.Pop()
		eles = append(eles, e)
	}
	v.stringStack.Push(fmt.Sprintf("[%s]", strings.Join(eles, " ")))
}
func (v *StringVisitor) MapLiteral(_ Visitor, l *MapLiteral) {
	eles := []string{}
	for key, value := range l.Values {
		key.Accept(v)
		k, _ := v.stringStack.Pop()
		value.Accept(v)
		val, _ := v.stringStack.Pop()
		eles = append(eles, fmt.Sprintf("%s:%s", k, val))
	}
	v.stringStack.Push(fmt.Sprintf("<|%s|>", strings.Join(eles, " ")))
}
func (v *StringVisitor) Indexing(_ Visitor, i *Indexing) {
	i.Left.Accept(v)
	collection, _ := v.stringStack.Pop()
	i.Right.Accept(v)
	index, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(_ %s %s)", collection, index))
}
func (v *StringVisitor) PrimitiveType(_ Visitor, t *PrimitiveType) {
	v.stringStack.Push(t.String())
}
func (v *StringVisitor) ListType(_ Visitor, t *ListType) {
	v.stringStack.Push(t.String())
}
func (v *StringVisitor) MapType(_ Visitor, t *MapType) {
	v.stringStack.Push(t.String())
}
