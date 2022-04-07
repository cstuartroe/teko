package checker

type GenericType struct {
	ttype TekoType
}

func (t GenericType) tekotypeToString() string {
	return "Generic[" + t.ttype.tekotypeToString() + "]"
}

func (t GenericType) allFields() map[string]TekoType {
	return t.ttype.allFields()
}

func (t *GenericType) addField(name string, ttype TekoType) {
	switch p := t.ttype.(type) {
	case *BasicType:
		p.setField(name, ttype)
	default:
		panic("Can't add field to a generic of non-object type")
	}
}

func newGenericType(name string) *GenericType {
	return &GenericType{
		ttype: newBasicType(name),
	}
}

func (c *Checker) resolveGeneric(generic *GenericType, resolution TekoType) {
	if !isTekoSubtype(resolution, generic.ttype) {
		panic("Can't resolve a generic to non-subtype")
	}

	c.generic_resolutions[generic] = resolution

	for name, ttype := range generic.allFields() {
		switch p := ttype.(type) {
		case *GenericType:
			c.resolveGeneric(p, getField(resolution, name))
		}
	}
}
