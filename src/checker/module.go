package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

type TypeModule struct {
	codeblock lexparse.Codeblock
	checker   Checker
}

func CheckTree(codeblock lexparse.Codeblock) TypeModule {
	main_checker := NewChecker(BaseChecker)

	for _, stmt := range codeblock.GetStatements() {
		main_checker.checkStatement(stmt)
	}

	return TypeModule{
		codeblock: codeblock,
		checker:   main_checker,
	}
}
