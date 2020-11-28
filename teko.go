package main

import (
  "github.com/cstuartroe/teko/src/lexparse"
)

func main() {
  l := lexparse.Line{
    Num: 45,
    Value: "let x = 3;",
    Filename: "/this/is/not/a/file.to",
  }

  lexparse.GrabTokens(l)
}
