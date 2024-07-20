package core

import "fmt"

func gsprint(args []Object) (Object, error) {
	interfaceArgs := make([]interface{}, len(args))
	for i, v := range args {
		str := v.String()
		interfaceArgs[i] = str.String().value
	}

	fmt.Println(interfaceArgs...)

	return &Nil{}, nil
}

func gslength(args []Object) (Object, error) {
	// TODO: build out builtin len function
	return &Nil{}, nil
}
