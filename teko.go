package main

import (
  "github.com/cstuartroe/teko/src/lexparse"
  "fmt"
  "os"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Please supply exactly one argument, the filename")
    os.Exit(0)
  }
  statements := lexparse.ParseFile(os.Args[1])
  for _, s := range(statements) {
    s.Printout(0)
  }
}
