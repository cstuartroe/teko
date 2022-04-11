package interpreter

type TekoMap struct {
	kvpairs     map[int]*TekoObject
	symbolTable SymbolTable
}

func tekoHash(o *TekoObject) int {
	hash_function := (*o).getFieldValue("hash")
	var hash *TekoObject

	switch fp := (*hash_function).(type) {
	case TekoFunction:
		hash = fp.executor(fp, map[string]*TekoObject{})

	default:
		panic("hash is not a function")
	}

	switch p := (*hash).(type) {
	case Integer:
		return p.value

	default:
		panic("hash did not return an integer")
	}
}

func (m *TekoMap) set(key *TekoObject, value *TekoObject) {
	m.kvpairs[tekoHash(key)] = value
}

func MapAtExecutor(kvpairs map[int]*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return kvpairs[tekoHash(evaluatedArgs["key"])]
	}
}

// TODO get argames from checker
func (m TekoMap) getFieldValue(name string) *TekoObject {
	return m.symbolTable.cached_get(name, func() *TekoObject {
		switch name {
		case "at":
			return tp(customExecutedFunction(MapAtExecutor(m.kvpairs), []string{"key"}))

		case "size":
			return tp(getInteger(len(m.kvpairs)))

		default:
			panic("Unknown array function")
		}
	})
}

func newTekoMap() TekoMap {
	return TekoMap{
		kvpairs:     map[int]*TekoObject{},
		symbolTable: newSymbolTable(nil),
	}
}
