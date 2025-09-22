package core

import (
	"fmt"
)

type Node interface {
	String() *String
}

type Object interface {
	Type() string
	Value() interface{}
	String() *String
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
	value int
}

func (i *Integer) Type() string {
	return "integer"
}

func (i *Integer) Value() interface{} {
	return i.value
}

func (i *Integer) String() *String {
	return &String{value: fmt.Sprintf("%d", i.value)}
}

func (i *Integer) Add(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.value + otherInt.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform addition operation with %s and %s", i.Type(), other.Type())
	}
}

func (i *Integer) Sub(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.value - otherInt.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform subtraction operation with %s and %s", i.Type(), other.Type())
	}
}

func (i *Integer) Multiply(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Integer{i.value * otherInt.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform multiplication operation with %s and %s", i.Type(), other.Type())
	}
}

func (i *Integer) Divide(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		if otherInt.value == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Integer{i.value / otherInt.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform division operation with %s and %s", i.Type(), other.Type())
	}
}

func (i *Integer) Modulo(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		if otherInt.value == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Integer{i.value % otherInt.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot perform modulo operation with %s and %s", i.Type(), other.Type())
	}
}

func (i *Integer) Equal(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value == otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using equal operator", i.Type(), other.Type())
}

func (i *Integer) NotEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value != otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using not equal operator", i.Type(), other.Type())
}

func (i *Integer) GreaterThan(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value > otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using greater than operator", i.Type(), other.Type())
}

func (i *Integer) LessThan(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value < otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using less than operator", i.Type(), other.Type())
}

func (i *Integer) GreaterThanOrEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value >= otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using greater than or equal operator", i.Type(), other.Type())
}

func (i *Integer) LessThanOrEqual(other Object) (Object, error) {
	if otherInt, ok := other.(*Integer); ok {
		return &Boolean{value: i.value <= otherInt.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot compare %s with %s using less than or equal operator", i.Type(), other.Type())
}

func (i *Integer) GetColumn() int {
	return 0
}
func (i *Integer) GetLine() int {
	return 0
}

type Float struct {
	value float64
}

func (f *Float) Type() string {
	return "float"
}

func (f *Float) Value() interface{} {
	return f.value
}

func (f *Float) String() *String {
	return &String{fmt.Sprintf("%f", f.value)}
}

func (f *Float) Add(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{value: f.value + otherFloat.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot add %s with %s", f.Type(), other.Type())
}

func (f *Float) Sub(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{value: f.value - otherFloat.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot subtract %s from %s", other.Type(), f.Type())
}

func (f *Float) Multiply(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Float{value: f.value * otherFloat.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot multiply %s with %s", f.Type(), other.Type())
}

func (f *Float) Divide(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		if otherFloat.value == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		return &Float{value: f.value / otherFloat.value}, nil
	}
	return nil, fmt.Errorf("Invalid type: cannot divide %s by %s", f.Type(), other.Type())
}

func (f *Float) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for %s", f.Type())
}

func (f *Float) Equal(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value == otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) NotEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value != otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) GreaterThan(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value > otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) LessThan(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value < otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) GreaterThanOrEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value >= otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) LessThanOrEqual(other Object) (Object, error) {
	if otherFloat, ok := other.(*Float); ok {
		return &Boolean{value: f.value <= otherFloat.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", f.Type(), other.Type())
	}
}

func (f *Float) GetColumn() int {
	return 0
}
func (f *Float) GetLine() int {
	return 0
}

type Boolean struct {
	value bool
}

func (b *Boolean) Type() string {
	return "boolean"
}

func (b *Boolean) Value() interface{} {
	return b.value
}

func (b *Boolean) String() *String {
	return &String{fmt.Sprintf("%t", b.value)}
}

func (b *Boolean) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for boolean")
}

func (b *Boolean) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for boolean")
}

func (b *Boolean) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for boolean")
}

func (b *Boolean) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for boolean")
}

func (b *Boolean) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for boolean")
}

func (b *Boolean) Equal(other Object) (Object, error) {
	if otherBool, ok := other.(*Boolean); ok {
		return &Boolean{value: b.value == otherBool.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", b.Type(), other.Type())
	}
}

func (b *Boolean) NotEqual(other Object) (Object, error) {
	if otherBool, ok := other.(*Boolean); ok {
		return &Boolean{value: b.value != otherBool.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", b.Type(), other.Type())
	}
}

func (b *Boolean) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b *Boolean) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b *Boolean) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b *Boolean) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for boolean")
}

func (b *Boolean) GetColumn() int {
	return 0
}
func (b *Boolean) GetLine() int {
	return 0
}

type String struct {
	value string
}

func (s *String) Type() string {
	return "string"
}

func (s *String) Value() interface{} {
	return s.value
}

func (s *String) String() *String {
	return s
}

func (s *String) Add(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &String{value: s.value + otherString.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot concatenate %s with %s", s.Type(), other.Type())
	}
}

func (s *String) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for string")
}

func (s *String) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for string")
}

func (s *String) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for string")
}

func (s *String) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for string")
}

func (s *String) Equal(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &Boolean{value: s.value == otherString.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", s.Type(), other.Type())
	}
}

func (s *String) NotEqual(other Object) (Object, error) {
	if otherString, ok := other.(*String); ok {
		return &Boolean{value: s.value != otherString.value}, nil
	} else {
		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", s.Type(), other.Type())
	}
}

func (s *String) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s *String) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s *String) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s *String) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for string")
}

func (s *String) GetColumn() int {
	return 0
}
func (s *String) GetLine() int {
	return 0
}

// type Array struct {
// 	Elements []Node
// }

// func (a *Array) Type() string {
// 	return "array"
// }

// func (a *Array) Value() interface{} {
// 	values := make([]interface{}, len(a.Elements))
// 	for i, element := range a.Elements {
// 		values[i] = element.String().value
// 	}
// 	return values
// }

// func (a *Array) String() *String {
// 	strValues := make([]string, len(a.Elements))
// 	for i, element := range a.Elements {
// 		x := element.String()

// 		strValues[i] = x.value

// 	}
// 	return &String{value: fmt.Sprintf("[%s]", strings.Join(strValues, ", "))}
// }

// func (a *Array) Add(other Object) (Object, error) {
// 	if otherArray, ok := other.(*Array); ok {
// 		newElements := make([]Node, len(a.Elements)+len(otherArray.Elements))
// 		copy(newElements, a.Elements)
// 		copy(newElements[len(a.Elements):], otherArray.Elements)
// 		return &Array{Elements: newElements}, nil
// 	} else {
// 		return nil, fmt.Errorf("Invalid type: cannot concatenate %s with %s", a.Type(), other.Type())
// 	}
// }

// func (a *Array) Sub(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Subtraction operation not supported for array")
// }

// func (a *Array) Multiply(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Multiplication operation not supported for array")
// }

// func (a *Array) Divide(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Division operation not supported for array")
// }

// func (a *Array) Modulo(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Modulo operation not supported for array")
// }

// func (a *Array) Equal(other Object) (Object, error) {
// 	if otherArray, ok := other.(*Array); ok {
// 		if len(a.Elements) != len(otherArray.Elements) {
// 			return &Boolean{value: false}, nil
// 		}
// 		for i := range a.Elements {
// 			equal, err := a.Elements[i].(Object).Equal(otherArray.Elements[i])
// 			if err != nil {
// 				return nil, err
// 			}
// 			if boolean, ok := equal.(*Boolean); ok && !boolean.value {
// 				return &Boolean{value: false}, nil
// 			}
// 		}
// 		return &Boolean{value: true}, nil
// 	} else {
// 		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", a.Type(), other.Type())
// 	}
// }

// func (a *Array) NotEqual(other Object) (Object, error) {
// 	equal, err := a.Equal(other)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if boolean, ok := equal.(*Boolean); ok {
// 		return &Boolean{value: !boolean.value}, nil
// 	} else {
// 		return nil, fmt.Errorf("Invalid type: cannot compare %s with %s", a.Type(), other.Type())
// 	}
// }

// func (a *Array) GreaterThan(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Comparison operation not supported for array")
// }

// func (a *Array) LessThan(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Comparison operation not supported for array")
// }

// func (a *Array) GreaterThanOrEqual(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Comparison operation not supported for array")
// }

// func (a *Array) LessThanOrEqual(other Object) (Object, error) {
// 	return nil, fmt.Errorf("Comparison operation not supported for array")
// }

// func (a *Array) GetColumn() int {
// 	return 0
// }
// func (a *Array) GetLine() int {
// 	return 0
// }

type Function struct {
	Name      string
	Arguments []*IdentifierLiteral
	Body      *BlockStatement
}

func (f *Function) Type() string {
	return "function"
}

func (f *Function) Value() interface{} {
	return ""
}

func (f *Function) GetName() string {
	return f.Name
}

func (f *Function) String() *String {
	return &String{value: fmt.Sprintf("<%s >", f.Name)}
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

func (f *Function) GetColumn() int {
	return 0
}
func (f *Function) GetLine() int {
	return 0
}

var _ Object = (*Nil)(nil)

type Nil struct {
}

func (n *Nil) Type() string {
	return "nil"
}

func (n *Nil) Inspect() string {
	return "nil"
}

func (n *Nil) Value() interface{} {
	return nil
}

func (n *Nil) String() *String {
	return &String{value: "nil"}
}

func (n *Nil) Add(other Object) (Object, error) {
	return nil, fmt.Errorf("Addition operation not supported for nil")
}

func (n *Nil) Sub(other Object) (Object, error) {
	return nil, fmt.Errorf("Subtraction operation not supported for nil")
}

func (n *Nil) Multiply(other Object) (Object, error) {
	return nil, fmt.Errorf("Multiplication operation not supported for nil")
}

func (n *Nil) Divide(other Object) (Object, error) {
	return nil, fmt.Errorf("Division operation not supported for nil")
}

func (n *Nil) Modulo(other Object) (Object, error) {
	return nil, fmt.Errorf("Modulo operation not supported for nil")
}

func (n *Nil) Equal(other Object) (Object, error) {
	return &Boolean{value: n == other}, nil
}

func (n *Nil) NotEqual(other Object) (Object, error) {
	return &Boolean{value: n != other}, nil
}

func (n *Nil) GreaterThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n *Nil) LessThan(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n *Nil) GreaterThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n *Nil) LessThanOrEqual(other Object) (Object, error) {
	return nil, fmt.Errorf("Comparison operation not supported for nil")
}

func (n *Nil) GetColumn() int {
	return 0
}
func (n *Nil) GetLine() int {
	return 0
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

func (f *GoFunction) String() *String {
	return &String{value: fmt.Sprintf("<%s >", f.Name)}
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

func (f *GoFunction) GetColumn() int {
	return 0
}
func (f *GoFunction) GetLine() int {
	return 0
}

// expressions below //

type Expression interface {
	GetLine() int
	GetColumn() int
}

type ReturnStatement struct {
	ReturnValue Node
	Line        int
	Column      int
}

func (rs *ReturnStatement) String() *String {
	return &String{fmt.Sprintf("return %s", rs.ReturnValue.String())}
}

func (rs *ReturnStatement) Value() interface{} {
	return rs
}

func (rs *ReturnStatement) GetLine() int {
	return rs.Line
}

func (rs *ReturnStatement) GetColumn() int {
	return rs.Column
}

type FunctionCall struct {
	Name      string
	Function  Node
	Arguments []Node
	Line      int
	Column    int
}

func (fc *FunctionCall) String() *String {
	return &String{fc.Name}
}

func (fc *FunctionCall) Value() interface{} {
	return fc
}

func (fc *FunctionCall) GetLine() int {
	return fc.Line
}

func (fc *FunctionCall) GetColumn() int {
	return fc.Column
}

type FunctionLiteral struct {
	Name      string
	Arguments []*IdentifierLiteral
	Body      *BlockStatement
	Line      int
	Column    int
}

func (fl *FunctionLiteral) String() *String {
	return &String{fl.Name}
}

func (fl *FunctionLiteral) Value() interface{} {
	return fl
}

func (fl *FunctionLiteral) GetLine() int {
	return fl.Line
}

func (fl *FunctionLiteral) GetColumn() int {
	return fl.Column
}

type IdentifierLiteral struct {
	value  string
	Type   int
	Line   int
	Column int
}

func (il *IdentifierLiteral) String() *String {
	return &String{il.value}
}

func (il *IdentifierLiteral) Value() interface{} {
	return il.value
}

func (il *IdentifierLiteral) GetLine() int {
	return il.Line
}

func (il *IdentifierLiteral) GetColumn() int {
	return il.Column
}

type InfixNode struct {
	Left     Node
	Operator string
	Right    Node
	Line     int
	Column   int
}

func (ie *InfixNode) String() *String {
	return &String{fmt.Sprintf("%s %s %s", ie.Left.String(), ie.Operator, ie.Right.String())}
}

func (ie *InfixNode) Value() interface{} {
	return ie.Operator
}

func (ie *InfixNode) GetLine() int {
	return ie.Line
}

func (ie *InfixNode) GetColumn() int {
	return ie.Column
}

type PrefixNode struct {
	Operator string
	Right    Node
	Line     int
	Column   int
}

func (pe *PrefixNode) String() *String {
	return &String{string(pe.Operator)}
}

func (pe *PrefixNode) Value() interface{} {
	return pe.Operator
}

func (pe *PrefixNode) GetLine() int {
	return pe.Line
}

func (pe *PrefixNode) GetColumn() int {
	return pe.Column
}

type SufixNode struct {
	Operator string
	Left     Node
	Line     int
	Column   int
}

func (se *SufixNode) String() *String {
	return &String{fmt.Sprintf("%s%s", se.Left.String(), se.Operator)}
}

func (se *SufixNode) Value() interface{} {
	return se.Operator
}

func (se *SufixNode) GetLine() int {
	return se.Line
}

func (se *SufixNode) GetColumn() int {
	return se.Column
}

type IfNode struct {
	Condition   Node
	Consequence Node
	Alternative Node
	Line        int
	Column      int
}

func (ie *IfNode) String() *String {
	return &String{"if"}
}

func (ie *IfNode) Value() interface{} {
	return ie
}

func (ie *IfNode) GetLine() int {
	return ie.Line
}

func (ie *IfNode) GetColumn() int {
	return ie.Column
}

type ForNode struct {
	Initialisation Node
	Condition      Node
	Updater        Node
	Body           Node
	Line           int
	Column         int
}

func (fe *ForNode) String() *String {
	return &String{"for"}
}

func (fe *ForNode) Value() interface{} {
	return fe
}

func (fe *ForNode) GetLine() int {
	return fe.Line
}

func (fe *ForNode) GetColumn() int {
	return fe.Column
}

type BlockStatement struct {
	Statements []Node
	Line       int
	Column     int
}

func (bs *BlockStatement) String() *String {
	return &String{"if"}
}

func (bs *BlockStatement) Value() interface{} {
	return bs
}

func (bs *BlockStatement) GetLine() int {
	return bs.Line
}

func (bs *BlockStatement) GetColumn() int {
	return bs.Column
}

func NewIfNode(condition Node, consequence Node, alternative Node, line, column int) *IfNode {
	return &IfNode{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
		Line:        line,
		Column:      column,
	}
}

func NewIdentifierLiteral(value string, line, column int) *IdentifierLiteral {
	return &IdentifierLiteral{
		value:  value,
		Line:   line,
		Column: column,
	}
}

func NewInfixNode(left Node, operator string, right Node, line, column int) *InfixNode {
	return &InfixNode{
		Left:     left,
		Operator: operator,
		Right:    right,
		Line:     line,
		Column:   column,
	}
}

func NewPrefixNode(operator string, right Node, line, column int) *PrefixNode {
	return &PrefixNode{
		Operator: operator,
		Right:    right,
		Line:     line,
		Column:   column,
	}
}

func NewFunctionLiteral(name string, parameters []*IdentifierLiteral, body *BlockStatement, line, column int) *FunctionLiteral {
	return &FunctionLiteral{
		Name:      name,
		Arguments: parameters,
		Body:      body,
		Line:      line,
		Column:    column,
	}
}

func NewFunctionCall(function Node, arguments []Node, line, column int) *FunctionCall {
	return &FunctionCall{
		Function:  function,
		Arguments: arguments,
		Line:      line,
		Column:    column,
	}
}

func NewReturnStatement(returnValue Node, line, column int) *ReturnStatement {
	return &ReturnStatement{
		ReturnValue: returnValue,
		Line:        line,
		Column:      column,
	}
}

type VariableDeclaration struct {
	Identifier *IdentifierLiteral
	Type       Token
	Line       int
	Column     int
}

func (vd *VariableDeclaration) String() *String {
	return &String{fmt.Sprintf("%s: %s", vd.Identifier.String(), vd.Type.Value)}
}

func (vd *VariableDeclaration) Value() interface{} { return vd }
func (vd *VariableDeclaration) GetLine() int       { return vd.Line }
func (vd *VariableDeclaration) GetColumn() int     { return vd.Column }

type StructLiteral struct {
	Name   *IdentifierLiteral
	Fields []*IdentifierLiteral
	Line   int
	Column int
}

func (s *StructLiteral) String() *String {
	return &String{fmt.Sprintf("struct %s", s.Name.String())}
}

func (s *StructLiteral) Value() interface{} { return s }
func (s *StructLiteral) GetLine() int       { return s.Line }
func (s *StructLiteral) GetColumn() int     { return s.Column }
