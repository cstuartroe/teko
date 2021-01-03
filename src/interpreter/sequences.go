package interpreter

type Array struct {
	elements []TekoObject
}

func (a Array) getFieldValue(name string) TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for arrays yet")
	}
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
