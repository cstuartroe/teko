package interpreter

type Array struct {
	elements    []TekoObject
	symbolTable *SymbolTable
}

func ArrayAddExecutor(receiverElements []TekoObject) executorType {
	return func(function *TekoFunction, evaluatedArgs map[string]TekoObject) TekoObject {
		switch p := (evaluatedArgs["other"]).(type) {
		case *Array:
			return &Array{
				append(append([]TekoObject{}, receiverElements...), p.elements...),
				newSymbolTable(nil),
			}
		default:
			panic("Non-array somehow made it past the type checker as an argument to array add!")
		}
	}
}

func ArrayAtExecutor(receiverElements []TekoObject) executorType {
	return func(function *TekoFunction, evaluatedArgs map[string]TekoObject) TekoObject {
		switch p := (evaluatedArgs["key"]).(type) {
		case *Integer:
			return receiverElements[p.value]
		default:
			panic("Non-integer somehow made it past the type checker as an argument to array at!")
		}
	}
}

func ArrayIncludesExecutor(receiverElements []TekoObject) executorType {
	return func(function *TekoFunction, evaluatedArgs map[string]TekoObject) TekoObject {
		for _, e := range receiverElements {
			if e == evaluatedArgs["element"] {
				return True
			}
		}
		return False
	}
}

func (a *Array) getFieldValue(name string) TekoObject {
	return cached_get(a.symbolTable, name, func() TekoObject {
		switch name {
		case "add":
			return customExecutedFunction(ArrayAddExecutor(a.elements), []string{"other"})

		case "at":
			return customExecutedFunction(ArrayAtExecutor(a.elements), []string{"key"})

		case "size":
			return getInteger(len(a.elements))

		case "includes":
			return customExecutedFunction(ArrayIncludesExecutor(a.elements), []string{"element"})

		default:
			panic("Unknown array function")
		}
	})
}

type Set struct {
	elements []TekoObject
}

func (s Set) getFieldValue(name string) TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for sets yet")
	}
}
