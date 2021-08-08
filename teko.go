package main

import (
	"fmt"
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/interpreter"
	"github.com/cstuartroe/teko/src/lexparse"
	"os"
)

func main() {
	interpreter.VerifyStdlibDeclarations()
	checker.SetupFunctionTypes()
	checker.SetupStringTypes()

	if len(os.Args) < 2 {
		fmt.Println("Please supply exactly one argument, the filename")
		os.Exit(1)
	}

	main_codeblock := lexparse.ParseFile(os.Args[1], true)

	if len(os.Args) >= 3 && os.Args[2] == "-p" {
		for _, stmt := range main_codeblock.Statements {
			lexparse.PrintNode(stmt)
		}
	}

	checker.CheckTree(main_codeblock)
	interpreter.ExecuteTree(&main_codeblock)
}
