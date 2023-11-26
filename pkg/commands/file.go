package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/hyperioxx/goscript/pkg/core/compiler"
	"github.com/hyperioxx/goscript/pkg/core/lexer"
	"github.com/hyperioxx/goscript/pkg/core/parser"
	"github.com/hyperioxx/goscript/pkg/core/virtualmachine"
)

type FileHandler struct {
	debugFlag *bool
}

func NewFileHandler(debugFlag *bool) *FileHandler {
	return &FileHandler{
		debugFlag: debugFlag,
	}
}

func (f *FileHandler) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing file name")
	}
	filename := args[len(args)-1]

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	fileContent := string(fileBytes)
	l := lexer.NewV1Lexer(fileContent)
	p := parser.NewV1Parser(l, *f.debugFlag)
	exp := p.ParseProgram()

	c := compiler.NewCompiler()
	instructions, err := c.Compile(exp, *f.debugFlag)
	if err != nil {
		return fmt.Errorf("error compiling expression: %v\n", err)
	}

	vm := virtualmachine.NewVirtualMachine(*f.debugFlag)
	vm.Run(instructions)

	return nil
}

func (f *FileHandler) Name() string {
	return "FileHandler"
}
