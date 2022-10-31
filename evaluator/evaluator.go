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

func (e *Evaluator) EvaluateListLiteral(literal *ast.ListLiteral, typeArgs ...value.ValueType) (*value.Value, error) {
	val := new(value.Value)
	valType := value.NewValueType("List")

	elementType := value.NewValueType("UNSET")
	inferType := len(typeArgs) < 1
	if !inferType {
		elementType = typeArgs[0]
		valType.AddTypeArg(elementType)
	}

	elements := []*value.Value{}

	for _, listVal := range literal.Values {
		lv := listVal
		v, err := e.EvaluateExpr(lv)
		if err != nil {
			return nil, err
		}
		if inferType && elementType.Name == "UNSET" {
			elementType = v.T
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
func (e *Evaluator) EvaluateMapLiteral(literal *ast.MapLiteral, typeArgs ...value.ValueType) (*value.Value, error) {
	val := new(value.Value)
	valType := value.NewValueType("Map")

	keyType, valueType := value.NewValueType("UNSET-KEY"), value.NewValueType("UNSET-VALUE")
	inferType := len(typeArgs) < 2
	if !inferType {
		keyType = typeArgs[0]
		valueType = typeArgs[1]
		valType.AddTypeArg(keyType)
		valType.AddTypeArg(valueType)
	}

	elements := make(map[*value.Value]*value.Value)

	for ke, valExpr := range literal.Values {
		kl := ke
		k, err := e.EvaluateExpr(kl)
		if err != nil {
			return nil, err
		}
		if inferType && keyType.Name == "UNSET-KEY" {
			keyType = k.T
		}
		if !k.T.Equal(keyType) {
			return nil, fmt.Errorf("incorrect element type: expected %v, got %v", keyType, k.T)
		}
		vl := valExpr
		v, err := e.EvaluateExpr(vl)
		if err != nil {
			return nil, err
		}
		if inferType && valueType.Name == "UNSET-VALUE" {
			valueType = v.T
		}
		if !v.T.Equal(valueType) {
			return nil, fmt.Errorf("incorrect element type: expected %v, got %v", valueType, v.T)
		}
		elements[k] = v
	}
	val.V = elements
	val.T = valType
	return val, nil
}

func (e *Evaluator) EvaluateLiteral(litExpr ast.Expr, typeArgs ...value.ValueType) (*value.Value, error) {
	var v *value.Value
	var err error
	// inferType := len(typeArgs) < 1
	// litType := value.NewValueType("UNSET")
	// if !inferType {
	// 	litType = typeArgs[0]
	// }

	switch litExpr := litExpr.(type) {
	case *ast.PrimitiveLiteral:
		v, err = e.EvaluatePrimitiveLiteral(litExpr)
	case *ast.ListLiteral:
		if len(typeArgs) < 1 {
			return nil, fmt.Errorf("not enough typeargs for List: got %v, want 1", len(typeArgs))
		}
		v, err = e.EvaluateListLiteral(litExpr, typeArgs...)
	case *ast.MapLiteral:
		if len(typeArgs) < 2 {
			return nil, fmt.Errorf("not enough typeargs for Map: got %v, want 2", len(typeArgs))
		}
		v, err = e.EvaluateMapLiteral(litExpr, typeArgs...)
	}

	if err != nil {
		return v, err
	}
	// if inferType && litType.Name == "UNSET" {
	// 	litType = v.T
	// }
	// // check desired type versus actual type
	// if !litType.Equal(v.T) {
	// 	return v, fmt.Errorf("types don't match: want %v, got %v", litType, v.T)
	// }
	return v, nil
}

func (e *Evaluator) EvaluateDeclaration(decl ast.Declaration) (*value.Value, error) {
	var v *value.Value
	var err error
	switch decl := decl.(type) {
	case *ast.VariableDeclaration:
		v, err = e.EvaluateExpr(decl.Value, decl.Type.ToValueType())
		if err != nil {
			return nil, err
		}
		e.Variables[decl.Name] = v
	case *ast.ConstantDeclaration:
		v, err = e.EvaluateExpr(decl.Value, decl.Type.ToValueType())
		if err != nil {
			return nil, err
		}
		e.Constants[decl.Name] = v
	default:
		return nil, fmt.Errorf("not declaration: %T", decl)
	}
	return v, nil
}

func (e *Evaluator) EvaluateExpr(expr ast.Expr, typeArgs ...value.ValueType) (*value.Value, error) {
	switch expr := expr.(type) {
	case *ast.PrimitiveLiteral, *ast.ListLiteral, *ast.MapLiteral:
		return e.EvaluateLiteral(expr, typeArgs...)
	case ast.Declaration:
		return e.EvaluateDeclaration(expr)
	}
	return nil, nil
}

func (e *Evaluator) EvaluateProgram(program ast.Program) error {
	var err error

	for _, stmt := range program.Statements {
		_, err = e.EvaluateExpr(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
