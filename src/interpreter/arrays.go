package interpreter

import "github.com/cstuartroe/teko/src/checker"

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
			// TODO we need to do better than pointer equality
			// It seems like maybe this shouldn't be a method of the array?
			if e == evaluatedArgs["element"] {
				return True
			}
		}
		return False
	}
}

func join(slices [][]rune, joiner []rune) []rune {
	out := []rune{}

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
		elements := [][]rune{}

		for _, e := range receiverElements {
			f := (*e).getFieldValue("to_str")

			switch fp := (*f).(type) {
			case TekoFunction:
				s := fp.executor(fp, map[string]*TekoObject{})

				switch sp := (*s).(type) {
				case String:
					elements = append(elements, sp.runes)

				default:
					panic("to_str did not return a string")
				}

			default:
				panic("to_str was not a function")
			}
		}

		return tp(newTekoString(
			append(
				[]rune{'['},
				append(
					join(
						elements,
						[]rune(", "),
					),
					']',
				)...,
			),
		))
	}
}

func ArrayForEachExecutor(receiveElements []*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		elements := []*TekoObject{}

		var f TekoFunction
		switch fp := (*evaluatedArgs["f"]).(type) {
		case TekoFunction:
			f = fp
		default:
			panic("something other than a function passed to forEach")
		}

		argname := f.argdefs[0].Name

		for _, e := range receiveElements {
			elements = append(elements, f.executor(f, map[string]*TekoObject{argname: e}))
		}

		return tp(newArray(elements))
	}
}

// TODO get argames from checker
func (a Array) getFieldValue(name string) *TekoObject {
	return a.symbolTable.cached_get(name, func() *TekoObject {
		switch name {
		case "add":
			return tp(customExecutedFunction(ArrayAddExecutor(a.elements), checker.NoDefaults("other")))

		case "at":
			return tp(customExecutedFunction(ArrayAtExecutor(a.elements), checker.NoDefaults("key")))

		case "size":
			return tp(getInteger(len(a.elements)))

		case "includes":
			return tp(customExecutedFunction(ArrayIncludesExecutor(a.elements), checker.NoDefaults("element")))

		case "to_str":
			return tp(customExecutedFunction(ArrayToStrExecutor(a.elements), checker.NoDefaults()))

		case "forEach":
			return tp(customExecutedFunction(ArrayForEachExecutor(a.elements), checker.NoDefaults("f")))

		default:
			panic("Unknown array function: " + name)
		}
	})
}

func newArray(elements []*TekoObject) Array {
	return Array{elements, newSymbolTable(nil)}
}
