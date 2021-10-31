package lexparse

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/cstuartroe/teko/src/shared"
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

func raiseTekoError(line *Line, col int, errorClass shared.TekoErrorClass, message string) {
	fmt.Fprintf(
		shared.ErrorDest,
		"%s (%s:%d:%d)\n%s\n%s%s\n%s\n",
		errorClass,
		line.Filename,
		line.Num+1,
		col+1,
		string(line.Value),
		strings.Repeat(" ", col),
		strings.Repeat("^", 1),
		message,
	)
	panic(shared.TekoErrorMessage)
}

func (token Token) Raise(errorClass shared.TekoErrorClass, message string) {
	raiseTekoError(token.Line, token.Col, errorClass, message)
}

type Lexer struct {
	Lines        []*Line
	LineN        int
	Col          int
	currentToken Token
}

func NewLexer(filename string, contents string) Lexer {
	lines_raw := strings.Split(contents, "\n")
	lines := make([]*Line, len(lines_raw))

	for i, s := range lines_raw {
		lines[i] = &Line{
			Num:      i,
			Value:    []rune(s),
			Filename: filename,
		}
	}

	lexer := Lexer{Lines: lines}
	lexer.passWhitespace()
	return lexer
}

func FromFile(filename string) Lexer {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	contents := string(dat)

	return NewLexer(filename, contents)
}

func (lexer *Lexer) Lex() []Token {
	tokens := []Token{}

	for lexer.HasMore() {
		tokens = append(tokens, lexer.GrabToken())
	}

	return tokens
}

func (lexer *Lexer) HasMore() bool {
	return lexer.LineN < len(lexer.Lines)
}

func (lexer *Lexer) newLine() {
	lexer.LineN += 1
	lexer.Col = 0
}

func (lexer *Lexer) Line() *Line {
	return lexer.Lines[lexer.LineN]
}

func (lexer *Lexer) newToken() {
	lexer.currentToken = Token{
		Line: lexer.Line(),
		Col:  lexer.Col,
	}
}

func (lexer *Lexer) next() rune {
	if !lexer.lineHasMore() {
		return rune(0)
	}
	return lexer.Line().Value[lexer.Col]
}

func (lexer *Lexer) nextFew(n int) string {
	end := lexer.Col + n
	if end > len(lexer.Line().Value) {
		end = len(lexer.Line().Value)
	}

	return string(lexer.Line().Value[lexer.Col:end])
}

func (lexer *Lexer) advance() {
	lexer.currentToken.Value = append(lexer.currentToken.Value, lexer.next())
	lexer.Col++
}

func (lexer *Lexer) lineHasMore() bool {
	return lexer.Col < len(lexer.Line().Value)
}

func (lexer *Lexer) passWhitespace() {
	for unicode.IsSpace(lexer.next()) {
		lexer.advance()
	}
	if !lexer.lineHasMore() {
		lexer.newLine()
		if lexer.HasMore() {
			lexer.passWhitespace()
		}
	}
}

func (lexer *Lexer) GrabToken() Token {
	lexer.newToken()

	c := lexer.next()

	if lexer.nextFew(2) == "//" {
		lexer.grabLineComment()
	} else if lexer.nextFew(2) == "/*" {
		lexer.grabBlockComment()
	} else if unicode.IsLetter(c) || (c == '_') {
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

	t := lexer.currentToken
	lexer.passWhitespace()
	return t
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
	"scope": ScopeT,

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

	if ttype, ok := keywords[string(lexer.currentToken.Value)]; ok {
		lexer.currentToken.TType = ttype
	} else {
		lexer.currentToken.TType = SymbolT
	}
}

func (lexer *Lexer) grabDecimalNumber() {
	for unicode.IsDigit(lexer.next()) {
		lexer.advance()
	}

	if lexer.next() != '.' {
		lexer.currentToken.TType = IntT
	} else {
		lexer.currentToken.TType = FloatT

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
	"=":  EqualT,
	"|>": PipeT,
}

func (lexer *Lexer) grabPunctuation() {
	if _, ok := punct_combos[lexer.nextFew(2)]; ok {
		lexer.advance()
	}

	lexer.advance()

	blob := string(lexer.currentToken.Value)

	if ttype, ok := uniquePuncts[blob]; ok {
		lexer.currentToken.TType = ttype
	} else if _, ok := binops[blob]; ok {
		lexer.currentToken.TType = BinopT
	} else if _, ok := comparisons[blob]; ok {
		lexer.currentToken.TType = ComparisonT
	} else if _, ok := updaters[blob]; ok {
		lexer.currentToken.TType = UpdaterT
	} else if _, ok := prefixes[blob]; ok {
		lexer.currentToken.TType = PrefixT
	} else if _, ok := suffixes[blob]; ok {
		lexer.currentToken.TType = SuffixT
	} else {
		raiseTekoError(lexer.Line(), lexer.Col-1, shared.LexerError, "Invalid start to token")
	}
}

func (lexer *Lexer) grabString() {
	if lexer.next() != '"' {
		panic("String didn't start with double quote")
	}
	lexer.currentToken.TType = StringT
	lexer.advance()

	value := []rune{}

	for lexer.next() != '"' {
		value = append(value, lexer.grabCharacter())
	}

	lexer.advance()

	lexer.currentToken.Value = value // a shameless hack
}

func (lexer *Lexer) grabChar() {
	if lexer.next() != '\'' {
		panic("Char didn't start with single quote")
	}
	lexer.currentToken.TType = CharT
	lexer.advance()

	value := []rune{lexer.grabCharacter()}

	if lexer.next() != '\'' {
		panic("Char too long")
	}
	lexer.advance()

	lexer.currentToken.Value = value // the same shameless hack - look any better this time?
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

func (lexer *Lexer) grabLineComment() {
	lexer.currentToken.TType = LineCommentT
	for lexer.lineHasMore() {
		lexer.advance()
	}
}

func (lexer *Lexer) grabBlockComment() {
	lexer.currentToken.TType = BlockCommentT
	lexer.advance()
	lexer.advance()

	comment_depth := 1

	for comment_depth > 0 {
		if !lexer.lineHasMore() {
			if lexer.LineN == len(lexer.Lines)-1 {
				raiseTekoError(lexer.Line(), lexer.Col, shared.LexerError, "EOF while parsing block comment")
			}

			lexer.newLine()
			lexer.currentToken.Value = append(lexer.currentToken.Value, '\n')

		} else if lexer.nextFew(2) == "*/" {
			lexer.advance()
			lexer.advance()
			comment_depth -= 1

		} else if lexer.nextFew(2) == "/*" {
			lexer.advance()
			lexer.advance()
			comment_depth += 1

		} else {
			lexer.advance()
		}
	}
}
