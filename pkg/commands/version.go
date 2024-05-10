package commands

import (
	"fmt"

	"goscript/pkg/version"
)

type VersionCommand struct {
	version version.Version
}

func NewVersionCommand(v version.Version) *VersionCommand {
	return &VersionCommand{version: v}
}

func (v *VersionCommand) Execute(args []string) error {
	fmt.Println(v.version.String())
	return nil
}

func (v *VersionCommand) Name() string {
	return "Version"
}
