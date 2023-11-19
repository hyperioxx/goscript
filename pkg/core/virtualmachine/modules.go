package virtualmachine

import "fmt"

type ModuleLoader interface {
	GetModule(name string) *Module
}

type DefaultModuleLoader struct {
	debug bool
}

func NewModuleLoader(debug bool) ModuleLoader {
	return &DefaultModuleLoader{debug: debug}
}

func (l *DefaultModuleLoader) GetModule(name string) *Module {
	if l.debug {
		fmt.Println("I GET HERE")
	}

	return &Module{Name: name}
}
