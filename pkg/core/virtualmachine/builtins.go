package virtualmachine

import "fmt"

func _print(args []Object) (Object, error) {
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

func _length(args []Object) (Object, error) {
	// TODO: build out builtin len function
	return Nil{}, nil
}

func _type(args []Object) (Object, error) {
	if len(args) != 1 {
		return Nil{}, fmt.Errorf("type function only takes one argument")
	}
	args[0].Type()
	return String{StringValue: args[0].Type()}, nil
}
