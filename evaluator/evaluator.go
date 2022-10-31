package evaluator

import (
	"fmt"
	"nice-expr/ast"
	"nice-expr/token/tokentype"
	"nice-expr/util"
	"nice-expr/value"
)

type Evaluator struct {
	Constants, Variables map[*ast.Identifier]*value.Value
	ValueStack           util.Stack[*value.Value]
}

func NewEvaluator() Evaluator {
	var evaluator Evaluator
	evaluator.Constants = make(map[*ast.Identifier]*value.Value)
	evaluator.Variables = make(map[*ast.Identifier]*value.Value)
	evaluator.ValueStack = util.Stack[*value.Value]{}
	return evaluator
}

func (e Evaluator) GetConstant(name string) (*ast.Identifier, *value.Value) {
	for ident, val := range e.Constants {
		if ident.Name.Lexeme == name {
			return ident, val
		}
	}
	return nil, nil
}
func (e Evaluator) GetVariable(name string) (*ast.Identifier, *value.Value) {
	for ident, val := range e.Variables {
		if ident.Name.Lexeme == name {
			return ident, val
		}
	}
	return nil, nil
}

func (e *Evaluator) EvaluatePrimitiveLiteral(literal *ast.PrimitiveLiteral) {
	val := new(value.Value)
	var valType = tokentype.LitToType[literal.Token.Tt]

	val.T = valType
	val.V = literal.Token.Value
	e.ValueStack.Push(val)
}

// func (e *Evaluator) EvaluateListLiteral(literal *ast.ListLiteral) {
// 	var value value.Value
// 	var valType value.ValueType

// 	// valType.Name = literal.Token.Tt.String()
// 	// value.V = literal.Token.value.Value
// 	value.T = valType
// 	e.ValueStack.Push(value)
// }

func (e *Evaluator) EvaluateLiteral(litType value.ValueType, litExpr ast.Expr) (*value.Value, error) {
	e.EvaluatePrimitiveLiteral(litExpr.(*ast.PrimitiveLiteral))
	v, err := e.ValueStack.Pop()
	if err != nil {
		return v, err
	}
	// check desired type versus actual type
	if !litType.Equal(v.T) {
		return v, fmt.Errorf("types don't match: %v and %v", litType, v.T)
	}
	return v, nil
}

func (e *Evaluator) EvaluateDeclaration(decl ast.Declaration) error {
	switch decl := decl.(type) {
	case *ast.VariableDeclaration:
		v, err := e.EvaluateLiteral(decl.Type.ToValueType(), decl.Value)
		if err != nil {
			return err
		}
		e.Variables[decl.Name] = v
	case *ast.ConstantDeclaration:
		v, err := e.EvaluateLiteral(decl.Type.ToValueType(), decl.Value)
		if err != nil {
			return err
		}
		e.Constants[decl.Name] = v
	default:
		return fmt.Errorf("not declaration: %T", decl)
	}
	return nil
}

func (e *Evaluator) EvaluateProgram(program ast.Program) error {
	var err error

	for _, stmt := range program.Statements {
		// TODO: Expand to all exprs
		err = e.EvaluateDeclaration(stmt.(ast.Declaration))
		if err != nil {
			return err
		}
	}
	return nil
}
