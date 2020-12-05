package main

import (
  "fmt"
  "os"
  "github.com/cstuartroe/teko/src/checker"
)

func main() {
  checker.SetupFunctionTypes()

  if len(os.Args) != 2 {
    fmt.Println("Please supply exactly one argument, the filename")
    os.Exit(0)
  }
  codeblock := checker.LoadFile(os.Args[1])
  for _,  stmt := range codeblock.GetStatements() {
    stmt.Printout(0)
  }
}
