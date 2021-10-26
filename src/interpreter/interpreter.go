package interpreter

import (
	"fmt"
	"github.com/cstuartroe/teko/src/lexparse"
	"strconv"
)

func ExecuteTree(codeblock *lexparse.Codeblock) {
	interpreter := InterpreterModule{
		codeblock: codeblock,
		scope: &BasicObject{
			newSymbolTable(BaseSymbolTable),
		},
	}

	interpreter.execute()
}

type InterpreterModule struct {
	codeblock *lexparse.Codeblock
	scope     *BasicObject
}

func (m InterpreterModule) execute() *TekoObject {
	var out *TekoObject = nil
	for _, stmt := range m.codeblock.Statements {
		out = m.executeStatement(stmt)
	}
	return out
}

func (m InterpreterModule) executeStatement(stmt lexparse.Statement) *TekoObject {
	switch p := stmt.(type) {
	case lexparse.ExpressionStatement:
		return m.evaluateExpression(p.Expression)

	case lexparse.TypeStatement:
		return nil // it's only for the type checker

	default:
		stmt.Token().Raise(lexparse.NotImplementedError, "Statement type not implemented")
		return nil
	}
}

func (m InterpreterModule) evaluateExpression(expr lexparse.Expression) *TekoObject {
	switch p := expr.(type) {

	case lexparse.SimpleExpression:
		return m.evaluateSimpleExpression(p)

	case lexparse.DeclarationExpression:
		return m.evaluateDeclaration(p)

	case lexparse.CallExpression:
		return m.evaluateFunctionCall(p)

	case lexparse.AttributeExpression:
		return m.evaluateAttributeExpression(p)

	case lexparse.IfExpression:
		return m.evaluateIfExpression(p)

	case lexparse.SequenceExpression:
		return m.evaluateSequenceExpression(p)

	case lexparse.ObjectExpression:
		return m.evaluateObjectExpression(p)

	case lexparse.FunctionExpression:
		return m.evaluateFunctionDefinition(p)

	case lexparse.DoExpression:
		return m.evaluateDoExpression(p)

	case lexparse.VarExpression:
		return m.evaluateExpression(p.Right)

	default:
		expr.Token().Raise(lexparse.NotImplementedError, "Intepretation of expression type not implemented")
		return nil
	}
}

func (m InterpreterModule) evaluateSimpleExpression(expr lexparse.SimpleExpression) *TekoObject {
	ttype := expr.Token().TType
	value := expr.Token().Value

	switch ttype {
	case lexparse.SymbolT:
		val := m.scope.getFieldValue(string(value))
		if val != nil {
			return val
		} else {
			expr.Token().Raise(lexparse.UnexpectedIssue, "Label not found")
			return nil
		}

	case lexparse.StringT:
		return tp(newTekoString(value))

	case lexparse.CharT:
		return nil // TODO

	case lexparse.IntT:
		n, ok := strconv.Atoi(string(value))
		if ok == nil {
			return tp(getInteger(n))
		} else {
			expr.Token().Raise(lexparse.UnexpectedIssue, "Invalid integer - how did this make it past the lexer?")
			return nil
		}

	case lexparse.FloatT:
		return nil // TODO

	case lexparse.BoolT:
		if string(value) == "true" {
			return True
		} else {
			return False
		}

	default:
		expr.Token().Raise(lexparse.NotImplementedError, fmt.Sprintf("Invalid or unimplemented simple expression type: %s", ttype))
		return nil
	}
}

func (m InterpreterModule) evaluateFunctionCall(call lexparse.CallExpression) *TekoObject {
	receiver := m.evaluateExpression(call.Receiver)
	switch p := (*receiver).(type) {
	case TekoFunction:
		return p.execute(m, call)

	default:
		call.Token().Raise(lexparse.UnexpectedIssue, "Non-function was the receiver of a call. Where was the type checker??")
		return nil
	}
}

func (m InterpreterModule) evaluateDeclaration(decl lexparse.DeclarationExpression) *TekoObject {
	name := string(decl.Symbol.Value)
	val := *m.evaluateExpression(decl.Right)
	m.scope.symbolTable.set(name, &val)
	return &val
}

func (m InterpreterModule) evaluateAttributeExpression(expr lexparse.AttributeExpression) *TekoObject {
	left := m.evaluateExpression(expr.Left)

	if string(expr.Symbol.Value) == "=" {
		return tp(TekoFunction{
			context: nil,
			owner: left,
			body: nil,
			argnames: []string{"value"},
			executor: updateExecutor,
		})
	} else {
		return (*left).getFieldValue(string(expr.Symbol.Value))
	}
}

func (m InterpreterModule) evaluateIfExpression(expr lexparse.IfExpression) *TekoObject {
	cond := m.evaluateExpression(expr.Condition)

	var cond_value bool
	switch p := (*cond).(type) {
	case Boolean:
		cond_value = p.value
	default:
		panic("How did a non-boolean slip in here?")
	}

	// It's lazy!
	if cond_value {
		return m.evaluateExpression(expr.Then)
	} else {
		return m.evaluateExpression(expr.Else)
	}
}

func (m InterpreterModule) evaluateSequenceExpression(expr lexparse.SequenceExpression) *TekoObject {
	switch expr.Stype {
	case lexparse.ArraySeqType:
		return tp(m.evaluateArray(expr))
	case lexparse.SetSeqType:
		return tp(m.evaluateSet(expr))
	default:
		panic("Unknown sequence type to interpret: " + expr.Stype)
	}
}

func (m InterpreterModule) evaluateArray(expr lexparse.SequenceExpression) Array {
	elements := []*TekoObject{}
	for _, e := range expr.Elements {
		o := m.evaluateExpression(e)
		elements = append(elements, o)
	}
	return newArray(elements)
}

func (m InterpreterModule) evaluateSet(expr lexparse.SequenceExpression) Set {
	elements := []*TekoObject{}
	for _, e := range expr.Elements {
		o := m.evaluateExpression(e)
		elements = append(elements, o)
	}
	return Set{elements}
}

func (m InterpreterModule) evaluateObjectExpression(expr lexparse.ObjectExpression) *TekoObject {
	symbolTable := newSymbolTable(nil)
	for _, of := range expr.Fields {
		o := m.evaluateExpression(of.Value)
		symbolTable.set(string(of.Symbol.Value), o)
	}
	return tp(BasicObject{symbolTable})
}

func (m *InterpreterModule) evaluateFunctionDefinition(expr lexparse.FunctionExpression) *TekoObject {
	argnames := []string{}
	for _, ad := range expr.Argdefs {
		argnames = append(argnames, string(ad.Symbol.Value))
	}

	f := tp(TekoFunction{
		context:  m,
		body:     expr.Right,
		argnames: argnames,
		executor: defaultFunctionExecutor,
	})

	if expr.Name != nil {
		m.scope.symbolTable.set(string(expr.Name.Value), f)
	}

	return f
}

func (m *InterpreterModule) evaluateDoExpression(expr lexparse.DoExpression) *TekoObject {
	interpreter := InterpreterModule{
		codeblock: &expr.Codeblock,
		scope: &BasicObject{
			newSymbolTable(&m.scope.symbolTable),
		},
	}

	return interpreter.execute()
}
