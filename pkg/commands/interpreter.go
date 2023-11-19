package commands

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hyperioxx/goscript/pkg/core/compiler"
	"github.com/hyperioxx/goscript/pkg/core/lexer"
	"github.com/hyperioxx/goscript/pkg/core/parser"
	"github.com/hyperioxx/goscript/pkg/core/virtualmachine"
	"github.com/hyperioxx/goscript/pkg/version"
)

type Interpreter struct {
	debugFlag *bool
	version   version.Version
}

func NewInterpreter(debugFlag *bool, ver version.Version) *Interpreter {
	return &Interpreter{
		debugFlag: debugFlag,
		version:   ver,
	}
}

func (i *Interpreter) Execute(args []string) error {
	i.printSystemInfo()
	scanner := bufio.NewScanner(os.Stdin)
	vm := virtualmachine.NewVirtualMachine(*i.debugFlag)

	var multiLine string
	isMultiLine := false
	prompt := ">> "

	for vm.IsRunning {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return fmt.Errorf("unable to scan input from stdin")
		}

		line := scanner.Text()

		if strings.Contains(line, "{") && !isMultiLine {
			isMultiLine = true
			prompt = "       "
		}

		if isMultiLine {
			multiLine += line

			if strings.Contains(line, "}") {
				isMultiLine = false
				line = multiLine
				multiLine = ""
				prompt = ">> "
			} else {
				continue
			}
		}

		if line == "" {
			continue
		}

		l := lexer.NewV1Lexer(line)
		p := parser.NewV1Parser(l, *i.debugFlag)
		exp := p.ParseProgram()
		c := compiler.NewCompiler()
		instructions, err := c.Compile(exp, *i.debugFlag)
		if err != nil {
			msg := fmt.Sprintf("error compiling expression: %v\n", err)
			vm.StackTrace(msg, virtualmachine.Instruction{}, 1, 1, 0)
			continue
		}

		result := vm.Run(instructions)
		if result != nil && result.Type() != "nil" {
			fmt.Println(result.Value())
		}
	}
	return nil
}

func (i *Interpreter) printSystemInfo() {
	fmt.Printf("goscript REPL (Version: %s)\n", i.version.String())
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Println("Enter 'exit' to quit the REPL.")
}

func (i *Interpreter) Name() string {
	return "Interpreter"
}
