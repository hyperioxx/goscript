package core

import "fmt"

type Evaluator struct {
	debug    bool
	scope    map[string]Object
	symtable map[string]*int
}

func NewEvaluator(debug bool) *Evaluator {
	return &Evaluator{debug: debug, scope: make(map[string]Object)}
}

func (e *Evaluator) Evaluate(exp Node) (Object, error) {
	switch n := exp.(type) {
	case *StringLiteral:
		return &String{StringValue: n.Value().(string)}, nil
	case *IntegerLiteral:
		return &Integer{IntValue: n.Value().(int)}, nil
	case *FloatLiteral:
		return &Float{FloatValue: n.Value().(float64)}, nil
	case *IdentifierLiteral:
		fmt.Println(n)
		return Nil{}, nil
	case *VariableDeclaration:
		// ident := n.Identifier.String()
		// e.scope[n.Identifier.String()]

		fmt.Println(n.Type.Value)
		return Nil{}, nil
	case *ForNode:
		fmt.Println(n)
		return Nil{}, nil
	case *FunctionLiteral:
		fmt.Println(n)
		return Nil{}, nil
	case *FunctionCall:
		fmt.Println(n)
		return Nil{}, nil
	case *InfixNode:
		left, err := e.Evaluate(n.Left)
		if err != nil {
			return Nil{}, err
		}
		right, err := e.Evaluate(n.Right)
		if err != nil {
			return Nil{}, err
		}
		switch n.Operator {
		case "+":
			return left.Add(right)
		case "-":
			return left.Sub(right)
		case "*":
			return left.Multiply(right)
		case "/":
			return left.Divide(right)
		default:
			return &Integer{IntValue: 0}, fmt.Errorf("unknown operator: %s", n.Operator)
		}
	default:
		return nil, fmt.Errorf("Unknown %T", n)
	}
}
