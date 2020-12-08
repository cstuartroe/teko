package main

import (
	"fmt"
	"github.com/cstuartroe/teko/src/checker"
	"os"
)

func main() {
	checker.SetupFunctionTypes()

	if len(os.Args) != 2 {
		fmt.Println("Please supply exactly one argument, the filename")
		os.Exit(1)
	}
	checker.LoadFile(os.Args[1])
}
