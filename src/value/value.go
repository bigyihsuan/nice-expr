package value

import (
	"fmt"
	"math/big"
	"nice-expr/src/token/tokentype"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	LitToType = func() map[tokentype.TokenType]ValueType {
		m := make(map[tokentype.TokenType]ValueType)
		for i := range tokentype.PrimitiveTypes {
			m[tokentype.PrimitiveLiterals[i]] = NewValueType(tokentype.PrimitiveTypes[i].String())
		}
		// manually add true and false
		m[tokentype.True] = BoolType
		m[tokentype.False] = BoolType
		return m
	}()
	NoneType = NewValueType("None")
	IntType  = NewValueType("Int")
	DecType  = NewValueType("Dec")
	StrType  = NewValueType("Str")
	BoolType = NewValueType("Bool")
	ListType = NewValueType("List")
	MapType  = NewValueType("Map")
)

type Value struct {
	T ValueType
	V interface{}
}

func (v Value) String() string {
	return fmt.Sprintf("{%s %v}", v.T, v.V)
}

func (v Value) Sprint() string {
	var b strings.Builder
	if v.T.Is(ListType) {
		b.WriteRune('[')
		for _, e := range v.V.([]*Value) {
			b.WriteString(e.Sprint())
			b.WriteRune(',')
		}
		b.WriteRune(']')
		return b.String()
	} else if v.T.Is(MapType) {
		b.WriteString("<|")
		m := v.V.(map[*Value]*Value)
		for key, val := range m {
			b.WriteString(key.Sprint())
			b.WriteRune(':')
			b.WriteString(val.Sprint())
			b.WriteRune(',')
		}
		b.WriteString("|>")
		return b.String()
	}
	return fmt.Sprint(v.V)
}

func (v Value) Int() (int64, error) {
	if val, ok := v.V.(*big.Int); ok {
		return val.Int64(), nil
	}
	return 0, fmt.Errorf("incorrect type for Int: %s", v.T.String())
}
func (v Value) Dec() (float64, error) {
	if val, ok := v.V.(*big.Float); ok {
		valf, _ := val.Float64()
		return valf, nil
	}
	return 0.0, fmt.Errorf("incorrect type for Dec: %s", v.T.String())
}
func (v Value) Str() (string, error) {
	if val, ok := v.V.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("incorrect type for Str: %s", v.T.String())
}

type ValueType struct {
	Name     string
	TypeArgs []ValueType
}

func NewValueType(name string, typeargs ...string) (vt ValueType) {
	vt.Name = name
	vt.TypeArgs = []ValueType{}
	for _, t := range typeargs {
		vt.TypeArgs = append(vt.TypeArgs, NewValueType(t))
	}
	return
}

func (vt *ValueType) AddTypeArg(typeargs ...ValueType) *ValueType {
	vt.TypeArgs = append(vt.TypeArgs, typeargs...)
	return vt
}

func (vt ValueType) String() string {
	switch len(vt.TypeArgs) {
	case 0:
		return vt.Name
	case 1:
		return fmt.Sprintf("%s[%s]", vt.Name, vt.TypeArgs[0])
	case 2:
		return fmt.Sprintf("%s[%s]%s", vt.Name, vt.TypeArgs[0], vt.TypeArgs[1])
	default:
		return fmt.Sprintf("%s%s", vt.Name, vt.TypeArgs)
	}
}

// compares to another type.
// checks for deep equality, i.e. the type arguments are also equal.
func (vt ValueType) Equal(other ValueType) bool {
	return vt.Name == other.Name &&
		slices.EqualFunc(
			vt.TypeArgs,
			other.TypeArgs,
			func(l, r ValueType) bool { return l.Equal(r) },
		)
}

// compares to another type.
// checks for deep equality, i.e. the type arguments are also equal.
func (vt ValueType) NotEqual(other ValueType) bool {
	return !vt.Equal(other)
}

// compares to another type.
// only checks if the base type is the same.
func (vt ValueType) Is(other ValueType) bool {
	return vt.Name == other.Name
}
