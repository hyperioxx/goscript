package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hyperioxx/goscript/pkg/core/lexer"
)

const debug bool = false // global toggle for debug in tests

const ModulePath string = "./"

// Test prefix integer literals
func TestPrefixIntegerLiteral(t *testing.T) {
	input := "5"
	l := lexer.NewV1Lexer(input)
	p := NewV1Parser(l, debug)
	exp := p.ParseNode(LOWEST)

	if exp == nil {
		t.Fatalf("exp is nil.")
	}

	intExp, ok := exp.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *IntegerLiteral. got=%T", exp)
	}

	if intExp.Value().(int) != 5 {
		t.Fatalf("intExp.Value not 5. got=%d", intExp.Value())
	}
}

// Test infix Nodes
func TestInfixNodes(t *testing.T) {
	input := "5 + 5"
	l := lexer.NewV1Lexer(input)
	p := NewV1Parser(l, debug)
	exp := p.ParseNode(LOWEST)

	if exp == nil {
		t.Fatalf("exp is nil.")
	}

	infixExp, ok := exp.(*InfixNode)
	if !ok {
		t.Fatalf("exp not *InfixNode. got=%T", exp)
	}

	left, ok := infixExp.Left.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp.Left not *IntegerLiteral. got=%T", infixExp.Left)
	}
	if left.Value() != 5 {
		t.Fatalf("left.Value not 5. got=%d", left.Value())
	}

	if infixExp.Operator != "+" {
		t.Fatalf("exp.Operator not '+'. got=%q", infixExp.Operator)
	}

	right, ok := infixExp.Right.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp.Right not *IntegerLiteral. got=%T", infixExp.Right)
	}
	if right.Value() != 5 {
		t.Fatalf("right.Value not 5. got=%d", right.Value())
	}
}

// Test function literal
func TestFunctionLiteral(t *testing.T) {
	input := `func add(x, y) { 
		x + y 
		}
		`
	l := lexer.NewV1Lexer(input)
	p := NewV1Parser(l, debug)
	exp := p.ParseNode(LOWEST)

	if exp == nil {
		t.Fatalf("exp is nil.")
	}

	fnExp, ok := exp.(*FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *FunctionLiteral. got=%T", exp)
	}

	if len(fnExp.Parameters) != 2 {
		t.Fatalf("fnExp.Parameters does not contain 2 items. got=%d", len(fnExp.Parameters))
	}

	if fnExp.Parameters[0].Value().(string) != "x" {
		t.Fatalf("fnExp.Parameters[0] not 'x'. got=%q", fnExp.Parameters[0])
	}

	if fnExp.Parameters[1].Value().(string) != "y" {
		t.Fatalf("fnExp.Parameters[1] not 'y'. got=%q", fnExp.Parameters[1])
	}

	body, ok := fnExp.Body[0].(*InfixNode)
	if !ok {
		t.Fatalf("fnExp.Body[0] not *InfixNode. got=%T", fnExp.Body[0])
	}

	if body.Operator != "+" {
		t.Fatalf("body.Operator not '+'. got=%q", body.Operator)
	}

	if body.Left.Value() != "x" {
		t.Fatalf("body.Left not 'x'. got=%q", body.Left)
	}

	if body.Right.Value() != "y" {
		t.Fatalf("body.Right not 'y'. got=%q", body.Right)
	}
}

// Test function literal with no arguments
func TestFunctionLiteralNoArgs(t *testing.T) {
	input := `func display() {  
		}
			
		`
	l := lexer.NewV1Lexer(input)
	p := NewV1Parser(l, debug)
	exp := p.ParseNode(LOWEST)

	if exp == nil {
		t.Fatalf("exp is nil.")
	}

	fnExp, ok := exp.(*FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *FunctionLiteral. got=%T", exp)
	}

	if len(fnExp.Parameters) != 0 {
		t.Fatalf("fnExp.Parameters should contain 0 items. got=%d", len(fnExp.Parameters))
	}

}

// Test function calls
func TestFunctionCall(t *testing.T) {
	input := `add(10, 20)`
	l := lexer.NewV1Lexer(input)
	p := NewV1Parser(l, debug)
	exp := p.ParseNode(LOWEST)

	if exp == nil {
		t.Fatalf("exp is nil.")
	}

	callExp, ok := exp.(*FunctionCall)
	if !ok {
		t.Fatalf("exp not *FunctionCall. got=%T", exp)
	}

	if callExp.Function.Value() != "add" {
		t.Fatalf("callExp.Function not 'add'. got=%q", callExp.Function)
	}

	if len(callExp.Arguments) != 2 {
		t.Fatalf("callExp.Arguments does not contain 2 items. got=%d", len(callExp.Arguments))
	}

	arg1, ok := callExp.Arguments[0].(*IntegerLiteral)
	if !ok {
		t.Fatalf("callExp.Arguments[0] not *IntegerLiteral. got=%T", callExp.Arguments[0])
	}
	if arg1.Value() != 10 {
		t.Fatalf("arg1.Value not 10. got=%d", arg1.Value())
	}

	arg2, ok := callExp.Arguments[1].(*IntegerLiteral)
	if !ok {
		t.Fatalf("callExp.Arguments[1] not *IntegerLiteral. got=%T", callExp.Arguments[1])
	}
	if arg2.Value() != 20 {
		t.Fatalf("arg2.Value not 20. got=%d", arg2.Value())
	}
}

// Test assignment Nodes for different types
func TestAssignNodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{"Integer assignment", "x = 5", &InfixNode{
			Left:     &IdentifierLiteral{value: "x"},
			Operator: "=",
			Right:    &IntegerLiteral{value: 5},
		}},
		{"Float assignment", "y = 3.14", &InfixNode{
			Left:     &IdentifierLiteral{value: "y"},
			Operator: "=",
			Right:    &FloatLiteral{value: 3.14},
		}},
		{"String assignment", "str = \"hello\"", &InfixNode{
			Left:     &IdentifierLiteral{value: "str"},
			Operator: "=",
			Right:    &StringLiteral{value: "hello"},
		}},
		{"Boolean assignment", "b = true", &InfixNode{
			Left:     &IdentifierLiteral{value: "b"},
			Operator: "=",
			Right:    &BooleanLiteral{value: true},
		}},
		{"Array assignment", "arr = [1, 2, 3]", &InfixNode{
			Left:     &IdentifierLiteral{value: "arr"},
			Operator: "=",
			Right: &ArrayLiteral{
				Elements: []Node{
					&IntegerLiteral{value: 1},
					&IntegerLiteral{value: 2},
					&IntegerLiteral{value: 3},
				},
			},
		}},
		{"Empty array assignment", "emptyArr = []", &InfixNode{
			Left:     &IdentifierLiteral{value: "emptyArr"},
			Operator: "=",
			Right: &ArrayLiteral{
				Elements: []Node{},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewV1Lexer(tt.input)
			p := NewV1Parser(l, debug)
			exp := p.ParseNode(LOWEST)

			if exp == nil {
				t.Fatalf("%s: exp is nil for input: %s", tt.name, tt.input)
			}

			assignExp := exp

			if !reflect.DeepEqual(assignExp.Value(), tt.expected.Value()) {
				fmt.Printf("assignExp: %v expected:  %v\n", assignExp, tt.expected)
				t.Fatalf("%s: assignExp not %T %v. got %v %T", tt.name, tt.expected, tt.expected, assignExp, assignExp)
			}

		})
	}
}

func TestFor(t *testing.T) {
	testCases := []struct {
		name                   string
		input                  string
		expectedInitialisation string
		expectedCondition      string
		expectedUpdater        string
		expectedBody           []string
	}{
		{
			name: "full for-loop",
			input: `for i = 0; i != 10; i = i + 1 {
				x = i * 2
				y = x + 1
			}`,
			expectedInitialisation: "i = 0",
			expectedCondition:      "i != 10",
			expectedUpdater:        "i = i + 1",
			expectedBody: []string{
				"x = i * 2",
				"y = x + 1",
			},
		},
	}

	for _, tc := range testCases {
		l := lexer.NewV1Lexer(tc.input)
		p := NewV1Parser(l, debug)
		exp := p.ParseNode(LOWEST)

		if exp == nil {
			t.Fatalf("exp is nil.")
		}

		forExp, ok := exp.(*ForNode)
		if !ok {
			t.Fatalf("exp not *ForNode. got=%T", exp)
		}

		// initialisation
		if initialisation := forExp.Initialisation.String(); initialisation != tc.expectedInitialisation {
			t.Fatalf("forExp.Initialisation should be '%s', got '%s'", tc.expectedInitialisation, initialisation)
		}

		// condition
		if condition := forExp.Condition.String(); condition != tc.expectedCondition {
			t.Fatalf("forExp.Initialisation should be '%s', got '%s'", tc.expectedCondition, condition)
		}

		// updater
		if updater := forExp.Updater.String(); updater != tc.expectedUpdater {
			t.Fatalf("forExp.Initialisation should be '%s', got '%s'", tc.expectedUpdater, updater)
		}

		// body
		if len(forExp.Body) != len(tc.expectedBody) {
			t.Fatalf("forExp.Body unexpectedly long. Expected %d Nodes, got %d", len(tc.expectedBody), len(forExp.Body))
		}
		for lineNo, expectedLine := range tc.expectedBody {
			gotLine := forExp.Body[lineNo]

			if gotLine.String() != expectedLine {
				t.Fatalf("forExp.Body line %d should be '%s', got '%s'", lineNo, expectedLine, gotLine.String())
			}
		}

		// TODO: better checking
		//		- are they Nodes?
		//		- etc.
	}
}
