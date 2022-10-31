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

func (e *Evaluator) EvaluatePrimitiveLiteral(literal *ast.PrimitiveLiteral) (*value.Value, error) {
	val := new(value.Value)
	valType, ok := tokentype.LitToType[literal.Token.Tt]
	if !ok {
		return nil, fmt.Errorf("unkown primitive literal %v", literal.Token.Tt)
	}

	val.T = valType
	val.V = literal.Token.Value
	return val, nil
}

func (e *Evaluator) EvaluateListLiteral(elementType value.ValueType, literal *ast.ListLiteral) (*value.Value, error) {
	val := new(value.Value)
	valType := value.NewValueType("List")
	valType.AddTypeArg(elementType)

	elements := []*value.Value{}

	for _, listVal := range literal.Values {
		// TODO: expand to all exprs
		lv := listVal.(*ast.PrimitiveLiteral)
		v, err := e.EvaluatePrimitiveLiteral(lv)
		if err != nil {
			return nil, err
		}
		if !v.T.Equal(elementType) {
			return nil, fmt.Errorf("incorrect element type: expected %v, got %v", elementType, v.T)
		}
		elements = append(elements, v)
	}
	val.V = elements
	val.T = valType
	return val, nil
}
func (e *Evaluator) EvaluateMapLiteral(keyType, elementType value.ValueType, literal *ast.MapLiteral) (*value.Value, error) {
	val := new(value.Value)
	valType := value.NewValueType("Map")
	valType.AddTypeArg(keyType, elementType)

	elements := make(map[*value.Value]*value.Value)

	for ke, valExpr := range literal.Values {
		// TODO: expand to all exprs
		kl := ke.(*ast.PrimitiveLiteral)
		k, err := e.EvaluatePrimitiveLiteral(kl)
		if err != nil {
			return nil, err
		}
		if !k.T.Equal(keyType) {
			return nil, fmt.Errorf("incorrect element type: expected %v, got %v", keyType, k.T)
		}
		vl := valExpr.(*ast.PrimitiveLiteral)
		v, err := e.EvaluatePrimitiveLiteral(vl)
		if err != nil {
			return nil, err
		}
		if !v.T.Equal(elementType) {
			return nil, fmt.Errorf("incorrect element type: expected %v, got %v", elementType, v.T)
		}
		elements[k] = v
	}
	val.V = elements
	val.T = valType
	return val, nil
}

func (e *Evaluator) EvaluateLiteral(litType value.ValueType, litExpr ast.Expr) (*value.Value, error) {
	var v *value.Value
	var err error
	switch litExpr := litExpr.(type) {
	case *ast.PrimitiveLiteral:
		v, err = e.EvaluatePrimitiveLiteral(litExpr)
	case *ast.ListLiteral:
		if len(litType.TypeArgs) < 1 {
			return nil, fmt.Errorf("not enough typeargs for List: got %v, want 1", len(litType.TypeArgs))
		}
		v, err = e.EvaluateListLiteral(litType.TypeArgs[0], litExpr)
	case *ast.MapLiteral:
		if len(litType.TypeArgs) < 2 {
			return nil, fmt.Errorf("not enough typeargs for Map: got %v, want 2", len(litType.TypeArgs))
		}
		v, err = e.EvaluateMapLiteral(litType.TypeArgs[0], litType.TypeArgs[1], litExpr)
	}

	if err != nil {
		return v, err
	}
	// check desired type versus actual type
	if !litType.Equal(v.T) {
		return v, fmt.Errorf("types don't match: want %v, got %v", litType, v.T)
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
