package lexparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const INDENT_AMOUNT int = 2

type Line struct {
	Num      int
	Value    []rune
	Filename string
}

type Token struct {
	Line  *Line
	Col   int
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

type TekoErrorClass string

const (
	LexerError          = "Lexer Error"
	SyntaxError         = "Syntax Error"
	NameError           = "Name Error"
	NotImplementedError = "Unimplemented (this is planned functionality for Teko)"
	TypeError           = "Type Error"
	ArgumentError       = "Argument Error"
	UnexpectedIssue     = "Unexpected issue (this is a mistake in the Teko implementation)"
	RuntimeError        = "Runtime Error"
)

func lexparsePanic(line *Line, col int, width int, errorClass TekoErrorClass, message string) {
	fmt.Printf(
		"%s (%s:%d:%d)\n%s\n%s%s\n%s\n",
		errorClass,
		line.Filename,
		line.Num,
		col,
		string(line.Value),
		strings.Repeat(" ", col),
		strings.Repeat("^", width),
		message,
	)
	os.Exit(0)
}

func (token Token) Raise(errorClass TekoErrorClass, message string) {
	lexparsePanic(token.Line, token.Col, len(token.Value), errorClass, message)
}

type Lexer struct {
	Line           *Line
	Col            int
	CurrentBlob    []rune
	InBlockComment bool
	CurrentTType   tokenType
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
			Num:      i,
			Value:    []rune(s),
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
	if !lexer.hasMore() {
		return rune(0)
	}
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
	return Token{
		Line:  lexer.Line,
		Col:   lexer.Col - len(lexer.CurrentBlob),
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

var keywords map[string]tokenType = map[string]tokenType{
	"true":  BoolT,
	"false": BoolT,

	"var": VarT,
	"fn":  FnT,

	"in":    InT,
	"for":   ForT,
	"while": WhileT,
	"if":    IfT,
	"then":  ThenT,
	"else":  ElseT,

	"set":  SetT,
	"map":  MapT,
	"data": DataT,
	"type": TypeT,
	"do":   DoT,
}

func (lexer *Lexer) grabSymbol() {
	c := lexer.next()
	for unicode.IsLetter(c) || unicode.IsDigit(c) || (c == '_') {
		lexer.advance()
		c = lexer.next()
	}

	if ttype, ok := keywords[string(lexer.CurrentBlob)]; ok {
		lexer.CurrentTType = ttype
	} else {
		lexer.CurrentTType = SymbolT
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

var uniquePuncts map[string]tokenType = map[string]tokenType{
	"(":  LParT,
	")":  RParT,
	"[":  LSquareBrT,
	"]":  RSquareBrT,
	"{":  LCurlyBrT,
	"}":  RCurlyBrT,
	".":  DotT,
	"?":  QMarkT,
	"..": EllipsisT,
	",":  CommaT,
	";":  SemicolonT,
	":":  ColonT,
	"<:": SubtypeT,
	"->": ArrowT,
}

func (lexer *Lexer) grabPunctuation() {
	lexer.advance()

	twochars := string(append(lexer.CurrentBlob, lexer.next()))
	if _, ok := punct_combos[twochars]; ok {
		lexer.advance()
	}

	blob := string(lexer.CurrentBlob)

	if ttype, ok := uniquePuncts[blob]; ok {
		lexer.CurrentTType = ttype
	} else if _, ok := binops[blob]; ok {
		lexer.CurrentTType = BinopT
	} else if _, ok := comparisons[blob]; ok {
		lexer.CurrentTType = ComparisonT
	} else if _, ok := definers[blob]; ok {
		lexer.CurrentTType = DefinerT
	} else if _, ok := updaters[blob]; ok {
		lexer.CurrentTType = UpdaterT
	} else if _, ok := prefixes[blob]; ok {
		lexer.CurrentTType = PrefixT
	} else if _, ok := suffixes[blob]; ok {
		lexer.CurrentTType = SuffixT
	} else {
		lexparsePanic(lexer.Line, lexer.Col-1, 1, LexerError, "Invalid start to token")
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
		c := lexer.next()
		lexer.advance()
		return c
	}
	lexer.advance()

	var c rune
	switch lexer.next() {
	case 't':
		c = '\t'
		lexer.advance()
	case 'n':
		c = '\n'
		lexer.advance()
	case 'r':
		c = '\r'
		lexer.advance()
	case '"':
		c = '"'
		lexer.advance()
	case '\'':
		c = '\''
		lexer.advance()
	default:
		panic("Ahh! We still need octal and unicode escapes!!!1!!one!")
	}

	return c
}
