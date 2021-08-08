package interpreter

type Array struct {
	elements    []*TekoObject
	symbolTable *SymbolTable
}

func ArrayAddExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		switch p := (*evaluatedArgs["other"]).(type) {
		case Array:
			return tp(Array{
				append(append([]*TekoObject{}, receiverElements...), p.elements...),
				newSymbolTable(nil),
			})
		default:
			panic("Non-array somehow made it past the type checker as an argument to array add!")
		}
	}
}

func ArrayAtExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		switch p := (*evaluatedArgs["key"]).(type) {
		case Integer:
			return receiverElements[p.value]
		default:
			panic("Non-integer somehow made it past the type checker as an argument to array at!")
		}
	}
}

func ArrayIncludesExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		for _, e := range receiverElements {
			if e == evaluatedArgs["element"] {
				return True
			}
		}
		return False
	}
}

func runeJoin(runeSlices [][]rune, joiner []rune) []rune {
	out := []rune{}

	for i, r := range runeSlices {
		if i > 0 {
			out = append(out, joiner...)
		}

		out = append(out, r...)
	}

	return out
}

func ArrayToStrExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		elementStrings := [][]rune{}

		for _, e := range receiverElements {
			f := (*e).getFieldValue("to_str")

			switch fp := (*f).(type) {
			case TekoFunction:
				s := fp.executor(fp, map[string]*TekoObject{})

				switch sp := (*s).(type) {
				case String:
					elementStrings = append(elementStrings, sp.value)

				default:
					panic("to_str did not return a string")
				}

			default:
				panic("to_str was not a function")
			}
		}

		return tp(String{
			value: append([]rune{'['}, append(runeJoin(elementStrings, []rune{',', ' '}), ']')...),
		})
	}
}

// TODO get argames from checker
func (a Array) getFieldValue(name string) *TekoObject {
	return cached_get(a.symbolTable, name, func() *TekoObject {
		switch name {
		case "add":
			return tp(customExecutedFunction(ArrayAddExecutor(a.elements), []string{"other"}))

		case "at":
			return tp(customExecutedFunction(ArrayAtExecutor(a.elements), []string{"key"}))

		case "size":
			return tp(getInteger(len(a.elements)))

		case "includes":
			return tp(customExecutedFunction(ArrayIncludesExecutor(a.elements), []string{"element"}))

		case "to_str":
			return tp(customExecutedFunction(ArrayToStrExecutor(a.elements), []string{}))

		default:
			panic("Unknown array function")
		}
	})
}

type Set struct {
	elements []*TekoObject
}

func (s Set) getFieldValue(name string) *TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for sets yet")
	}
}
