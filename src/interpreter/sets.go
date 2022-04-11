package interpreter

type Set struct {
	elements []*TekoObject
}

func (s Set) getFieldValue(name string) *TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for sets yet")
	}
}
