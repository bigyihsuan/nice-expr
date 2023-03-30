package util

import "fmt"

type Stack[T any] struct {
	stack []T
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
	// fmt.Println("Push", s.Len(), s.Peek())
}

func (s *Stack[T]) Pop() (T, error) {
	var v T
	if s.Len() <= 0 {
		return v, fmt.Errorf("no more elements to pop")
	}
	v, s.stack = s.stack[s.Len()-1], s.stack[:s.Len()-1]
	// fmt.Println("Pop", s.Len(), s.Peek())
	return v, nil
}
func (s *Stack[T]) Peek() *T {
	if s.Len() <= 0 {
		return nil
	}
	return &s.stack[s.Len()-1]
}

func (s Stack[T]) Len() int {
	return len(s.stack)
}
