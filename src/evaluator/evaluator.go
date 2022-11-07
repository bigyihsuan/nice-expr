package evaluator

import (
	"fmt"
	"math/big"
	"nice-expr/src/ast"
	TT "nice-expr/src/token/tokentype"
	"nice-expr/src/util"
	"nice-expr/src/value"
	"strings"

	"golang.org/x/exp/maps"
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
			valType.AddTypeArg(elementType)
		}
		if !v.EqualsType(elementType) {
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
			valType.AddTypeArg(keyType)
		}
		if !k.EqualsType(keyType) {
			return nil, fmt.Errorf("incorrect key type: expected %v, got %v", keyType, k.T)
		}
		vl := valExpr
		v, err := e.EvaluateExpr(vl)
		if err != nil {
			return nil, err
		}
		if inferType && valueType.Name == "UNSET-VALUE" {
			valueType = v.T
			valType.AddTypeArg(valueType)
		}
		if !v.EqualsType(valueType) {
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
		if !idVal.EqualsType(typeArgs[0]) {
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
	switch funcCall.Ident.Name.Lexeme {
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
		case collection.EqualsType(value.StrType):
			val := collection.V.(string)
			return value.NewValue(value.IntType, big.NewInt(int64(len([]rune(val))))), nil
		case collection.T.Is(value.ListType):
			return value.NewValue(value.IntType, big.NewInt(int64(len(collection.V.([]*value.Value))))), nil
		case collection.T.Is(value.MapType):
			return value.NewValue(value.IntType, big.NewInt(int64(len(collection.V.(map[*value.Value]*value.Value))))), nil
		default:
			return nil, fmt.Errorf("invalid type for `len`: %s", collection.T.Name)
		}
	}
	return nil, nil
}

func (e *Evaluator) EvaluateFunctionCall(funcCall *ast.FunctionCall, typeArgs ...value.ValueType) (*value.Value, error) {
	if slices.Contains(BuiltinFunctionNames, funcCall.Ident.Name.Lexeme) {
		return e.EvaluateBuiltinFunction(funcCall, typeArgs...)
	}
	return nil, nil
}

func (e *Evaluator) EvaluateUnaryMinusExpr(unary ast.Expr, typeArgs ...value.ValueType) (*value.Value, error) {
	switch unary := unary.(type) {
	case *ast.UnaryMinusExpr:
		{
			val, err := e.EvaluateUnaryMinusExpr(unary.Right)
			if err != nil {
				return val, err
			}
			if unary.Op != nil {
				switch unary.Op.Tt {
				case TT.Not:
					if !val.EqualsType(value.BoolType) {
						return val, fmt.Errorf("incompatible type for %s: %s", unary.Op.Tt, val.T.Name)
					}
					val.V = !val.V.(bool)
					return val, nil
				case TT.Minus:
					switch {
					case val.EqualsType(value.IntType):
						val.V.(*big.Int).Neg(val.V.(*big.Int))
					case val.EqualsType(value.DecType):
						val.V.(*big.Float).Neg(val.V.(*big.Float))
					default:
						return val, fmt.Errorf("incompatible type for %s: %s", unary.Op.Tt, val.T.Name)
					}
				}
			}
			return val, nil
		}
	default:
		return e.EvaluateExpr(unary)
	}
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
	switch {
	// test
	case slices.Contains(TT.BinLogOps, binary.Op.Tt):
		// bin ops work on 2 bools
		if left.EqualsType(value.BoolType) && right.EqualsType(value.BoolType) {
			l, err := left.Bool()
			if err != nil {
				return nil, err
			}
			r, err := right.Bool()
			if err != nil {
				return nil, err
			}
			var val bool
			switch binary.Op.Tt {
			case TT.And:
				val = l && r
			case TT.Or:
				val = l || r
			default:
				return nil, fmt.Errorf("invalid logical op: got `%s`", binary.Op.Lexeme)
			}
			return value.NewValue(value.BoolType, val), nil
		} else {
			return nil, fmt.Errorf("mismatched types for %s: expected two Bools, got %s and %s", binary.Op.Tt.String(), left.T, right.T)
		}
	// comparison
	case slices.Contains(TT.ComparisonOps, binary.Op.Tt):
		// both sides need to have the same type
		// returns a bool
		if !left.EqualsValueType(right) {
			return nil, fmt.Errorf("mismatched types for %s: got %s and %s", binary.Op.Tt.String(), left.T, right.T)
		}
		val := value.NewValue(value.BoolType, nil)
		switch {
		case left.EqualsType(value.IntType):
			l, _ := left.BigInt()
			r, _ := right.BigInt()
			switch binary.Op.Tt {
			case TT.Equal:
				val.V = l.Cmp(r) == 0
			case TT.Greater:
				val.V = l.Cmp(r) > 0
			case TT.GreaterEqual:
				val.V = l.Cmp(r) >= 0
			case TT.Less:
				val.V = l.Cmp(r) < 0
			case TT.LessEqual:
				val.V = l.Cmp(r) <= 0
			}
		case left.EqualsType(value.DecType):
			l, _ := left.BigDec()
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Equal:
				val.V = l.Cmp(r) == 0
			case TT.Greater:
				val.V = l.Cmp(r) > 0
			case TT.GreaterEqual:
				val.V = l.Cmp(r) >= 0
			case TT.Less:
				val.V = l.Cmp(r) < 0
			case TT.LessEqual:
				val.V = l.Cmp(r) <= 0
			}
		case left.EqualsType(value.StrType):
			l, _ := left.Str()
			r, _ := right.Str()
			switch binary.Op.Tt {
			case TT.Equal:
				val.V = l == r
			case TT.Greater:
				val.V = l > r
			case TT.GreaterEqual:
				val.V = l >= r
			case TT.Less:
				val.V = l < r
			case TT.LessEqual:
				val.V = l <= r
			}
		case left.IsType(value.ListType):
			l := left.V.([]*value.Value)
			r := right.V.([]*value.Value)
			switch binary.Op.Tt {
			case TT.Equal:
				val.V = slices.Equal(l, r)
			default:
				return nil, fmt.Errorf("operation %s not supported on list", binary.Op.Tt)
			}
		case left.IsType(value.MapType):
			l := left.V.(map[*value.Value]*value.Value)
			r := right.V.(map[*value.Value]*value.Value)
			switch binary.Op.Tt {
			case TT.Equal:
				val.V = maps.Equal(l, r)
			default:
				return nil, fmt.Errorf("operation %s not supported on map", binary.Op.Tt)
			}
		}
		return val, nil
	// addExpr
	case slices.Contains(TT.AddOps, binary.Op.Tt):
		// add works on int and dec
		// int + int = int
		// int + dec = dec
		// dec + int = dec
		// dec + dec = dec
		// str + str = str
		// list + list = list
		val := new(value.Value)
		switch {
		case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
			val.T = value.IntType
			l, _ := left.BigInt()
			r, _ := right.BigInt()
			switch binary.Op.Tt {
			case TT.Plus:
				val.V = big.NewInt(0).Add(l, r)
			case TT.Minus:
				val.V = big.NewInt(0).Sub(l, r)
			}
		case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
			val.T = value.DecType
			l, _ := left.BigDec()
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Plus:
				val.V = big.NewFloat(0).Add(l, r)
			case TT.Minus:
				val.V = big.NewFloat(0).Sub(l, r)
			}
		case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
			val.T = value.DecType
			l, _ := left.BigDec()
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Plus:
				val.V = big.NewFloat(0).Add(l, r)
			case TT.Minus:
				val.V = big.NewFloat(0).Sub(l, r)
			}
		case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
			val.T = value.DecType
			l, _ := left.BigDec()
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Plus:
				val.V = big.NewFloat(0).Add(l, r)
			case TT.Minus:
				val.V = big.NewFloat(0).Sub(l, r)
			}
		case left.EqualsType(value.StrType) && right.EqualsType(value.StrType):
			val.T = value.StrType
			l, _ := left.Str()
			r, _ := right.Str()
			switch binary.Op.Tt {
			case TT.Plus:
				val.V = l + r
			case TT.Minus:
				val.V = strings.ReplaceAll(l, r, "")
			}
		case left.IsType(value.ListType) && right.IsType(value.ListType):
			// check subtypes
			if !left.EqualsValueType(right) {
				return nil, fmt.Errorf("mismatched types for `%s`: got %s and %s", binary.Op.Tt.String(), left.T, right.T)
			}
			values := left.V.([]*value.Value)
			r := right.V.([]*value.Value)
			switch binary.Op.Tt {
			case TT.Plus:
				return value.NewValue(left.T, append(values, r...)), nil
			case TT.Minus:
				val := new(value.Value)
				val.T = left.T
				var diff []*value.Value
				for _, x := range values {
					for _, y := range r {
						if !x.Equal(y) {
							diff = append(diff, x)
						}
					}
				}
				val.V = diff
				return val, nil
			}
		default:
			return nil, fmt.Errorf("invalid types for `%s`: got %s and %s", binary.Op.Tt.String(), left.T, right.T)
		}
		return val, nil

	// mulExpr
	case slices.Contains(TT.MulOps, binary.Op.Tt):
		// mul works on int and dec
		// int * int = int
		// int * dec = dec
		// dec * int = dec
		// dec * dec = dec
		// additional check: div and mod cannot have the right == 0
		// additional check: mod only works on int and int
		val := new(value.Value)
		switch {
		case left.EqualsType(value.IntType) && right.EqualsType(value.IntType):
			val.T = value.IntType
			l, _ := left.BigInt()
			r, _ := right.BigInt()
			switch binary.Op.Tt {
			case TT.Star:
				val.V = big.NewInt(0).Mul(l, r)
			case TT.Slash:
				if r.Cmp(big.NewInt(0)) == 0 {
					return nil, fmt.Errorf("division by 0: %s / %s", left.V, right.V)
				}
				val.V = big.NewInt(0).Div(l, r)
			case TT.Percent:
				if r.Cmp(big.NewInt(0)) == 0 {
					return nil, fmt.Errorf("modulo by 0: %s %% %s", left.V, right.V)
				}
				val.V = big.NewInt(0).Mod(l, r)
			}
		case left.EqualsType(value.DecType) && right.EqualsType(value.IntType):
			val.T = value.DecType
			l, _ := left.BigDec()
			d, _ := right.BigInt()
			r := new(big.Float).SetInt(d)
			switch binary.Op.Tt {
			case TT.Star:
				val.V = big.NewFloat(0.0).Mul(l, r)
			case TT.Slash:
				if r.Cmp(big.NewFloat(0)) == 0 {
					return nil, fmt.Errorf("division by 0: %s / %s", left.V, right.V)
				}
				val.V = big.NewFloat(0.0).Quo(l, r)
			case TT.Percent:
				return nil, fmt.Errorf("modulo not allowed for Dec")
			}
		case left.EqualsType(value.IntType) && right.EqualsType(value.DecType):
			val.T = value.DecType
			d, _ := left.BigInt()
			l := new(big.Float).SetInt(d)
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Star:
				val.V = big.NewFloat(0.0).Mul(l, r)
			case TT.Slash:
				if r.Cmp(big.NewFloat(0)) == 0 {
					return nil, fmt.Errorf("division by 0: %s / %s", left.V, right.V)
				}
				val.V = big.NewFloat(0.0).Quo(l, r)
			case TT.Percent:
				return nil, fmt.Errorf("modulo not allowed for Dec")
			}
		case left.EqualsType(value.DecType) && right.EqualsType(value.DecType):
			val.T = value.DecType
			l, _ := left.BigDec()
			r, _ := right.BigDec()
			switch binary.Op.Tt {
			case TT.Star:
				val.V = big.NewFloat(0.0).Mul(l, r)
			case TT.Slash:
				if r.Cmp(big.NewFloat(0)) == 0 {
					return nil, fmt.Errorf("division by 0: %s / %s", left.V, right.V)
				}
				val.V = big.NewFloat(0.0).Quo(l, r)
			case TT.Percent:
				return nil, fmt.Errorf("modulo not allowed for Dec")
			}
		default:
			return nil, fmt.Errorf("invalid types for `%s`: got %s and %s", binary.Op.Tt.String(), left.T, right.T)
		}
		return val, nil
	// indexing
	case binary.Op.Is(TT.Underscore):
		if !left.T.IsIndexable() {
			return nil, fmt.Errorf("left-hand side is not indexable: %s", left.T)
		}
		switch {
		case left.IsType(value.StrType):
			if right.IsNotType(value.IntType) {
				return value.NewZeroValue(left.T), fmt.Errorf("invalid type for indexing Str: expected Int, got %v", right.T)
			}
			// right is an int
			s, _ := left.Str()
			rs := []rune(s)
			i, _ := right.Int()
			idx := int(i)
			if idx >= len(rs) {
				return value.NewZeroValue(left.T), fmt.Errorf("index out of range: idx %v > len %v", idx, len(rs))
			}
			return value.NewValue(left.T, string(rs[idx])), nil
		case left.IsType(value.ListType):
			if right.IsNotType(value.IntType) {
				return value.NewZeroValue(left.T), fmt.Errorf("invalid type for indexing List: expected Int, got %v", right.T)
			}
			// right is an int
			l, _ := left.List()
			i, _ := right.Int()
			idx := int(i)
			if idx >= len(l) {
				return value.NewZeroValue(left.T), fmt.Errorf("index out of range: idx %v > len %v", idx, len(l))
			}
			return value.NewValue(left.T.TypeArgs[0], l[idx].V), nil
		case left.IsType(value.MapType):
			// check if right is the same type as the key
			if right.IsNotType(left.T.TypeArgs[0]) {
				return value.NewZeroValue(left.T.TypeArgs[1]), fmt.Errorf("invalid type for indexing %v: expected %v, got %v", left.T, left.T.TypeArgs[0], right.T)
			}
			// right is a key type
			m, _ := left.Map()
			for key, val := range m {
				if key.Equal(right) {
					return value.NewValue(left.T.TypeArgs[1], val.V), nil
				}
			}
			// return zero value if not found
			return value.NewZeroValue(left.T.TypeArgs[1]), nil
		}
	default:
		return nil, fmt.Errorf("unknown op: got `%s` (%s)", binary.Op.Lexeme, binary.Op.Tt.String())
	}
	return nil, fmt.Errorf("some other error: %s", binary.String())
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
	case *ast.UnaryMinusExpr:
		return e.EvaluateUnaryMinusExpr(expr, typeArgs...)
	case *ast.BinaryExpr:
		return e.EvaluateBinaryExpr(expr, typeArgs...)
	case *ast.ConstantDeclaration, *ast.VariableDeclaration:
		return e.EvaluateDeclaration(expr)
	}
	return nil, fmt.Errorf("unknown expr %v of type %T", expr, expr)
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
