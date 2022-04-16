package interpreter

import (
	"fmt"
	"strconv"

	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

type InterpreterModule struct {
	scope *BasicObject
}

func New(parent *InterpreterModule) InterpreterModule {
	return InterpreterModule{
		scope: &BasicObject{
			newSymbolTable(&parent.scope.symbolTable),
		},
	}
}

func (m InterpreterModule) Execute(codeblock *lexparse.Codeblock) *TekoObject {
	var out *TekoObject = nil
	for _, stmt := range codeblock.Statements {
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
		stmt.Token().Raise(shared.NotImplementedError, "Statement type not implemented")
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

	case lexparse.MapExpression:
		return m.evaluateMapExpression(p)

	case lexparse.ObjectExpression:
		return m.evaluateObjectExpression(p)

	case lexparse.FunctionExpression:
		return m.evaluateFunctionDefinition(p)

	case lexparse.DoExpression:
		return m.evaluateDoExpression(p)

	case lexparse.VarExpression:
		return m.evaluateExpression(p.Right)

	case lexparse.WhileExpression:
		return m.evaluateWhile(p)

	case lexparse.ForExpression:
		return m.evaluateFor(p)

	case lexparse.ScopeExpression:
		return m.evaluateScope(p)

	case lexparse.ComparisonExpression:
		return m.evaluateComparisonExpression(p)

	default:
		expr.Token().Raise(shared.NotImplementedError, "Intepretation of expression type not implemented")
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
			expr.Token().Raise(shared.UnexpectedIssue, "Label not found")
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
			expr.Token().Raise(shared.UnexpectedIssue, "Invalid integer - how did this make it past the lexer?")
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
		expr.Token().Raise(shared.NotImplementedError, fmt.Sprintf("Invalid or unimplemented simple expression type: %s", ttype))
		return nil
	}
}

func (m InterpreterModule) evaluateFunctionCall(call lexparse.CallExpression) *TekoObject {
	receiver := m.evaluateExpression(call.Receiver)
	switch p := (*receiver).(type) {
	case TekoFunction:
		return p.execute(m, call)

	default:
		call.Token().Raise(shared.UnexpectedIssue, "Non-function was the receiver of a call. Where was the type checker??")
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
			owner:   left,
			body:    nil,
			argdefs: []checker.FunctionArgDef{
				{
					Name:    "value",
					Default: nil,
				},
			},
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

func (m InterpreterModule) evaluateMapExpression(expr lexparse.MapExpression) *TekoObject {
	tmap := newTekoMap()

	for _, kvpair := range expr.KVPairs {
		tmap.set(m.evaluateExpression(kvpair.Key), m.evaluateExpression(kvpair.Value))
	}

	return tp(tmap)
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
	argdefs := []checker.FunctionArgDef{}
	for _, ad := range expr.Argdefs {
		ad := checker.FunctionArgDef{
			Name:    string(ad.Symbol.Value),
			Default: ad.Default,
		}
		argdefs = append(argdefs, ad)
	}

	f := tp(TekoFunction{
		context:  m,
		body:     expr.Right,
		argdefs:  argdefs,
		executor: defaultFunctionExecutor,
	})

	if expr.Name != nil {
		m.scope.symbolTable.set(string(expr.Name.Value), f)
	}

	return f
}

func (m *InterpreterModule) evaluateDoExpression(expr lexparse.DoExpression) *TekoObject {
	return New(m).Execute(&expr.Codeblock)
}

func (m *InterpreterModule) isTrue(expr lexparse.Expression) bool {
	switch p := (*m.evaluateExpression(expr)).(type) {
	case Boolean:
		return p.value
	default:
		expr.Token().Raise(shared.UnexpectedIssue, "Not a boolean")
		return false
	}
}

func (m InterpreterModule) evaluateWhile(expr lexparse.WhileExpression) *TekoObject {
	elements := []*TekoObject{}

	for m.isTrue(expr.Condition) {
		elements = append(elements, m.evaluateExpression(expr.Body))
	}

	return tp(newArray(elements))
}

func (m *InterpreterModule) evaluateFor(expr lexparse.ForExpression) *TekoObject {
	elements := []*TekoObject{}

	iterator := *m.evaluateExpression(expr.Iterator)

	var size int
	sizeObj := *iterator.getFieldValue("size")
	switch p := sizeObj.(type) {
	case Integer:
		size = p.value
	default:
		expr.Iterator.Token().Raise(shared.UnexpectedIssue, "Did not evaluate to an int object")
	}

	var at TekoFunction
	atObj := *iterator.getFieldValue("at")

	switch fp := atObj.(type) {
	case TekoFunction:
		at = fp
	default:
		panic("Iterator does not have at function")
	}

	for i := 0; i < size; i++ {
		el := at.executor(at, map[string]*TekoObject{at.argdefs[0].Name: tp(getInteger(i))})

		blockInterpreter := New(m)

		blockInterpreter.scope.symbolTable.set(string(expr.Iterand.Value), el)

		elements = append(elements, blockInterpreter.evaluateExpression(expr.Body))
	}

	return tp(newArray(elements))
}

func (m *InterpreterModule) evaluateScope(expr lexparse.ScopeExpression) *TekoObject {
	interpreter := New(m)
	interpreter.Execute(&expr.Codeblock)
	return tp(*interpreter.scope)
}

var comparisons map[string]func(int) bool = map[string]func(int) bool{
	"==": func(n int) bool { return n == 0 },
	"!=": func(n int) bool { return n != 0 },
	"<":  func(n int) bool { return n < 0 },
	"<=": func(n int) bool { return n <= 0 },
	">":  func(n int) bool { return n > 0 },
	">=": func(n int) bool { return n >= 0 },
}

func (m *InterpreterModule) evaluateComparisonExpression(expr lexparse.ComparisonExpression) *TekoObject {
	comparison_int := m.evaluateFunctionCall(lexparse.ComparisonCallExpression(expr))

	switch p := (*comparison_int).(type) {
	case Integer:
		f := comparisons[string(expr.Comparator.Value)]

		if f(p.value) {
			return True
		} else {
			return False
		}
	default:
		panic("Not an int. What trickery is this?")
	}
}
