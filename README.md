# GoScript


GoScript is a dynamically typed, interpreted language created out of curiosity to answer that question we ask as programmers: "How do you make a programming language from scratch?" So, I've given it a try.

## Requirements:
- Go 1.21

## Install

```bash
go install github.com/hyperioxx/goscript/cmd/goscript@latest 
```

Note: This is still a work in progress 

Broken things:
- call stack not fully implemented
- scoping 
- native goscript function calls
- native function arguments 


Example syntax:
```
// variable declaration
myInt = 1
myString = "foo"
myFloat = 1.0
myArray = [1,2,3,4,"bar"]



// conditionals
if myInt > 1 {
    print(x) // builtin function 
}
```

