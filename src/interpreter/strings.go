package interpreter

import (
	"math/big"

	"github.com/cstuartroe/teko/src/checker"
)

type String struct {
	runes       []rune
	symbolTable SymbolTable
}

func (s String) getUnderlyingType() checker.TekoType {
	return checker.NewConstantStringType(s.runes)
}

func newTekoString(runes []rune) String {
	return String{
		runes:       runes,
		symbolTable: newSymbolTable(nil),
	}
}

func StringAddExecutor(receiverRunes []rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		switch p := (*evaluatedArgs["other"]).(type) {
		case String:
			return tp(newTekoString(
				append(append([]rune{}, receiverRunes...), p.runes...),
			))
		default:
			panic("Non-string somehow made it past the type checker as an argument to string add!")
		}
	}
}

func StringAtExecutor(receiverRunes []rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		switch p := (*evaluatedArgs["key"]).(type) {
		case Integer:
			return tp(newChar(receiverRunes[p.value.Int64()]))
		default:
			panic("Non-integer somehow made it past the type checker as an argument to string at!")
		}
	}
}

func StringToStrExecutor(receiverRunes []rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		runes := []rune{'"'}

		for _, c := range receiverRunes {
			runes = append(runes, []rune(tekoEscape(c, true, false))...)
		}

		runes = append(runes, '"')

		return tp(newTekoString(runes))
	}
}

func intExp(base int64, exp int64) int64 {
	var out int64 = 1

	for exp > 0 {
		if exp%2 == 0 {
			base = base * base
			exp = exp / 2
		} else {
			out = out * base
			exp = exp - 1
		}
	}

	return out
}

const GENERATOR int64 = 30922493929929101
const OFFSET int64 = 7920982390293083290

// This is not a cryptographic hash. It would probably be pretty straightforward to find a collision if you wanted to,
// but I tried 26^6 random strings and didn't find one so it should be fine for anyone not actively trying to break it.
func hash(runes []rune) int64 {
	var n int64 = OFFSET

	for i, c := range runes {
		n += intExp(GENERATOR, int64(i)+1) * int64(c)
	}

	return n
}

func StringHashExecutor(receiverRunes []rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		// TODO on 32-bit systems this cast to int will lead to a possibility of collision
		// but I won't bother converting everything to int64s because I should ultimately use variable-width integers anyway
		return tp(getInteger(big.NewInt(hash(receiverRunes))))
	}
}

// TODO get argames from checker
func (s String) getFieldValue(name string) *TekoObject {
	return s.symbolTable.cached_get(name, func() *TekoObject {
		switch name {
		case "add":
			return tp(customExecutedFunction(StringAddExecutor(s.runes), checker.NoDefaults("other")))

		case "at":
			return tp(customExecutedFunction(StringAtExecutor(s.runes), checker.NoDefaults("key")))

		case "size":
			return tp(int2Integer(len(s.runes)))

		case "to_str":
			return tp(customExecutedFunction(StringToStrExecutor(s.runes), checker.NoDefaults()))

		case "hash":
			return tp(customExecutedFunction(StringHashExecutor(s.runes), checker.NoDefaults()))

		default:
			panic("Unknown string function: " + name)
		}
	})
}
