package router

// Stack is a generic stack data structure.
type stack[T any] []T

// IsEmpty returns true if the stack is empty, false otherwise.
func (s *stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds an element to the top of the stack.
func (s *stack[T]) Push(in T) {
	*s = append(*s, in)
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns the zero value of the element type and false.
func (s *stack[T]) Pop() (T, bool) {
	var zero T
	if s.IsEmpty() {
		return zero, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
