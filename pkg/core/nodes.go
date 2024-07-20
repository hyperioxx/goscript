package core

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	String() string
	Value() interface{}
	GetLine() int
	GetColumn() int
}

type StringLiteral struct {
	Line   int
	Column int
	value  string
}

func (sl *StringLiteral) String() string     { return sl.value }
func (sl *StringLiteral) Value() interface{} { return sl.value }
func (sl *StringLiteral) GetLine() int       { return sl.Line }
func (sl *StringLiteral) GetColumn() int     { return sl.Column }

type FloatLiteral struct {
	Line   int
	Column int
	value  float64
}

func (fl *FloatLiteral) String() string     { return fmt.Sprintf("%f", fl.Value()) }
func (fl *FloatLiteral) Value() interface{} { return fl.value }
func (fl *FloatLiteral) GetLine() int       { return fl.Line }
func (fl *FloatLiteral) GetColumn() int     { return fl.Column }

type ArrayLiteral struct {
	Line     int
	Column   int
	Elements []Node
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (al *ArrayLiteral) Value() interface{} { return al.Elements }
func (al *ArrayLiteral) GetLine() int       { return al.Line }
func (al *ArrayLiteral) GetColumn() int     { return al.Column }

type HashLiteral struct {
	Line   int
	Column int
	Pairs  map[Node]Node
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (hl *HashLiteral) Value() interface{} { return hl.Pairs }
func (hl *HashLiteral) GetLine() int       { return hl.Line }
func (hl *HashLiteral) GetColumn() int     { return hl.Column }

type ReturnStatement struct {
	ReturnValue Node
	Line        int
	Column      int
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", rs.ReturnValue.String())
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

func (fc *FunctionCall) String() string {
	return fc.Name
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
	Name       string
	Parameters []Node
	Body       *BlockStatement
	Line       int
	Column     int
}

func (fl *FunctionLiteral) String() string {
	return fl.Name
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

type IntegerLiteral struct {
	value  int
	Line   int
	Column int
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprint(il.value)
}

func (il *IntegerLiteral) Value() interface{} {
	return il.value
}

func (il *IntegerLiteral) GetLine() int {
	return il.Line
}

func (il *IntegerLiteral) GetColumn() int {
	return il.Column
}

type IdentifierLiteral struct {
	value  string
	Type   int
	Line   int
	Column int
}

func (il *IdentifierLiteral) String() string {
	return il.value
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

func (ie *InfixNode) String() string {
	return fmt.Sprintf("%s %s %s", ie.Left.String(), ie.Operator, ie.Right.String())
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

func (pe *PrefixNode) String() string {
	return string(pe.Operator)
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

type IfNode struct {
	Condition   Node
	Consequence Node
	Alternative Node
	Line        int
	Column      int
}

func (ie *IfNode) String() string {
	return "if"
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

func (fe *ForNode) String() string {
	return "for"
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

func (bs *BlockStatement) String() string {
	return "if"
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

func NewIntegerLiteral(value int, line, column int) *IntegerLiteral {
	return &IntegerLiteral{
		value:  value,
		Line:   line,
		Column: column,
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

func NewFunctionLiteral(name string, parameters []Node, body *BlockStatement, line, column int) *FunctionLiteral {
	return &FunctionLiteral{
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Line:       line,
		Column:     column,
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

func NewFloatLiteral(value float64, line, column int) *FloatLiteral {
	return &FloatLiteral{
		value:  value,
		Line:   line,
		Column: column,
	}
}

func NewStringLiteral(value string, line, column int) *StringLiteral {
	return &StringLiteral{
		value:  value,
		Line:   line,
		Column: column,
	}
}

func NewBooleanLiteral(value bool, line, column int) *BooleanLiteral {
	return &BooleanLiteral{
		value:  value,
		Line:   line,
		Column: column,
	}
}

type BooleanLiteral struct {
	value  bool
	Line   int
	Column int
}

func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value())
}

func (bl *BooleanLiteral) Value() interface{} {
	return bl.value
}

func (bl *BooleanLiteral) GetLine() int {
	return bl.Line
}

func (bl *BooleanLiteral) GetColumn() int {
	return bl.Column
}

type ModuleLiteral struct {
	Name   string
	Nodes  []Node
	Line   int
	Column int
}

func (m *ModuleLiteral) String() string {
	return m.Name
}

func (m *ModuleLiteral) Value() interface{} {
	return m.Nodes
}

func (m *ModuleLiteral) GetLine() int {
	return m.Line
}

func (m *ModuleLiteral) GetColumn() int {
	return m.Column
}

func NewModuleLiteral(name string, Nodes []Node, line, column int) *ModuleLiteral {
	return &ModuleLiteral{
		Name:   name,
		Nodes:  Nodes,
		Line:   line,
		Column: column,
	}
}

type ModuleListNode struct {
	Modules []Node
	Line    int
	Column  int
}

func (mle *ModuleListNode) String() string {
	var out bytes.Buffer

	modules := []string{}
	for _, module := range mle.Modules {
		modules = append(modules, module.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(modules, ", "))
	out.WriteString(")")

	return out.String()
}

func (mle *ModuleListNode) Value() interface{} {
	return mle.Modules
}

func (mle *ModuleListNode) GetLine() int {
	return mle.Line
}

func (mle *ModuleListNode) GetColumn() int {
	return mle.Column
}

func NewModuleListNode(modules []Node, line, column int) *ModuleListNode {
	return &ModuleListNode{
		Modules: modules,
		Line:    line,
		Column:  column,
	}
}

type DotNotationNode struct {
	Line   int
	Column int
	Left   Node
	Right  Node
}

func (dn *DotNotationNode) String() string     { return dn.Left.String() + "." + dn.Right.String() }
func (dn *DotNotationNode) Value() interface{} { return dn } // DotNotationNode itself can be the value.
func (dn *DotNotationNode) GetLine() int       { return dn.Line }
func (dn *DotNotationNode) GetColumn() int     { return dn.Column }

type IncrementNode struct {
	Line    int
	Column  int
	Operand *IdentifierLiteral
}

func (in *IncrementNode) String() string {
	return fmt.Sprintf("%s++", in.Operand.String())
}

func (i *IncrementNode) Value() interface{} { return i }
func (i *IncrementNode) GetLine() int       { return i.Line }
func (i *IncrementNode) GetColumn() int     { return i.Column }

type DecrementNode struct {
	Line    int
	Column  int
	Operand *IdentifierLiteral
}

func (dn *DecrementNode) String() string {
	return fmt.Sprintf("%s--", dn.Operand.String())
}

func (d *DecrementNode) Value() interface{} { return d }
func (d *DecrementNode) GetLine() int       { return d.Line }
func (d *DecrementNode) GetColumn() int     { return d.Column }

type StructLiteral struct {
	Fields map[string]Node
	Line   int
	Column int
}

func NewStructLiteral(fields map[string]Node, line, column int) *StructLiteral {
	return &StructLiteral{
		Fields: fields,
		Line:   line,
		Column: column,
	}
}

func (sl *StructLiteral) String() string {
	fieldStrings := make([]string, 0, len(sl.Fields))
	for name, value := range sl.Fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", name, value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(fieldStrings, ", "))
}

func (sl *StructLiteral) Value() interface{} {
	return sl.Fields
}

func (sl *StructLiteral) GetLine() int {
	return sl.Line
}

func (sl *StructLiteral) GetColumn() int {
	return sl.Column
}

type VariableDeclaration struct {
	Identifier *IdentifierLiteral
	Type       Token
	Line       int
	Column     int
}

func (vd *VariableDeclaration) String() string {
	return fmt.Sprintf("%s: %s", vd.Identifier.String(), vd.Type.Value)
}

func (vd *VariableDeclaration) Value() interface{} { return vd }
func (vd *VariableDeclaration) GetLine() int       { return vd.Line }
func (vd *VariableDeclaration) GetColumn() int     { return vd.Column }
