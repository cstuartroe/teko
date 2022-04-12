package interpreter

type String struct {
	runes       []rune
	symbolTable SymbolTable
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
			return tp(newChar(receiverRunes[p.value]))
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

const GENERATOR uint64 = 30922101

func intExp(base uint64, exp uint64) uint64 {
	var out uint64 = 1

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

func StringHashExecutor(receiverRunes []rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		var n uint64 = 0

		for i, c := range receiverRunes {
			n += intExp(GENERATOR, uint64(i)) * uint64(int(c))
		}

		return tp(getInteger(int(n)))
	}
}

// TODO get argames from checker
func (s String) getFieldValue(name string) *TekoObject {
	return s.symbolTable.cached_get(name, func() *TekoObject {
		switch name {
		case "add":
			return tp(customExecutedFunction(StringAddExecutor(s.runes), []string{"other"}))

		case "at":
			return tp(customExecutedFunction(StringAtExecutor(s.runes), []string{"key"}))

		case "size":
			return tp(getInteger(len(s.runes)))

		case "to_str":
			return tp(customExecutedFunction(StringToStrExecutor(s.runes), []string{}))

		case "hash":
			return tp(customExecutedFunction(StringHashExecutor(s.runes), []string{}))

		default:
			panic("Unknown string function: " + name)
		}
	})
}
