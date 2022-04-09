package cmd

import (
	"os"

	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/interpreter"
	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

func ExecuteFile(filename string) {
	p := lexparse.Parser{
		Transform: true,
	}
	p.ParseFile(filename)

	c := checker.NewBaseChecker()
	c.CheckTree(p.Codeblock)

	i := interpreter.New(&interpreter.StdLibModule)
	i.Execute(&p.Codeblock)
}

func ExecuteFileSafe(filename string) {
	defer func() {
		if r := recover(); r == shared.TekoErrorMessage {
			os.Exit(1)
		} else {
			panic(r)
		}
	}()

	ExecuteFile(filename)
}
