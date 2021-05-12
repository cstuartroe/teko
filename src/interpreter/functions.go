package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
)

type executorType func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject

type TekoFunction struct {
	context  *InterpreterModule
	body     *lexparse.Codeblock
	argnames []string
	executor executorType
}

func (f TekoFunction) getFieldValue(name string) *TekoObject { return nil }

func (t TekoFunction) execute(callingContext InterpreterModule, call lexparse.CallExpression) *TekoObject {
	resolvedArgs := checker.ResolveArgs(t.argnames, call)
	evaluatedArgs := t.evaluateArgs(callingContext, resolvedArgs)
	return t.executor(t, evaluatedArgs)
}

// This will eventually use the function's own interpreter to evaluate default args
func (t TekoFunction) evaluateArgs(callingContext InterpreterModule, resolvedArgs map[string]lexparse.Expression) map[string]*TekoObject {
	out := map[string]*TekoObject{}
	for name, expr := range resolvedArgs {
		// everything should be pass by value, so we need to generate a shallow copy
		out[name] = sc(callingContext.evaluateExpression(expr))
	}
	return out
}

func defaultFunctionExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	interpreter := InterpreterModule{
		codeblock: function.body,
		scope: &BasicObject{
			symbolTable: &SymbolTable{
				parent: function.context.scope.symbolTable,
				table:  evaluatedArgs,
			},
		},
	}

	return interpreter.execute()
}

func customExecutedFunction(executor executorType, argnames []string) TekoFunction {
	return TekoFunction{
		context:  nil,
		body:     nil,
		argnames: argnames,
		executor: executor,
	}
}