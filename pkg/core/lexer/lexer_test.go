package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `func add(x, y){
		return x + y
	}
    
	func area(width, height){
		return width * height
	}


	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{FUNC, "func"},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENT, "x"},
		{ADD, "+"},
		{IDENT, "y"},
		{RBRACE, "}"},
		{FUNC, "func"},
		{IDENT, "area"},
		{LPAREN, "("},
		{IDENT, "width"},
		{COMMA, ","},
		{IDENT, "height"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENT, "width"},
		{MUL, "*"},
		{IDENT, "height"},
		{RBRACE, "}"},
		{EOF, "EOF"},
	}

	l := NewV1Lexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s, expectedValue=%q, value=%q", i, TokenTypeStr[tt.expectedType], TokenTypeStr[tok.Type], tt.expectedLiteral, tok.Value)
		}

		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}
}

func TestAssignmentAndArithmetic(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedTokens   []TokenType
		expectedLiterals []string
	}{
		{
			name:  "Integer assignment",
			input: `x = 10`,
			expectedTokens: []TokenType{
				IDENT, ASSIGN, INT, EOF,
			},
			expectedLiterals: []string{
				"x", "=", "10", "EOF",
			},
		},
		{
			name: "Arithmetic operations",
			input: `y = 20 + 30
			        z = y - x
			        a = x * y
			        b = y / x
			        c = 10 % 2`,
			expectedTokens: []TokenType{
				IDENT, ASSIGN, INT, ADD, INT, IDENT, ASSIGN, IDENT, SUB, IDENT,
				IDENT, ASSIGN, IDENT, MUL, IDENT, IDENT, ASSIGN, IDENT, DIV, IDENT,
				IDENT, ASSIGN, INT, REM, INT, EOF,
			},
			expectedLiterals: []string{
				"y", "=", "20", "+", "30", "z", "=", "y", "-", "x",
				"a", "=", "x", "*", "y", "b", "=", "y", "/", "x",
				"c", "=", "10", "%", "2", "EOF",
			},
		},
		{
			name: "Assignment of multiple types",
			input: `name = "John"
			        age = 25
			        pi = 3.14
			        isStudent = true
			        numbers = [1, 2, 3, 4, 5]`,
			expectedTokens: []TokenType{
				IDENT, ASSIGN, STRING, IDENT, ASSIGN, INT, IDENT, ASSIGN, FLOAT, IDENT, ASSIGN, TRUE,
				IDENT, ASSIGN, LBRACKET, INT, COMMA, INT, COMMA, INT, COMMA, INT, COMMA, INT, RBRACKET, EOF,
			},
			expectedLiterals: []string{
				"name", "=", "John", "age", "=", "25", "pi", "=", "3.14", "isStudent", "=", "true",
				"numbers", "=", "[", "1", ",", "2", ",", "3", ",", "4", ",", "5", "]", "EOF",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewV1Lexer(tt.input)

			for i, expectedToken := range tt.expectedTokens {
				token := l.NextToken()

				if token.Type != expectedToken {
					t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s, expectedValue=%q, value=%q", i, TokenTypeStr[expectedToken], TokenTypeStr[token.Type], tt.expectedLiterals[i], token.Value)
				}

				expectedLiteral := tt.expectedLiterals[i]
				if token.Value != expectedLiteral {
					t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, expectedLiteral, token.Value)
				}
			}
		})
	}
}

func TestControlFlowIf(t *testing.T) {
	input := `if x > 10 {
    return true
} else {
    return false
}
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IF, "if"},
		{IDENT, "x"},
		{GT, ">"},
		{INT, "10"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{RBRACE, "}"},
		{EOF, "EOF"},
	}

	l := NewV1Lexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q, expectedValue=%q, value=%q", i, TokenTypeStr[tt.expectedType], TokenTypeStr[tok.Type], tt.expectedLiteral, tok.Value)
		}

		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}
}

func TestControlFlowFor(t *testing.T) {
	input := `for i = 0; i < 10; i = i + 1 {
		x = i * 2
		y = x + 1
	}`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{FOR, "for"},
		{IDENT, "i"},
		{ASSIGN, "="},
		{INT, "0"},
		{SEMICOLON, ";"},
		{IDENT, "i"},
		{LT, "<"},
		{INT, "10"},
		{SEMICOLON, ";"},
		{IDENT, "i"},
		{ASSIGN, "="},
		{IDENT, "i"},
		{ADD, "+"},
		{INT, "1"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{IDENT, "i"},
		{MUL, "*"},
		{INT, "2"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{IDENT, "x"},
		{ADD, "+"},
		{INT, "1"},
		{RBRACE, "}"},
		{EOF, "EOF"},
	}

	l := NewV1Lexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q, expectedValue=%q, value=%q", i, TokenTypeStr[tt.expectedType], TokenTypeStr[tok.Type], tt.expectedLiteral, tok.Value)
		}

		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}
}

func TestFunctionCalls(t *testing.T) {
	input := `add(10, 20)
               area(5, 10)`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "add"},
		{LPAREN, "("},
		{INT, "10"},
		{COMMA, ","},
		{INT, "20"},
		{RPAREN, ")"},
		{IDENT, "area"},
		{LPAREN, "("},
		{INT, "5"},
		{COMMA, ","},
		{INT, "10"},
		{RPAREN, ")"},
		{EOF, "EOF"},
	}

	l := NewV1Lexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s, expectedValue=%q, value=%q", i, TokenTypeStr[tt.expectedType], TokenTypeStr[tok.Type], tt.expectedLiteral, tok.Value)
		}

		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}
}
