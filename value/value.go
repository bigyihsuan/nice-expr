package value

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Value struct {
	T ValueType
	V interface{}
}

func (v Value) String() string {
	return fmt.Sprintf("{%s %v}", v.T, v.V)
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

func (vt ValueType) Equal(other ValueType) bool {
	return vt.Name == other.Name &&
		slices.EqualFunc(
			vt.TypeArgs,
			other.TypeArgs,
			func(l, r ValueType) bool { return l.Equal(r) },
		)
}

func (vt ValueType) NotEqual(other ValueType) bool {
	return !vt.Equal(other)
}
