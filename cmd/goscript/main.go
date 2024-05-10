package main

import (
	"os"

	"goscript/pkg/commands"
)

func main() {
	app := commands.NewApplication(os.Args, os.Exit)
	app.Run()
}
