package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []Object
	}{
		{
			name:     "test increment",
			input:    []string{"i = 0", "i++"},
			expected: []Object{&Nil{}, &Integer{value: 1}},
		},
	}

	for _, test := range cases {
		fmt.Printf("running test %s", test.name)
		evaluator := NewEvaluator(false)
		for i, line := range test.input {
			fmt.Printf("parsing line: %s\n %d", line, i)
			lexer := NewV1Lexer(line)
			parser := NewV1Parser(lexer, false)
			node, err := parser.ParseNode(0)
			if err != nil {
				fmt.Println(err)
			}

			result, err := evaluator.Evaluate(node)
			if err != nil {
				fmt.Println(err)
			}

			if !reflect.DeepEqual(result, test.expected[i]) {
				t.Fatalf("expected %v, got %v", test.expected[i], result)
			}

		}

	}

}
