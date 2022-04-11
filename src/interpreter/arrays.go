package interpreter

type Array struct {
	elements    []*TekoObject
	symbolTable SymbolTable
}

func ArrayAddExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		switch p := (*evaluatedArgs["other"]).(type) {
		case Array:
			return tp(newArray(
				append(append([]*TekoObject{}, receiverElements...), p.elements...),
			))
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

func join(slices [][]*TekoObject, joiner []*TekoObject) []*TekoObject {
	out := []*TekoObject{}

	for i, r := range slices {
		if i > 0 {
			out = append(out, joiner...)
		}

		out = append(out, r...)
	}

	return out
}

func ArrayToStrExecutor(receiverElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		elements := [][]*TekoObject{}

		for _, e := range receiverElements {
			f := (*e).getFieldValue("to_str")

			switch fp := (*f).(type) {
			case TekoFunction:
				s := fp.executor(fp, map[string]*TekoObject{})

				switch sp := (*s).(type) {
				case Array:
					elements = append(elements, sp.elements)

				default:
					panic("to_str did not return an array")
				}

			default:
				panic("to_str was not a function")
			}
		}

		return tp(newArray(
			append(
				[]*TekoObject{tp(newChar('['))},
				append(
					join(
						elements,
						[]*TekoObject{
							tp(newChar(',')),
							tp(newChar(' ')),
						},
					),
					tp(newChar(']')),
				)...,
			),
		))
	}
}

// TODO get argames from checker
func (a Array) getFieldValue(name string) *TekoObject {
	return a.symbolTable.cached_get(name, func() *TekoObject {
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

func newArray(elements []*TekoObject) Array {
	return Array{elements, newSymbolTable(nil)}
}

func ArrayMapExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	l, ok := evaluatedArgs["l"]
	if !ok {
		panic("No array parameter passed to map")
	}

	switch lp := (*l).(type) {
	case Array:
		f, ok2 := evaluatedArgs["f"]
		if !ok2 {
			panic("No function parameter passed to map")
		}

		switch fp := (*f).(type) {
		case TekoFunction:
			elements := []*TekoObject{}

			for _, e := range lp.elements {
				elements = append(elements, fp.executor(fp, map[string]*TekoObject{"e": e}))
			}

			return tp(newArray(elements))

		default:
			panic("Non-function made it past the type checker as an argument to map!")
		}
	default:
		panic("Non-array somehow made it past the type checker as an argument to map!")
	}
	return nil
}

var ArrayMap TekoFunction = customExecutedFunction(ArrayMapExecutor, []string{"f", "l"})

func newTekoString(runes []rune) Array {
	chars := []*TekoObject{}

	for _, c := range runes {
		chars = append(chars, tp(newChar(c)))
	}

	return newArray(chars)
}
