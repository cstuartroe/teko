package checker

import "fmt"

type GenericType struct {
	ttype TekoType
}

func (t *GenericType) isDeferred() bool {
	// TODO I don't think we want generic types to be deferrable, but needs more thought
	return false
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
		ttype: nil, // TODO?: newBasicType(name),
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

func degenericize(ttype TekoType, generic_resolutions map[*GenericType]TekoType, ancestors map[TekoType]bool) TekoType {
	if ancestors == nil {
		ancestors = map[TekoType]bool{}
	}

	if ancestors[ttype] {
		return ttype
	}

	ancestors[ttype] = true

	var out TekoType

	switch p := ttype.(type) {
	case *GenericType:
		if res, ok := generic_resolutions[p]; ok { // TODO what about declared generics?
			out = res
		} else {
			out = ttype
		}

	case *TemplateType:
		out = p.degenericize(generic_resolutions)

	case *UnionType:
		types := []TekoType{}
		for _, ttype := range p.types {
			types = append(types, degenericize(ttype, generic_resolutions, ancestors))
		}

		out = &UnionType{
			types: types,
		}

	case *FunctionType:
		argdefs := []FunctionArgDef{}

		for _, argdef := range p.argdefs {
			argdefs = append(argdefs, FunctionArgDef{
				Name:  argdef.Name,
				ttype: degenericize(argdef.ttype, generic_resolutions, ancestors),
			})
		}

		out = &FunctionType{
			rtype:   degenericize(p.rtype, generic_resolutions, ancestors),
			argdefs: argdefs,
		}

	case *BasicType:
		if _, ok := primitives[p]; ok {
			out = ttype
		} else {
			fields := map[string]TekoType{}

			for name, ttype := range p.fields {
				fields[name] = degenericize(ttype, generic_resolutions, ancestors)
			}

			out = &BasicType{
				name:   p.name,
				fields: fields,
			}
		}

	case *ArrayType:
		out = newArrayType(degenericize(p.etype, generic_resolutions, ancestors))

	case *VarType:
		out = newVarType(degenericize(p.ttype, generic_resolutions, ancestors))

	default:
		out = ttype
	}

	delete(ancestors, ttype)

	return out
}
