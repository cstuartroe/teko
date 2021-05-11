package interpreter

import (
	"fmt"
	"github.com/cstuartroe/teko/src/lexparse"
)

var blankInterpreter *InterpreterModule = &InterpreterModule{}

func TekoPrintExecutor(function *TekoFunction, evaluatedArgs map[string]TekoObject) TekoObject {
	s, ok := evaluatedArgs["s"]
	if !ok {
		panic("No parameter passed to print")
	}

	switch p := (s).(type) {
	case *String:
		fmt.Printf(string(p.value))
	default:
		panic("Non-string somehow made it past the type checker as an argument to print!")
	}
	return nil
}

var TekoPrint *TekoFunction = &TekoFunction{
	context:  blankInterpreter,
	body:     lexparse.Codeblock{},
	argnames: []string{"s"},
	executor: TekoPrintExecutor,
}

var BaseInterpreterFieldValues map[string]TekoObject = map[string]TekoObject{
	"print": TekoPrint,
}

var BaseSymbolTable *SymbolTable = &SymbolTable{
	parent: nil,
	table:  BaseInterpreterFieldValues,
}
