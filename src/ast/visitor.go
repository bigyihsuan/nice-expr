package ast

type Visitor interface {
	UnaryExpr(v Visitor, e *UnaryExpr)
	BinaryExpr(v Visitor, e *BinaryExpr)
	// program
	Program(v Visitor, p *Program)
	Expr(v Visitor, p Expr)
	// "statements"

	// TODO: Assignment(v Visitor, s *Assignment)

	Declaration(v Visitor, d Declaration)
	VariableDeclaration(v Visitor, s *VariableDeclaration)
	ConstantDeclaration(v Visitor, s *ConstantDeclaration)
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
	PrimitiveType(v Visitor, t *PrimitiveType)
	ListType(v Visitor, t *ListType)
	MapType(v Visitor, t *MapType)
}

type DefaultVisitor struct{}

func (v *DefaultVisitor) UnaryExpr(_ Visitor, e *UnaryExpr) {
	v.Expr(v, e.Right)
}
func (v *DefaultVisitor) BinaryExpr(_ Visitor, e *BinaryExpr) {
	v.Expr(v, e.Left)
	v.Expr(v, e.Right)
}

func (v *DefaultVisitor) Program(_ Visitor, p *Program) {
	for _, e := range p.Statements {
		v.Expr(v, e)
	}
}
func (v *DefaultVisitor) Expr(_ Visitor, p Expr) {
	switch p := p.(type) {
	case *Indexing:
		v.Indexing(v, p)
	case Declaration:
		v.Declaration(v, p)
		// TODO: case *Assignment:
		// v.Assignment(v,p)
	case Test:
		v.Test(v, p)
	}
}

// TODO: func (v *DefaultVisitor) Assignment(_ Visitor, s *Assignment) {}

func (v *DefaultVisitor) Declaration(_ Visitor, d Declaration) {
	switch d := d.(type) {
	case *VariableDeclaration:
		v.VariableDeclaration(v, d)
	case *ConstantDeclaration:
		v.ConstantDeclaration(v, d)
	}
}
func (v *DefaultVisitor) VariableDeclaration(_ Visitor, s *VariableDeclaration) {}
func (v *DefaultVisitor) ConstantDeclaration(_ Visitor, s *ConstantDeclaration) {}
func (v *DefaultVisitor) Test(_ Visitor, t Test) {
	switch t := t.(type) {
	case *AndTest:
		v.AndTest(v, t)
	case *OrTest:
		v.OrTest(v, t)
	}
}
func (v *DefaultVisitor) AndTest(_ Visitor, t *AndTest) {}
func (v *DefaultVisitor) OrTest(_ Visitor, t *OrTest)   {}
func (v *DefaultVisitor) NotTest(_ Visitor, t *NotTest) {}
func (v *DefaultVisitor) Comparison(_ Visitor, c Comparison) {
	switch c := c.(type) {
	case *Equal:
		v.Equal(v, c)
	case *Greater:
		v.Greater(v, c)
	case *Less:
		v.Less(v, c)
	case *GreaterEqual:
		v.GreaterEqual(v, c)
	case *LessEqual:
		v.LessEqual(v, c)
	}
}
func (v *DefaultVisitor) Equal(_ Visitor, c *Equal)               {}
func (v *DefaultVisitor) Greater(_ Visitor, c *Greater)           {}
func (v *DefaultVisitor) Less(_ Visitor, c *Less)                 {}
func (v *DefaultVisitor) GreaterEqual(_ Visitor, c *GreaterEqual) {}
func (v *DefaultVisitor) LessEqual(_ Visitor, c *LessEqual)       {}
func (v *DefaultVisitor) AddExpr(_ Visitor, a AddExpr) {
	switch a := a.(type) {
	case *Add:
		v.Add(v, a)
	case *Sub:
		v.Sub(v, a)
	}
}
func (v *DefaultVisitor) Add(_ Visitor, a *Add) {}
func (v *DefaultVisitor) Sub(_ Visitor, a *Sub) {}
func (v *DefaultVisitor) MulExpr(_ Visitor, m MulExpr) {
	switch m := m.(type) {
	case *Mul:
		v.Mul(v, m)
	case *Div:
		v.Div(v, m)
	case *Mod:
		v.Mod(v, m)
	}
}
func (v *DefaultVisitor) Mul(_ Visitor, m *Mul)               {}
func (v *DefaultVisitor) Div(_ Visitor, m *Div)               {}
func (v *DefaultVisitor) Mod(_ Visitor, m *Mod)               {}
func (v *DefaultVisitor) UnaryMinus(_ Visitor, t *UnaryMinus) {}
func (v *DefaultVisitor) Primary(_ Visitor, p Primary) {
	switch p := p.(type) {
	case *Identifier:
		v.Identifier(v, p)
	case *FunctionCall:
		v.FunctionCall(v, p)
	case Literal:
		v.Literal(v, p)
	}
}
func (v *DefaultVisitor) Identifier(_ Visitor, i *Identifier)     {}
func (v *DefaultVisitor) FunctionCall(_ Visitor, f *FunctionCall) {}
func (v *DefaultVisitor) Literal(_ Visitor, l Literal) {
	switch l := l.(type) {
	case *PrimitiveLiteral:
		v.PrimitiveLiteral(v, l)
	case CompoundLiteral:
		v.CompoundLiteral(v, l)
	}
}
func (v *DefaultVisitor) PrimitiveLiteral(_ Visitor, l *PrimitiveLiteral) {}
func (v *DefaultVisitor) CompoundLiteral(_ Visitor, l CompoundLiteral) {
	switch l := l.(type) {
	case *ListLiteral:
		v.ListLiteral(v, l)
	case *MapLiteral:
		v.MapLiteral(v, l)
	}
}
func (v *DefaultVisitor) ListLiteral(_ Visitor, l *ListLiteral) {}
func (v *DefaultVisitor) MapLiteral(_ Visitor, l *MapLiteral)   {}
func (v *DefaultVisitor) Indexing(_ Visitor, i *Indexing)       {}
func (v *DefaultVisitor) Type(_ Visitor, t Type) {
	switch t := t.(type) {
	case *PrimitiveType:
		v.PrimitiveType(v, t)
	case *ListType:
		v.ListType(v, t)
	case *MapType:
		v.MapType(v, t)
	}
}
func (v *DefaultVisitor) PrimitiveType(_ Visitor, t *PrimitiveType) {}
func (v *DefaultVisitor) ListType(_ Visitor, t *ListType)           {}
func (v *DefaultVisitor) MapType(_ Visitor, t *MapType)             {}
