package ast

type Visitor interface {
	UnaryExpr(v Visitor, e *UnaryExpr)
	BinaryExpr(v Visitor, e *BinaryExpr)
	// program
	Program(v Visitor, p *Program)
	Expr(v Visitor, p Expr)
	// "statements"
	Declaration(v Visitor, d Declaration)
	VariableDeclaration(v Visitor, s *VariableDeclaration)
	ConstantDeclaration(v Visitor, s *ConstantDeclaration)
	Assignment(v Visitor, s *Assignment)
	Return(v Visitor, r *Return)
	Block(v Visitor, b *Block)
	// tests
	Test(v Visitor, t Test)
	AndTest(v Visitor, t *AndTest)
	OrTest(v Visitor, t *OrTest)
	NotTest(v Visitor, t *NotTest)
	// comparisons
	Comparison(v Visitor, c Comparison)
	Equal(v Visitor, c *Equal)
	Greater(v Visitor, c *Greater)
	Less(v Visitor, c *Less)
	GreaterEqual(v Visitor, c *GreaterEqual)
	LessEqual(v Visitor, c *LessEqual)
	// arithmetic
	AddExpr(v Visitor, a AddExpr)
	Add(v Visitor, a *Add)
	Sub(v Visitor, a *Sub)
	MulExpr(v Visitor, m MulExpr)
	Mul(v Visitor, m *Mul)
	Div(v Visitor, m *Div)
	Mod(v Visitor, m *Mod)
	UnaryMinus(v Visitor, t *UnaryMinus)
	// primaries
	Primary(v Visitor, p Primary)
	Identifier(v Visitor, i *Identifier)
	FunctionCall(v Visitor, f *FunctionCall)
	Literal(v Visitor, l Literal)
	PrimitiveLiteral(v Visitor, l *PrimitiveLiteral)
	CompoundLiteral(v Visitor, l CompoundLiteral)
	ListLiteral(v Visitor, l *ListLiteral)
	MapLiteral(v Visitor, l *MapLiteral)
	// other exprs
	Indexing(v Visitor, i *Indexing)
	// types
	Type(v Visitor, t Type)
	PrimitiveType(v Visitor, t *PrimitiveType)
	ListType(v Visitor, t *ListType)
	MapType(v Visitor, t *MapType)
}

type DefaultVisitor struct{}

/* DEFAULT IMPLEMENTATIONS */

func (v *DefaultVisitor) UnaryExpr(o Visitor, e *UnaryExpr) {
	e.Right.Accept(v)
}
func (v *DefaultVisitor) BinaryExpr(o Visitor, e *BinaryExpr) {
	e.Left.Accept(v)
	e.Right.Accept(v)
}

func (v *DefaultVisitor) Program(o Visitor, p *Program) {
	for _, e := range p.Statements {
		e.Accept(v)
	}
}
func (v *DefaultVisitor) Expr(o Visitor, p Expr) {
	switch p := p.(type) {
	case *Indexing:
		p.Accept(v)
	case *Assignment:
		p.Accept(v)
	case Declaration:
		p.Accept(v)
	case Test:
		p.AcceptTest(v)
	}
}

func (v *DefaultVisitor) Declaration(o Visitor, d Declaration) {
	switch d := d.(type) {
	case *VariableDeclaration:
		o.VariableDeclaration(o, d)
	case *ConstantDeclaration:
		o.ConstantDeclaration(o, d)
	}
}
func (v *DefaultVisitor) Test(o Visitor, t Test) {
	switch t := t.(type) {
	case *AndTest:
		o.AndTest(o, t)
	case *OrTest:
		o.OrTest(o, t)
	}
}
func (v *DefaultVisitor) Comparison(o Visitor, c Comparison) {
	c.AcceptCompare(o)
}
func (v *DefaultVisitor) AddExpr(o Visitor, a AddExpr) {
	a.AcceptAddExpr(o)
}
func (v *DefaultVisitor) MulExpr(o Visitor, m MulExpr) {
	m.AcceptMulExpr(o)
}
func (v *DefaultVisitor) Primary(o Visitor, p Primary) {
	switch p := p.(type) {
	case *Identifier:
		p.Accept(o)
	case *FunctionCall:
		p.Accept(o)
	case Literal:
		p.Accept(o)
	}
}
func (v *DefaultVisitor) Literal(o Visitor, l Literal) {
	switch l := l.(type) {
	case *PrimitiveLiteral:
		l.Accept(o)
	case CompoundLiteral:
		l.AcceptCompoundLiteral(o)
	}
}
func (v *DefaultVisitor) CompoundLiteral(o Visitor, l CompoundLiteral) {
	l.AcceptCompoundLiteral(o)
}
func (v *DefaultVisitor) Type(o Visitor, t Type) {
	switch t := t.(type) {
	case *PrimitiveType:
		v.PrimitiveType(o, t)
	case *ListType:
		v.ListType(o, t)
	case *MapType:
		v.MapType(o, t)
	}
}

/* METHODS TO IMPLEMENT */

func (v *DefaultVisitor) VariableDeclaration(_ Visitor, s *VariableDeclaration) {}
func (v *DefaultVisitor) ConstantDeclaration(_ Visitor, s *ConstantDeclaration) {}
func (v *DefaultVisitor) Assignment(_ Visitor, s *Assignment)                   {}
func (v *DefaultVisitor) Return(_ Visitor, r *Return)                           {}
func (v *DefaultVisitor) Block(_ Visitor, b *Block)                             {}
func (v *DefaultVisitor) AndTest(_ Visitor, t *AndTest)                         {}
func (v *DefaultVisitor) OrTest(_ Visitor, t *OrTest)                           {}
func (v *DefaultVisitor) NotTest(_ Visitor, t *NotTest)                         {}
func (v *DefaultVisitor) Equal(_ Visitor, c *Equal)                             {}
func (v *DefaultVisitor) Greater(_ Visitor, c *Greater)                         {}
func (v *DefaultVisitor) Less(_ Visitor, c *Less)                               {}
func (v *DefaultVisitor) GreaterEqual(_ Visitor, c *GreaterEqual)               {}
func (v *DefaultVisitor) LessEqual(_ Visitor, c *LessEqual)                     {}
func (v *DefaultVisitor) Add(_ Visitor, a *Add)                                 {}
func (v *DefaultVisitor) Sub(_ Visitor, a *Sub)                                 {}
func (v *DefaultVisitor) Mul(_ Visitor, m *Mul)                                 {}
func (v *DefaultVisitor) Div(_ Visitor, m *Div)                                 {}
func (v *DefaultVisitor) Mod(_ Visitor, m *Mod)                                 {}
func (v *DefaultVisitor) UnaryMinus(_ Visitor, t *UnaryMinus)                   {}
func (v *DefaultVisitor) Identifier(_ Visitor, i *Identifier)                   {}
func (v *DefaultVisitor) FunctionCall(_ Visitor, f *FunctionCall)               {}
func (v *DefaultVisitor) PrimitiveLiteral(_ Visitor, l *PrimitiveLiteral)       {}
func (v *DefaultVisitor) ListLiteral(_ Visitor, l *ListLiteral)                 {}
func (v *DefaultVisitor) MapLiteral(_ Visitor, l *MapLiteral)                   {}
func (v *DefaultVisitor) Indexing(_ Visitor, i *Indexing)                       {}
func (v *DefaultVisitor) PrimitiveType(_ Visitor, t *PrimitiveType)             {}
func (v *DefaultVisitor) ListType(_ Visitor, t *ListType)                       {}
func (v *DefaultVisitor) MapType(_ Visitor, t *MapType)                         {}
