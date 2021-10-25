package interpreter

type TekoChar struct {
	value rune
}

func (c TekoChar) getFieldValue(name string) *TekoObject {
	panic("Char attributes not yet implemented")
}
