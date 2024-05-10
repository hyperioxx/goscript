package commands

import (
	"fmt"
	"os"

	"goscript/pkg/core"
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

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	fileContent := string(fileBytes)
	l := core.NewV1Lexer(fileContent)
	p := core.NewV1Parser(l, *f.debugFlag)
	e := core.NewEvaluator(*f.debugFlag)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, err = e.Evaluate(program)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileHandler) Name() string {
	return "FileHandler"
}
