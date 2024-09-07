package core

import "fmt"

const (
	INT_TYPE int = iota
	FLOAT_TYPE
	STRING_TYPE
	FUNC_TYPE
)

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
	case Object:
		return n, nil
	case *IdentifierLiteral:
		variable, err := e.getIdentifier(n.String().value)
		if err != nil {
			return &Nil{}, fmt.Errorf("variable '%s' is not defined", n.String())
		}
		return variable, nil
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
		e.callStack[e.framePointer].scope[n.Name] = &Function{Name: n.Name, Body: n.Body, Arguments: n.Arguments}
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
		fn, err := e.getIdentifier(n.Name)
		if err != nil {
			return &Nil{}, fmt.Errorf("function '%s' is not defined", n.Name)
		}
		e.pushFrame()
		defer e.popFrame()
		switch fn := fn.(type) {
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
			if len(n.Arguments) != len(fn.Arguments) {
				return &Nil{}, fmt.Errorf("function '%s' takes %d arguments only %d was given", fn.Name, len(fn.Arguments), len(n.Arguments))
			}
			for i, argIdent := range fn.Arguments {
				arg, err := e.Evaluate(n.Arguments[i])
				if err != nil {
					return &Nil{}, err
				}
				e.callStack[e.framePointer].scope[argIdent.value] = arg
			}
			return e.Evaluate(fn.Body)
		}

		return &Nil{}, fmt.Errorf("function '%s' is not defined", n.Name)
	case *IfNode:
		condition, err := e.Evaluate(n.Condition)
		if err != nil {
			return &Nil{}, err
		}
		boolean, _ := condition.(*Boolean)
		if boolean.value {
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
			e.callStack[e.framePointer].scope[n.Left.String().value] = right
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

func (e *Evaluator) pushFrame() {
	// create new scope and copy old scope to new scope
	frame := Frame{scope: make(map[string]Object)}

	e.framePointer++
	e.callStack[e.framePointer] = frame
}

func (e *Evaluator) popFrame() {
	e.callStack = e.callStack[:len(e.callStack)-1]
	e.framePointer--
}

func (e *Evaluator) getIdentifier(name string) (Object, error) {
	for i := e.framePointer; i >= 0; i-- {
		if variable, ok := e.callStack[i].scope[name]; ok {
			return variable, nil
		}
	}
	return nil, fmt.Errorf("unable to find reference")
}

func isTruthy(obj Node) bool {
	switch obj := obj.(type) {
	case *Integer:
		return obj.value != 0
	case *Float:
		return obj.value != 0.0
	case *String:
		return obj.value != ""
	case *Boolean:
		return obj.value
	case *Nil:
		return false
	default:
		return true
	}
}
