package util

import "nice-expr/src/value"

type TestCase struct {
	Code         string
	Expected     interface{}
	ExpectedType value.ValueType
}
