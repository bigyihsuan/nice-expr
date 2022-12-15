package value

import (
	"fmt"
	"nice-expr/src/token/tokentype"

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

	IndexableTypes  = []ValueType{StrType, ListType, MapType}
	ComparableTypes = []ValueType{IntType, DecType, StrType}
)

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

// compares to another type.
// only checks if the base type is the same.
func (vt ValueType) IsNot(other ValueType) bool {
	return vt.Name != other.Name
}

// is this type an indexable type?
func (vt ValueType) IsIndexable() bool {
	for _, t := range IndexableTypes {
		if vt.Is(t) {
			return true
		}
	}
	return false
}

// does this type allow for non-equal comparison?
func (vt ValueType) IsComparable() bool {
	for _, t := range ComparableTypes {
		if vt.Is(t) {
			return true
		}
	}
	return false
}
