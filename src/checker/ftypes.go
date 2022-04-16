package checker

import (
	"fmt"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

type FunctionArgDef struct {
	Name    string
	ttype   TekoType
	Default lexparse.Expression
}

func NoDefaults(argnames ...string) []FunctionArgDef {
	out := []FunctionArgDef{}

	for _, name := range argnames {
		out = append(out, FunctionArgDef{
			Name: name,
		})
	}

	return out
}

type FunctionType struct {
	rtype   TekoType
	argdefs []FunctionArgDef
}

func (ftype FunctionType) tekotypeToString() string {
	out := "fn("
	for _, argdef := range ftype.argdefs {
		out += argdef.Name + ": " + argdef.ttype.tekotypeToString() + ", "
	}

	out += "): "
	out += ftype.rtype.tekotypeToString()

	return out
}

func (ftype FunctionType) allFields() map[string]TekoType {
	return map[string]TekoType{} // TODO should actually be all fields shared by types
}

// func (ftype FunctionType) argnames() []string {
// 	var out []string = []string{}

// 	for _, argdef := range ftype.argdefs {
// 		out = append(out, argdef.name)
// 	}

// 	return out
// }

func containsName(argdefs []FunctionArgDef, name string) bool {
	for _, ad := range argdefs {
		if name == ad.Name {
			return true
		}
	}
	return false
}

func ResolveArgs(argdefs []FunctionArgDef, expr lexparse.CallExpression) map[string]lexparse.Expression {
	args_by_name := map[string]lexparse.Expression{}

	if len(expr.Args) > len(argdefs) {
		expr.Args[len(argdefs)].Token().Raise(
			shared.ArgumentError,
			fmt.Sprintf("Too many arguments (%d expected, %d given)", len(argdefs), len(expr.Args)),
		)
	}

	for i, arg := range expr.Args {
		args_by_name[argdefs[i].Name] = arg
	}

	for _, kwarg := range expr.Kwargs {
		name := string(kwarg.Symbol.Value)

		if _, ok := args_by_name[name]; ok {
			kwarg.Token().Raise(
				shared.ArgumentError,
				fmt.Sprintf("Argument already passed: %s", name),
			)
		}

		if containsName(argdefs, name) {
			args_by_name[name] = kwarg.Value
		} else {
			kwarg.Symbol.Raise(
				shared.ArgumentError,
				fmt.Sprintf("Function doesn't take argument %s", name),
			)
		}
	}

	for _, ad := range argdefs {
		if _, ok := args_by_name[ad.Name]; !ok {
			args_by_name[ad.Name] = ad.Default
		}
	}

	return args_by_name
}
