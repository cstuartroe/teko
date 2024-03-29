package lexparse

import "github.com/cstuartroe/teko/src/shared"

var StdLibParser Parser = Parser{}

func ParseStldLib() {
	lexer := NewLexer("lib/stdlib.to", shared.StdLibContents)
	tokens := lexer.Lex()
	StdLibParser.Parse(tokens, true)
}
