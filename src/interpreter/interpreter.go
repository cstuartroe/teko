package interpreter

import (
	"fmt"
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
	"strconv"
)

func ExecuteTree(codeblock lexparse.Codeblock) {
	interpreter := InterpreterModule{
		codeblock: codeblock,
		scope: BasicObject{
			newSymbolTable(&BaseSymbolTable),
		},
	}

	interpreter.execute()
}

type InterpreterModule struct {
	codeblock lexparse.Codeblock
	scope     BasicObject
}

func (m InterpreterModule) execute() TekoObject {
	for _, stmt := range m.codeblock.GetStatements() {
		m.executeStatement(stmt)
	}
	return nil
}

func (m InterpreterModule) executeStatement(stmt lexparse.Statement) {
	switch p := stmt.(type) {
	case lexparse.ExpressionStatement:
		m.evaluateExpression(p.Expression)
	default:
		lexparse.TokenPanic(stmt.Token(), "Statement type not implemented")
	}
}

func (m InterpreterModule) evaluateExpression(expr lexparse.Expression) TekoObject {
	switch p := expr.(type) {

	case lexparse.SimpleExpression:
		return m.evaluateSimpleExpression(p)

	case lexparse.DeclarationExpression:
		return m.evaluateDeclaration(p)

	case lexparse.CallExpression:
		return m.evaluateFunctionCall(p)

	case lexparse.AttributeExpression:
		return m.evaluateAttributeExpression(p)

	default:
		lexparse.TokenPanic(expr.Token(), "Intepretation of expression type not implemented: "+expr.Ntype())
		return nil
	}
}

func (m InterpreterModule) evaluateSimpleExpression(expr lexparse.SimpleExpression) TekoObject {
	ttype := expr.Token().TType
	value := expr.Token().Value

	switch ttype {
	case lexparse.SymbolT:
		val := m.scope.getFieldValue(string(value))
		if val != nil {
			return val
		} else {
			lexparse.TokenPanic(expr.Token(), "Label not found")
		}

	case lexparse.StringT:
		return String{value}

	case lexparse.CharT:
		return nil // TODO

	case lexparse.IntT:
		n, ok := strconv.Atoi(string(value))
		if ok == nil {
			return getInteger(n)
		} else {
			lexparse.TokenPanic(expr.Token(), "Invalid integer - how did this make it past the lexer?")
		}

	case lexparse.FloatT:
		return nil // TODO

	case lexparse.BoolT:
		b := false
		if string(value) == "true" {
			b = true
		}
		return Boolean{b}

	default:
		lexparse.TokenPanic(expr.Token(), fmt.Sprintf("Invalid or unimplemented simple expression type: %s", ttype))
	}

	return nil
}

func (m InterpreterModule) evaluateFunctionCall(call lexparse.CallExpression) TekoObject {
	receiver := m.evaluateExpression(call.Receiver)
	switch p := receiver.(type) {
	case TekoFunction:
		resolved_args := checker.ResolveArgs(p.ftype.GetArgdefs(), call)
		return p.execute(m, resolved_args)

	default:
		lexparse.TokenPanic(call.Token(), "Non-function was the receiver of a call. Where was the type checker??")
		return nil
	}
}

func (m InterpreterModule) evaluateDeclaration(decl lexparse.DeclarationExpression) TekoObject {
	for _, declared := range decl.Declareds {
		name := string(declared.Symbol.Value)
		m.scope.symbolTable.set(name, m.evaluateExpression(declared.Right))
	}

	return nil // TODO: tuple?
}

func (m InterpreterModule) evaluateAttributeExpression(expr lexparse.AttributeExpression) TekoObject {
	left := m.evaluateExpression(expr.Left)
	return left.getFieldValue(string(expr.Symbol.Value))
}
