package commands

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hyperioxx/goscript/pkg/core/lexer"
	"github.com/hyperioxx/goscript/pkg/core/parser"
)

type Linter struct {
	debugFlag *bool
}

func NewLinter(debugFlag *bool) *Linter {
	return &Linter{
		debugFlag: debugFlag,
	}
}

func (l *Linter) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing file name")
	}

	filename := args[2]

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}

	fileContent := string(fileBytes)

	lex := lexer.NewV1Lexer(fileContent)
	p := parser.NewV1Parser(lex, *l.debugFlag)
	expr := p.ParseProgram()
	var formattedContent string
	for _, e := range expr {
		formattedContent += formatExpression(e, 0) + "\n\n"
	}

	if formattedContent == fileContent {
		fmt.Printf("File %s is already correctly formatted\n", filename)
		return nil
	}

	err = ioutil.WriteFile(filename, []byte(formattedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", filename, err)
	}

	fmt.Printf("File %s has been successfully linted\n", filename)
	return nil
}

func formatExpression(expr parser.Node, indentLevel int) string {
	indent := strings.Repeat("\t", indentLevel)
	switch e := expr.(type) {
	case *parser.ReturnStatement:
		return fmt.Sprintf("%sreturn %s", indent, formatExpression(e.ReturnValue, indentLevel))
	case *parser.FunctionCall:
		var args []string
		for _, a := range e.Arguments {
			args = append(args, formatExpression(a, indentLevel))
		}
		return fmt.Sprintf("%s%s(%s)", indent, formatExpression(e.Function, indentLevel), strings.Join(args, ", "))
	case *parser.FunctionLiteral:
		var params []string
		for _, p := range e.Parameters {
			params = append(params, formatExpression(p, indentLevel))
		}
		var body []string
		for _, b := range e.Body {
			body = append(body, formatExpression(b, indentLevel+1))
		}
		return fmt.Sprintf("%sfunc %s(%s) {\n%s\n%s}", indent, e.Name, strings.Join(params, ", "), strings.Join(body, "\n"), indent)
	case *parser.IntegerLiteral:
		return fmt.Sprintf("%s%d", indent, e.Value())
	case *parser.IdentifierLiteral:
		return fmt.Sprintf("%s%s", indent, e.Value())
	case *parser.PrefixNode:
		return fmt.Sprintf("%s%s%s", indent, e.Operator, formatExpression(e.Right, indentLevel))

	case *parser.IfNode:
		var consequence []string
		for _, c := range e.Consequence {
			consequence = append(consequence, formatExpression(c, indentLevel+1))
		}
		var alternative []string
		for _, a := range e.Alternative {
			alternative = append(alternative, formatExpression(a, indentLevel+1))
		}
		return fmt.Sprintf("%sif %s {\n%s\n%s} else {\n%s\n%s}", indent, formatExpression(e.Condition, indentLevel), strings.Join(consequence, "\n"), indent, strings.Join(alternative, "\n"), indent)
	case *parser.BlockStatement:
		var stmts []string
		for _, stmt := range e.Statements {
			stmts = append(stmts, formatExpression(stmt, indentLevel+1))
		}
		return fmt.Sprintf("%s{\n%s\n%s}", indent, strings.Join(stmts, "\n"), indent)
	case *parser.InfixNode:
		left := formatExpression(e.Left, indentLevel)
		right := formatExpression(e.Right, indentLevel)
		if e.Operator == "=" && right[len(right)-2:] == "()" {
			right = right[:len(right)-2]
		}
		return fmt.Sprintf("%s%s %s %s", indent, left, e.Operator, right)

	default:
		return ""
	}
}

func (l *Linter) Name() string {
	return "Linter"
}
