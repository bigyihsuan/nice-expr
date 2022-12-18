package visitor

import (
	"fmt"
	"nice-expr/src/ast"
	"nice-expr/src/evaluator"
	"nice-expr/src/token/tokentype"
	"nice-expr/src/util"
	"nice-expr/src/value"
)

type TypeChecker struct {
	ast.DefaultVisitor
	typeStack      util.Stack[value.ValueType]
	errors         util.Stack[error]
	currentContext *evaluator.Context[value.ValueType]
}

func NewTypeChecker() *TypeChecker {
	tc := new(TypeChecker)
	tc.typeStack = util.Stack[value.ValueType]{}
	tc.errors = util.Stack[error]{}
	tc.currentContext = evaluator.NewContext[value.ValueType]()
	return tc
}

func (v TypeChecker) TypeStack() util.Stack[value.ValueType] {
	return v.typeStack
}

func (v TypeChecker) Errors() util.Stack[error] {
	return v.errors
}

func (v TypeChecker) Identifiers() map[string]evaluator.IdentifierEntry[value.ValueType] {
	return v.currentContext.Identifiers
}

func (v TypeChecker) HasErrors() bool {
	return v.errors.Len() > 0
}

// make a new context with optional parent
func (v TypeChecker) NewContext(parent ...*evaluator.Context[value.ValueType]) *evaluator.Context[value.ValueType] {
	return evaluator.NewContext(parent...)
}

func (v *TypeChecker) AddError(format string, args ...any) {
	v.errors.Push(fmt.Errorf(format, args...))
}

func (v *TypeChecker) UnaryExpr(_ ast.Visitor, e *ast.UnaryExpr) {
	e.Right.Accept(v)
}
func (v *TypeChecker) BinaryExpr(_ ast.Visitor, e *ast.BinaryExpr) {
	e.Left.Accept(v)
	e.Right.Accept(v)
}
func (v *TypeChecker) Program(_ ast.Visitor, p *ast.Program) {
	for _, e := range p.Statements {
		e.Accept(v)
	}
}
func (v *TypeChecker) Expr(_ ast.Visitor, p ast.Expr) {
	switch p := p.(type) {
	case *ast.Indexing:
		p.Accept(v)
	case *ast.Assignment:
		p.Accept(v)
	case *ast.VariableDeclaration, *ast.ConstantDeclaration:
		p.Accept(v)
	case *ast.Return:
		p.Accept(v)
	case *ast.Block:
		p.Accept(v)
	case ast.Test:
		p.AcceptTest(v)
	}
}

func (v *TypeChecker) Declaration(_ ast.Visitor, d ast.Declaration) {
	switch d := d.(type) {
	case *ast.VariableDeclaration:
		v.VariableDeclaration(v, d)
	case *ast.ConstantDeclaration:
		v.ConstantDeclaration(v, d)
	}
}
func (v *TypeChecker) Test(_ ast.Visitor, t ast.Test) {
	switch t := t.(type) {
	case *ast.AndTest:
		v.AndTest(v, t)
	case *ast.OrTest:
		v.OrTest(v, t)
	}
}
func (v *TypeChecker) Comparison(_ ast.Visitor, c ast.Comparison) {
	c.AcceptCompare(v)
}
func (v *TypeChecker) AddExpr(_ ast.Visitor, a ast.AddExpr) {
	a.AcceptAddExpr(v)
}
func (v *TypeChecker) MulExpr(_ ast.Visitor, m ast.MulExpr) {
	m.AcceptMulExpr(v)
}
func (v *TypeChecker) Primary(_ ast.Visitor, p ast.Primary) {
	switch p := p.(type) {
	case *ast.Identifier:
		p.Accept(v)
	case *ast.FunctionCall:
		p.Accept(v)
	case ast.Literal:
		p.Accept(v)
	}
}
func (v *TypeChecker) Literal(_ ast.Visitor, l ast.Literal) {
	switch l := l.(type) {
	case *ast.PrimitiveLiteral:
		l.Accept(v)
	case ast.CompoundLiteral:
		l.AcceptCompoundLiteral(v)
	}
}
func (v *TypeChecker) CompoundLiteral(_ ast.Visitor, l ast.CompoundLiteral) {
	l.AcceptCompoundLiteral(v)
}
func (v *TypeChecker) Type(_ ast.Visitor, t ast.Type) {
	switch t := t.(type) {
	case *ast.PrimitiveType:
		v.PrimitiveType(v, t)
	case *ast.ListType:
		v.ListType(v, t)
	case *ast.MapType:
		v.MapType(v, t)
	}
}

func (v *TypeChecker) VariableDeclaration(_ ast.Visitor, s *ast.VariableDeclaration) {
	s.Type.Accept(v)
	varType, _ := v.typeStack.Pop()
	s.Value.Accept(v)
	valType, _ := v.typeStack.Pop()
	if varType.NotEqual(valType) {
		v.AddError("mismatched types: got %v and %v at %s", varType, valType, s)
		return
	}

	name := s.Name.Name()
	// overwrite any entries in the current context
	v.typeStack.Push(varType)
	v.currentContext.AddIdentifier(
		name,
		evaluator.IdentifierEntry[value.ValueType]{
			Ident:   s.Name,
			Value:   varType,
			VarType: evaluator.Var,
		},
	)
}
func (v *TypeChecker) ConstantDeclaration(_ ast.Visitor, s *ast.ConstantDeclaration) {
	s.Type.Accept(v)
	varType, _ := v.typeStack.Pop()
	s.Value.Accept(v)
	valType, _ := v.typeStack.Pop()
	if varType.NotEqual(valType) {
		v.AddError("mismatched types: got %v and %v at %s", varType, valType, s)
		return
	}

	name := s.Name.Name()
	// overwrite any entries in the current context
	v.typeStack.Push(varType)
	v.currentContext.AddIdentifier(
		name,
		evaluator.IdentifierEntry[value.ValueType]{
			Ident:   s.Name,
			Value:   varType,
			VarType: evaluator.Const,
		},
	)
}
func (v *TypeChecker) Assignment(_ ast.Visitor, s *ast.Assignment) {
	var ident *ast.Identifier
	var kind evaluator.VariableType
	name := s.Name.Name()
	_, identType, _, source := v.currentContext.GetIdentifier(name)
	if source != v.currentContext {
		// make a copy for this context
		ident, identType, kind, _ = source.GetIdentifier(name)
		v.currentContext.AddIdentifier(
			name,
			evaluator.IdentifierEntry[value.ValueType]{
				Ident:   ident,
				Value:   identType,
				VarType: kind,
			},
		)
	}

	s.Value.Accept(v)
	valType, _ := v.typeStack.Pop()
	if identType.NotEqual(valType) {
		v.AddError("mismatched types: got %v and %v at %s", identType, valType, s)
		return
	}
	v.typeStack.Push(identType)
}
func (v *TypeChecker) Return(_ ast.Visitor, r *ast.Return) {
	r.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	v.typeStack.Push(right)
}
func (v *TypeChecker) Block(_ ast.Visitor, b *ast.Block) {
	// when entering a block, make a new context, with outside variables available
	blockContext := evaluator.CopyContext(v.currentContext, v.currentContext)
	// move contexts into this block
	v.currentContext = blockContext
	// blocks have the type of the first return statement
	for _, e := range b.Statements {
		e.Accept(v)
		t, err := v.typeStack.Pop()
		if err != nil {
			v.AddError("%s at %s in %s", err, e, b)
		}
		_, isReturn := e.(*ast.Return)
		if isReturn {
			v.typeStack.Push(t)
			break
		}
	}
	// exit this context
	v.currentContext = v.currentContext.Parent
}

func (v *TypeChecker) AndTest(_ ast.Visitor, t *ast.AndTest) {
	t.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.AddError("operation And only allowed on Bool and Bool: got %v and %v at %s", left, right, t)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) OrTest(_ ast.Visitor, t *ast.OrTest) {
	t.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.AddError("operation Or only allowed on Bool and Bool: got %v and %v at %s", left, right, t)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) NotTest(_ ast.Visitor, t *ast.NotTest) {
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.AddError("operation Not only allowed on Bool: got %v at %s", right, t)
	}
	v.typeStack.Push(value.BoolType)
}

func (v *TypeChecker) Equal(_ ast.Visitor, c *ast.Equal) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.AddError("mismatched types for Equal: got %v and %v at %s", left, right, c)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) Greater(_ ast.Visitor, c *ast.Greater) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.AddError("mismatched types for Greater: got %v and %v at %s", left, right, c)
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.AddError("types are not comparable for Greater: got %v and %v at %s", left, right, c)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) Less(_ ast.Visitor, c *ast.Less) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.AddError("mismatched types for Less: got %v and %v at %s", left, right, c)
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.AddError("types are not comparable for Less: got %v and %v at %s", left, right, c)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) GreaterEqual(_ ast.Visitor, c *ast.GreaterEqual) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.AddError("mismatched types for GreaterEqual: got %v and %v at %s", left, right, c)
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.AddError("types are not comparable for GreaterEqual: got %v and %v at %s", left, right, c)
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) LessEqual(_ ast.Visitor, c *ast.LessEqual) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.AddError("mismatched types for LessEqual: got %v and %v at %s", left, right, c)
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.AddError("types are not comparable for LessEqual: got %v and %v at %s", left, right, c)
	}
	v.typeStack.Push(value.BoolType)
}

func (v *TypeChecker) Add(_ ast.Visitor, a *ast.Add) {
	a.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	a.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	switch {
	case left.Equal(value.IntType) && right.Equal(value.IntType) ||
		left.Equal(value.IntType) && right.Equal(value.DecType) ||
		left.Equal(value.DecType) && right.Equal(value.IntType) ||
		left.Equal(value.DecType) && right.Equal(value.DecType) ||
		left.Equal(value.StrType) && right.Equal(value.StrType) ||
		left.Is(value.ListType) && right.Is(value.ListType) && left.Equal(right):
		// nop
	default:
		v.AddError("mismatched types for Add: got %v and %v at %s", left, right, a)
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) Sub(_ ast.Visitor, a *ast.Sub) {
	a.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	a.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	switch {
	case left.Equal(value.IntType) && right.Equal(value.IntType) ||
		left.Equal(value.IntType) && right.Equal(value.DecType) ||
		left.Equal(value.DecType) && right.Equal(value.IntType) ||
		left.Equal(value.DecType) && right.Equal(value.DecType) ||
		left.Equal(value.StrType) && right.Equal(value.StrType) ||
		left.Is(value.ListType) && right.Is(value.ListType) && left.Equal(right):
		// nop
	default:
		v.AddError("mismatched types for Sub: got %v and %v at %s", left, right, a)
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) Mul(_ ast.Visitor, m *ast.Mul) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	switch {
	case left.Equal(value.IntType) && right.Equal(value.IntType) ||
		left.Equal(value.IntType) && right.Equal(value.DecType) ||
		left.Equal(value.DecType) && right.Equal(value.IntType) ||
		left.Equal(value.DecType) && right.Equal(value.DecType):
		// nop
	default:
		v.AddError("mismatched types for Mul: got %v and %v at %s", left, right, m)
	}
	v.typeStack.Push(left)
}
func (v *TypeChecker) Div(_ ast.Visitor, m *ast.Div) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	switch {
	case left.Equal(value.IntType) && right.Equal(value.IntType) ||
		left.Equal(value.IntType) && right.Equal(value.DecType) ||
		left.Equal(value.DecType) && right.Equal(value.IntType) ||
		left.Equal(value.DecType) && right.Equal(value.DecType):
		// nop
	default:
		v.AddError("mismatched types for Div: got %v and %v at %s", left, right, m)
	}
	v.typeStack.Push(left)
}
func (v *TypeChecker) Mod(_ ast.Visitor, m *ast.Mod) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.IsNot(value.IntType) || right.IsNot(value.IntType) {
		v.AddError("only Int allowed for Mod, got %v and %v at %s", left, right, m)
	}
	if left.NotEqual(right) {
		v.AddError("mismatched types for Mod: got %v and %v at %s", left, right, m)
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) UnaryMinus(_ ast.Visitor, t *ast.UnaryMinus) {
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.IntType) || right.IsNot(value.DecType) {
		v.AddError("UnaryMinus only allowed on Int or Dec: got %v at %s", right, t)
	}
	v.typeStack.Push(right)
}

func (v *TypeChecker) Identifier(_ ast.Visitor, i *ast.Identifier) {
	entry, ok := v.currentContext.Identifiers[i.Name()]
	if !ok {
		v.AddError("identifier not found: %s at %s", i.Name(), i)
		v.typeStack.Push(value.NoneType)
	} else {
		v.typeStack.Push(entry.Value)
	}
}
func (v *TypeChecker) FunctionCall(_ ast.Visitor, f *ast.FunctionCall) {}

func (v *TypeChecker) PrimitiveLiteral(_ ast.Visitor, l *ast.PrimitiveLiteral) {
	switch l.Token.Tt {
	case tokentype.Integer:
		v.typeStack.Push(value.IntType)
	case tokentype.Floating:
		v.typeStack.Push(value.DecType)
	case tokentype.String:
		v.typeStack.Push(value.StrType)
	case tokentype.True, tokentype.False:
		v.typeStack.Push(value.BoolType)
	default:
		v.AddError("unknown primitive literal: %s", l.Token)
	}
}

func (v *TypeChecker) ListLiteral(_ ast.Visitor, l *ast.ListLiteral) {
	t := value.NewValueType("List")
	elementTypes := []value.ValueType{}
	for _, element := range l.Values {
		element.Accept(v)
		eleType, _ := v.typeStack.Pop()
		elementTypes = append(elementTypes, eleType)
	}
	var first value.ValueType
	if len(elementTypes) > 0 {
		first = elementTypes[0]
		if util.Any(elementTypes, func(e value.ValueType) bool { return e.NotEqual(first) }) {
			v.AddError("mixed types in list: %v at %s", elementTypes, l)
			return
		}
		t.AddTypeArg(first)
		v.typeStack.Push(t)
	} else {
		t.AddTypeArg(value.NoneType)
		v.typeStack.Push(t)
	}
}
func (v *TypeChecker) MapLiteral(_ ast.Visitor, l *ast.MapLiteral) {
	t := value.NewValueType("Map")

	keys := make([]ast.Expr, 0, len(l.Values))
	values := make([]ast.Expr, 0, len(l.Values))
	for k, val := range l.Values {
		keys = append(keys, k)
		values = append(values, val)
	}

	keyTypes := []value.ValueType{}
	valueTypes := []value.ValueType{}
	for i := range keys {
		keys[i].Accept(v)
		keyType, _ := v.typeStack.Pop()
		values[i].Accept(v)
		valueType, _ := v.typeStack.Pop()
		keyTypes = append(keyTypes, keyType)
		valueTypes = append(valueTypes, valueType)
	}
	var firstKey, firstValue value.ValueType
	if len(valueTypes) > 0 {
		firstKey, firstValue = keyTypes[0], valueTypes[0]
		if util.Any(keyTypes, func(e value.ValueType) bool { return e.NotEqual(firstKey) }) {
			v.AddError("mixed key types in map: %v at %s", keyTypes, l)
			return
		}
		if util.Any(valueTypes, func(e value.ValueType) bool { return e.NotEqual(firstValue) }) {
			v.AddError("mixed value in map: %v at %s", valueTypes, l)
			return
		}
		t.AddTypeArg(firstKey)
		t.AddTypeArg(firstValue)
		v.typeStack.Push(t)
	} else {
		t.AddTypeArg(value.NoneType)
		t.AddTypeArg(value.NoneType)
		v.typeStack.Push(t)
	}
}
func (v *TypeChecker) Indexing(_ ast.Visitor, i *ast.Indexing) {
	i.Left.Accept(v)
	collection, _ := v.typeStack.Pop()
	if util.All(value.IndexableTypes, func(t value.ValueType) bool { return t.IsNot(collection) }) {
		v.AddError("type not indexable: %v at %s", collection, i)
	}
	i.Right.Accept(v)
	index, _ := v.typeStack.Pop()
	switch {
	case collection.Is(value.StrType):
		// int index only
		if index.IsNot(value.IntType) {
			v.AddError("collection Str is not indexable by %v, expect Int at %s", index, i)
		}
		v.typeStack.Push(value.StrType)
	case collection.Is(value.ListType):
		// int index only
		if index.IsNot(value.IntType) {
			v.AddError("collection List is not indexable by %v, expect Int at %s", index, i)
		}
		v.typeStack.Push(collection.TypeArgs[0])
	case collection.Is(value.MapType):
		// match keytype
		keyType := collection.TypeArgs[0]
		if index.IsNot(keyType) {
			v.AddError("collection Map is not indexable by %v, expect %v at %s", index, keyType, i)
		}
		v.typeStack.Push(collection.TypeArgs[1])
	}
}

func (v *TypeChecker) PrimitiveType(_ ast.Visitor, t *ast.PrimitiveType) {
	v.typeStack.Push(t.ToValueType())
}
func (v *TypeChecker) ListType(_ ast.Visitor, t *ast.ListType) {
	v.typeStack.Push(t.ToValueType())
}
func (v *TypeChecker) MapType(_ ast.Visitor, t *ast.MapType) {
	v.typeStack.Push(t.ToValueType())
}
