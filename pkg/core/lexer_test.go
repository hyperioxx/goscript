package core

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "test for loop",
			input: "for i = 0 ; i < 10; i = i + 1 {}",
			expected: []Token{
				{Value: "for", Type: FOR, Line: 1, Column: 1},
				{Value: "i", Type: IDENT, Line: 1, Column: 5},
				{Value: "=", Type: ASSIGN, Line: 1, Column: 7},
				{Value: "0", Type: INT, Line: 1, Column: 9},
				{Value: ";", Type: SEMICOLON, Line: 1, Column: 11},
				{Value: "i", Type: IDENT, Line: 1, Column: 13},
				{Value: "<", Type: LT, Line: 1, Column: 15},
				{Value: "10", Type: INT, Line: 1, Column: 17},
				{Value: ";", Type: SEMICOLON, Line: 1, Column: 19},
				{Value: "i", Type: IDENT, Line: 1, Column: 21},
				{Value: "=", Type: ASSIGN, Line: 1, Column: 23},
				{Value: "i", Type: IDENT, Line: 1, Column: 25},
				{Value: "+", Type: ADD, Line: 1, Column: 27},
				{Value: "1", Type: INT, Line: 1, Column: 29},
				{Value: "{", Type: LBRACE, Line: 1, Column: 31},
				{Value: "}", Type: RBRACE, Line: 1, Column: 32},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewV1Lexer(test.input)
			for i, expectedToken := range test.expected {
				token := lexer.NextToken()
				if !reflect.DeepEqual(token, expectedToken) {
					t.Errorf("test %s failed: token %d - expected %v, got %v", test.name, i, expectedToken, token)
				}
			}
		})
	}
}
