package value

import (
	"fmt"
	"math/big"
	"strings"
)

type Value struct {
	T ValueType
	V interface{}
}

func NewValue(vt ValueType, val interface{}) *Value {
	v := new(Value)
	v.T = vt
	v.V = val
	return v
}

func NewZeroValue(vt ValueType) *Value {
	switch {
	case vt.Is(IntType):
		return NewValue(vt, big.NewInt(0))
	case vt.Is(DecType):
		return NewValue(vt, big.NewFloat(0.0))
	case vt.Is(StrType):
		return NewValue(vt, "")
	case vt.Is(ListType):
		return NewValue(vt, []*Value{})
	case vt.Is(MapType):
		return NewValue(vt, make(map[*Value]*Value))
	case vt.Is(NoneType):
		return NewValue(vt, "None")
	default:
		return nil
	}
}

// if the value is a return/break value,
// unwrap the type so that it can be used elsewhere
func (v *Value) UnwrapReturn() *Value {
	if v.IsType(ReturnedType) || v.IsType(BrokeType) {
		v.T = v.T.TypeArgs[0]
	}
	return v
}

func (v Value) String() string {
	return fmt.Sprintf("{%s %v}", v.T, v.V)
}

func (v Value) Sprint() string {
	var b strings.Builder
	if v.IsType(ListType) {
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

func (v Value) Int64() (int64, error) {
	if val, ok := v.V.(*big.Int); ok {
		return val.Int64(), nil
	}
	return NewZeroValue(v.T).V.(int64), fmt.Errorf("incorrect type for Int: %s", v.T.String())
}
func (v Value) Dec() (float64, error) {
	if val, ok := v.V.(*big.Float); ok {
		valf, _ := val.Float64()
		return valf, nil
	}
	return NewZeroValue(v.T).V.(float64), fmt.Errorf("incorrect type for Dec: %s", v.T.String())
}
func (v Value) BigInt() (*big.Int, error) {
	if val, ok := v.V.(*big.Int); ok {
		return val, nil
	}
	// v is a float
	i, _ := v.V.(*big.Float).Int(nil)
	return i, fmt.Errorf("incorrect type for Int: %s", v.T.String())
}
func (v Value) BigDec() (*big.Float, error) {
	if val, ok := v.V.(*big.Float); ok {
		return val, nil
	}
	// v is an int
	f := big.NewFloat(0).SetInt(v.V.(*big.Int))
	return f, fmt.Errorf("incorrect type for Dec: %s", v.T.String())
}
func (v Value) Str() (string, error) {
	if val, ok := v.V.(string); ok {
		return val, nil
	}
	return NewZeroValue(v.T).V.(string), fmt.Errorf("incorrect type for Str: %s", v.T.String())
}
func (v Value) Bool() (bool, error) {
	if val, ok := v.V.(bool); ok {
		return val, nil
	}
	return NewZeroValue(v.T).V.(bool), fmt.Errorf("incorrect type for Bool: %s", v.T.String())
}
func (v Value) List() ([]*Value, error) {
	if val, ok := v.V.([]*Value); ok {
		return val, nil
	}
	return NewZeroValue(v.T).V.([]*Value), fmt.Errorf("incorrect type for List: %s", v.T.String())
}
func (v Value) Map() (map[*Value]*Value, error) {
	if val, ok := v.V.(map[*Value]*Value); ok {
		return val, nil
	}
	return NewZeroValue(v.T).V.(map[*Value]*Value), fmt.Errorf("incorrect type for Map: %s", v.T.String())
}

// is this equal to another value?
// returns false if the types don't match.
func (v Value) Equal(other *Value) bool {
	if !v.EqualsValueType(other) {
		return false
	}
	// different ways of comparing based on the base type
	switch {
	case v.IsType(IntType):
		l, _ := v.BigInt()
		r, _ := other.BigInt()
		return l.Cmp(r) == 0
	case v.IsType(DecType):
		l, _ := v.BigDec()
		r, _ := other.BigDec()
		return l.Cmp(r) == 0
	case v.IsType(StrType):
		l, _ := v.Str()
		r, _ := other.Str()
		return l == r
	case v.IsType(BoolType):
		l, _ := v.Bool()
		r, _ := v.Bool()
		return l == r
	case v.IsType(ListType):
		l, _ := v.List()
		r, _ := v.List()
		if len(l) != len(r) {
			return false
		}
		for i := range l {
			if !l[i].Equal(r[i]) {
				return false
			}
		}
		return true
	case v.IsType(MapType):
		l, _ := v.Map()
		r, _ := v.Map()
		if len(l) != len(r) {
			return false
		}
		for key := range l {
			if !l[key].Equal(r[key]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// checks if another value has a type **exactly*** equal to this value's type.
func (v Value) EqualsValueType(other *Value) bool {
	return v.EqualsType(other.T)
}

// checks if another value has the same base type as this value's base type.
func (v Value) IsValueType(other *Value) bool {
	return v.IsType(other.T)
}

// checks if this value has a type **exactly** equal to a given type.
func (v Value) EqualsType(vt ValueType) bool {
	return v.T.Equal(vt)
}

// checks if this value has a **base type** equal to a given type.
func (v Value) IsType(vt ValueType) bool {
	return v.T.Is(vt)
}

// checks if this value has a **base type** not equal to a given type.
func (v Value) IsNotType(vt ValueType) bool {
	return v.T.IsNot(vt)
}
