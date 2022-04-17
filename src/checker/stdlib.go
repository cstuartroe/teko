package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

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
	"print": PrintType,
	"map":   arrayMapType,
	"null":  NullType,
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
	process_type := &BasicType{
		name: "",
		fields: map[string]TekoType{
			"args": emptyChecker.newArrayType(StringType, hashable_sequence),
		},
	}

	BaseChecker.ctype.setField("process", process_type)

	BaseChecker.CheckTree(lexparse.StdLibParser.Codeblock)
}
