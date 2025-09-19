package main

import (
	"os"

	"github.com/hyperioxx/goscript/pkg/commands"
)

func main() {
	app := commands.NewApplication(os.Args, os.Exit)
	app.Run()
}
