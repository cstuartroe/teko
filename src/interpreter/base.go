package interpreter

import (
	"fmt"
)

func TekoPrintExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	s, ok := evaluatedArgs["s"]
	if !ok {
		panic("No parameter passed to print")
	}

	switch p := (*s).(type) {
	case Array:
		for _, c := range p.elements {
			switch cp := (*c).(type) {
			case TekoChar:
				fmt.Printf(string(cp.value))
			default:
				panic("Not a TekoChar")
			}
		}
	default:
		panic("Non-string somehow made it past the type checker as an argument to print!")
	}
	return nil
}

var TekoPrint TekoFunction = customExecutedFunction(TekoPrintExecutor, []string{"s"})

var BaseInterpreterFieldValues map[string]*TekoObject = map[string]*TekoObject{
	"print": tp(TekoPrint),
}

var BaseSymbolTable *SymbolTable = &SymbolTable{
	parent: nil,
	table:  BaseInterpreterFieldValues,
}
