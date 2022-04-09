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

func (tt *TypeTable) isDeclared(g *GenericType) bool {
	_, ok := tt.declared_generics[g]

	if ok {
		return true
	} else if tt.parent == nil {
		return false
	} else {
		return tt.parent.isDeclared(g)
	}
}

func (c *Checker) isDeclared(g *GenericType) bool {
	return c.typeTable.isDeclared(g)
}

func (c *Checker) declareGeneric(g *GenericType) {
	if c.isDeclared(g) {
		panic("Cannot redeclare generic")
	}

	c.typeTable.declared_generics[g] = true
}
