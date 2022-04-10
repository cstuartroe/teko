package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

var stdlibTypeTable *TypeTable = &TypeTable{
	parent: nil,
	table: map[string]TekoType{
		"int":  IntType,
		"bool": BoolType,
		"str":  StringType,
		"char": CharType,
	},
}

var StdlibSymbolTable map[string]TekoType = map[string]TekoType{
	"print": PrintType,
	"map":   arrayMapType,
}

var stdLibType *CheckerType = &CheckerType{
	fields: StdlibSymbolTable,
	parent: nil,
}

var BaseChecker *Checker = &Checker{
	typeTable: stdlibTypeTable,
	ctype:     stdLibType,
}

func CheckTekoLangStdLib() {
	BaseChecker.CheckTree(lexparse.StdLibParser.Codeblock)
}
