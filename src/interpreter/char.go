package interpreter

type TekoChar struct {
	value       rune
	symbolTable SymbolTable
}

func newChar(c rune) TekoChar {
	return TekoChar{
		value:       c,
		symbolTable: newSymbolTable(nil),
	}
}

func tekoEscape(c rune, escapeDoubleQuotes bool, escapeSingleQuotes bool) string {
	if c == '"' && escapeDoubleQuotes {
		return "\\\""
	} else if c == '\'' && escapeSingleQuotes {
		return "\\'"
	} else {
		return string(c)
	}
}

func CharToStrExecutor(c rune) executorType {
	return func(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
		return tp(newTekoString([]rune("'" + tekoEscape(c, false, true) + "'")))
	}
}

func (c TekoChar) getFieldValue(name string) *TekoObject {
	return c.symbolTable.cached_get(name, func() *TekoObject {
		switch name {

		case "to_str":
			return tp(customExecutedFunction(CharToStrExecutor(c.value), []string{}))

		default:
			panic("Operation not implemented: " + name)
		}
	})
}
