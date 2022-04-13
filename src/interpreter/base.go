package interpreter

import (
	"fmt"
	"os"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

func TekoPrintExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	s, ok := evaluatedArgs["s"]
	if !ok {
		panic("No parameter passed to print")
	}

	switch p := (*s).(type) {
	case String:
		fmt.Fprint(shared.PrintDest, string(p.runes))
	default:
		panic("Non-string somehow made it past the type checker as an argument to print!")
	}
	return nil
}

var TekoPrint TekoFunction = customExecutedFunction(TekoPrintExecutor, []string{"s"})

func getProcessArgs() []*TekoObject {
	out := []*TekoObject{}

	for _, arg := range os.Args[2:] {
		out = append(out, tp(newTekoString([]rune(arg))))
	}

	return out
}

var Process BasicObject = BasicObject{
	symbolTable: SymbolTable{
		parent: nil,
		table: map[string]*TekoObject{
			"args": tp(newArray(getProcessArgs())),
		},
	},
}

var StdLibFieldValues map[string]*TekoObject = map[string]*TekoObject{
	"print":   tp(TekoPrint),
	"map":     tp(ArrayMap),
	"process": tp(Process),
	"null":    Null,
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
