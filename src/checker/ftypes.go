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

func (ftype FunctionType) allFields() map[string]TekoType {
	return map[string]TekoType{}
}

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
		lexparse.TokenPanic(
			expr.Args[len(argnames)].Token(),
			fmt.Sprintf("Too many arguments (%d expected, %d given)", len(argnames), len(expr.Args)),
		)
	}

	for i, arg := range expr.Args {
		args_by_name[argnames[i]] = arg
	}

	for _, kwarg := range expr.Kwargs {
		name := string(kwarg.Symbol.Value)

		if _, ok := args_by_name[name]; ok {
			lexparse.TokenPanic(
				kwarg.Token(),
				fmt.Sprintf("Argument already passed: %s", name),
			)
		}

		if contains(argnames, name) {
			args_by_name[name] = kwarg.Value
		} else {
			lexparse.TokenPanic(
				kwarg.Symbol,
				fmt.Sprintf("Function doesn't take argument %s", name),
			)
		}
	}

	return args_by_name
}
