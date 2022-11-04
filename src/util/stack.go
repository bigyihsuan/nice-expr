package util

import "fmt"

type Stack[T any] struct {
	stack []T
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var v T
	if s.Len() <= 0 {
		return v, fmt.Errorf("no more elements to pop")
	}
	v, s.stack = s.stack[s.Len()-1], s.stack[:s.Len()-1]
	return v, nil
}

func (s Stack[T]) Len() int {
	return len(s.stack)
}
