package lexparse

type tokenType int

const (
  SymbolT tokenType = iota

  StringT tokenType = iota
  CharT tokenType = iota
  IntT tokenType = iota
  FloatT tokenType = iota
  BoolT tokenType = iota
  NullT tokenType = iota

  ForT tokenType = iota
  WhileT tokenType = iota
  InT tokenType = iota
  TypeT tokenType = iota

  BinopT tokenType = iota // + - * / ^ % & |
  ComparisonT tokenType = iota // == != < > <= >=
  SetterT tokenType = iota // = += -= *= /= ^= %= &= |= ->
  PrefixT tokenType = iota // ! ~ ?
  SuffixT tokenType = iota // $ # .

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
  SubtypeT tokenType = iota
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
