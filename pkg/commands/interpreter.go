package commands

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hyperioxx/goscript/pkg/core"
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
	e := core.NewEvaluator(*i.debugFlag)
	
	var multiLine string
	isMultiLine := false
	prompt := ">> "

	for {
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

		l := core.NewV1Lexer(line)
		p := core.NewV1Parser(l, *i.debugFlag)
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println(err)
			return nil
		}
        for _, exp := range program {
			value, err := e.Evaluate(exp)
			if err != nil {
				fmt.Println(err)
			}
			if _, ok := value.(*core.Nil); !ok {
				fmt.Println(value.Value())
			}
		}
	}
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
