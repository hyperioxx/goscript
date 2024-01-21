package core

const CALL_STACK_SIZE = 10000

type Frame struct {
	scope map[string]Object
}

func NewFrame() *Frame {
	return &Frame{map[string]Object{}}
}

func NewCallStack() *[]Frame {
	stack := make([]Frame, 0, CALL_STACK_SIZE)
	return &stack
}
