package interpreter

func updateExecutor(function TekoFunction, evaluatedArgs map[string]*TekoObject) *TekoObject {
	s, ok := evaluatedArgs["value"]
	if !ok {
		panic("No parameter passed to print")
	}

	*(function.owner) = *s

  return s
}
