package interpreter

import "github.com/cstuartroe/teko/src/checker"

type Boolean struct {
	value       bool
	symbolTable SymbolTable
}

func (b Boolean) getUnderlyingType() checker.TekoType {
	return checker.BoolType
}

var True *TekoObject = tp(Boolean{
	value:       true,
	symbolTable: newSymbolTable(nil),
})

var False *TekoObject = tp(Boolean{
	value:       false,
	symbolTable: newSymbolTable(nil),
})

func getBool(b bool) *TekoObject {
	if b {
		return True
	} else {
		return False
	}
}

type boolOpType func(n1 bool, n2 bool) bool

var boolOps map[string]boolOpType = map[string]boolOpType{
	"and": func(b1 bool, b2 bool) bool { return b1 && b2 },
	"or":  func(b1 bool, b2 bool) bool { return b1 || b2 },
}

func BoolBinopExecutor(receiverValue bool, op boolOpType) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		other, ok := evaluatedArgs["other"]
		if !ok {
			panic("No parameter passed to int arithmetic function")
		}

		return getBool(op(receiverValue, (*other).(Boolean).value))
	}
}

func NotExecutor(receiverValue bool) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return getBool(!receiverValue)
	}
}

func BoolToStrExecutor(receiverValue bool) executorType {
	var s string = "false"

	if receiverValue {
		s = "true"
	}

	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return tp(newTekoString([]rune(s)))
	}
}

func (b Boolean) getFieldValue(name string) *TekoObject {
	return b.symbolTable.cached_get(name, func() *TekoObject {
		if op, ok := boolOps[name]; ok {
			return tp(customExecutedFunction(BoolBinopExecutor(b.value, op), checker.NoDefaults("other")))
		}

		switch name {

		case "not":
			return tp(customExecutedFunction(NotExecutor(b.value), checker.NoDefaults()))

		case "to_str":
			return tp(customExecutedFunction(BoolToStrExecutor(b.value), checker.NoDefaults()))

		default:
			panic("Operation not implemented: " + name)
		}
	})
}
