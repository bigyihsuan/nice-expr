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
	// block-related
	Return(v Visitor, r *Return)
	Break(v Visitor, r *Break)
	Block(v Visitor, b *Block)
	If(v Visitor, i *If)
	For(v Visitor, f *For)
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
