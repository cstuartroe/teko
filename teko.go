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
  tokens := lexparse.LexFile(os.Args[1])
  for _, t := range(tokens) {
    fmt.Printf("col: %d type: %d value: %s\n", t.Col, t.TType, string(t.Value))
  }
}
