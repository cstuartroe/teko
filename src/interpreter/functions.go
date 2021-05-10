package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
)

type executorType func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject

type TekoFunction struct {
	context  InterpreterModule
	body     lexparse.Codeblock
	ftype    *checker.FunctionType
	executor executorType
}

func (f TekoFunction) getFieldValue(name string) *TekoObject { return nil }

func (t TekoFunction) execute(callingContext InterpreterModule, resolvedArgs map[string]lexparse.Expression) *TekoObject {
	evaluatedArgs := t.evaluateArgs(callingContext, resolvedArgs)
	return t.executor(t, evaluatedArgs)
}

// This will eventually use the function's own interpreter to evaluate default args
func (t TekoFunction) evaluateArgs(callingContext InterpreterModule, resolvedArgs map[string]lexparse.Expression) map[string]*TekoObject {
	out := map[string]*TekoObject{}
	for name, expr := range resolvedArgs {
		out[name] = callingContext.evaluateExpression(expr)
	}
	return out
}

func defaultFunctionExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	interpreter := InterpreterModule{
		codeblock: function.body,
		scope: BasicObject{
			symbolTable: SymbolTable{
				parent: &function.context.scope.symbolTable,
				table:  evaluatedArgs,
			},
		},
	}

	return interpreter.execute()
}
