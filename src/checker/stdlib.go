package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

var ProcessType TekoType = &BasicType{
	name: "",
	fields: map[string]TekoType{
		"args": newArrayType(StringType),
	},
}

var stdlibTypeTable *TypeTable = &TypeTable{
	parent: nil,
	table: map[string]TekoType{
		"int":      IntType,
		"bool":     BoolType,
		"str":      StringType,
		"char":     CharType,
		"null":     NullType,
		"Hashable": Hashable,
	},
}

var StdlibSymbolTable map[string]TekoType = map[string]TekoType{
	"print":   PrintType,
	"process": ProcessType,
	"null":    NullType,
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
