package checker

import (
	"fmt"
	"github.com/cstuartroe/teko/src/lexparse"
)

type FunctionArgDef struct {
	name    string
	mutable bool
	ttype   TekoType
}

type FunctionType struct {
	rtype   TekoType
	argdefs []FunctionArgDef
}

func (ftype FunctionType) tekotype() {}

func (ftype FunctionType) argnames() []string {
	var out []string = []string{}

	for _, argdef := range ftype.argdefs {
		out = append(out, argdef.name)
	}

	return out
}

func contains(strings []string, s string) bool {
	for _, e := range strings {
		if s == e {
			return true
		}
	}
	return false
}

func ResolveArgs(argnames []string, expr lexparse.CallExpression) map[string]lexparse.Expression {
	args_by_name := map[string]lexparse.Expression{}

	if len(expr.Args) > len(argnames) {
		expr.Args[len(argnames)].Token().Raise(
			lexparse.ArgumentError,
			fmt.Sprintf("Too many arguments (%d expected, %d given)", len(argnames), len(expr.Args)),
		)
	}

	for i, arg := range expr.Args {
		args_by_name[argnames[i]] = arg
	}

	for _, kwarg := range expr.Kwargs {
		name := string(kwarg.Symbol.Value)

		if _, ok := args_by_name[name]; ok {
			kwarg.Token().Raise(
				lexparse.ArgumentError,
				fmt.Sprintf("Argument already passed: %s", name),
			)
		}

		if contains(argnames, name) {
			args_by_name[name] = kwarg.Value
		} else {
			kwarg.Symbol.Raise(
				lexparse.ArgumentError,
				fmt.Sprintf("Function doesn't take argument %s", name),
			)
		}
	}

	return args_by_name
}
