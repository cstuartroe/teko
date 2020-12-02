package lexparse

type tokenType int

const (
  SymbolT tokenType = iota

  StringT
  CharT
  IntT
  FloatT
  BoolT
  NullT

  ForT
  WhileT
  InT
  TypeT
  LetT

  BinopT // + - * / ^ % & |
  ComparisonT // == != < > <= >=
  SetterT // = += -= *= /= ^= %= &= |= ->
  PrefixT // ! ~ ?
  SuffixT // $ # .

  LParT
  RParT
  LSquareBrT
  RSquareBrT
  LCurlyBrT
  RCurlyBrT
  LAngleT
  RAngleT

  DotT
  QMarkT
  EllipsisT
  CommaT
  SemicolonT
  ColonT
  SubtypeT
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
