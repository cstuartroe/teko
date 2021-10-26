
package interpreter

import (
	"math"
	"strconv"
)

func getInteger(n int) Integer {
	return Integer{
		value:       n,
		symbolTable: newSymbolTable(nil),
	}
}

type intOpType func(n1 int, n2 int) int

func IntBinopExecutor(receiverValue int, op intOpType) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		other, ok := evaluatedArgs["other"]
		if !ok {
			panic("No parameter passed to int arithmetic function")
		}

		switch p := (*other).(type) {
		case Integer:
			return tp(getInteger(op(receiverValue, p.value)))
		default:
			panic("Non-int somehow made it past the type checker as an argument to int arithmetic!")
		}
	}
}

var intOps map[string]intOpType = map[string]intOpType{
	"add":  func(n1 int, n2 int) int { return n1 + n2 },
	"sub":  func(n1 int, n2 int) int { return n1 - n2 },
	"mult": func(n1 int, n2 int) int { return n1 * n2 },
	"div":  func(n1 int, n2 int) int { return n1 / n2 },
	"exp":  func(n1 int, n2 int) int { return int(math.Pow(float64(n1), float64(n2))) },
}

func IntToStrExecutor(receiverValue int) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return tp(newTekoString([]rune(strconv.Itoa(receiverValue))))
	}
}

func (n Integer) getFieldValue(name string) *TekoObject {
	return n.symbolTable.cached_get(name, func() *TekoObject {
		if op, ok := intOps[name]; ok {
			return tp(customExecutedFunction(IntBinopExecutor(n.value, op), []string{"other"}))
		}

		switch name {

		case "to_str":
			return tp(customExecutedFunction(IntToStrExecutor(n.value), []string{}))

		default:
			panic("Operation not implemented: " + name)
		}
	})
}
