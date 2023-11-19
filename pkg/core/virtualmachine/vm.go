package virtualmachine

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/hyperioxx/goscript/pkg/utils"
)

type StackFrame struct {
	FunctionName       string
	LineNumber         int
	InstructionPointer int
	Instructions       []Instruction
	Args               []Object
	Scope              map[string]Object
	ReturnValue        Object
}

type CallStack []StackFrame

type OpCode int

const (
	OpAdd OpCode = iota
	OpSub
	OpMultiply
	OpDivide
	OpModulo
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpGreaterThanOrEqual
	OpLessThanOrEqual
	OpJump
	OpPush
	OpAssign
	OpGet
	OpCall
	OpStoreFunc
	OpReturn
	OpJumpIfFalse
	OpJumpIfTrue
	OpCreateThread
	OpCreateArray
	OpPop
	OpImport    // import "foo"
	OpIncrement // foo++
	OpDecrement // foo--
)

var OpCodeStrings = map[OpCode]string{
	OpAdd:                "Add",
	OpSub:                "Sub",
	OpMultiply:           "Multiply",
	OpDivide:             "Divide",
	OpModulo:             "Modulo",
	OpEqual:              "Equal",
	OpNotEqual:           "NotEqual",
	OpGreaterThan:        "GreaterThan",
	OpLessThan:           "LessThan",
	OpGreaterThanOrEqual: "GreaterThanOrEqual",
	OpLessThanOrEqual:    "LessThanOrEqual",
	OpJump:               "Jump",
	OpPush:               "Push",
	OpAssign:             "Assign",
	OpGet:                "Get",
	OpCall:               "Call",
	OpStoreFunc:          "StoreFunc",
	OpReturn:             "Return",
	OpJumpIfFalse:        "JumpIfFalse",
	OpJumpIfTrue:         "JumpIfTrue",
	OpCreateThread:       "CreateThread",
	OpCreateArray:        "CreateArray",
	OpPop:                "Pop",
	OpImport:             "Import",
	OpIncrement:          "Increment",
}

type Instruction struct {
	Opcode OpCode
	Value  Object
	Line   int
	Column int
}

type VM struct {
	Instructions []Instruction
	Stack        []Object
	Variables    map[string]Object
	CallStack    CallStack
	callStackPtr int
	Debug        bool
	Ticks        int
	MaxStackSize int
	IP           int
	IsRunning    bool
	ModuleLoader ModuleLoader
}

func NewVirtualMachine(debug bool) *VM {
	vm := &VM{
		Stack:        make([]Object, 0),
		CallStack:    CallStack{StackFrame{FunctionName: "main", LineNumber: 1, Scope: map[string]Object{}}},
		Variables:    make(map[string]Object),
		Debug:        debug,
		MaxStackSize: 100000,
		IsRunning:    true,
		ModuleLoader: NewModuleLoader(debug),
	}
	vm.registerBuiltin("print", gsprint)
	return vm
}

func (vm *VM) StackTrace(msg string, curInstr Instruction, instrPointer int, line int, column int) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("\nError: %s\n", msg)
	fmt.Printf("Timestamp: %s\n", time.Now().Format(time.RFC3339Nano))
	fmt.Printf("Allocated memory: %v bytes\n", mem.Alloc)
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Line: %d, Column: %d\n", line, column)

	for depth := len(vm.CallStack) - 1; depth >= 0; depth-- {
		trace := vm.CallStack[depth]
		fmt.Printf("%d in %s\n", trace.LineNumber, trace.FunctionName)
		fmt.Printf("\tOpcode: %s\n", OpCodeStrings[curInstr.Opcode])
		fmt.Printf("\tArguments: %v\n", trace.Args)
	}
}

func (vm *VM) Run(instructions []Instruction) Object {

	vm.CurrentFrame().Instructions = append(vm.CurrentFrame().Instructions, instructions...)

	for vm.IsRunning {
		// check if we have returned from the main function
		if len(vm.CallStack) == 0 {
			break
		}

		if vm.CurrentFrame().InstructionPointer >= len(vm.CurrentFrame().Instructions) {
			break
		}

		instruction := vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer]

		if len(vm.Stack) >= vm.MaxStackSize {
			vm.StackTrace("Error: stack overflow", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			os.Exit(1)
		}

		if vm.Debug {
			instruction := vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer]
			fmt.Printf("%s=============== %sVCPU Cycle %s%d%s  ===============%s\n", utils.ColorCyan, utils.ColorYellow, utils.ColorWhite, vm.Ticks, utils.ColorCyan, utils.ColorReset)
			fmt.Printf("%sExecuting Instruction %d:%s %s %v%s\n", utils.ColorCyan, vm.CurrentFrame().InstructionPointer+1, utils.ColorYellow, OpCodeStrings[instruction.Opcode], instruction.Value, utils.ColorReset)
			fmt.Printf("%sInstruction Pointer: %s%d %s\n", utils.ColorCyan, utils.ColorYellow, vm.CurrentFrame().InstructionPointer, utils.ColorReset)
			fmt.Printf("%sStack:%s %+v%s\n", utils.ColorCyan, utils.ColorYellow, vm.Stack, utils.ColorReset)
			fmt.Printf("%sCall Stack:%s %+v%s\n", utils.ColorCyan, utils.ColorYellow, vm.CallStack, utils.ColorReset)
			fmt.Printf("%sFunction:%s %s%s\n", utils.ColorCyan, utils.ColorYellow, vm.CurrentFrame().FunctionName, utils.ColorReset)
			fmt.Printf("%sScoped:%s %+v%s\n", utils.ColorCyan, utils.ColorYellow, vm.CurrentFrame().Scope, utils.ColorReset)
		}

		switch instruction.Opcode {
		case OpAdd, OpSub, OpMultiply, OpDivide, OpModulo:
			if len(vm.Stack) < 2 {
				vm.StackTrace("not enough operands in the stack", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			a, ok := vm.Stack[len(vm.Stack)-2].(Object)
			if !ok {
				vm.StackTrace("the first operand is not an Object", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			b, ok := vm.Stack[len(vm.Stack)-1].(Object)
			if !ok {
				vm.StackTrace(fmt.Sprintf("the second operand is %s not an Object", b.Type()), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			var result Object
			var err error
			switch instruction.Opcode {
			case OpAdd:
				result, err = a.Add(b)
			case OpSub:
				result, err = a.Sub(b)
			case OpMultiply:
				result, err = a.Multiply(b)
			case OpDivide:
				result, err = a.Divide(b)
			case OpModulo:
				result, err = a.Modulo(b)
			}
			if err != nil {
				vm.StackTrace(fmt.Sprintf("Error: %s", err.Error()), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			vm.Stack = vm.Stack[:len(vm.Stack)-2]
			vm.Stack = append(vm.Stack, result)

		case OpEqual, OpNotEqual, OpGreaterThan, OpLessThan, OpGreaterThanOrEqual, OpLessThanOrEqual:
			if len(vm.Stack) < 2 {
				vm.StackTrace("Error: not enough operands in the stack", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			a, ok := vm.Stack[len(vm.Stack)-2].(Object)
			if !ok {
				vm.StackTrace("Error: the first operand is not an Object", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			b, ok := vm.Stack[len(vm.Stack)-1].(Object)
			if !ok {
				vm.StackTrace("Error: the second operand is not an Object", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			var result Object
			var err error
			switch instruction.Opcode {
			case OpEqual:
				result, err = a.Equal(b)
			case OpNotEqual:
				result, err = a.NotEqual(b)
			case OpGreaterThan:
				result, err = a.GreaterThan(b)
			case OpLessThan:
				result, err = a.LessThan(b)
			case OpGreaterThanOrEqual:
				result, err = a.GreaterThanOrEqual(b)
			case OpLessThanOrEqual:
				result, err = a.LessThanOrEqual(b)
			}
			if err != nil {
				vm.StackTrace(fmt.Sprintf("Error: %s", err.Error()), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			vm.Stack = vm.Stack[:len(vm.Stack)-2]
			vm.Stack = append(vm.Stack, result)
		case OpJumpIfTrue:
			condition, ok := vm.Stack[len(vm.Stack)-1].(*Boolean)
			if !ok {
				vm.StackTrace("Error: the condition is not a Boolean", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}
			// pop off the stack
			vm.Stack = vm.Stack[:len(vm.Stack)-1]

			if condition.BoolValue {
				ipOffset, ok := instruction.Value.(*Integer)
				if !ok {
					vm.StackTrace("Error: jump offset is not an Integer", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}
				// -1 because after the switch, the pointer is incremented
				vm.CurrentFrame().InstructionPointer += ipOffset.IntValue - 1
			}
		case OpJumpIfFalse:
			condition, ok := vm.Stack[len(vm.Stack)-1].(*Boolean)
			if !ok {
				vm.StackTrace("Error: the condition is not a Boolean", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			// pop off the stack
			vm.Stack = vm.Stack[:len(vm.Stack)-1]

			if !condition.BoolValue {
				ipOffset, ok := instruction.Value.(*Integer)
				if !ok {
					vm.StackTrace("Error: jump offset is not an Integer", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}
				// -1 because after the switch, the pointer is incremented
				vm.CurrentFrame().InstructionPointer += ipOffset.IntValue - 1
			}
		case OpJump:
			ipOffset, ok := instruction.Value.(*Integer)
			if !ok {
				vm.StackTrace("Error: jump offset is not an Integer", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}
			// -1 because after the switch, the pointer is incremented
			vm.CurrentFrame().InstructionPointer += ipOffset.IntValue - 1

			// in the OpCall case of the switch
		case OpCall:
			name, ok := instruction.Value.(*String)
			if !ok {
				vm.StackTrace("Error: function name is not a String", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			// Check if the function is a goscript function
			if function, ok := vm.CurrentFrame().Scope[name.StringValue].(*Function); ok {
				expectedArgs := len(function.Parameters)
				// Check if the stack contains enough arguments
				if len(vm.Stack) < expectedArgs {
					vm.StackTrace(fmt.Sprintf("Error: function %s requires %d arguments but stack contains only %d", name.StringValue, expectedArgs, len(vm.Stack)-1), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}

				numArgs := vm.Stack[len(vm.Stack)-1].(*Integer).IntValue
				vm.Stack = vm.Stack[:len(vm.Stack)-1]

				args := vm.Stack[len(vm.Stack)-numArgs:]
				vm.Stack = vm.Stack[:len(vm.Stack)-numArgs]

				// Assign arguments to function parameters
				scope := make(map[string]Object)
				for i, arg := range args {
					scope[function.Parameters[i]] = arg
				}

				// Create a new frame for the called goscript function
				newFrame := StackFrame{
					FunctionName:       name.StringValue,
					LineNumber:         instruction.Line,
					InstructionPointer: 0,
					Instructions:       function.Body,
					Scope:              scope,
				}

				// Push the new frame onto the call stack
				vm.CallStack = append(vm.CallStack, newFrame)
				continue
			} else if goFunc, ok := vm.CurrentFrame().Scope[name.StringValue].(*GoFunction); ok {
				// Check if the stack contains the number of arguments
				if len(vm.Stack) < 1 {
					vm.StackTrace(fmt.Sprintf("Go function %s requires the number of arguments on the stack", name.StringValue), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}

				numArgs := vm.Stack[len(vm.Stack)-1].(*Integer).IntValue
				vm.Stack = vm.Stack[:len(vm.Stack)-1]

				// Check if the stack contains enough arguments
				if len(vm.Stack) < numArgs {
					vm.StackTrace(fmt.Sprintf("Go function %s requires %d arguments but stack contains only %d", name.StringValue, numArgs, len(vm.Stack)), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}

				args := vm.Stack[len(vm.Stack)-numArgs:]
				vm.Stack = vm.Stack[:len(vm.Stack)-numArgs]

				// Call the Go function and push the return value onto the stack
				returnValue, err := goFunc.Call(args)
				if err != nil {
					vm.StackTrace(fmt.Sprintf("Error calling Go function %s: %v", name.StringValue, err), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
				}
				switch returnValue.(type) {
				case Nil:
				default:
					vm.Stack = append(vm.Stack, returnValue)

				}
			} else {
				vm.StackTrace(fmt.Sprintf("Error: function %s is not defined", name.StringValue), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

		case OpReturn:
			var returnValue Object
			if len(vm.Stack) > 0 {
				returnValue = vm.Stack[len(vm.Stack)-1]
				vm.Stack = vm.Stack[:len(vm.Stack)-1]
				vm.CurrentFrame().ReturnValue = returnValue
			}

			vm.CallStack = vm.CallStack[:len(vm.CallStack)-1]

			if len(vm.CallStack) == 0 {
				break
			}

			if returnValue != nil {
				vm.Stack = append(vm.Stack, returnValue)
				vm.CurrentFrame().ReturnValue = returnValue
			}

		case OpPush:
			if vm.Debug {
				fmt.Printf("%s", utils.ColorYellow)
				fmt.Printf("Pushing %v on to stack...\n", instruction.Value.Value())
				fmt.Printf("%s", utils.ColorReset)
			}
			value := instruction.Value
			vm.Stack = append(vm.Stack, value)

			popInstruction := Instruction{
				Opcode: OpPop,
			}
			vm.CurrentFrame().Instructions = append(vm.CurrentFrame().Instructions, popInstruction)
		case OpAssign:
			if len(vm.Stack) < 2 {
				vm.StackTrace("Error: not enough operands in the stack", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			name, ok := vm.Stack[len(vm.Stack)-2].(*String)
			if !ok {
				vm.StackTrace("Error: the variable name is not a String", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			value := vm.Stack[len(vm.Stack)-1]

			vm.CurrentFrame().Scope[name.StringValue] = value

			vm.Stack = vm.Stack[:len(vm.Stack)-2]
		case OpStoreFunc:
			function := instruction.Value.(Callable)
			vm.CurrentFrame().Scope[function.GetName()] = function

		case OpGet:
			name, ok := instruction.Value.(*String)
			if !ok {
				vm.StackTrace("Error: variable name is not a String", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			var value Object
			var found bool
			for i := len(vm.CallStack) - 1; i >= 0; i-- {
				value, found = vm.CallStack[i].Scope[name.StringValue]
				if found {
					break
				}
			}

			if !found {
				vm.StackTrace(fmt.Sprintf("Error: variable %s is not defined", name.StringValue), vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)
			}

			vm.Stack = append(vm.Stack, value)

		case OpImport:
			name, ok := instruction.Value.(*String)
			if !ok {
				vm.StackTrace("Error: variable name is not a String", vm.CurrentFrame().Instructions[vm.CurrentFrame().InstructionPointer], vm.CurrentFrame().InstructionPointer, instruction.Line, instruction.Column)

			}
			module := vm.ModuleLoader.GetModule(name.StringValue)
			vm.CurrentFrame().Scope[module.Name] = module
		}

		vm.CurrentFrame().InstructionPointer++
		vm.Ticks++

	}

	if len(vm.Stack) < 1 || vm.Stack[0] == nil {
		return &Nil{} // Return a pointer to the Nil object
	}

	result := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:0]

	return result
}

func (vm *VM) CurrentFrame() *StackFrame {
	return &vm.CallStack[len(vm.CallStack)-1]
}

func (vm *VM) registerBuiltin(funcName string, fn func([]Object) (Object, error)) {
	goFunc := &GoFunction{Name: funcName, Func: fn}
	vm.CurrentFrame().Scope[funcName] = goFunc
}
