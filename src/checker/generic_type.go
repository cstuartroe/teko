package checker

type GenericType struct {
	ttype TekoType
}

func (t GenericType) tekotypeToString() string {
	var s string
	if t.ttype == nil {
		s = ""
	} else {
		s = t.ttype.tekotypeToString()
	}

	return "Generic[" + s + "]"
}

func (t GenericType) allFields() map[string]TekoType {
	if t.ttype == nil {
		return map[string]TekoType{}
	} else {
		return t.ttype.allFields()
	}
}

func (t *GenericType) addField(name string, ttype TekoType) {
	switch p := t.ttype.(type) {
	case *BasicType:
		p.setField(name, ttype)
	case nil:
		bt := newBasicType("")
		t.ttype = bt
		bt.setField(name, ttype)
	default:
		panic("Can't add field to a generic of non-object type")
	}
}

func newGenericType(name string) *GenericType {
	return &GenericType{
		ttype: nil,
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
