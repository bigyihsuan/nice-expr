package visitor

import (
	"fmt"
	"nice-expr/src/ast"
	"nice-expr/src/token/tokentype"
	"nice-expr/src/util"
	"nice-expr/src/value"
)

type TypeChecker struct {
	ast.DefaultVisitor
	identifiers map[string]IdentifierEntry[value.ValueType]
	typeStack   util.Stack[value.ValueType]
	errors      util.Stack[error]
}

func NewTypeChecker() *TypeChecker {
	tc := new(TypeChecker)
	tc.identifiers = make(map[string]IdentifierEntry[value.ValueType])
	tc.typeStack = util.Stack[value.ValueType]{}
	tc.errors = util.Stack[error]{}
	return tc
}

func (v TypeChecker) TypeStack() util.Stack[value.ValueType] {
	return v.typeStack
}

func (v TypeChecker) Errors() util.Stack[error] {
	return v.errors
}

func (v TypeChecker) Identifiers() map[string]IdentifierEntry[value.ValueType] {
	return v.identifiers
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
	case ast.Declaration:
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
		v.errors.Push(fmt.Errorf("mismatched types in VariableDeclaration: got %v and %v", varType, valType))
	}
	v.typeStack.Push(varType)
	v.identifiers[s.Name.Name()] = IdentifierEntry[value.ValueType]{s.Name, varType, Var}
}
func (v *TypeChecker) ConstantDeclaration(_ ast.Visitor, s *ast.ConstantDeclaration) {
	s.Type.Accept(v)
	varType, _ := v.typeStack.Pop()
	s.Value.Accept(v)
	valType, _ := v.typeStack.Pop()
	if varType.NotEqual(valType) {
		v.errors.Push(fmt.Errorf("mismatched types in ConstantDeclaration: got %v and %v", varType, valType))
	}
	v.typeStack.Push(varType)
	v.identifiers[s.Name.Name()] = IdentifierEntry[value.ValueType]{s.Name, varType, Const}
}
func (v *TypeChecker) Assignment(_ ast.Visitor, s *ast.Assignment) {
	entry, ok := v.identifiers[s.Name.Name()]
	if !ok {
		v.errors.Push(fmt.Errorf("identifier `%s` not defined before use", s.Name.Name()))
	}
	varType := entry.Value
	s.Value.Accept(v)
	valType, _ := v.typeStack.Pop()
	if varType.NotEqual(valType) {
		v.errors.Push(fmt.Errorf("mismatched types in Assignment: got %v and %v", varType, valType))
	}
	v.typeStack.Push(varType)
}

func (v *TypeChecker) AndTest(_ ast.Visitor, t *ast.AndTest) {
	t.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.errors.Push(fmt.Errorf("operation And only allowed on Bool and Bool: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) OrTest(_ ast.Visitor, t *ast.OrTest) {
	t.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.errors.Push(fmt.Errorf("operation Or only allowed on Bool and Bool: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) NotTest(_ ast.Visitor, t *ast.NotTest) {
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.BoolType) {
		v.errors.Push(fmt.Errorf("operation Not only allowed on Bool: got %v", right))
	}
	v.typeStack.Push(value.BoolType)
}

func (v *TypeChecker) Equal(_ ast.Visitor, c *ast.Equal) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Equal: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) Greater(_ ast.Visitor, c *ast.Greater) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Greater: got %v and %v", left, right))
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.errors.Push(fmt.Errorf("types are not comparable for Greater: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) Less(_ ast.Visitor, c *ast.Less) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Less: got %v and %v", left, right))
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.errors.Push(fmt.Errorf("types are not comparable for Less: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) GreaterEqual(_ ast.Visitor, c *ast.GreaterEqual) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for GreaterEqual: got %v and %v", left, right))
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.errors.Push(fmt.Errorf("types are not comparable for GreaterEqual: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}
func (v *TypeChecker) LessEqual(_ ast.Visitor, c *ast.LessEqual) {
	c.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	c.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for LessEqual: got %v and %v", left, right))
	}
	if !left.IsComparable() || !right.IsComparable() {
		v.errors.Push(fmt.Errorf("types are not comparable for LessEqual: got %v and %v", left, right))
	}
	v.typeStack.Push(value.BoolType)
}

func (v *TypeChecker) Add(_ ast.Visitor, a *ast.Add) {
	a.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	a.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Add: got %v and %v", left, right))
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) Sub(_ ast.Visitor, a *ast.Sub) {
	a.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	a.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Sub: got %v and %v", left, right))
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) Mul(_ ast.Visitor, m *ast.Mul) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Mul: got %v and %v", left, right))
	}
	v.typeStack.Push(left)
}
func (v *TypeChecker) Div(_ ast.Visitor, m *ast.Div) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Div: got %v and %v", left, right))
	}
	v.typeStack.Push(left)
}
func (v *TypeChecker) Mod(_ ast.Visitor, m *ast.Mod) {
	m.Left.Accept(v)
	left, _ := v.typeStack.Pop()
	m.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if left.NotEqual(right) {
		v.errors.Push(fmt.Errorf("mismatched types for Mod: got %v and %v", left, right))
	}
	if left.IsNot(value.IntType) || right.IsNot(value.IntType) {
		v.errors.Push(fmt.Errorf("only Int allowed for Mod, got %v and %v", left, right))
	}
	v.typeStack.Push(left)
}

func (v *TypeChecker) UnaryMinus(_ ast.Visitor, t *ast.UnaryMinus) {
	t.Right.Accept(v)
	right, _ := v.typeStack.Pop()
	if right.IsNot(value.IntType) || right.IsNot(value.DecType) {
		v.errors.Push(fmt.Errorf("UnaryMinus only allowed on Int or Dec: got %v", right))
	}
	v.typeStack.Push(right)
}

func (v *TypeChecker) Identifier(_ ast.Visitor, i *ast.Identifier) {
	entry, ok := v.identifiers[i.Name()]
	if !ok {
		v.errors.Push(fmt.Errorf("identifier not found: %s", i.Name()))
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
		v.errors.Push(fmt.Errorf("unknown primitive literal: %s", l.Token))
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
			v.errors.Push(fmt.Errorf("mixed types in list: %v", elementTypes))
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
			v.errors.Push(fmt.Errorf("mixed key types in map: %v", keyTypes))
			return
		}
		if util.Any(valueTypes, func(e value.ValueType) bool { return e.NotEqual(firstValue) }) {
			v.errors.Push(fmt.Errorf("mixed value in map: %v", valueTypes))
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
		v.errors.Push(fmt.Errorf("type not indexable: %v", collection))
	}
	i.Right.Accept(v)
	index, _ := v.typeStack.Pop()
	switch {
	case collection.Is(value.StrType):
		// int index only
		if index.IsNot(value.IntType) {
			v.errors.Push(fmt.Errorf("collection Str is not indexable by %v, expect Int", index))
		}
		v.typeStack.Push(value.StrType)
	case collection.Is(value.ListType):
		// int index only
		if index.IsNot(value.IntType) {
			v.errors.Push(fmt.Errorf("collection List is not indexable by %v, expect Int", index))
		}
		v.typeStack.Push(collection.TypeArgs[0])
	case collection.Is(value.MapType):
		// match keytype
		keyType := collection.TypeArgs[0]
		if index.IsNot(keyType) {
			v.errors.Push(fmt.Errorf("collection Map is not indexable by %v, expect %v", index, keyType))
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
