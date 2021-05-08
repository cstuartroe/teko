package interpreter

import (
	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/lexparse"
	"math"
	"strconv"
)

var INT_CACHE map[int]Integer = map[int]Integer{}

func getInteger(n int) Integer {
	if obj, ok := INT_CACHE[n]; ok {
		return obj
	} else {
		obj = Integer{
			value:       n,
			symbolTable: newSymbolTable(nil),
		}
		INT_CACHE[n] = obj
		return obj
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
			var n TekoObject = getInteger(op(receiverValue, p.value))
			return &n
		default:
			panic("Non-int somehow made it past the type checker as an argument to int arithmetic!")
		}
	}
}

func IntBinop(receiverValue int, op intOpType) TekoFunction {
	return TekoFunction{
		context:  blankInterpreter,
		body:     lexparse.Codeblock{},
		ftype:    checker.IntBinopType,
		executor: IntBinopExecutor(receiverValue, op),
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
		var s TekoObject = String{
			value: []rune(strconv.Itoa(receiverValue)),
		}
		return &s
	}
}

func IntToStr(receiverValue int) TekoFunction {
	return TekoFunction{
		context:  blankInterpreter,
		body:     lexparse.Codeblock{},
		ftype:    checker.ToStrType,
		executor: IntToStrExecutor(receiverValue),
	}
}

func set_and_return(n Integer, name string, f *TekoObject) *TekoObject {
	n.symbolTable.set(name, f)
	return f
}

func (n Integer) getFieldValue(name string) *TekoObject {
	if attr := n.symbolTable.get(name); attr != nil {
		return attr
	} else if op, ok := intOps[name]; ok {
		var f TekoObject = IntBinop(n.value, op)
		return set_and_return(n, name, &f)
	}

	switch name {

	case "to_str":
		var f TekoObject = IntToStr(n.value)
		return set_and_return(n, name, &f)

	default:
		panic("Operation not implemented: " + name)
	}
}
