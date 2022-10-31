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
	return fmt.Sprint(v.V)
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

func (vt ValueType) Equal(other ValueType) bool {
	return vt.Name == other.Name && slices.EqualFunc(vt.TypeArgs, other.TypeArgs, func(l, r ValueType) bool {
		return l.Equal(r)
	})
}
