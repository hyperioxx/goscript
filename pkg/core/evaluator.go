package core

import "fmt"

const (
	INT_TYPE int = iota
	FLOAT_TYPE
	STRING_TYPE
	FUNC_TYPE
)

var typeMapping = map[string]int{
	"int":    INT_TYPE,
	"float":  FLOAT_TYPE,
	"string": STRING_TYPE,
	"func":   FUNC_TYPE,
}

type Evaluator struct {
	debug        bool
	callStack    []Frame
	framePointer int
}

func NewEvaluator(debug bool) *Evaluator {
	evaluator := &Evaluator{debug: debug}
	frame := Frame{scope: map[string]Object{}}

	// setup builtin functions in root scope
	frame.scope["print"] = &GoFunction{Name: "print", Func: gsprint}
	frame.scope["length"] = &GoFunction{Name: "length", Func: gslength}
	
	evaluator.callStack = append(evaluator.callStack, frame)
	return evaluator
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
		if variable, ok := e.callStack[e.framePointer].scope[n.String()]; ok {
			return variable, nil
		}
		return &Nil{}, fmt.Errorf("variable '%s' is not defined", n.String())
	case *ArrayLiteral:
		return &Nil{}, nil
	case *ForNode:
		return &Nil{}, nil
	case *FunctionLiteral:
		e.callStack[e.framePointer].scope[n.Name] = &Function{Body: n.Body}
		return &Nil{}, nil
	case *BlockStatement:
		for _, exp := range n.Statements {
			_, err := e.Evaluate(exp)
			if err != nil {
				return nil, err
			}
		}
		return &Nil{}, nil
	case *FunctionCall:

		switch fn := e.callStack[e.framePointer].scope[n.Name].(type) {
		case *GoFunction:
			var args []Object
			for _, arg := range n.Arguments {
				val, err := e.Evaluate(arg)
				if err != nil {
					return &Nil{}, err
				}
				args = append(args, val)
			}

			return fn.Call(args)
		case *Function:
			return e.Evaluate(fn.Body)
		}
		
		return &Nil{}, fmt.Errorf("function '%s' is not defined", n.Name)
	case *IfNode:
		condition, err := e.Evaluate(n.Condition)
		if err != nil {
            return &Nil{}, err
		}
		if _, ok := condition.Value().(bool); ok {
			return e.Evaluate(n.Consequence)
		} else {
			return e.Evaluate(n.Alternative)
		}
	case *InfixNode:
		switch n.Operator {
		case "+":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.Add(right)
		case "-":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.Sub(right)
		case "*":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.Multiply(right)
		case "/":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.Divide(right)
		case "=":
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			e.callStack[e.framePointer].scope[n.Left.String()] = right
			return &Nil{}, nil
		case ">":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.GreaterThan(right)
		default:
			return &Nil{}, fmt.Errorf("unknown operator: %s", n.Operator)
		}
	default:
		return nil, fmt.Errorf("Unknown %T", n)
	}
}
