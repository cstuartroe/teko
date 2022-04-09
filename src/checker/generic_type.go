package checker

import "fmt"

type GenericType struct {
	ttype TekoType
}

func (t *GenericType) tekotypeToString() string {
	var s string
	if t.ttype == nil {
		s = fmt.Sprintf("%p", t)
	} else {
		s = t.ttype.tekotypeToString()
	}

	return "Generic[" + s + "]"
}

func (t *GenericType) allFields() map[string]TekoType {
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

func (c *Checker) degenericizeReturnType(ttype TekoType) TekoType {
	switch p := ttype.(type) {
	case *GenericType:
		return c.generic_resolutions[p] // TODO what about declared generics?

	case *UnionType:
		types := []TekoType{}
		for _, ttype := range p.types {
			types = append(types, c.degenericizeReturnType(ttype))
		}

		return &UnionType{
			types: types,
		}

	case *FunctionType:
		argdefs := []FunctionArgDef{}

		for _, argdef := range p.argdefs {
			argdefs = append(argdefs, FunctionArgDef{
				name:  argdef.name,
				ttype: c.degenericizeReturnType(ttype),
			})
		}

		return &FunctionType{
			rtype:   c.degenericizeReturnType(p.rtype),
			argdefs: argdefs,
		}

	case *BasicType:
		if _, ok := primitives[p]; ok {
			return ttype
		}

		fields := map[string]TekoType{}

		for name, ttype := range p.fields {
			fields[name] = c.degenericizeReturnType(ttype)
		}

		return &BasicType{
			name:   p.name,
			fields: fields,
		}

	case *ArrayType:
		return newArrayType(c.degenericizeReturnType(p.etype))

	default:
		return ttype
		// panic("Unknown kind of type")
	}
}
