package interpreter

import (
	"math/big"

	"github.com/cstuartroe/teko/src/checker"
)

type Integer struct {
	value       *big.Int // TODO: arbitrary-precision integer
	symbolTable SymbolTable
}

func (i Integer) getUnderlyingType() checker.TekoType {
	return checker.NewConstantIntType(i.value)
}

func int2Integer(n int) Integer {
	return getInteger(big.NewInt(int64(n)))
}

func getInteger(n *big.Int) Integer {
	return Integer{
		value:       n,
		symbolTable: newSymbolTable(nil),
	}
}

type intOpType func(n1 *big.Int, n2 *big.Int) *big.Int

func IntBinopExecutor(receiverValue *big.Int, op intOpType) executorType {
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
	"add":     func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Add(n1, n2) },
	"sub":     func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Sub(n1, n2) },
	"mult":    func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Mul(n1, n2) },
	"div":     func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Div(n1, n2) },
	"exp":     func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Exp(n1, n2, nil) },
	"mod":     func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Mod(n1, n2) },
	"compare": func(n1 *big.Int, n2 *big.Int) *big.Int { return big.NewInt(0).Sub(n1, n2) },
}

func IntToStrExecutor(receiverValue *big.Int) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return tp(newTekoString([]rune(receiverValue.String())))
	}
}

func IntHashExecutor(receiver Integer) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return tp(receiver)
	}
}

func (n Integer) getFieldValue(name string) *TekoObject {
	return n.symbolTable.cached_get(name, func() *TekoObject {
		if op, ok := intOps[name]; ok {
			return tp(customExecutedFunction(IntBinopExecutor(n.value, op), checker.NoDefaults("other")))
		}

		switch name {

		case "to_str":
			return tp(customExecutedFunction(IntToStrExecutor(n.value), checker.NoDefaults()))

		case "hash":
			return tp(customExecutedFunction(IntHashExecutor(n), checker.NoDefaults()))

		default:
			panic("Operation not implemented: " + name)
		}
	})
}
