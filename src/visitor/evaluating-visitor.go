package visitor

import (
	"fmt"
	"math/big"
	"nice-expr/src/ast"
	"nice-expr/src/evaluator"
	"nice-expr/src/util"
	"nice-expr/src/value"
	"strings"

	"golang.org/x/exp/slices"
)

type EvaluatingVisitor struct {
	ast.DefaultVisitor
	valueStack     util.Stack[*value.Value]
	errors         util.Stack[error]
	currentContext *evaluator.Context[*value.Value]
}

func NewEvaluatingVisitor() *EvaluatingVisitor {
	ev := new(EvaluatingVisitor)
	ev.valueStack = util.Stack[*value.Value]{}
	ev.errors = util.Stack[error]{}
	ev.currentContext = evaluator.NewContext[*value.Value]()
	return ev
}

func (v EvaluatingVisitor) ValueStack() util.Stack[*value.Value] {
	return v.valueStack
}

func (v EvaluatingVisitor) Errors() util.Stack[error] {
	return v.errors
}

func (v EvaluatingVisitor) Identifiers() map[string]evaluator.IdentifierEntry[*value.Value] {
	return v.currentContext.Identifiers
}

func (v EvaluatingVisitor) NewContext(parent ...*evaluator.Context[*value.Value]) *evaluator.Context[*value.Value] {
	return evaluator.NewContext(parent...)
}

func (v *EvaluatingVisitor) PrepareUnary(u ast.UnaryExpr) (right *value.Value) {
	var err error
	u.Accept(v)
	right, err = v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s (right)", err, u))
		return
	}
	return right
}
func (v *EvaluatingVisitor) PrepareBinary(b ast.BinaryExpr) (left, right *value.Value) {
	var err error
	b.Accept(v)
	right, err = v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s (right)", err, b))
		return
	}
	left, err = v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s (left)", err, b))
		return
	}
	return left, right
}

func (v *EvaluatingVisitor) UnaryExpr(_ ast.Visitor, e *ast.UnaryExpr) {
	e.Right.Accept(v)
}
func (v *EvaluatingVisitor) BinaryExpr(_ ast.Visitor, e *ast.BinaryExpr) {
	e.Left.Accept(v)
	e.Right.Accept(v)
}

func (v *EvaluatingVisitor) Program(_ ast.Visitor, p *ast.Program) {
	for _, e := range p.Statements {
		e.Accept(v)
	}
}
func (v *EvaluatingVisitor) Expr(_ ast.Visitor, p ast.Expr) {
	switch p := p.(type) {
	case *ast.Indexing:
		p.Accept(v)
	case *ast.Block:
		p.Accept(v)
	case *ast.Assignment:
		p.Accept(v)
	case *ast.VariableDeclaration, *ast.ConstantDeclaration:
		v.Declaration(v, p)
	case ast.Test:
		v.Test(v, p)
	}
}

func (v *EvaluatingVisitor) Declaration(_ ast.Visitor, d ast.Declaration) {
	switch d := d.(type) {
	case *ast.VariableDeclaration:
		v.VariableDeclaration(v, d)
	case *ast.ConstantDeclaration:
		v.ConstantDeclaration(v, d)
	}
}
func (v *EvaluatingVisitor) Test(_ ast.Visitor, t ast.Test) {
	switch t := t.(type) {
	case *ast.AndTest:
		v.AndTest(v, t)
	case *ast.OrTest:
		v.OrTest(v, t)
	}
}
func (v *EvaluatingVisitor) Comparison(_ ast.Visitor, c ast.Comparison) {
	c.AcceptCompare(v)
}
func (v *EvaluatingVisitor) AddExpr(_ ast.Visitor, a ast.AddExpr) {
	a.AcceptAddExpr(v)
}
func (v *EvaluatingVisitor) MulExpr(_ ast.Visitor, m ast.MulExpr) {
	m.AcceptMulExpr(v)
}
func (v *EvaluatingVisitor) Primary(_ ast.Visitor, p ast.Primary) {
	switch p := p.(type) {
	case *ast.Identifier:
		p.Accept(v)
	case *ast.FunctionCall:
		p.Accept(v)
	case ast.Literal:
		p.Accept(v)
	}
}
func (v *EvaluatingVisitor) Literal(_ ast.Visitor, l ast.Literal) {
	switch l := l.(type) {
	case *ast.PrimitiveLiteral:
		l.Accept(v)
	case ast.CompoundLiteral:
		l.AcceptCompoundLiteral(v)
	}
}
func (v *EvaluatingVisitor) CompoundLiteral(_ ast.Visitor, l ast.CompoundLiteral) {
	l.AcceptCompoundLiteral(v)
}
func (v *EvaluatingVisitor) Type(_ ast.Visitor, t ast.Type) {
	switch t := t.(type) {
	case *ast.PrimitiveType:
		v.PrimitiveType(v, t)
	case *ast.ListType:
		v.ListType(v, t)
	case *ast.MapType:
		v.MapType(v, t)
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
	// if existingValue, alreadyExists := v.currentContext.Identifiers[name]; alreadyExists {
	// 	v.errors.Push(fmt.Errorf("name `%s` already exists as `%s` with value `%v` at %s", name, existingValue.VarType, existingValue.Value.V, s))
	// }
	v.currentContext.AddIdentifier(
		name,
		evaluator.IdentifierEntry[*value.Value]{
			Ident:   s.Name,
			Value:   varVal.UnwrapReturn(),
			VarType: evaluator.Var,
		},
	)
	v.valueStack.Push(varVal)
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
	// if existingValue, alreadyExists := v.currentContext.Identifiers[name]; alreadyExists {
	// 	v.errors.Push(fmt.Errorf("name `%s` already exists as `%s` with value `%v` at %s", name, existingValue.VarType, existingValue.Value.V, s))
	// }
	v.currentContext.AddIdentifier(
		name,
		evaluator.IdentifierEntry[*value.Value]{
			Ident:   s.Name,
			Value:   varVal.UnwrapReturn(),
			VarType: evaluator.Const,
		},
	)
	v.valueStack.Push(varVal)
}
func (v *EvaluatingVisitor) Assignment(_ ast.Visitor, s *ast.Assignment) {
	name := s.Name.Name()
	s.Value.Accept(v)
	assVal, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, s))
		return
	}
	entry, exists := v.currentContext.Identifiers[name]
	typesMatch := entry.Value.T.Equal(assVal.T)
	isVariable := entry.VarType == evaluator.Var
	switch {
	case !exists:
		v.errors.Push(fmt.Errorf("variable `%s` used before defintion at %s", name, s))
		return
	case !isVariable:
		v.errors.Push(fmt.Errorf("cannot assign `%v` to name `%s` because it is not a variable at %s", assVal, name, s))
		return
	case !typesMatch:
		v.errors.Push(fmt.Errorf("cannot assign `%v` to variable `%s` of type %s at %s", assVal, name, entry.Value.T.String(), s))
		return
	}
	// unwrap return/break
	assVal = assVal.UnwrapReturn()

	entry.Value = assVal
	v.currentContext.Identifiers[name] = entry
	v.valueStack.Push(entry.Value)
}
func (v *EvaluatingVisitor) Return(_ ast.Visitor, r *ast.Return) {
	if r.UnaryExpr.Right == nil {
		return // don't push a value when no right expr
	}
	right := v.PrepareUnary(r.UnaryExpr)
	// wrap in a Return[] signalling type
	returnedType := right.T
	right.T = value.NewValueType("Return")
	right.T.AddTypeArg(returnedType)
	v.valueStack.Push(right)
}
func (v *EvaluatingVisitor) Break(_ ast.Visitor, r *ast.Break) {
	if r.UnaryExpr.Right == nil {
		return // don't push a value when no right expr
	}
	right := v.PrepareUnary(r.UnaryExpr)
	// wrap in a Break[] signalling type
	returnedType := right.T
	right.T = value.NewValueType("Break")
	right.T.AddTypeArg(returnedType)
	v.valueStack.Push(right)
}
func (v *EvaluatingVisitor) Block(_ ast.Visitor, b *ast.Block) {
	for _, e := range b.Statements {
		e.Accept(v)
		returnVal, err := v.valueStack.Pop()
		var isReturnOrBroke bool
		if returnVal != nil {
			isReturnOrBroke = returnVal.IsType(value.ReturnedType) || returnVal.IsType(value.BrokeType)
		}
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, e))
		} else if isReturnOrBroke {
			v.valueStack.Push(returnVal)
			break
		}
		v.valueStack.Push(returnVal)
	}
}
func (v *EvaluatingVisitor) If(_ ast.Visitor, i *ast.If) {
	i.Condition.Accept(v)
	condition, err := v.valueStack.Pop()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, i))
		return
	}
	conditionResult, err := condition.Bool()
	if err != nil {
		v.errors.Push(fmt.Errorf("%s at %s", err, i))
		return
	}
	// when entering a block, make a new context, with outside variables available
	blockContext := evaluator.CopyContext(v.currentContext, v.currentContext)
	// move contexts into this block
	v.currentContext = blockContext
	if conditionResult {
		i.Then.Accept(v)
	} else if i.ElseIf != nil {
		i.ElseIf.Accept(v)
	} else if i.Else != nil {
		i.Else.Accept(v)
	}
	v.currentContext = v.currentContext.Parent
}
func (v *EvaluatingVisitor) For(_ ast.Visitor, f *ast.For) {
	// let the variables declare themselves
	for _, local := range f.LocalVariables {
		local.Accept(v)
	}
	// when entering a block, make a new context, with outside variables available
	blockContext := evaluator.CopyContext(v.currentContext, v.currentContext)
	// move contexts into this block
	v.currentContext = blockContext
	for {
		f.Body.Accept(v)
		maybeBreak, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, f))
		}
		if maybeBreak != nil && (maybeBreak.IsType(value.ReturnedType) || maybeBreak.IsType(value.BrokeType)) {
			// push back the return value of this loop
			v.valueStack.Push(maybeBreak.UnwrapReturn())
			break
		}
	}
	// remove the for local variables from the parent context
	for _, local := range f.LocalVariables {
		var name string
		switch local := local.(type) {
		case *ast.VariableDeclaration:
			name = local.Name.Name()
		case *ast.ConstantDeclaration:
			name = local.Name.Name()
		}
		if name != "" {
			v.currentContext.DeleteIdentifier(name)
		}
	}
	v.currentContext = v.currentContext.Parent
}

func (v *EvaluatingVisitor) AndTest(_ ast.Visitor, t *ast.AndTest) {
	leftValue, rightValue := v.PrepareBinary(t.BinaryExpr)
	v.valueStack.Push(value.NewValue(value.BoolType, leftValue.V.(bool) && rightValue.V.(bool)))
}
func (v *EvaluatingVisitor) OrTest(_ ast.Visitor, t *ast.OrTest) {
	leftValue, rightValue := v.PrepareBinary(t.BinaryExpr)
	v.valueStack.Push(value.NewValue(value.BoolType, leftValue.V.(bool) || rightValue.V.(bool)))
}
func (v *EvaluatingVisitor) NotTest(_ ast.Visitor, t *ast.NotTest) {
	rightValue := v.PrepareUnary(t.UnaryExpr)
	if rightValue.IsNotType(value.BoolType) {
		v.errors.Push(fmt.Errorf("cannot take logical not of value %s of type %s at %s", rightValue.V, rightValue.T, t))
		return
	}
	v.valueStack.Push(value.NewValue(value.BoolType, !rightValue.V.(bool)))
}

func (v *EvaluatingVisitor) Equal(_ ast.Visitor, c *ast.Equal) {
	left, right := v.PrepareBinary(c.BinaryExpr)
	switch {
	case left.EqualsValueType(right):
		v.valueStack.Push(value.NewValue(value.BoolType, left.Equal(right)))
	default:
		v.errors.Push(fmt.Errorf("mismatched types for `=`: %s and %s at %s", left.T, right.T, c))
	}
}
func (v *EvaluatingVisitor) Greater(_ ast.Visitor, c *ast.Greater) {
	left, right := v.PrepareBinary(c.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) > 0))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) > 0))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(value.BoolType, l > r))
	default:
		v.errors.Push(fmt.Errorf("mismatched or unsupported types for `>`: %s and %s at %s", left.T, right.T, c))
	}
}
func (v *EvaluatingVisitor) Less(_ ast.Visitor, c *ast.Less) {
	left, right := v.PrepareBinary(c.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) < 0))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) < 0))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(value.BoolType, l < r))
	default:
		v.errors.Push(fmt.Errorf("mismatched or unsupported types for `<`: %s and %s at %s", left.T, right.T, c))
	}
}
func (v *EvaluatingVisitor) GreaterEqual(_ ast.Visitor, c *ast.GreaterEqual) {
	left, right := v.PrepareBinary(c.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) >= 0))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) >= 0))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(value.BoolType, l >= r))
	default:
		v.errors.Push(fmt.Errorf("mismatched or unsupported types for `>=`: %s and %s at %s", left.T, right.T, c))
	}
}
func (v *EvaluatingVisitor) LessEqual(_ ast.Visitor, c *ast.LessEqual) {
	left, right := v.PrepareBinary(c.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) <= 0))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(value.BoolType, l.Cmp(r) <= 0))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(value.BoolType, l <= r))
	default:
		v.errors.Push(fmt.Errorf("mismatched or unsupported types for `<=`: %s and %s at %s", left.T, right.T, c))
	}
}

func (v *EvaluatingVisitor) Add(_ ast.Visitor, a *ast.Add) {
	left, right := v.PrepareBinary(a.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(left.T, big.NewInt(0).Add(l, r)))
	case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Add(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Add(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Add(l, r)))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(left.T, l+r))
	case left.IsType(value.ListType) && right.IsType(value.ListType) && left.EqualsValueType(right):
		l, _ := left.List()
		r, _ := right.List()
		v.valueStack.Push(value.NewValue(left.T, append(l, r...)))
	default:
		v.errors.Push(fmt.Errorf("invalid type combo `%s` + `%s` at %s", left.T, right.T, a))
	}
}
func (v *EvaluatingVisitor) Sub(_ ast.Visitor, a *ast.Sub) {
	left, right := v.PrepareBinary(a.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(left.T, big.NewInt(0).Sub(l, r)))
	case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Sub(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Sub(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Sub(l, r)))
	case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
		l, _ := left.Str()
		r, _ := right.Str()
		v.valueStack.Push(value.NewValue(left.T, strings.ReplaceAll(l, r, "")))
	case left.IsType(value.ListType) && right.IsType(value.ListType) && left.EqualsValueType(right):
		l, _ := left.List()
		r, _ := right.List()
		var diff []*value.Value
		for _, x := range l {
			for _, y := range r {
				if !x.Equal(y) {
					diff = append(diff, x)
				}
			}
		}
		v.valueStack.Push(value.NewValue(left.T, diff))
	default:
		v.errors.Push(fmt.Errorf("invalid type combo `%s` - `%s` at %s", left.T, right.T, a))
	}
}

func (v *EvaluatingVisitor) Mul(_ ast.Visitor, m *ast.Mul) {
	left, right := v.PrepareBinary(m.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		v.valueStack.Push(value.NewValue(left.T, big.NewInt(0).Mul(l, r)))
	case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Mul(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Mul(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Mul(l, r)))
	default:
		v.errors.Push(fmt.Errorf("invalid type combo `%s` * `%s` at %s", left.T, right.T, m))
	}
}
func (v *EvaluatingVisitor) Div(_ ast.Visitor, m *ast.Div) {
	left, right := v.PrepareBinary(m.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		if r.Cmp(big.NewInt(0)) == 0 {
			v.errors.Push(fmt.Errorf("division by 0 at %s", m))
		}
		v.valueStack.Push(value.NewValue(left.T, big.NewInt(0).Div(l, r)))
	case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		if r.Cmp(big.NewFloat(0.0)) == 0 {
			v.errors.Push(fmt.Errorf("division by 0 at %s", m))
		}
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Quo(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		if r.Cmp(big.NewFloat(0.0)) == 0 {
			v.errors.Push(fmt.Errorf("division by 0 at %s", m))
		}
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Quo(l, r)))
	case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
		l, _ := left.BigDec()
		r, _ := right.BigDec()
		if r.Cmp(big.NewFloat(0.0)) == 0 {
			v.errors.Push(fmt.Errorf("division by 0 at %s", m))
		}
		v.valueStack.Push(value.NewValue(left.T, big.NewFloat(0.0).Quo(l, r)))
	default:
		v.errors.Push(fmt.Errorf("invalid type combo `%s` / `%s` at %s", left.T, right.T, m))
	}
}
func (v *EvaluatingVisitor) Mod(_ ast.Visitor, m *ast.Mod) {
	left, right := v.PrepareBinary(m.BinaryExpr)
	switch {
	case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
		l, _ := left.BigInt()
		r, _ := right.BigInt()
		if r.Cmp(big.NewInt(0)) == 0 {
			v.errors.Push(fmt.Errorf("division by 0 at %s", m))
		}
		v.valueStack.Push(value.NewValue(left.T, big.NewInt(0).Mod(l, r)))
	default:
		v.errors.Push(fmt.Errorf("invalid type combo `%s` %% `%s` at %s", left.T, right.T, m))
	}
}

func (v *EvaluatingVisitor) UnaryMinus(_ ast.Visitor, t *ast.UnaryMinus) {
	rightVal := v.PrepareUnary(t.UnaryExpr)
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
	ident, ok := v.currentContext.Identifiers[name]
	if !ok {
		v.errors.Push(fmt.Errorf("name `%s` does not exist", name))
		return
	}
	v.valueStack.Push(ident.Value)
}
func (v *EvaluatingVisitor) FunctionCall(_ ast.Visitor, f *ast.FunctionCall) {
	if slices.Contains(evaluator.BuiltinFunctionNames, f.Ident.Name()) {
		v.BuiltinFunction(f)
		return
	} else {
		// TODO: call user-defined functions
		fmt.Println(f)
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

func (v *EvaluatingVisitor) Indexing(_ ast.Visitor, i *ast.Indexing) {
	left, right := v.PrepareBinary(i.BinaryExpr)
	if !left.T.IsIndexable() {
		v.errors.Push(fmt.Errorf("type `%s` is not indexable at %s", left.T, i))
		return
	}
	switch {
	case left.IsType(value.StrType) && right.IsType(value.IntType):
		str, _ := left.Str()
		runes := []rune(str)
		id, _ := right.Int64()
		index := int(id)
		if index >= len(runes) {
			v.errors.Push(fmt.Errorf("index `%d` out of range for `%s` (length %d) at %s", index, str, len(runes), i))
			return
		}
		v.valueStack.Push(value.NewValue(left.T, string(runes[index])))
	case left.IsType(value.ListType) && right.IsType(value.IntType):
		list, _ := left.List()
		id, _ := right.Int64()
		index := int(id)
		if index >= len(list) {
			v.errors.Push(fmt.Errorf("index `%d` out of range for `%s` (length %d) at %s", index, list, len(list), i))
			return
		}
		v.valueStack.Push(list[index])
	case left.IsType(value.MapType) && right.EqualsType(left.T.TypeArgs[0]):
		m, _ := left.Map()
		for key, value := range m {
			if key.Equal(right) {
				v.valueStack.Push(value)
				return
			}
		}
		// v.errors.Push(fmt.Errorf("key `%s` does not exist in map `%s` at %s", right, m, i))
		v.valueStack.Push(value.NewZeroValue(left.T.TypeArgs[1]))
	default:
		v.errors.Push(fmt.Errorf("type `%s` cannot be indexed by type `%s` at %s", left.T, right.T, i))
		return
	}
}

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
		// return none type
		// v.valueStack.Push(value.NewZeroValue(value.NoneType))
	case "println":
		if len(arguments) < 1 {
			fmt.Println()
		} else {
			for _, ex := range arguments {
				ex.Accept(v)
				val, err := v.valueStack.Pop()
				if err != nil {
					v.errors.Push(fmt.Errorf("%s at %s", err, f))
					return
				}
				fmt.Println(val.Sprint())
			}
		}
		// return none type
		// v.valueStack.Push(value.NewZeroValue(value.NoneType))
	case "len":
		if len(arguments) != 1 {
			v.errors.Push(fmt.Errorf("incorrect number of arguments for `len`: got %d, want %d", len(arguments), 1))
			return
		}
		arguments[0].Accept(v)
		collection, err := v.valueStack.Pop()
		if err != nil {
			v.errors.Push(fmt.Errorf("%s at %s", err, f))
			return
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
