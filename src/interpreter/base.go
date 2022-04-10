package interpreter

import (
	"fmt"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
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
				fmt.Fprint(shared.PrintDest, string(cp.value))
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

var StdLibFieldValues map[string]*TekoObject = map[string]*TekoObject{
	"print": tp(TekoPrint),
	"map":   tp(ArrayMap),
}

var StdLibSymbolTable SymbolTable = SymbolTable{
	parent: nil,
	table:  StdLibFieldValues,
}

var StdLibModule InterpreterModule = InterpreterModule{
	scope: &BasicObject{
		StdLibSymbolTable,
	},
}

func SetupStdLibModule() {
	StdLibModule.Execute(&lexparse.StdLibParser.Codeblock)
}
