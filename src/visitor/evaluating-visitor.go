package visitor

import (
	"fmt"
	"math/big"
	"nice-expr/src/ast"
	"nice-expr/src/evaluator"
	"nice-expr/src/util"
	"nice-expr/src/value"

	"golang.org/x/exp/slices"
)

type EvaluatingVisitor struct {
	ast.DefaultVisitor
	identifiers map[string]IdentifierEntry[*value.Value]
	valueStack  util.Stack[*value.Value]
	errors      util.Stack[error]
}

func NewEvaluatingVisitor() *EvaluatingVisitor {
	ev := new(EvaluatingVisitor)
	ev.identifiers = make(map[string]IdentifierEntry[*value.Value])
	ev.valueStack = util.Stack[*value.Value]{}
	ev.errors = util.Stack[error]{}
	return ev
}

func (v EvaluatingVisitor) ValueStack() util.Stack[*value.Value] {
	return v.valueStack
}

func (v EvaluatingVisitor) Errors() util.Stack[error] {
	return v.errors
}

func (v EvaluatingVisitor) Identifiers() map[string]IdentifierEntry[*value.Value] {
	return v.identifiers
}

func (v *EvaluatingVisitor) UnaryExpr(o ast.Visitor, e *ast.UnaryExpr) {
	e.Right.Accept(v)
}
func (v *EvaluatingVisitor) BinaryExpr(o ast.Visitor, e *ast.BinaryExpr) {
	e.Left.Accept(v)
	e.Right.Accept(v)
}

func (v *EvaluatingVisitor) Program(o ast.Visitor, p *ast.Program) {
	for _, e := range p.Statements {
		e.Accept(v)
	}
}
func (v *EvaluatingVisitor) Expr(o ast.Visitor, p ast.Expr) {
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

func (v *EvaluatingVisitor) Declaration(o ast.Visitor, d ast.Declaration) {
	switch d := d.(type) {
	case *ast.VariableDeclaration:
		o.VariableDeclaration(o, d)
	case *ast.ConstantDeclaration:
		o.ConstantDeclaration(o, d)
	}
}
func (v *EvaluatingVisitor) Test(o ast.Visitor, t ast.Test) {
	switch t := t.(type) {
	case *ast.AndTest:
		o.AndTest(o, t)
	case *ast.OrTest:
		o.OrTest(o, t)
	}
}
func (v *EvaluatingVisitor) Comparison(o ast.Visitor, c ast.Comparison) {
	c.AcceptCompare(o)
}
func (v *EvaluatingVisitor) AddExpr(o ast.Visitor, a ast.AddExpr) {
	a.AcceptAddExpr(o)
}
func (v *EvaluatingVisitor) MulExpr(o ast.Visitor, m ast.MulExpr) {
	m.AcceptMulExpr(o)
}
func (v *EvaluatingVisitor) Primary(o ast.Visitor, p ast.Primary) {
	switch p := p.(type) {
	case *ast.Identifier:
		p.Accept(o)
	case *ast.FunctionCall:
		p.Accept(o)
	case ast.Literal:
		p.Accept(o)
	}
}
func (v *EvaluatingVisitor) Literal(o ast.Visitor, l ast.Literal) {
	switch l := l.(type) {
	case *ast.PrimitiveLiteral:
		l.Accept(o)
	case ast.CompoundLiteral:
		l.AcceptCompoundLiteral(o)
	}
}
func (v *EvaluatingVisitor) CompoundLiteral(o ast.Visitor, l ast.CompoundLiteral) {
	l.AcceptCompoundLiteral(o)
}
func (v *EvaluatingVisitor) Type(o ast.Visitor, t ast.Type) {
	switch t := t.(type) {
	case *ast.PrimitiveType:
		v.PrimitiveType(o, t)
	case *ast.ListType:
		v.ListType(o, t)
	case *ast.MapType:
		v.MapType(o, t)
	}
}

func (v *EvaluatingVisitor) VariableDeclaration(_ ast.Visitor, s *ast.VariableDeclaration) {
	name := s.Name.Name()
	// ignore type
	s.Value.Accept(v)
	varVal, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, s))
		return
	}
	if existingValue, alreadyExists := v.identifiers[name]; alreadyExists {
		v.errors.Push(fmt.Errorf("name `%s` already exists as `%s` with value `%v`", name, existingValue.VarType, existingValue.Value.V))
	}
	v.identifiers[name] = IdentifierEntry[*value.Value]{s.Name, varVal, Var}
}
func (v *EvaluatingVisitor) ConstantDeclaration(_ ast.Visitor, s *ast.ConstantDeclaration) {
	name := s.Name.Name()
	// ignore type
	s.Value.Accept(v)
	varVal, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, s))
		return
	}
	if existingValue, alreadyExists := v.identifiers[name]; alreadyExists {
		v.errors.Push(fmt.Errorf("name `%s` already exists as `%s` with value `%v`", name, existingValue.VarType, existingValue.Value.V))
	}
	v.identifiers[name] = IdentifierEntry[*value.Value]{s.Name, varVal, Const}
}
func (v *EvaluatingVisitor) Assignment(_ ast.Visitor, s *ast.Assignment) {
	name := s.Name.Name()
	s.Value.Accept(v)
	assVal, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, s))
		return
	}
	entry, exists := v.identifiers[name]
	typesMatch := entry.Value.T.Equal(assVal.T)
	isVariable := entry.VarType == Var
	switch {
	case !exists:
		v.errors.Push(fmt.Errorf("variable `%s` used before defintion", name))
		return
	case !isVariable:
		v.errors.Push(fmt.Errorf("cannot assign `%v` to name `%s` because it is not a variable", assVal, name))
		return
	case !typesMatch:
		v.errors.Push(fmt.Errorf("cannot assign `%v` to variable `%s` of type %s", assVal, name, entry.Value.T.String()))
		return
	}
	// TODO: check s.Op for what operation to do. `is` is plain reassignment. call binaryexpr code on left OP value.
	entry.Value = assVal
	v.identifiers[name] = entry
}

func (v *EvaluatingVisitor) AndTest(_ ast.Visitor, t *ast.AndTest) {
	t.Left.Accept(v)
	leftValue, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("left: %s at %s", err, t))
		return
	}
	t.Right.Accept(v)
	rightValue, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("right: %s at %s", err, t))
		return
	}
	v.valueStack.Push(value.NewValue(value.BoolType, leftValue.V.(bool) && rightValue.V.(bool)))
}
func (v *EvaluatingVisitor) OrTest(_ ast.Visitor, t *ast.OrTest) {
	t.Left.Accept(v)
	leftValue, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("left: %s at %s", err, t))
		return
	}
	t.Right.Accept(v)
	rightValue, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("right: %s at %s", err, t))
		return
	}
	v.valueStack.Push(value.NewValue(value.BoolType, leftValue.V.(bool) || rightValue.V.(bool)))
}
func (v *EvaluatingVisitor) NotTest(_ ast.Visitor, t *ast.NotTest) {
	t.Right.Accept(v)
	rightValue, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("right: %s at %s", err, t))
		return
	} else if rightValue.IsNotType(value.BoolType) {
		v.errors.Push(fmt.Errorf("cannot take logical not of value %s of type %s at %s", rightValue.V, rightValue.T, t))
		return
	}
	v.valueStack.Push(value.NewValue(value.BoolType, !rightValue.V.(bool)))
}

func (v *EvaluatingVisitor) Equal(_ ast.Visitor, c *ast.Equal)               {}
func (v *EvaluatingVisitor) Greater(_ ast.Visitor, c *ast.Greater)           {}
func (v *EvaluatingVisitor) Less(_ ast.Visitor, c *ast.Less)                 {}
func (v *EvaluatingVisitor) GreaterEqual(_ ast.Visitor, c *ast.GreaterEqual) {}
func (v *EvaluatingVisitor) LessEqual(_ ast.Visitor, c *ast.LessEqual)       {}

func (v *EvaluatingVisitor) Add(_ ast.Visitor, a *ast.Add) {}
func (v *EvaluatingVisitor) Sub(_ ast.Visitor, a *ast.Sub) {}
func (v *EvaluatingVisitor) Mul(_ ast.Visitor, m *ast.Mul) {}
func (v *EvaluatingVisitor) Div(_ ast.Visitor, m *ast.Div) {}
func (v *EvaluatingVisitor) Mod(_ ast.Visitor, m *ast.Mod) {}

func (v *EvaluatingVisitor) UnaryMinus(_ ast.Visitor, t *ast.UnaryMinus) {
	t.Accept(v)
	rightVal, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, t))
		return
	}
	if rightVal.IsNotType(value.IntType) && rightVal.IsNotType(value.DecType) {
		v.errors.Push(fmt.Errorf("cannot use unary minus on value %s of type %s at %s", rightVal.V, rightVal.T, t))
		return
	}
	resultType := rightVal.T
	resultVal := rightVal.V
	switch {
	case rightVal.EqualsType(value.IntType):
		resultVal.(*big.Int).Neg(resultVal.(*big.Int))
	case rightVal.EqualsType(value.DecType):
		resultVal.(*big.Float).Neg(resultVal.(*big.Float))
	}
	v.valueStack.Push(value.NewValue(resultType, resultVal))
}

func (v *EvaluatingVisitor) Identifier(_ ast.Visitor, i *ast.Identifier) {
	name := i.Name()
	ident, ok := v.identifiers[name]
	if !ok {
		v.errors.Push(fmt.Errorf("name `%s` does not exist", name))
		return
	}
	v.valueStack.Push(ident.Value)
}
func (v *EvaluatingVisitor) FunctionCall(_ ast.Visitor, f *ast.FunctionCall) {
	if slices.Contains(evaluator.BuiltinFunctionNames, f.Ident.Name()) {
		v.BuiltinFunction(f)
	}
}

func (v *EvaluatingVisitor) PrimitiveLiteral(_ ast.Visitor, l *ast.PrimitiveLiteral) {
	val := new(value.Value)
	valType, ok := value.LitToType[l.Token.Tt]
	if !ok {
		v.errors.Push(fmt.Errorf("unkown primitive literal %s at %s", l.Token.Tt, l))
		return
	}
	val.T = valType
	val.V = l.Token.Value
	v.valueStack.Push(val)
}

func (v *EvaluatingVisitor) ListLiteral(_ ast.Visitor, l *ast.ListLiteral) {
	val := new(value.Value)
	valType := value.NewValueType("List")

	elements := []*value.Value{}

	for _, listVal := range l.Values {
		lv := listVal
		lv.Accept(v)
		value, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, l))
			return
		}
		elementType := value.T
		if len(valType.TypeArgs) < 1 {
			valType.AddTypeArg(elementType)
		}
		if !value.EqualsType(valType.TypeArgs[0]) {
			v.errors.Push(fmt.Errorf("incorrect element type: expected %v, got %v", elementType, value.T))
			return
		}
		elements = append(elements, value)
	}
	val.V = elements
	val.T = valType
	v.valueStack.Push(val)
}

func (v *EvaluatingVisitor) MapLiteral(_ ast.Visitor, l *ast.MapLiteral) {
	val := new(value.Value)
	valType := value.NewValueType("Map")

	elements := make(map[*value.Value]*value.Value)

	for keyExpr, valExpr := range l.Values {
		keyExpr.Accept(v)
		keyValue, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, l))
			return
		}
		keyType := keyValue.T
		if len(valType.TypeArgs) < 1 {
			valType.AddTypeArg(keyType)
		}
		if !keyValue.EqualsType(keyType) {
			v.errors.Push(fmt.Errorf("incorrect key type: expected %v, got %v", keyType, keyValue.T))
			return
		}
		valExpr.Accept(v)
		valueValue, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, l))
			return
		}
		valueType := valueValue.T
		if len(valType.TypeArgs) < 2 {
			valType.AddTypeArg(valueType)
		}
		if !valueValue.EqualsType(valueType) {
			v.errors.Push(fmt.Errorf("incorrect value type: expected %v, got %v", valueType, valueValue.T))
			return
		}
		elements[keyValue] = valueValue
	}
	val.V = elements
	val.T = valType
	v.valueStack.Push(val)
}

func (v *EvaluatingVisitor) Indexing(_ ast.Visitor, i *ast.Indexing) {}

func (v *EvaluatingVisitor) PrimitiveType(_ ast.Visitor, t *ast.PrimitiveType) { /* nop */ }
func (v *EvaluatingVisitor) ListType(_ ast.Visitor, t *ast.ListType)           { /* nop */ }
func (v *EvaluatingVisitor) MapType(_ ast.Visitor, t *ast.MapType)             { /* nop */ }

func (v *EvaluatingVisitor) BuiltinFunction(f *ast.FunctionCall) {
	name, arguments := f.Ident.Name(), f.Arguments
	switch name {
	case "print":
		if len(arguments) < 1 {
			fmt.Print()
		} else {
			for _, ex := range arguments {

				ex.Accept(v)
				val, err := v.valueStack.Pop()
				if err != nil {
					v.errors.Push(fmt.Errorf("%s at %s", err, f))
				}
				fmt.Print(val.Sprint())
			}
		}
	case "println":
		if len(arguments) < 1 {
			fmt.Println()
		} else {
			for _, ex := range arguments {
				ex.Accept(v)
				val, err := v.valueStack.Pop()
				if err != nil {
					v.errors.Push(fmt.Errorf("%s at %s", err, f))
				}
				fmt.Println(val.Sprint())
			}
		}
	case "len":
		if len(arguments) != 1 {
			v.errors.Push(fmt.Errorf("incorrect number of arguments for `len`: got %d, want %d", len(arguments), 1))
		}
		arguments[0].Accept(v)
		collection, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, f))
		}
		switch {
		case collection.EqualsType(value.StrType):
			val := collection.V.(string)
			v.valueStack.Push(value.NewValue(value.IntType, big.NewInt(int64(len([]rune(val))))))
			return
		case collection.T.Is(value.ListType):
			v.valueStack.Push(value.NewValue(value.IntType, big.NewInt(int64(len(collection.V.([]*value.Value))))))
			return
		case collection.T.Is(value.MapType):
			v.valueStack.Push(value.NewValue(value.IntType, big.NewInt(int64(len(collection.V.(map[*value.Value]*value.Value))))))
			return
		default:
			v.errors.Push(fmt.Errorf("invalid type for `len`: %s", collection.T.Name))
			return
		}
	default:
		v.errors.Push(fmt.Errorf("function `%s` does not exist at %s", name, f))
	}
}
