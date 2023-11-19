package compiler

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/hyperioxx/goscript/pkg/core/lexer"
	"github.com/hyperioxx/goscript/pkg/core/parser"
	"github.com/hyperioxx/goscript/pkg/core/virtualmachine"
)

const debug bool = false // global toggle for debug in tests

func TestCompiler_Compile(t *testing.T) {
	c := NewCompiler()

	testCases := []struct {
		name        string
		Node        parser.Node
		expected    []virtualmachine.Instruction
		expectError bool
	}{
		{
			name:     "IntegerLiteral",
			Node:     parser.NewIntegerLiteral(123, 1, 0),
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 123}}},
		},
		{
			name:        "Invalid InfixNode - Left side not an Identifier",
			Node:        &parser.InfixNode{Left: parser.NewIntegerLiteral(123, 1, 0), Operator: lexer.TokenTypeStr[lexer.ASSIGN], Right: parser.NewIntegerLiteral(456, 1, 0)},
			expectError: true,
		},
		{
			name:     "IdentifierLiteral",
			Node:     parser.NewIdentifierLiteral("abc", 1, 0),
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpGet, Value: &virtualmachine.String{StringValue: "abc"}}},
		},
		{
			name:     "InfixNode - Add",
			Node:     &parser.InfixNode{Left: parser.NewIntegerLiteral(1, 1, 0), Operator: lexer.TokenTypeStr[lexer.ADD], Right: parser.NewIntegerLiteral(2, 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 2}}, {Opcode: virtualmachine.OpAdd}},
		},
		{
			name:     "InfixNode - Subtract",
			Node:     &parser.InfixNode{Left: parser.NewIntegerLiteral(3, 1, 0), Operator: lexer.TokenTypeStr[lexer.SUB], Right: parser.NewIntegerLiteral(2, 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 3}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 2}}, {Opcode: virtualmachine.OpSub}},
		},
		{
			name:     "Assignment - Integer",
			Node:     &parser.InfixNode{Left: parser.NewIdentifierLiteral("x", 1, 0), Operator: lexer.TokenTypeStr[lexer.ASSIGN], Right: parser.NewIntegerLiteral(5, 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "x"}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 5}}, {Opcode: virtualmachine.OpAssign}},
		},
		{
			name:     "Assignment - Float",
			Node:     &parser.InfixNode{Left: parser.NewIdentifierLiteral("y", 1, 0), Operator: lexer.TokenTypeStr[lexer.ASSIGN], Right: parser.NewFloatLiteral(3.14, 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "y"}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.Float{FloatValue: 3.14}}, {Opcode: virtualmachine.OpAssign}},
		},
		{
			name:     "Assignment - String",
			Node:     &parser.InfixNode{Left: parser.NewIdentifierLiteral("str", 1, 0), Operator: lexer.TokenTypeStr[lexer.ASSIGN], Right: parser.NewStringLiteral("hello", 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "str"}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "hello"}}, {Opcode: virtualmachine.OpAssign}},
		},
		{
			name:     "Assignment - Boolean",
			Node:     &parser.InfixNode{Left: parser.NewIdentifierLiteral("b", 1, 0), Operator: lexer.TokenTypeStr[lexer.ASSIGN], Right: parser.NewBooleanLiteral(true, 1, 0)},
			expected: []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "b"}}, {Opcode: virtualmachine.OpPush, Value: &virtualmachine.Boolean{BoolValue: true}}, {Opcode: virtualmachine.OpAssign}},
		},
		{
			name: "IfStatement - Basic",
			Node: &parser.IfNode{
				Condition: &parser.InfixNode{
					Left:     parser.NewIntegerLiteral(1, 1, 0),
					Operator: lexer.TokenTypeStr[lexer.EQ],
					Right:    parser.NewIntegerLiteral(1, 1, 0),
				},
				Consequence: []parser.Node{
					&parser.InfixNode{
						Left:     parser.NewIdentifierLiteral("x", 1, 0),
						Operator: lexer.TokenTypeStr[lexer.ASSIGN],
						Right:    parser.NewIntegerLiteral(1, 1, 0),
					},
				},
			},
			expected: []virtualmachine.Instruction{
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpEqual},
				{Opcode: virtualmachine.OpJumpIfFalse, Value: &virtualmachine.Integer{IntValue: 4}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "x"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpAssign},
			},
		},
		{
			name: "ForStatement - while true",
			Node: &parser.ForNode{
				Body: []parser.Node{
					&parser.InfixNode{
						Left:     parser.NewIdentifierLiteral("x", 1, 0),
						Operator: lexer.TokenTypeStr[lexer.ASSIGN],
						Right:    parser.NewIntegerLiteral(1, 1, 0),
					},
				},
			},
			expected: []virtualmachine.Instruction{
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Boolean{BoolValue: true}},
				{Opcode: virtualmachine.OpJumpIfFalse, Value: &virtualmachine.Integer{IntValue: 5}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "x"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpAssign},
				{Opcode: virtualmachine.OpJump, Value: &virtualmachine.Integer{IntValue: -5}},
			},
		},
		{
			name: "ForStatement - while condition",
			Node: &parser.ForNode{
				Condition: &parser.InfixNode{
					Left:     parser.NewIdentifierLiteral("x", 1, 0),
					Operator: lexer.TokenTypeStr[lexer.LT],
					Right:    parser.NewIntegerLiteral(10, 1, 0),
				},
				Body: []parser.Node{
					&parser.InfixNode{
						Left:     parser.NewIdentifierLiteral("y", 1, 0),
						Operator: lexer.TokenTypeStr[lexer.ASSIGN],
						Right:    parser.NewIntegerLiteral(1, 1, 0),
					},
				},
			},
			expected: []virtualmachine.Instruction{
				{Opcode: virtualmachine.OpGet, Value: &virtualmachine.String{StringValue: "x"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 10}},
				{Opcode: virtualmachine.OpLessThan},
				{Opcode: virtualmachine.OpJumpIfFalse, Value: &virtualmachine.Integer{IntValue: 5}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "y"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpAssign},
				{Opcode: virtualmachine.OpJump, Value: &virtualmachine.Integer{IntValue: -7}},
			},
		},
		{
			name: "ForStatement - full whack",
			Node: &parser.ForNode{
				Initialisation: &parser.InfixNode{
					Left:     parser.NewIdentifierLiteral("i", 1, 0),
					Operator: lexer.TokenTypeStr[lexer.ASSIGN],
					Right:    parser.NewIntegerLiteral(0, 1, 0),
				},
				Condition: &parser.InfixNode{
					Left:     parser.NewIdentifierLiteral("i", 1, 0),
					Operator: lexer.TokenTypeStr[lexer.LT],
					Right:    parser.NewIntegerLiteral(10, 1, 0),
				},
				Updater: &parser.InfixNode{
					Left:     parser.NewIdentifierLiteral("i", 1, 0),
					Operator: lexer.TokenTypeStr[lexer.ASSIGN],
					Right: &parser.InfixNode{
						Left:     parser.NewIdentifierLiteral("i", 1, 0),
						Operator: lexer.TokenTypeStr[lexer.ADD],
						Right:    parser.NewIntegerLiteral(1, 1, 0),
					},
				},
				Body: []parser.Node{
					&parser.InfixNode{
						Left:     parser.NewIdentifierLiteral("y", 1, 0),
						Operator: lexer.TokenTypeStr[lexer.ASSIGN],
						Right:    parser.NewIntegerLiteral(1, 1, 0),
					},
				},
			},
			expected: []virtualmachine.Instruction{
				// initialisation
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "i"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 0}},
				{Opcode: virtualmachine.OpAssign},
				// condition
				{Opcode: virtualmachine.OpGet, Value: &virtualmachine.String{StringValue: "i"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 10}},
				{Opcode: virtualmachine.OpLessThan},
				// jump if false
				{Opcode: virtualmachine.OpJumpIfFalse, Value: &virtualmachine.Integer{IntValue: 10}},
				// body
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "y"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpAssign},
				// updater
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: "i"}},
				{Opcode: virtualmachine.OpGet, Value: &virtualmachine.String{StringValue: "i"}},
				{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
				{Opcode: virtualmachine.OpAdd},
				{Opcode: virtualmachine.OpAssign},
				// jump to loop
				{Opcode: virtualmachine.OpJump, Value: &virtualmachine.Integer{IntValue: -12}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := make([]parser.Node, 0)
			x = append(x, tc.Node)
			result, err := c.Compile(x, debug)
			if (err != nil) != tc.expectError {
				t.Fatalf("unexpected error status: got error %v, want error status %v", err, tc.expectError)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("unexpected result:\ngot:\n%v\nwant:\n%v", result, tc.expected)
			}
		})
	}
}

func printInstructions(instructions []virtualmachine.Instruction) string {
	var buf bytes.Buffer
	for _, inst := range instructions {
		buf.WriteString(fmt.Sprintf("%#v\n", virtualmachine.OpCodeStrings[inst.Opcode]))
	}
	return buf.String()
}

func compareInstructions(a, b []virtualmachine.Instruction) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Opcode != b[i].Opcode {
			return false
		}
		switch va := a[i].Value.(type) {
		case *virtualmachine.Integer:
			if vb, ok := b[i].Value.(*virtualmachine.Integer); !ok || va.IntValue != vb.IntValue {
				return false
			}
		case *virtualmachine.String:
			if vb, ok := b[i].Value.(*virtualmachine.String); !ok || va.StringValue != vb.StringValue {
				return false
			}
		// add other cases for other types...
		default:
			if a[i].Value != b[i].Value {
				return false
			}
		}
	}
	return true
}
