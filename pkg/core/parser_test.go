package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []Node
	}{
		{
			name:  "test increment",
			input: "i++",
			expected: []Node{
				&SufixNode{
					Left:     &IdentifierLiteral{value: "i"},
					Operator: "++",
				},
			},
		},
		{
			name:  "test decrement",
			input: "i--",
			expected: []Node{
				&SufixNode{
					Left:     &IdentifierLiteral{value: "i"},
					Operator: "--",
				},
			},
		},
		{
			name:  "test greater than",
			input: "i > 10",
			expected: []Node{
				&InfixNode{
					Left:     &IdentifierLiteral{value: "i"},
					Operator: ">",
					Right:    &Integer{value: 10},
				},
			},
		},
		{
			name: "test for loop increment",
			input: "for i = 0 ; i < 10; i++ {}",
			expected: []Node{
				&ForNode{
					Initialisation: &InfixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: "=",
						Right:    &Integer{value: 0},
					},
					Condition: &InfixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: "<",
						Right:    &Integer{value: 10},
					},
					Updater: &SufixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: "++",
					},
					Body: &BlockStatement{Statements: []Node{}},
				},
			},	
		},
		{
			name: "test for loop decrement",
			input: "for i = 10 ; i > 10; i-- {}",
			expected: []Node{
				&ForNode{
					Initialisation: &InfixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: "=",
						Right:    &Integer{value: 10},
					},
					Condition: &InfixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: ">",
						Right:    &Integer{value: 10},
					},
					Updater: &SufixNode{
						Left:     &IdentifierLiteral{value: "i"},
						Operator: "--",
					},
					Body: &BlockStatement{Statements: []Node{}},
				},
			},	
		},
	}

	for _, test := range cases {
		fmt.Printf("running test %s", test.name)
		lexer := NewV1Lexer(test.input)
		parser := NewV1Parser(lexer, false)
		node, err := parser.ParseNode(0)
		if err != nil {
			fmt.Println(err)
		}

		if !reflect.DeepEqual(node, test.expected[0]) {
			t.Errorf("test %s failed: expected %v, got %v", test.name, test.expected[0], node)
		}
	}
}
