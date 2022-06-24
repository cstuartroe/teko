package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/cmd"
	"github.com/cstuartroe/teko/src/interpreter"
	"github.com/cstuartroe/teko/src/lexparse"
)

var commands map[string]func() = map[string]func(){
	"test": cmd.RunTests,
}

func commandsList() []string {
	cl := []string{}
	for key := range commands {
		cl = append(cl, key)
	}
	return cl
}

func main() {
	lexparse.ParseStldLib()
	checker.SetupFunctionTypes()
	checker.SetupSequenceTypes()
	checker.SetupStringTypes()
	checker.CheckTekoLangStdLib()
	interpreter.SetupStdLibModule()
	interpreter.VerifyStdlibDeclarations()

	if len(os.Args) < 2 {
		fmt.Printf("Please supply a command or filename (available commands: %s)\n", strings.Join(commandsList(), ", "))
		os.Exit(1)
	}

	command_or_filename := os.Args[1]

	if f, ok := commands[command_or_filename]; ok {
		f()
	} else {
		cmd.ExecuteFileSafe(command_or_filename)
	}
}
