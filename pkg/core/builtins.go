package core

import "fmt"

func gsprint(args []Object) (Object, error) {
	for _, v := range args {
		str := v.String()
		fmt.Print(str.String().value)
	}
	fmt.Print("\n")
	return &Nil{}, nil
}

func gslength(args []Object) (Object, error) {
	// TODO: build out builtin len function
	return &Nil{}, nil
}
