package lexparse

type tokenType string

const (
  SymbolT tokenType = "symbol"

  StringT = "string"
  CharT = "char"
  IntT = "int"
  FloatT = "float"
  BoolT = "bool"
  NullT = "null"

  ForT = "for"
  WhileT = "while"
  InT = "in"
  TypeT = "type"
  LetT = "let"

  BinopT = "binary operation" // + - * / ^ % & |
  ComparisonT = "comparison" // == != < > <= >=
  SetterT = "setter" // = += -= *= /= ^= %= &= |= ->
  PrefixT = "prefix" // ! ~ ?
  SuffixT = "suffix" // $ # .

  LParT = "("
  RParT = ")"
  LSquareBrT = "["
  RSquareBrT = "]"
  LCurlyBrT = "{"
  RCurlyBrT = "}"
  LAngleT = "<"
  RAngleT = ">"

  DotT = "."
  QMarkT = "?"
  EllipsisT = ".."
  CommaT = ","
  SemicolonT = ";"
  ColonT = ":"
  SubtypeT = "<:"
)

// Go doesn't have sets, which is dumb.
var punct_combos map[string]bool = map[string]bool {
  "==": true,
  "!=": true,
  "<=": true,
  ">=": true,
  "+=": true,
  "-=": true,
  "*=": true,
  "/=": true,
  "^=": true,
  "%=": true,
  "&=": true,
  "|=": true,
  "->": true,
  "<-": true,
  "<:": true,
  "..": true,
}

var binops map[string]string = map[string]string {
  "+": "add",
  "-": "sub",
  "*": "mult",
  "/": "div",
  "^": "exp",
  "%": "mod",
  "&": "and",
  "|": "or",
}

var comparisons map[string]string = map[string]string {
  "==": "eq",
  "!": "neq",
  "<": "lt",
  ">": "gt",
  "<=": "leq",
  ">=": "geq",
}

var setters map[string]string = map[string]string {
  "=": "=",
  "+=": "add",
  "-=": "sub",
  "*=": "mult",
  "/=": "div",
  "^=": "exp",
  "%=": "mod",
  "&=": "and",
  "|=": "or",
  "->": "->",
  "<-": "<-",
}

var prefixes map[string]string = map[string]string {
  "!": "not",
  "?": "to_bool",
  "~": "explode",
}

var suffixes map[string]string = map[string]string {
  "$": "to_str",
  "#": "to_int",
  ".": "to_float",
}
