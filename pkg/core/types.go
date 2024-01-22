package core

import (
	"fmt"
	"strings"
)

type Object interface {
	Type() string
	Value() interface{}
	String() (String, error)
	Add(other Object) (Object, error)
	Sub(other Object) (Object, error)
	Multiply(other Object) (Object, error)
	Divide(other Object) (Object, error)
	Modulo(other Object) (Object, error)
	Equal(other Object) (Object, error)
	NotEqual(other Object) (Object, error)
	GreaterThan(other Object) (Object, error)
	LessThan(other Object) (Object, error)
	GreaterThanOrEqual(other Object) (Object, error)
	LessThanOrEqual(other Object) (Object, error)
}

type Callable interface {
	Object
	GetName() string
	Call(args []Object) (Object, error)
}

type Error interface {
	Object
	Error() string
}

type Integer struct {
	IntValue int
}

func (i Integer) Type() string {
	return "integer"
}

func (i Integer) Value() interface{} {
	return i.IntValue
}

func (i Integer) String() (String, error) {
	return String{StringValue: fmt.Sprintf("%d", i.IntValue)}, nil
}

func (i Integer) Add(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.IntValue + otherInt.IntValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform addition operation with %s and %s", i.Type(), other.Type())
	}
}

func (i Integer) Sub(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.IntValue - otherInt.IntValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform subtraction operation with %s and %s", i.Type(), other.Type())
	}
}

func (i Integer) Multiply(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.IntValue * otherInt.IntValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform multiplication operation with %s and %s", i.Type(), other.Type())
	}
}

func (i Integer) Divide(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		if otherInt.IntValue == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Integer{i.IntValue / otherInt.IntValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform division operation with %s and %s", i.Type(), other.Type())
	}
}

func (i Integer) Modulo(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		if otherInt.IntValue == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Integer{i.IntValue % otherInt.IntValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform modulo operation with %s and %s", i.Type(), other.Type())
	}
}

func (i Integer) Equal(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue == otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using equal operator", i.Type(), other.Type())
}

func (i Integer) NotEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue != otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using not equal operator", i.Type(), other.Type())
}

func (i Integer) GreaterThan(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue > otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using greater than operator", i.Type(), other.Type())
}

func (i Integer) LessThan(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue < otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using less than operator", i.Type(), other.Type())
}

func (i Integer) GreaterThanOrEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue >= otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using greater than or equal operator", i.Type(), other.Type())
}

func (i Integer) LessThanOrEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{BoolValue: i.IntValue <= otherInt.IntValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using less than or equal operator", i.Type(), other.Type())
}

type Float struct {
	FloatValue float64
}

func (f Float) Type() string {
	return "float"
}

func (f Float) Value() interface{} {
	return f.FloatValue
}

func (f Float) String() (String, error) {
	return String{fmt.Sprintf("%f", f.FloatValue)}, nil
}

func (f Float) Add(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{FloatValue: f.FloatValue + otherFloat.FloatValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot add %s with %s", f.Type(), other.Type())
}

func (f Float) Sub(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{FloatValue: f.FloatValue - otherFloat.FloatValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot subtract %s from %s", other.Type(), f.Type())
}

func (f Float) Multiply(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{FloatValue: f.FloatValue * otherFloat.FloatValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot multiply %s with %s", f.Type(), other.Type())
}

func (f Float) Divide(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		if otherFloat.FloatValue == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Float{FloatValue: f.FloatValue / otherFloat.FloatValue}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot divide %s by %s", f.Type(), other.Type())
}

func (f Float) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for %s", f.Type())
}

func (f Float) Equal(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue == otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f Float) NotEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue != otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f Float) GreaterThan(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue > otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f Float) LessThan(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue < otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f Float) GreaterThanOrEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue >= otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f Float) LessThanOrEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{BoolValue: f.FloatValue <= otherFloat.FloatValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

type Boolean struct {
	BoolValue bool
}

func (b Boolean) Type() string {
	return "boolean"
}

func (b Boolean) Value() interface{} {
	return b.BoolValue
}

func (b Boolean) String() (String, error) {
	return String{fmt.Sprintf("%t", b.BoolValue)}, nil
}

func (b Boolean) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for boolean")
}

func (b Boolean) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for boolean")
}

func (b Boolean) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for boolean")
}

func (b Boolean) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for boolean")
}

func (b Boolean) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for boolean")
}

func (b Boolean) Equal(other Object) (Object, error) {
	if otherBool, ok := other.(*Boolean); ok {
		return &Boolean{BoolValue: b.BoolValue == otherBool.BoolValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", b.Type(), other.Type())
	}
}

func (b Boolean) NotEqual(other Object) (Object, error) {
	if otherBool, ok := other.(*Boolean); ok {
		return &Boolean{BoolValue: b.BoolValue != otherBool.BoolValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", b.Type(), other.Type())
	}
}

func (b Boolean) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b Boolean) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b Boolean) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b Boolean) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

type String struct {
	StringValue string
}

func (s String) Type() string {
	return "string"
}

func (s String) Value() interface{} {
	return s.StringValue
}

func (s String) String() (String, error) {
	return String{StringValue: s.StringValue}, nil
}

func (s String) Add(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &String{StringValue: s.StringValue + otherString.StringValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot concatenate %s with %s", s.Type(), other.Type())
	}
}

func (s String) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for string")
}

func (s String) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for string")
}

func (s String) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for string")
}

func (s String) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for string")
}

func (s String) Equal(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &Boolean{BoolValue: s.StringValue == otherString.StringValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", s.Type(), other.Type())
	}
}

func (s String) NotEqual(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &Boolean{BoolValue: s.StringValue != otherString.StringValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", s.Type(), other.Type())
	}
}

func (s String) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s String) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s String) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s String) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

type Array struct {
	Elements []Object
}

func (a Array) Type() string {
	return "array"
}

func (a Array) Value() interface{} {
	values := make([]interface{}, len(a.Elements))
	for i, element := range a.Elements {
		values[i] = element.Value()
	}
	return values
}

func (a Array) String() (String, error) {
	strValues := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		x, _ := element.String()

		strValues[i] = x.StringValue

	}
	return String{StringValue: fmt.Sprintf("[%s]", strings.Join(strValues, ", "))}, nil
}

func (a Array) Add(other Object) (Object, error) {
	if otherArray, ok := other.(Array); ok {
		newElements := make([]Object, len(a.Elements)+len(otherArray.Elements))
		copy(newElements, a.Elements)
		copy(newElements[len(a.Elements):], otherArray.Elements)
		return Array{Elements: newElements}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot concatenate %s with %s", a.Type(), other.Type())
	}
}

func (a Array) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for array")
}

func (a Array) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for array")
}

func (a Array) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for array")
}

func (a Array) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for array")
}

func (a Array) Equal(other Object) (Object, error) {
	if otherArray, ok := other.(Array); ok {
		if len(a.Elements) != len(otherArray.Elements) {
			return &Boolean{BoolValue: false}, nil
		}
		for i := range a.Elements {
			equal, err := a.Elements[i].Equal(otherArray.Elements[i])
			if err != nil {
				return nil, err
			}
			if boolean, ok := equal.(*Boolean); ok && !boolean.BoolValue {
				return &Boolean{BoolValue: false}, nil
			}
		}
		return &Boolean{BoolValue: true}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", a.Type(), other.Type())
	}
}

func (a Array) NotEqual(other Object) (Object, error) {
	equal, err := a.Equal(other)
	if err != nil {
		return nil, err
	}
	if boolean, ok := equal.(*Boolean); ok {
		return &Boolean{BoolValue: !boolean.BoolValue}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", a.Type(), other.Type())
	}
}

func (a Array) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for array")
}

func (a Array) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for array")
}

func (a Array) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for array")
}

func (a Array) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for array")
}

type Map struct {
	mapValue map[Object]Object
}

func (m Map) Type() string {
	return "map"
}

func (m Map) Value() interface{} {
	return m.mapValue
}

func (m Map) String() Object {
	return String{StringValue: fmt.Sprintf("%v", m)}
}

type Function struct {
	Name       string
	Parameters []string
	Body       *BlockStatement
}

func (f *Function) Type() string {
	return "function"
}

func (f *Function) Value() interface{} {
	return struct {
		parameters []string
	}{f.Parameters}
}

func (f *Function) GetName() string {
	return f.Name
}

func (f *Function) String() (String, error) {
	return String{StringValue: fmt.Sprintf("<%s >", f.Name)}, nil
}

func (f *Function) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for function")
}

func (f *Function) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for function")
}

func (f *Function) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for function")
}

func (f *Function) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for function")
}

func (f *Function) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for function")
}

func (f *Function) Equal(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) NotEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *Function) Call(args []Object) (Object, error) {
	return nil, fmt.Errorf("currently not supported")
}

var _ Object = (*Nil)(nil)

type Nil struct {
}

func (n Nil) Type() string {
	return "nil"
}

func (n *Nil) Inspect() string {
	return "nil"
}

func (n Nil) Value() interface{} {
	return nil
}

func (n Nil) String() (String, error) {
	return String{StringValue: "nil"}, nil
}

func (n Nil) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for nil")
}

func (n Nil) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for nil")
}

func (n Nil) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for nil")
}

func (n Nil) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for nil")
}

func (n Nil) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for nil")
}

func (n Nil) Equal(other Object) (Object, error) {
	return &Boolean{BoolValue: n == other}, nil
}

func (n Nil) NotEqual(other Object) (Object, error) {
	return &Boolean{BoolValue: n != other}, nil
}

func (n Nil) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n Nil) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n Nil) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n Nil) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

type Module struct {
	Name  string
	Scope map[string]Object
}

func (m *Module) Type() string {
	return "module"
}

func (m *Module) Value() interface{} {
	return &m.Scope
}

func (m *Module) String() (String, error) {
	return String{StringValue: "nil"}, nil
}

func (m *Module) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for nil")
}

func (m *Module) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for nil")
}

func (m *Module) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for nil")
}

func (m *Module) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for nil")
}

func (m *Module) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for nil")
}

func (m *Module) Equal(other Object) (Object, error) {
	return &Boolean{BoolValue: m == other}, nil
}

func (m *Module) NotEqual(other Object) (Object, error) {
	return &Boolean{BoolValue: m != other}, nil
}

func (m *Module) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (m *Module) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (m *Module) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (m *Module) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

type GoFunction struct {
	Name string
	Func func([]Object) (Object, error) // The actual Go function
}

func (f *GoFunction) Type() string {
	return "gofunction"
}

func (f *GoFunction) Value() interface{} {
	return 1
}

func (f *GoFunction) GetName() string {
	return f.Name
}

func (f *GoFunction) String() (String, error) {
	return String{StringValue: fmt.Sprintf("<%s >", f.Name)}, nil
}

func (f *GoFunction) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for function")
}

func (f *GoFunction) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for function")
}

func (f *GoFunction) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for function")
}

func (f *GoFunction) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for function")
}

func (f *GoFunction) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for function")
}

func (f *GoFunction) Equal(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) NotEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for function")
}

func (f *GoFunction) Call(args []Object) (Object, error) {
	// Call the actual Go function here
	return f.Func(args)
}
