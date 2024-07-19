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
	frame := Frame{scope: map[string]Object{}} // global scope

	// setup builtin functions in root scope
	frame.scope["print"] = &GoFunction{Name: "print", Func: gsprint}
	frame.scope["length"] = &GoFunction{Name: "length", Func: gslength}
    evaluator.callStack = make([]Frame, 10000) // call stack of 10000
	evaluator.callStack[evaluator.framePointer] = frame

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
		e.pushFrame()
		defer e.popFrame()
		if _, err := e.Evaluate(n.Initialisation); err != nil {
			return &Nil{}, err
		}

		for {
			cond, err := e.Evaluate(n.Condition)
			if err != nil {
				return &Nil{}, err
			}
			if !isTruthy(cond) {
				break
			}

			if _, err := e.Evaluate(n.Body); err != nil {
				return &Nil{}, err
			}

			if _, err := e.Evaluate(n.Updater); err != nil {
				return &Nil{}, err
			}
		}

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
        e.pushFrame()
		defer e.popFrame()
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
		boolean, _ := condition.(*Boolean)
		if boolean.BoolValue {
			return e.Evaluate(n.Consequence)
		}

		if n.Alternative != nil {
			return e.Evaluate(n.Alternative)
		}
		return &Nil{}, nil
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
		case "<":
			left, err := e.Evaluate(n.Left)
			if err != nil {
				return &Nil{}, err
			}
			right, err := e.Evaluate(n.Right)
			if err != nil {
				return &Nil{}, err
			}
			return left.LessThan(right)
		default:
			return &Nil{}, fmt.Errorf("unknown operator: %s", n.Operator)
		}
	default:
		return nil, fmt.Errorf("Unknown %T", n)
	}
}


func (e *Evaluator) pushFrame(){
	// create new scope and copy old scope to new scope 
	frame := Frame{scope: make(map[string]Object, len(e.callStack[e.framePointer].scope))}
	for k, v := range e.callStack[e.framePointer].scope {
		frame.scope[k] = v
	}

	e.framePointer++
	e.callStack[e.framePointer] = frame
}

func (e *Evaluator) popFrame(){
	e.callStack = e.callStack[:len(e.callStack)-1]
	e.framePointer--
}

func isTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *Integer:
		return obj.IntValue != 0
	case *Float:
		return obj.FloatValue != 0.0
	case *String:
		return obj.StringValue != ""
	case *Boolean:
		return obj.BoolValue
	case *Nil:
		return false
	default:
		return true
	}
}