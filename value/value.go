package value

import "golang.org/x/exp/slices"

type Value struct {
	T ValueType
	V interface{}
}

type ValueType struct {
	Name     string
	TypeArgs []ValueType
}

func NewValueType(name string) (vt ValueType) {
	vt.Name = name
	vt.TypeArgs = []ValueType{}
	return
}

func (vt ValueType) Equal(other ValueType) bool {
	return vt.Name == other.Name && slices.EqualFunc(vt.TypeArgs, other.TypeArgs, func(l, r ValueType) bool {
		return l.Equal(r)
	})
}
