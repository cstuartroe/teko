package lexparse

import (
  "fmt"
  "strings"
)

type tokenType int

const (
  SymbolT tokenType = iota

  StringT tokenType = iota
  CharT tokenType = iota
  IntT tokenType = iota
  FloatT tokenType = iota
  BoolT tokenType = iota
  NullT tokenType = iota

  BinopT tokenType = iota // + - * / ^ % == != <= >=
  SetterT tokenType = iota // = += -= *= /= ^= %=
  PrefixT tokenType = iota // ! ~
  Suffix tokenType = iota // $ #

  LParT tokenType = iota
  RParT tokenType = iota
  LSquareBrT tokenType = iota
  RSquareBrT tokenType = iota
  LCurlyBrT tokenType = iota
  RCurlyBrT tokenType = iota
  LAngleT tokenType = iota
  RAngleT tokenType = iota

  DotT tokenType = iota
  QMarkT tokenType = iota
  EllipsisT tokenType = iota
  CommaT tokenType = iota
  SemicolonT tokenType = iota
  ColonT tokenType = iota
)

type Line struct {
  Num int
  Value string
  Filename string
}

type Token struct {
  Line Line
  Col int
  TType tokenType
  Value string
}

func lexerPanic(token Token, message string) {
  s := fmt.Sprintf(
    "Teko Parser Error (%s:%d)\n%s\n%s%s\n%s",
    token.Line.Filename,
    token.Col,
    token.Line.Value,
    strings.Repeat(" ",token.Col),
    strings.Repeat("^",len(token.Value)),
    message,
  )
  panic(s)
}

func GrabTokens(line Line) []Token {
  t := Token{line, 4, SymbolT, line.Value[4:5]}
  lexerPanic(t, "Oh noes! ur a bad programmer! I'm gonna tell ur mom!")
  return []Token{t}
}

// func GrabToken(line Line, col uint8) {
// }
