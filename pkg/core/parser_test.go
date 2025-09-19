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
	}

	for _, test := range cases {
		fmt.Printf("running test %s", test.name)
		lexer := NewV1Lexer(test.input)
		parser := NewV1Parser(lexer, true)
		node, err := parser.ParseNode(0)
		if err != nil {
			fmt.Println(err)
		}

		if !reflect.DeepEqual(node, test.expected[0]) {
			t.Errorf("test %s failed: expected %v, got %v", test.name, test.expected[0], node)
		}
	}
}
