package core

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name: "test for loop",
			input: "for i = 0 ; i < 10; i++ {}",
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
		
		fmt.Printf("%+v", node.String())
	}
}
