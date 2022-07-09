package interpreter

import "github.com/cstuartroe/teko/src/checker"

type TekoMap struct {
	kvpairs     map[int64]*TekoObject
	symbolTable SymbolTable
}

func (m TekoMap) getUnderlyingType() checker.TekoType {
	return nil // TODO
}

func tekoHash(o *TekoObject) int64 {
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
		return p.value.Int64()

	default:
		panic("hash did not return an integer")
	}
}

func (m *TekoMap) set(key *TekoObject, value *TekoObject) {
	m.kvpairs[tekoHash(key)] = value
}

func MapAtExecutor(kvpairs map[int64]*TekoObject) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return kvpairs[tekoHash(evaluatedArgs["key"])]
	}
}

// TODO get argames from checker
func (m TekoMap) getFieldValue(name string) *TekoObject {
	return m.symbolTable.cached_get(name, func() *TekoObject {
		switch name {
		case "at":
			return tp(customExecutedFunction(MapAtExecutor(m.kvpairs), checker.NoDefaults("key")))

		case "size":
			return tp(int2Integer(len(m.kvpairs)))

		default:
			panic("Unknown map function")
		}
	})
}

func newTekoMap() TekoMap {
	return TekoMap{
		kvpairs:     map[int64]*TekoObject{},
		symbolTable: newSymbolTable(nil),
	}
}
