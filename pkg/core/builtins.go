package core

import "fmt"

func gsprint(args []Object) (Object, error) {
	interfaceArgs := make([]interface{}, len(args))
	for i, v := range args {
		str, err := v.String()
		if err != nil {
			fmt.Printf("Error converting argument %d to string: %v\n", i, err)
			continue
		}
		interfaceArgs[i] = str.StringValue
	}

	fmt.Println(interfaceArgs...)

	return Nil{}, nil
}

func gslength(args []Object) (Object, error) {
	// TODO: build out builtin len function
	return Nil{}, nil
}
