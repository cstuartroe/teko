package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
)

type executorType func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject

type TekoFunction struct {
	context  *InterpreterModule
	owner    *TekoObject
	body     lexparse.Expression
	argdefs  []checker.FunctionArgDef
	executor executorType
}

func (f TekoFunction) getUnderlyingType() checker.TekoType {
	return nil // TODO
}

func (f TekoFunction) getFieldValue(name string) *TekoObject { return nil }

func (t TekoFunction) execute(callingContext InterpreterModule, call *lexparse.CallExpression) *TekoObject {
	resolvedArgs := checker.ResolveArgs(t.argdefs, call)
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
		scope: &BasicObject{
			symbolTable: SymbolTable{
				parent: &function.context.scope.symbolTable,
				table:  evaluatedArgs,
			},
		},
	}

	codeblock := &lexparse.Codeblock{
		Statements: []lexparse.Statement{
			&lexparse.ExpressionStatement{
				Expression: function.body,
			},
		},
	}

	return interpreter.Execute(codeblock)
}

func customExecutedFunction(executor executorType, argdefs []checker.FunctionArgDef) TekoFunction {
	return TekoFunction{
		context:  nil,
		owner:    nil,
		body:     nil,
		argdefs:  argdefs,
		executor: executor,
	}
}
