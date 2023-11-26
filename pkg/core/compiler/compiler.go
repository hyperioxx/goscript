package compiler

import (
	"fmt"

	"github.com/hyperioxx/goscript/pkg/core/lexer"
	"github.com/hyperioxx/goscript/pkg/core/parser"
	"github.com/hyperioxx/goscript/pkg/core/virtualmachine"
)

// Compiler is the interface for compilers so we can create different compilers in future
type Compiler interface {
	Compile(program []parser.Node, debug bool) ([]virtualmachine.Instruction, error)
}

// SymbolEntry is an entry into a symbol table
type SymbolEntry struct {
    Type       string
    ScopeLevel int
} 

// V1Compiler is the initial compiler build
type V1Compiler struct{
	symbolTable map[string] SymbolEntry
}

// NewCompiler creates a new instance on the V1Compiler
func NewCompiler() Compiler {
	return &V1Compiler{}
}

// Compile is the main interface for the compiler
func (c *V1Compiler) Compile(program []parser.Node, debug bool) ([]virtualmachine.Instruction, error) {
	var compiled []virtualmachine.Instruction

	for _, Node := range program {
		instructions, err := c.compileNode(Node, debug)
		if err != nil {
			return nil, fmt.Errorf("error compiling Node: %v", err)
		}
		compiled = append(compiled, instructions...)
	}

	return compiled, nil
}

func (c *V1Compiler) compileNode(node parser.Node, debug bool) ([]virtualmachine.Instruction, error) {
	switch exp := node.(type) {
	case *parser.VariableDeclaration:
		var _type virtualmachine.Object
		switch exp.Type.Type { // TODO: rethink this
		case lexer.INT:
			_type = &virtualmachine.Integer{}
		case lexer.STRING:
			_type = &virtualmachine.String{}
		case lexer.FLOAT:
			_type = &virtualmachine.Float{}
		default:
			return nil, fmt.Errorf("unknown type")
		}
		// TODO: create instructions to create a variable
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: exp.Identifier.String()}},
			{Opcode: virtualmachine.OpPush, Value: _type},
			{Opcode: virtualmachine.OpDeclare},
		}
		if debug {
			fmt.Println("Compiled IntegerLiteral:", instructions)
		}
		return instructions, nil
	case *parser.IntegerLiteral:
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: exp.Value().(int)}},
		}
		if debug {
			fmt.Println("Compiled IntegerLiteral:", instructions)
		}
		return instructions, nil
	case *parser.BooleanLiteral:
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Boolean{BoolValue: exp.Value().(bool)}},
		}
		if debug {
			fmt.Println("Compiled BooleanLiteral:", instructions)
		}
		return instructions, nil
	case *parser.FloatLiteral:
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Float{FloatValue: exp.Value().(float64)}},
		}
		if debug {
			fmt.Println("Compiled FloatLiteral:", instructions)
		}
		return instructions, nil

	case *parser.StringLiteral:
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: exp.Value().(string)}},
		}
		if debug {
			fmt.Println("Compiled StringLiteral:", instructions)
		}
		return instructions, nil

	case *parser.InfixNode:
		right, err := c.compileNode(exp.Right, debug)
		if err != nil {
			return nil, err
		}
		if debug {
			fmt.Println(err)
		}
		left, err := c.compileNode(exp.Left, debug)
		if err != nil {
			return nil, err
		}
		var op virtualmachine.OpCode
		switch exp.Operator {
		case lexer.TokenTypeStr[lexer.ADD]:
			op = virtualmachine.OpAdd
		case lexer.TokenTypeStr[lexer.SUB]:
			op = virtualmachine.OpSub
		case lexer.TokenTypeStr[lexer.ASSIGN]:
			// check if left is ident then must be assignment??
			if id, ok := exp.Left.(*parser.IdentifierLiteral); ok {
				left = []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: id.Value().(string)}}}
				op = virtualmachine.OpAssign
			} else if id, ok := exp.Left.(*parser.VariableDeclaration); ok {
				left = []virtualmachine.Instruction{{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: id.Value().(string)}}}
				op = virtualmachine.OpAssign
			}
		case lexer.TokenTypeStr[lexer.EQ]:
			op = virtualmachine.OpEqual
		case lexer.TokenTypeStr[lexer.GT_EQ]:
			op = virtualmachine.OpGreaterThanOrEqual
		case lexer.TokenTypeStr[lexer.GT]:
			op = virtualmachine.OpGreaterThan
		case lexer.TokenTypeStr[lexer.LT_EQ]:
			op = virtualmachine.OpLessThanOrEqual
		case lexer.TokenTypeStr[lexer.LT]:
			op = virtualmachine.OpLessThan
		case lexer.TokenTypeStr[lexer.NOT_EQ]:
			op = virtualmachine.OpNotEqual
		case lexer.TokenTypeStr[lexer.RETURN]:
			op = virtualmachine.OpReturn
		default:
			return nil, fmt.Errorf("unknown operator: %s", exp.Operator)
		}
		instructions := append(append(left, right...), virtualmachine.Instruction{Opcode: op})
		if debug {
			fmt.Println("Compiled InfixNode:", instructions)
		}
		return instructions, nil

	case *parser.IdentifierLiteral:
		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpGet, Value: &virtualmachine.String{StringValue: exp.Value().(string)}},
		}
		if debug {
			fmt.Println("Compiled IdentifierLiteral:", instructions)
		}
		return instructions, nil
	case *parser.FunctionCall:
		var instructions []virtualmachine.Instruction
		for _, arg := range exp.Arguments {
			argInstructions, err := c.compileNode(arg, debug)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, argInstructions...)
		}
		// Push the number of arguments onto the stack
		instructions = append(instructions, virtualmachine.Instruction{
			Opcode: virtualmachine.OpPush,
			Value:  &virtualmachine.Integer{IntValue: len(exp.Arguments)},
		})

		// Create a call instruction
		callInstruction := virtualmachine.Instruction{
			Opcode: virtualmachine.OpCall,
			Value:  &virtualmachine.String{StringValue: exp.Name},
		}

		instructions = append(instructions, callInstruction)

		if debug {
			fmt.Println("Compiled parser.FunctionCall:", instructions)
		}
		return instructions, nil
	case *parser.FunctionLiteral:
		bodyInstructions := make([]virtualmachine.Instruction, 0)
		for _, bodyExp := range exp.Body {
			inst, err := c.compileNode(bodyExp, debug)
			if err != nil {
				return nil, err
			}
			bodyInstructions = append(bodyInstructions, inst...)
		}

		params := make([]string, len(exp.Parameters))
		for i, p := range exp.Parameters {
			params[i] = p.Identifier.String()
		}
		function := &virtualmachine.Function{
			Name:       exp.Name,
			Parameters: params,
			Body:       bodyInstructions,
		}

		instructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpStoreFunc, Value: function},
		}

		if debug {
			fmt.Println("Compiled FunctionLiteral:", instructions)
		}

		return instructions, nil

	case *parser.ReturnStatement:
		if exp.ReturnValue == nil {
			instructions := []virtualmachine.Instruction{
				{Opcode: virtualmachine.OpReturn, Value: nil},
			}
			return instructions, nil
		}
		returnValue, err := c.compileNode(exp.ReturnValue, debug)
		if err != nil {
			return nil, err
		}
		instructions := append(returnValue, virtualmachine.Instruction{Opcode: virtualmachine.OpReturn})
		if debug {
			fmt.Println("Compiled ReturnStatement:", instructions)
		}
		return instructions, nil
	case *parser.IfNode:
		conditionInstructions, err := c.compileNode(exp.Condition, debug)
		if err != nil {
			return nil, err
		}

		thenInstructions, err := c.compileNodeList(exp.Consequence, debug)
		if err != nil {
			return nil, err
		}

		elseInstructions, err := c.compileNodeList(exp.Alternative, debug)
		if err != nil {
			return nil, err
		}

		conditionInstructions = append(conditionInstructions, virtualmachine.Instruction{
			Opcode: virtualmachine.OpJumpIfFalse,
			Value:  &virtualmachine.Integer{IntValue: len(thenInstructions) + 1},
		})

		instructions := append(conditionInstructions, thenInstructions...)

		if len(elseInstructions) > 0 {
			instructions = append(instructions, virtualmachine.Instruction{
				Opcode: virtualmachine.OpJump,
				Value:  &virtualmachine.Integer{IntValue: len(elseInstructions)},
			})
			instructions = append(instructions, elseInstructions...)
		}

		if debug {
			fmt.Println("Compiled IfNode:", instructions)
		}

		return instructions, nil
	case *parser.ForNode:
		var initialisationInstructions []virtualmachine.Instruction
		var conditionInstructions []virtualmachine.Instruction
		var updaterInstructions []virtualmachine.Instruction
		var bodyInstructions []virtualmachine.Instruction

		var err error

		if exp.Initialisation != nil {
			initialisationInstructions, err = c.compileNode(exp.Initialisation, debug)
			if err != nil {
				return nil, err
			}
		}

		if exp.Condition != nil {
			conditionInstructions, err = c.compileNode(exp.Condition, debug)
			if err != nil {
				return nil, err
			}
		} else {
			// this is our "while True" equivalent
			conditionInstructions = []virtualmachine.Instruction{
				{
					Opcode: virtualmachine.OpPush,
					Value:  &virtualmachine.Boolean{BoolValue: true},
				},
			}
		}

		if exp.Updater != nil {
			updaterInstructions, err = c.compileNode(exp.Updater, debug)
			if err != nil {
				return nil, err
			}
		}

		bodyInstructions, err = c.compileNodeList(exp.Body, debug)
		if err != nil {
			return nil, err
		}

		// assemble the instructions
		// rough layout:
		//		INITIALISATION
		//		CONDITION
		//		JUMP_IF_FALSE (end)
		//		BODY
		//		UPDATER
		//		JUMP (condition)
		//		(end)

		var instructions []virtualmachine.Instruction

		// any of these that are empty don't change `instructions`, no need to check
		instructions = append(instructions, initialisationInstructions...)
		instructions = append(instructions, conditionInstructions...)
		instructions = append(instructions, virtualmachine.Instruction{
			Opcode: virtualmachine.OpJumpIfFalse,
			// +2 to skip the final JUMP
			Value: &virtualmachine.Integer{IntValue: len(bodyInstructions) + len(updaterInstructions) + 2},
		})
		instructions = append(instructions, bodyInstructions...)
		instructions = append(instructions, updaterInstructions...)
		instructions = append(instructions, virtualmachine.Instruction{
			Opcode: virtualmachine.OpJump,
			// +1 to skip the JUMP_IF_FALSE
			//                                       v note the `-`!
			Value: &virtualmachine.Integer{IntValue: -(len(updaterInstructions) + len(bodyInstructions) + 1 + len(conditionInstructions))},
		})

		if debug {
			fmt.Println("Compiled ForNode:", instructions)
		}

		return instructions, nil
	case *parser.ArrayLiteral:
		var elementInstructions []virtualmachine.Instruction
		for _, elem := range exp.Elements {

			elemInstr, err := c.compileNode(elem, debug)
			if err != nil {
				return nil, err
			}
			elementInstructions = append(elementInstructions, elemInstr...)
		}

		for i := 0; i < len(exp.Elements); i++ {
			elementInstructions = append(elementInstructions, virtualmachine.Instruction{
				Opcode: virtualmachine.OpPop,
			})
		}

		arrayInstr := virtualmachine.Instruction{
			Opcode: virtualmachine.OpCreateArray,
			Value:  virtualmachine.Integer{IntValue: len(exp.Elements)},
		}
		elementInstructions = append(elementInstructions, arrayInstr)

		if debug {
			fmt.Println("Compiled ArrayLiteral:", elementInstructions)
		}
		return elementInstructions, nil

	case *parser.ModuleListNode:
		var importInstructions []virtualmachine.Instruction
		for _, module := range exp.Modules {
			moduleName := virtualmachine.String{StringValue: module.String()}
			importInstructions = append(importInstructions, virtualmachine.Instruction{Opcode: virtualmachine.OpImport, Value: &moduleName})
		}
		if debug {
			fmt.Println("Compiled ModuleListNode:", importInstructions)
		}
		return importInstructions, nil
	case *parser.DotNotationNode:
		var instructions []virtualmachine.Instruction

		// Compile the left node (object)
		objectInstructions, err := c.compileNode(exp.Left, debug)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, objectInstructions...)

		// Compile the right node (property)
		propertyInstructions, err := c.compileNode(exp.Right, debug)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, propertyInstructions...)

		// Create a get property instruction
		getPropertyInstruction := virtualmachine.Instruction{
			Opcode: virtualmachine.OpGet,
			// Assuming the Value of the property is its name
			Value: &virtualmachine.String{StringValue: exp.Right.Value().(string)},
		}

		instructions = append(instructions, getPropertyInstruction)

		if debug {
			fmt.Println("Compiled DotNotationNode:", instructions)
		}
		return instructions, nil
	case *parser.IncrementNode:
		// Compile operand to get identifierInstructions
		identifierInstructions, err := c.compileNode(exp.Operand, debug)
		if err != nil {
			return nil, err
		}

		// Define incrementInstructions
		incrementInstructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: exp.Operand.String()}},
		}

		// Combine incrementInstructions and identifierInstructions
		instructions := append(incrementInstructions, identifierInstructions...)

		// Continue with additional instructions
		moreInstructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
			{Opcode: virtualmachine.OpAdd},
			{Opcode: virtualmachine.OpAssign},
		}

		// Append additional instructions
		instructions = append(instructions, moreInstructions...)

		if debug {
			fmt.Println("Compiled IncrementNode:", instructions)
		}
		return instructions, nil

	case *parser.DecrementNode:
		// Compile operand to get identifierInstructions
		identifierInstructions, err := c.compileNode(exp.Operand, debug)
		if err != nil {
			return nil, err
		}

		// Define decrementInstructions
		decrementInstructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.String{StringValue: exp.Operand.String()}},
		}

		// Combine decrementInstructions and identifierInstructions
		instructions := append(decrementInstructions, identifierInstructions...)

		// Continue with additional instructions
		moreInstructions := []virtualmachine.Instruction{
			{Opcode: virtualmachine.OpPush, Value: &virtualmachine.Integer{IntValue: 1}},
			{Opcode: virtualmachine.OpSub},
			{Opcode: virtualmachine.OpAssign},
		}

		// Append additional instructions
		instructions = append(instructions, moreInstructions...)

		if debug {
			fmt.Println("Compiled IncrementNode:", instructions)
		}
		return instructions, nil

	default:
		if debug {
			fmt.Println("unknown Node type: %T", node)
		}
		return nil, fmt.Errorf("unknown Node type: %T\n", node)
	}
}

func (c *V1Compiler) compileNodeList(Nodes []parser.Node, debug bool) ([]virtualmachine.Instruction, error) {
	var instructions []virtualmachine.Instruction

	for _, expr := range Nodes {
		instr, err := c.compileNode(expr, debug)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, instr...)
	}

	return instructions, nil
}
