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
			name:  "test increment",
			input: "i++",
			expected: []Token{
				{Value: "i", Type: IDENT, Line: 1, Column: 1},
				{Value: "++", Type: INC, Line: 1, Column: 2},
			},
		},
		{
			name:  "test decrement",
			input: "i--",
			expected: []Token{
				{Value: "i", Type: IDENT, Line: 1, Column: 1},
				{Value: "--", Type: DEC, Line: 1, Column: 2},
			},
		},
		{
			name:  "test greater than",
			input: "i > 10",
			expected: []Token{
				{Value: "i", Type: IDENT, Line: 1, Column: 1},
				{Value: ">", Type: GT, Line: 1, Column: 3},
				{Value: "10", Type: INT, Line: 1, Column: 5},
			},
		},
		{
			name:  "test less than",
			input: "i < 10",
			expected: []Token{
				{Value: "i", Type: IDENT, Line: 1, Column: 1},
				{Value: "<", Type: LT, Line: 1, Column: 3},
				{Value: "10", Type: INT, Line: 1, Column: 5},
			},
		},
		{
			name:  "test for loop",
			input: "for i = 0 ; i < 10; i = i++ {}",
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
				{Value: "++", Type: INC, Line: 1, Column: 26},
				{Value: "{", Type: LBRACE, Line: 1, Column: 29},
				{Value: "}", Type: RBRACE, Line: 1, Column: 30},
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
