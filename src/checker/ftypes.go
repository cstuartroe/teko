package checker

import (
	"fmt"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

type FunctionArgDef struct {
	Name  string
	ttype TekoType
	// TODO the typechecker doesn't really need to know the default value,
	// it just needs to know whether there is one. Change this to a boolean
	// and get the interpreter module its own arg def struct that doesn't have
	// a ttype field.
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

func (ftype FunctionType) isDeferred() bool {
	// TODO function types should be deferrable
	return false
}

func (ftype FunctionType) tekotypeToString() string {
	out := "fn("
	for i, argdef := range ftype.argdefs {
		if i > 0 {
			out += ", "
		}
		out += argdef.Name + ": " + argdef.ttype.tekotypeToString()
	}

	out += "): "
	out += ftype.rtype.tekotypeToString()

	return out
}

func (ftype FunctionType) allFields() map[string]TekoType {
	return map[string]TekoType{} // TODO should actually be all fields shared by types
}

func containsName(argdefs []FunctionArgDef, name string) bool {
	for _, ad := range argdefs {
		if name == ad.Name {
			return true
		}
	}
	return false
}

func ResolveArgs(argdefs []FunctionArgDef, expr *lexparse.CallExpression) map[string]lexparse.Expression {
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
