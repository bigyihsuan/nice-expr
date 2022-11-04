package evaluator

import (
	"fmt"
	"math/big"
	"nice-expr/src/ast"
	TT "nice-expr/src/token/tokentype"
	"nice-expr/src/util"
	"nice-expr/src/value"

	"golang.org/x/exp/slices"
)

//go:generate stringer -type=IdentifierType
type IdentifierType int

const (
	InvalidIdentifier IdentifierType = iota - 1
	ConstantIdentifier
	VariableIdentifier
	FunctionIdentifier
)

var (
	BuiltinFunctionNames = []string{"print", "println", "len"}
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

func (e Evaluator) GetIdentifier(name string) (*ast.Identifier, *value.Value, IdentifierType) {
	if con, conval := e.GetConstant(name); con != nil {
		return con, conval, ConstantIdentifier
	}
	if v, vval := e.GetVariable(name); v != nil {
		return v, vval, VariableIdentifier
	}
	return nil, nil, InvalidIdentifier
}

func (e *Evaluator) EvaluatePrimitiveLiteral(literal *ast.PrimitiveLiteral, typeArgs ...value.ValueType) (*value.Value, error) {
	val := new(value.Value)
	valType, ok := value.LitToType[literal.Token.Tt]
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
		if len(typeArgs[0].TypeArgs) < 1 {
			return nil, fmt.Errorf("not enough typeargs for List: want 1, got %v", len(typeArgs[0].TypeArgs))
		}
		elementType = typeArgs[0].TypeArgs[0]
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
		if v.T.NotEqual(elementType) {
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

	keyType := value.NewValueType("UNSET-KEY")
	valueType := value.NewValueType("UNSET-VALUE")

	inferType := len(typeArgs) < 1
	if !inferType {
		if len(typeArgs[0].TypeArgs) < 2 {
			return nil, fmt.Errorf("not enough typeargs for Map: got %v, want 2", len(typeArgs))
		}
		keyType = typeArgs[0].TypeArgs[0]
		valueType = typeArgs[0].TypeArgs[1]
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
		if k.T.NotEqual(keyType) {
			return nil, fmt.Errorf("incorrect key type: expected %v, got %v", keyType, k.T)
		}
		vl := valExpr
		v, err := e.EvaluateExpr(vl)
		if err != nil {
			return nil, err
		}
		if inferType && valueType.Name == "UNSET-VALUE" {
			valueType = v.T
		}
		if v.T.NotEqual(valueType) {
			return nil, fmt.Errorf("incorrect value type: expected %v, got %v", valueType, v.T)
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

	switch litExpr := litExpr.(type) {
	case *ast.PrimitiveLiteral:
		v, err = e.EvaluatePrimitiveLiteral(litExpr, typeArgs...)
	case *ast.ListLiteral:
		v, err = e.EvaluateListLiteral(litExpr, typeArgs...)
	case *ast.MapLiteral:
		v, err = e.EvaluateMapLiteral(litExpr, typeArgs...)
	}
	if err != nil {
		return v, err
	}
	return v, nil
}

func (e *Evaluator) EvaluateIdentifier(ident *ast.Identifier, typeArgs ...value.ValueType) (*value.Value, error) {
	id, idVal, idKind := e.GetIdentifier(ident.Name.Lexeme)
	if idKind == InvalidIdentifier {
		return nil, fmt.Errorf("identifier not found: %s", ident.Name.Lexeme)
	}
	if len(typeArgs) > 0 {
		// check type
		if idVal.T.NotEqual(typeArgs[0]) {
			return nil, fmt.Errorf("value type of `%s` and desired type don't match: got %s, want %s", id.Name.Lexeme, idVal.T, typeArgs[0])
		}
	}
	return idVal, nil
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

func (e *Evaluator) EvaluateBuiltinFunction(funcCall *ast.FunctionCall, typeArgs ...value.ValueType) (*value.Value, error) {
	switch funcCall.Name.Lexeme {
	case "print":
		if len(funcCall.Arguments) < 1 {
			fmt.Print()
		} else {
			for _, ex := range funcCall.Arguments {
				val, err := e.EvaluateExpr(ex)
				if err != nil {
					return val, err
				}
				fmt.Print(val.Sprint())
			}
		}
	case "println":
		if len(funcCall.Arguments) < 1 {
			fmt.Println()
		} else {
			for _, ex := range funcCall.Arguments {
				val, err := e.EvaluateExpr(ex)
				if err != nil {
					return val, err
				}
				fmt.Println(val.Sprint())
			}
		}
	case "len":
		if len(funcCall.Arguments) != 1 {
			return nil, fmt.Errorf("incorrect number of arguments for `len`: got %d, want %d", len(funcCall.Arguments), 1)
		}
		collection, err := e.EvaluateExpr(funcCall.Arguments[0])
		if err != nil {
			return nil, err
		}
		switch {
		case collection.T.Equal(value.StrType):
			val := collection.V.(string)
			return &value.Value{
				T: value.IntType,
				V: big.NewInt(int64(len([]rune(val)))),
			}, nil
		case collection.T.Is(value.ListType):
			return &value.Value{
				T: value.IntType,
				V: big.NewInt(int64(len(collection.V.([]*value.Value)))),
			}, nil
		case collection.T.Is(value.MapType):
			return &value.Value{
				T: value.IntType,
				V: big.NewInt(int64(len(collection.V.(map[*value.Value]*value.Value)))),
			}, nil
		default:
			return nil, fmt.Errorf("invalid type for `len`: %s", collection.T.Name)
		}
	}
	return nil, nil
}

func (e *Evaluator) EvaluateFunctionCall(funcCall *ast.FunctionCall, typeArgs ...value.ValueType) (*value.Value, error) {
	if slices.Contains(BuiltinFunctionNames, funcCall.Name.Lexeme) {
		return e.EvaluateBuiltinFunction(funcCall, typeArgs...)
	}
	return nil, nil
}

func (e *Evaluator) EvaluateUnaryExpr(unary *ast.UnaryExpr, typeArgs ...value.ValueType) (*value.Value, error) {
	val, err := e.EvaluateExpr(unary.Right)
	if err != nil {
		return val, err
	}
	if unary.Op != nil {
		switch unary.Op.Tt {
		case TT.Not:
			if val.T.NotEqual(value.BoolType) {
				return val, fmt.Errorf("incompatible type for %s: %s", unary.Op.Tt, val.T.Name)
			}
			val.V = !val.V.(bool)
			return val, nil
		case TT.Minus:
			switch {
			case val.T.Equal(value.IntType):
				val.V.(*big.Int).Neg(val.V.(*big.Int))
			case val.T.Equal(value.DecType):
				val.V.(*big.Float).Neg(val.V.(*big.Float))
			default:
				return val, fmt.Errorf("incompatible type for %s: %s", unary.Op.Tt, val.T.Name)
			}
		}
	}
	return val, nil
}
func (e *Evaluator) EvaluateBinaryExpr(binary *ast.BinaryExpr, typeArgs ...value.ValueType) (*value.Value, error) {
	left, err := e.EvaluateExpr(binary.Left)
	if err != nil {
		return left, err
	}
	right, err := e.EvaluateExpr(binary.Right)
	if err != nil {
		return right, err
	}
	if binary.Op == nil {
		return nil, fmt.Errorf("missing operation for binary expr")
	}
	switch binary.Op.Tt {
	// arithmetic
	case TT.Plus, TT.Minus, TT.Star, TT.Slash, TT.Percent:
	// comparisons
	case TT.Equal, TT.Greater, TT.GreaterEqual, TT.Less, TT.LessEqual:
		// logical
	case TT.And, TT.Or:
		// indexing
	case TT.Underscore:
	default:
		return nil, fmt.Errorf("unknown op: got `%s`", binary.Op.Lexeme)
	}
	return nil, fmt.Errorf("unknown op: got `%s`", binary.Op.Lexeme)
}

func (e *Evaluator) EvaluateExpr(expr ast.Expr, typeArgs ...value.ValueType) (*value.Value, error) {
	switch expr := expr.(type) {
	// leaves
	case *ast.PrimitiveLiteral, *ast.ListLiteral, *ast.MapLiteral:
		return e.EvaluateLiteral(expr, typeArgs...)
	case *ast.Identifier:
		return e.EvaluateIdentifier(expr, typeArgs...)
	// branches
	case *ast.FunctionCall:
		return e.EvaluateFunctionCall(expr, typeArgs...)
	case *ast.UnaryExpr:
		return e.EvaluateUnaryExpr(expr, typeArgs...)
	case *ast.BinaryExpr:
		return e.EvaluateBinaryExpr(expr, typeArgs...)
	case *ast.ConstantDeclaration, *ast.VariableDeclaration:
		return e.EvaluateDeclaration(expr)
	}
	return nil, fmt.Errorf("unknown expr %v", expr)
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
