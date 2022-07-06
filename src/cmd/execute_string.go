package cmd

import "github.com/cstuartroe/teko/src/lexparse"

func ExecuteString(contents string) {
	p := lexparse.Parser{}
	p.ParseString("<stdin>", contents, true)

	ExecuteParser(p)
}
