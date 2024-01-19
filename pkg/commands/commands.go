package commands

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/hyperioxx/goscript/pkg/version"
)

const ModulePath string = "./"

type Command interface {
	Execute(args []string) error
	Name() string
}

type CommandFactory func(debugFlag *bool) (Command, error)

type Application struct {
	debugFlag   *bool
	versionFlag *bool
	commands    map[string]CommandFactory
	args        []string
	exit        func(int)
}

func NewApplication(args []string, exit func(int)) *Application {
	debugFlag := flag.Bool("v", false, "verbose")
	versionFlag := flag.Bool("version", false, "Print version information")

	commands := map[string]CommandFactory{}

	return &Application{
		debugFlag:   debugFlag,
		versionFlag: versionFlag,
		commands:    commands,
		args:        args,
		exit:        exit,
	}
}

func (app *Application) Run() {
	flag.Parse() // Parse flags early to read debug flag

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			fmt.Println("Caught interrupt signal, stopping gracefully...")
			app.exit(0) // Exit the application gracefully
		}
	}()

	if *app.versionFlag {
		NewVersionCommand(version.GetVersion()).Execute(nil)
		return
	}

	if *app.debugFlag {
		cpuFile, err := os.Create("cpu.prof")
		if err != nil {
			fmt.Printf("Could not create CPU Profile: %v\n", err)
		}
		defer cpuFile.Close()
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			fmt.Printf("Could not start CPU Profile: %v\n", err)
		}
		defer pprof.StopCPUProfile()

		memFile, err := os.Create("mem.prof")
		if err != nil {
			fmt.Printf("Could not create memory profile: %v\n", err)
		}
		defer memFile.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			fmt.Printf("Could not write memory profile: %v\n", err)
		}
	}

	if len(app.args) <= 2 && !strings.HasSuffix(app.args[len(app.args)-1], ".gs") {

		interpreter := NewInterpreter(app.debugFlag, version.GetVersion())
		err := interpreter.Execute(nil)
		if err != nil {
			fmt.Printf("Interpreter failed: %s\n", err.Error())
			app.exit(1)
		}
		return
	}

	if cmdFactory, ok := app.commands[app.args[1]]; ok {
		cmd, err := cmdFactory(app.debugFlag)
		if err != nil {
			fmt.Printf("Command creation failed: %s\n", err.Error())
			app.exit(1)
		}
		err = cmd.Execute(app.args)
		if err != nil {
			fmt.Printf("Command failed: %s\n", err.Error())
			app.exit(1)
		}
	} else {
		if strings.HasSuffix(app.args[len(app.args)-1], ".gs") {
			fileHandler := NewFileHandler(app.debugFlag)
			err := fileHandler.Execute(app.args)
			if err != nil {
				fmt.Printf("File handler failed: %s\n", err.Error())
				app.exit(1)
			}
			app.exit(0)
		} else {
			fmt.Printf("Unknown command: %s\n", app.args[1])
			app.exit(1)
		}
	}
}
