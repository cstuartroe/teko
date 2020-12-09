package checker

import (
	"fmt"
	"github.com/cstuartroe/teko/src/lexparse"
)

type FunctionArgDef struct {
	name    string
	mutable bool
	byref   bool
	ttype   TekoType
}

type FunctionType struct {
	rtype   TekoType
	argdefs []FunctionArgDef
}

func (ftype FunctionType) allFields() map[string]TekoType {
	return map[string]TekoType{}
}

func (ftype FunctionType) GetArgdefs() []FunctionArgDef {
	return ftype.argdefs
}

func getArgTypeByName(argdefs []FunctionArgDef, name string) TekoType {
	for _, argdef := range argdefs {
		if argdef.name == name {
			return argdef.ttype
		}
	}
	return nil
}

func ResolveArgs(argdefs []FunctionArgDef, expr lexparse.CallExpression) map[string]lexparse.Expression {
	args_by_name := map[string]lexparse.Expression{}

	if len(expr.Args) > len(argdefs) {
		lexparse.TokenPanic(
			expr.Args[len(argdefs)].Token(),
			fmt.Sprintf("Too many arguments (%d expected, %d given)", len(argdefs), len(expr.Args)),
		)
	}

	for i, arg := range expr.Args {
		args_by_name[argdefs[i].name] = arg
	}

	for _, kwarg := range expr.Kwargs {
		name := string(kwarg.Symbol.Value)

		if _, ok := args_by_name[name]; ok {
			lexparse.TokenPanic(
				kwarg.Token(),
				fmt.Sprintf("Argument already passed: %s", name),
			)
		}

		argtype := getArgTypeByName(argdefs, name)
		if argtype != nil {
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
