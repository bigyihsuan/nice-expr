package visitor

import (
	"fmt"
	"nice-expr/src/ast"
	"nice-expr/src/util"
	"strings"
)

type StringVisitor struct {
	ast.DefaultVisitor
	stringStack *util.Stack[string]
}

func NewStringVisitor() *StringVisitor {
	sv := new(StringVisitor)
	sv.stringStack = new(util.Stack[string])
	return sv
}

func (v StringVisitor) String() string {
	// fmt.Println(v.stringStack.Len())
	s := *v.stringStack.Peek()
	return s
}

func (v *StringVisitor) UnaryExpr(_ ast.Visitor, e *ast.UnaryExpr) {
	e.Right.Accept(v)
}
func (v *StringVisitor) BinaryExpr(_ ast.Visitor, e *ast.BinaryExpr) {
	e.Left.Accept(v)
	e.Right.Accept(v)
}
func (v *StringVisitor) Program(_ ast.Visitor, p *ast.Program) {
	stmts := []string{}
	for _, e := range p.Statements {
		e.Accept(v)
		str, err := v.stringStack.Pop()
		if err != nil {
			panic(err)
		}
		stmts = append(stmts, str)
	}
	v.stringStack.Push("(" + strings.Join(stmts, " ") + ")")
}
func (v *StringVisitor) Expr(_ ast.Visitor, p ast.Expr) {
	switch p := p.(type) {
	case *ast.Indexing:
		p.Accept(v)
	case *ast.Assignment:
		p.Accept(v)
	case ast.Declaration:
		p.Accept(v)
	case ast.Test:
		p.AcceptTest(v)
	}
}

func (v *StringVisitor) Declaration(_ ast.Visitor, d ast.Declaration) {
	switch d := d.(type) {
	case *ast.VariableDeclaration:
		v.VariableDeclaration(v, d)
	case *ast.ConstantDeclaration:
		v.ConstantDeclaration(v, d)
	}
}
func (v *StringVisitor) Test(_ ast.Visitor, t ast.Test) {
	switch t := t.(type) {
	case *ast.AndTest:
		v.AndTest(v, t)
	case *ast.OrTest:
		v.OrTest(v, t)
	}
}
func (v *StringVisitor) Comparison(_ ast.Visitor, c ast.Comparison) {
	c.AcceptCompare(v)
}
func (v *StringVisitor) AddExpr(_ ast.Visitor, a ast.AddExpr) {
	a.AcceptAddExpr(v)
}
func (v *StringVisitor) MulExpr(_ ast.Visitor, m ast.MulExpr) {
	m.AcceptMulExpr(v)
}
func (v *StringVisitor) Primary(_ ast.Visitor, p ast.Primary) {
	switch p := p.(type) {
	case *ast.Identifier:
		p.Accept(v)
	case *ast.FunctionCall:
		p.Accept(v)
	case ast.Literal:
		p.Accept(v)
	}
}
func (v *StringVisitor) Literal(_ ast.Visitor, l ast.Literal) {
	switch l := l.(type) {
	case *ast.PrimitiveLiteral:
		l.Accept(v)
	case ast.CompoundLiteral:
		l.AcceptCompoundLiteral(v)
	}
}
func (v *StringVisitor) CompoundLiteral(_ ast.Visitor, l ast.CompoundLiteral) {
	l.AcceptCompoundLiteral(v)
}
func (v *StringVisitor) Type(_ ast.Visitor, t ast.Type) {
	switch t := t.(type) {
	case *ast.PrimitiveType:
		v.PrimitiveType(v, t)
	case *ast.ListType:
		v.ListType(v, t)
	case *ast.MapType:
		v.MapType(v, t)
	}
}

func (v *StringVisitor) VariableDeclaration(_ ast.Visitor, s *ast.VariableDeclaration) {
	s.Name.Accept(v)
	name, _ := v.stringStack.Pop()
	s.Type.Accept(v)
	typ, _ := v.stringStack.Pop()
	s.Value.Accept(v)
	expr, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(var %s is %s %s)", name, typ, expr)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) ConstantDeclaration(_ ast.Visitor, s *ast.ConstantDeclaration) {
	s.Name.Accept(v)
	name, _ := v.stringStack.Pop()
	s.Type.Accept(v)
	typ, _ := v.stringStack.Pop()
	s.Value.Accept(v)
	expr, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(const %s is %s %s)", name, typ, expr)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Assignment(_ ast.Visitor, s *ast.Assignment) {
	s.Name.Accept(v)
	name, _ := v.stringStack.Pop()
	s.Value.Accept(v)
	expr, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(set %s is %s)", name, expr))
}
func (v *StringVisitor) Return(_ ast.Visitor, r *ast.Return) {
	r.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	v.stringStack.Push(fmt.Sprintf("(return %s)", right))
}
func (v *StringVisitor) Block(_ ast.Visitor, b *ast.Block) {
	stmts := []string{}
	for _, e := range b.Statements {
		e.Accept(v)
		str, err := v.stringStack.Pop()
		if err != nil {
			panic(err)
		}
		stmts = append(stmts, str)
	}
	v.stringStack.Push("{" + strings.Join(stmts, " ") + "}")
}

func (v *StringVisitor) AndTest(_ ast.Visitor, t *ast.AndTest) {
	t.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(and %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) OrTest(_ ast.Visitor, t *ast.OrTest) {
	t.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(or %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) NotTest(_ ast.Visitor, t *ast.NotTest) {
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(not %s)", right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Equal(_ ast.Visitor, c *ast.Equal) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(= %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Greater(_ ast.Visitor, c *ast.Greater) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(> %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Less(_ ast.Visitor, c *ast.Less) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(< %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) GreaterEqual(_ ast.Visitor, c *ast.GreaterEqual) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(>= %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) LessEqual(_ ast.Visitor, c *ast.LessEqual) {
	c.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	c.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(<= %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Add(_ ast.Visitor, a *ast.Add) {
	a.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	a.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(+ %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Sub(_ ast.Visitor, a *ast.Sub) {
	a.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	a.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(- %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Mul(_ ast.Visitor, m *ast.Mul) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(* %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Div(_ ast.Visitor, m *ast.Div) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(/ %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Mod(_ ast.Visitor, m *ast.Mod) {
	m.Left.Accept(v)
	left, _ := v.stringStack.Pop()
	m.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(%% %s %s)", left, right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) UnaryMinus(_ ast.Visitor, t *ast.UnaryMinus) {
	t.Right.Accept(v)
	right, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(- %s)", right)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Identifier(_ ast.Visitor, i *ast.Identifier) {
	str := i.Tok.Lexeme
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) FunctionCall(_ ast.Visitor, f *ast.FunctionCall) {
	f.Ident.Accept(v)
	name, _ := v.stringStack.Pop()
	args := []string{}
	for _, a := range f.Arguments {
		a.Accept(v)
		arg, _ := v.stringStack.Pop()
		args = append(args, arg)
	}
	str := fmt.Sprintf("(%s %s)", name, strings.Join(args, " "))
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) PrimitiveLiteral(_ ast.Visitor, l *ast.PrimitiveLiteral) {
	v.stringStack.Push(l.String())
}
func (v *StringVisitor) ListLiteral(_ ast.Visitor, l *ast.ListLiteral) {
	eles := []string{}
	for _, expr := range l.Values {
		expr.Accept(v)
		e, _ := v.stringStack.Pop()
		eles = append(eles, e)
	}
	str := fmt.Sprintf("[%s]", strings.Join(eles, " "))
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) MapLiteral(_ ast.Visitor, l *ast.MapLiteral) {
	eles := []string{}
	for key, value := range l.Values {
		key.Accept(v)
		k, _ := v.stringStack.Pop()
		value.Accept(v)
		val, _ := v.stringStack.Pop()
		eles = append(eles, fmt.Sprintf("%s:%s", k, val))
	}
	str := fmt.Sprintf("<|%s|>", strings.Join(eles, " "))
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) Indexing(_ ast.Visitor, i *ast.Indexing) {
	i.Left.Accept(v)
	collection, _ := v.stringStack.Pop()
	i.Right.Accept(v)
	index, _ := v.stringStack.Pop()
	str := fmt.Sprintf("(_ %s %s)", collection, index)
	// fmt.Println(str)
	v.stringStack.Push(str)
}
func (v *StringVisitor) PrimitiveType(_ ast.Visitor, t *ast.PrimitiveType) {
	v.stringStack.Push(t.String())
}
func (v *StringVisitor) ListType(_ ast.Visitor, t *ast.ListType) {
	v.stringStack.Push(t.String())
}
func (v *StringVisitor) MapType(_ ast.Visitor, t *ast.MapType) {
	v.stringStack.Push(t.String())
}
