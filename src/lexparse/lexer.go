package lexparse

import (
  "fmt"
  "strings"
  "os"
  "io/ioutil"
  "unicode"
)

const INDENT_AMOUNT int = 2

type Line struct {
  Num int
  Value []rune
  Filename string
}

type Token struct {
  Line *Line
  Col int
  TType tokenType
  Value []rune
}

func (t Token) to_str() string {
  return fmt.Sprintf(
    "{line: %d, col: %d, type: %s, value: %s}",
    t.Line.Num,
    t.Col,
    t.TType,
    string(t.Value),
  )
}

func (t Token) to_indented_str(indent int) string {
  return strings.Repeat(" ", indent*INDENT_AMOUNT) + t.to_str() + "\n"
}

func lexparsePanic(line *Line, col int, width int, message string) {
  fmt.Printf(
    "Teko Parser Error (%s:%d:%d)\n%s\n%s%s\n%s\n",
    line.Filename,
    line.Num,
    col,
    string(line.Value),
    strings.Repeat(" ",col),
    strings.Repeat("^",width),
    message,
  )
  os.Exit(0)
}

func TokenPanic(token Token, message string) {
  lexparsePanic(token.Line, token.Col, len(token.Value), message)
}

type Lexer struct {
  Line *Line
  Col int
  CurrentBlob []rune
  InBlockComment bool
  CurrentTType tokenType
}

func LexFile(filename string) []Token {
  dat, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  contents := string(dat)

  lines_raw := strings.Split(contents, "\n")
  lines := make([]Line, len(lines_raw))

  for i, s := range lines_raw {
    lines[i] = Line{
      Num: i,
      Value: []rune(s),
      Filename: filename}
  }

  tokens := []Token{}
  lexer := Lexer{InBlockComment: false}

  for _, line := range lines {
    lexer.startline(line)
    tokens = append(tokens, lexer.grabTokens()...)
  }

  return tokens
}

func (lexer *Lexer) startline(line Line) {
  lexer.Line = &line
  lexer.Col = 0
}

func (lexer *Lexer) newToken() {
  lexer.CurrentBlob = []rune{}
  lexer.CurrentTType = ""
}

func (lexer *Lexer) next() rune {
  if !lexer.hasMore() { return rune(0) }
  return lexer.Line.Value[lexer.Col]
}

func (lexer *Lexer) advance() {
  lexer.CurrentBlob = append(lexer.CurrentBlob, lexer.next())
  lexer.Col++
}

func (lexer *Lexer) hasMore() bool {
  return lexer.Col < len(lexer.Line.Value)
}

func (lexer *Lexer) passWhitespace() {
  for unicode.IsSpace(lexer.next()) {
    lexer.advance()
  }
}

func (lexer *Lexer) currentToken() Token {
  return Token {
    Line: lexer.Line,
    Col: lexer.Col - len(lexer.CurrentBlob),
    TType: lexer.CurrentTType,
    Value: lexer.CurrentBlob,
  }
}

func (lexer *Lexer) grabTokens() []Token {
  tokens := []Token{}
  lexer.passWhitespace()

  for lexer.hasMore() {
    tokens = append(tokens, lexer.grabToken())
    lexer.passWhitespace()
  }

  return tokens
}

func (lexer *Lexer) grabToken() Token {
  lexer.newToken()

  c := lexer.next()
  if unicode.IsLetter(c) || (c == '_') {
    lexer.grabSymbol()
  } else if unicode.IsDigit(c) {
    lexer.grabDecimalNumber()
  } else if c == '"' {
    lexer.grabString()
  } else if c == '\'' {
    lexer.grabChar()
  } else {
    lexer.grabPunctuation()
  }

  return lexer.currentToken()
}

func (lexer *Lexer) grabSymbol() {
  c := lexer.next()
  for unicode.IsLetter(c) || unicode.IsDigit(c) || (c == '_') {
    lexer.advance()
    c = lexer.next()
  }

  switch symbol := string(lexer.CurrentBlob); symbol {
    case "true":  lexer.CurrentTType = BoolT
    case "false": lexer.CurrentTType = BoolT
    case "for":   lexer.CurrentTType = ForT
    case "while": lexer.CurrentTType = WhileT
    case "in":    lexer.CurrentTType = InT
    case "type":  lexer.CurrentTType = TypeT
    case "data": lexer.CurrentTType = DataT
    case "let":   lexer.CurrentTType = LetT
    default:      lexer.CurrentTType = SymbolT
  }
}

func (lexer *Lexer) grabDecimalNumber() {
  for unicode.IsDigit(lexer.next()) {
    lexer.advance()
  }

  if lexer.next() != '.' {
    lexer.CurrentTType = IntT
  } else {
    lexer.CurrentTType = FloatT

    lexer.advance()

    for unicode.IsDigit(lexer.next()) {
      lexer.advance()
    }
  }
}

func (lexer *Lexer) grabPunctuation() {
  lexer.advance()

  twochars := string(append(lexer.CurrentBlob, lexer.next()))
  if _, ok := punct_combos[twochars]; ok {
    lexer.advance()
  }

  blob := string(lexer.CurrentBlob)

  // Need to check specific cases first because some
  // (such as angle brackets) can be interpreted as
  // multiple categories and may be reassigned a token
  // type later based on context
  switch blob {
    case "(":  lexer.CurrentTType = LParT
    case ")":  lexer.CurrentTType = RParT
    case "[":  lexer.CurrentTType = LSquareBrT
    case "]":  lexer.CurrentTType = RSquareBrT
    case "{":  lexer.CurrentTType = LCurlyBrT
    case "}":  lexer.CurrentTType = RCurlyBrT
    case "<":  lexer.CurrentTType = LAngleT
    case ">":  lexer.CurrentTType = RAngleT
    case ".":  lexer.CurrentTType = DotT
    case "?":  lexer.CurrentTType = QMarkT
    case "..": lexer.CurrentTType = EllipsisT
    case ",":  lexer.CurrentTType = CommaT
    case ";":  lexer.CurrentTType = SemicolonT
    case ":":  lexer.CurrentTType = ColonT
    case "<:": lexer.CurrentTType = SubtypeT
    default:
  }

  if lexer.CurrentTType == "" {
    if _, ok := binops[blob]; ok {
      lexer.CurrentTType = BinopT
    } else if _, ok := comparisons[blob]; ok {
      lexer.CurrentTType = ComparisonT
    } else if _, ok := setters[blob]; ok {
      lexer.CurrentTType = SetterT
    } else if _, ok := prefixes[blob]; ok {
      lexer.CurrentTType = PrefixT
    } else if _, ok := suffixes[blob]; ok {
      lexer.CurrentTType = SuffixT
    } else {
      lexparsePanic(lexer.Line, lexer.Col-1, 1, "Invalid start to token")
    }
  }
}

func (lexer *Lexer) grabString() {
  if lexer.next() != '"' {
    panic("String didn't start with double quote")
  }
  lexer.CurrentTType = StringT
  lexer.advance()

  value := []rune{}

  for lexer.next() != '"' {
    value = append(value, lexer.grabCharacter())
  }

  lexer.advance()

  lexer.CurrentBlob = value // a shameless hack
}

func (lexer *Lexer) grabChar() {
  if lexer.next() != '\'' {
    panic("Char didn't start with single quote")
  }
  lexer.CurrentTType = CharT
  lexer.advance()

  value := []rune{lexer.grabCharacter()}

  if lexer.next() != '\'' {
    panic("Char too long")
  }
  lexer.advance()

  lexer.CurrentBlob = value // the same shameless hack - look any better this time?
}

func (lexer *Lexer) grabCharacter() rune {
  if lexer.next() != '\\' {
    c := lexer.next(); lexer.advance(); return c
  }
  lexer.advance()

  var c rune
  switch lexer.next() {
  case 't': c = '\t'; lexer.advance()
  case 'n': c = '\n'; lexer.advance()
  case 'r': c = '\r'; lexer.advance()
  case '"': c = '"'; lexer.advance()
  case '\'': c = '\''; lexer.advance()
  default: panic("Ahh! We still need octal and unicode escapes!!!1!!one!")
  }

  return c
}
