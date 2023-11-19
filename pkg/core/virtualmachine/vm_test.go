package virtualmachine

import (
	"fmt"
	"reflect"
	"testing"
)

const debug bool = false // global toggle for debug in tests

func TestVM_Run(t *testing.T) {
	testCases := []struct {
		name         string
		instructions []Instruction
		expected     Object
	}{
		{
			name: "Add Operation",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &Integer{IntValue: 5}},
				{Opcode: OpPush, Value: &Integer{IntValue: 3}},
				{Opcode: OpAdd},
			},
			expected: &Integer{IntValue: 8},
		},
		{
			name: "Subtract Operation",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &Integer{IntValue: 5}},
				{Opcode: OpPush, Value: &Integer{IntValue: 3}},
				{Opcode: OpSub},
			},
			expected: &Integer{IntValue: 2},
		},
		{
			name: "Equal Operation",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &Integer{IntValue: 5}},
				{Opcode: OpPush, Value: &Integer{IntValue: 5}},
				{Opcode: OpEqual},
			},
			expected: &Boolean{BoolValue: true},
		},
		{
			name: "NotEqual Operation",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &Integer{IntValue: 5}},
				{Opcode: OpPush, Value: &Integer{IntValue: 3}},
				{Opcode: OpNotEqual},
			},
			expected: &Boolean{BoolValue: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vm := NewVirtualMachine(debug)
			result := vm.Run(tc.instructions)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("unexpected result: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestFunctionCall(t *testing.T) {
	vm := NewVirtualMachine(debug)

	instructions := []Instruction{
		{Opcode: OpStoreFunc, Value: &Function{
			Name:       "add",
			Parameters: []string{"a", "b"},
			Body: []Instruction{
				{Opcode: OpGet, Value: &String{StringValue: "a"}},
				{Opcode: OpGet, Value: &String{StringValue: "b"}},
				{Opcode: OpAdd},
				{Opcode: OpReturn},
			},
		}},
		{Opcode: OpPush, Value: &Integer{IntValue: 5}},
		{Opcode: OpPush, Value: &Integer{IntValue: 10}},
		{Opcode: OpPush, Value: &Integer{IntValue: 2}},
		{Opcode: OpCall, Value: &String{StringValue: "add"}},
	}

	result := vm.Run(instructions)

	if result.Type() != "integer" {
		t.Fatalf("Expected Integer, got %s", result.Type())
	}

	if result.(*Integer).IntValue != 15 {
		t.Fatalf("Expected 15, got %d", result.(*Integer).IntValue)
	}
}

func TestGoFunctionCall(t *testing.T) {
	vm := NewVirtualMachine(debug)

	instructions := []Instruction{
		{Opcode: OpStoreFunc, Value: &GoFunction{
			Name: "add",
			Func: func(args []Object) (Object, error) {
				a, ok1 := args[0].(*Integer)
				b, ok2 := args[1].(*Integer)
				if !ok1 || !ok2 {
					return nil, fmt.Errorf("Invalid arguments, expected integers")
				}
				return &Integer{IntValue: a.IntValue + b.IntValue}, nil
			},
		}},
		{Opcode: OpPush, Value: &Integer{IntValue: 5}},
		{Opcode: OpPush, Value: &Integer{IntValue: 10}},
		{Opcode: OpPush, Value: &Integer{IntValue: 2}},
		{Opcode: OpCall, Value: &String{StringValue: "add"}},
	}

	result := vm.Run(instructions)

	if result.Type() != "integer" {
		t.Fatalf("Expected Integer, got %s", result.Type())
	}

	if result.(*Integer).IntValue != 15 {
		t.Fatalf("Expected 15, got %d", result.(*Integer).IntValue)
	}
}

func TestJump(t *testing.T) {
	vm := NewVirtualMachine(debug)

	instructions := []Instruction{
		{Opcode: OpPush, Value: &Integer{IntValue: 5}},
		{Opcode: OpPush, Value: &Integer{IntValue: 3}},
		{Opcode: OpGreaterThan},
		{Opcode: OpJump, Value: &Integer{IntValue: 2}},
		{Opcode: OpPush, Value: &Integer{IntValue: 10}},
		{Opcode: OpPush, Value: &Integer{IntValue: 7}},
	}

	result := vm.Run(instructions)

	if result.Type() != "integer" {
		t.Fatalf("Expected Integer, got %s", result.Type())
	}

	if result.(*Integer).IntValue != 7 {
		t.Fatalf("Expected 7, got %d", result.(*Integer).IntValue)
	}
}

func TestVariableAssignment(t *testing.T) {
	vm := NewVirtualMachine(debug)

	testCases := []struct {
		name         string
		instructions []Instruction
		expected     Object
	}{
		{
			name: "Assign String",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpPush, Value: &String{StringValue: "Hello, world!"}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &String{StringValue: "Hello, world!"},
		},
		{
			name: "Assign Integer",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 42}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &Integer{IntValue: 42},
		},
		{
			name: "Assign Float",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpPush, Value: &Float{FloatValue: 3.14}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &Float{FloatValue: 3.14},
		},
		{
			name: "Assign Boolean",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpPush, Value: &Boolean{BoolValue: true}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &Boolean{BoolValue: true},
		},
		{
			name: "Assign Array",
			instructions: []Instruction{
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpPush, Value: &Array{Elements: []Object{
					&Integer{IntValue: 1},
					&Integer{IntValue: 2},
					&Integer{IntValue: 3},
				}}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &Array{Elements: []Object{
				&Integer{IntValue: 1},
				&Integer{IntValue: 2},
				&Integer{IntValue: 3},
			}},
		},
		{
			name: "Assign Function",
			instructions: []Instruction{
				{Opcode: OpStoreFunc, Value: &Function{
					Name:       "add",
					Parameters: []string{"a", "b"},
					Body: []Instruction{
						{Opcode: OpGet, Value: &String{StringValue: "a"}},
						{Opcode: OpGet, Value: &String{StringValue: "b"}},
						{Opcode: OpAdd},
						{Opcode: OpReturn},
					},
				}},
				{Opcode: OpPush, Value: &String{StringValue: "x"}},
				{Opcode: OpGet, Value: &String{StringValue: "add"}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "x"}},
			},
			expected: &Function{
				Name:       "add",
				Parameters: []string{"a", "b"},
				Body: []Instruction{
					{Opcode: OpGet, Value: &String{StringValue: "a"}},
					{Opcode: OpGet, Value: &String{StringValue: "b"}},
					{Opcode: OpAdd},
					{Opcode: OpReturn},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := vm.Run(tc.instructions)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("unexpected result: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestForLoop(t *testing.T) {
	testCases := []struct {
		name         string
		instructions []Instruction
		expected     Object
	}{
		{
			name: "For loop - full",
			// adapted from output of compiler
			instructions: []Instruction{
				// y = 0
				{Opcode: OpPush, Value: &String{StringValue: "y"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 0}},
				{Opcode: OpAssign},
				// for i = 0; i < 10; i = i + 1 { y = y + 42 }
				{Opcode: OpPush, Value: &String{StringValue: "i"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 0}},
				{Opcode: OpAssign},
				{Opcode: OpGet, Value: &String{StringValue: "i"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 10}},
				{Opcode: OpLessThan},
				{Opcode: OpJumpIfFalse, Value: &Integer{IntValue: 12}},
				{Opcode: OpPush, Value: &String{StringValue: "y"}},
				{Opcode: OpGet, Value: &String{StringValue: "y"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 42}},
				{Opcode: OpAdd},
				{Opcode: OpAssign},
				{Opcode: OpPush, Value: &String{StringValue: "i"}},
				{Opcode: OpGet, Value: &String{StringValue: "i"}},
				{Opcode: OpPush, Value: &Integer{IntValue: 1}},
				{Opcode: OpAdd},
				{Opcode: OpAssign},
				{Opcode: OpJump, Value: &Integer{IntValue: -14}},
				// put result of `y` at top of stack
				{Opcode: OpGet, Value: &String{StringValue: "y"}},
			},
			expected: &Integer{IntValue: 420},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vm := NewVirtualMachine(debug)
			result := vm.Run(tc.instructions)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("unexpected result: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func BenchmarkAddition(b *testing.B) {
	vm := NewVirtualMachine(debug)
	instructions := []Instruction{
		{Opcode: OpPush, Value: &Integer{IntValue: 2}},
		{Opcode: OpPush, Value: &Integer{IntValue: 3}},
		{Opcode: OpAdd},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vm.Run(instructions)
	}
}
