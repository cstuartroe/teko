package checker

func (c *Checker) translateType(ttype TekoType) TekoType {
	switch p := ttype.(type) {
	case *GenericType:
		return c.generic_resolutions[p]

	case *UnionType:
		types := []TekoType{}
		for _, ttype := range p.types {
			types = append(types, c.translateType(ttype))
		}

		return &UnionType{
			types: types,
		}

	case *FunctionType:
		argdefs := []FunctionArgDef{}

		for _, argdef := range p.argdefs {
			argdefs = append(argdefs, FunctionArgDef{
				name:  argdef.name,
				ttype: c.translateType(ttype),
			})
		}

		return &FunctionType{
			rtype:   c.translateType(p.rtype),
			argdefs: argdefs,
		}

	default:
		return ttype
	}
}

func (c *Checker) isTekoSubtype(sub TekoType, sup TekoType) bool {
	if sub == sup {
		return true
	}

	switch psub := sub.(type) {
	case *GenericType:
		if c.isDeclared(psub) && psub.ttype == nil {
			psub.ttype = sup
			return true
		}
	}

	switch psup := sup.(type) {
	case *GenericType:
		if !(psup.ttype == nil || c.isTekoSubtype(sub, psup.ttype)) {
			return false
		}

		if !c.isDeclared(psup) {
			_, resolved := c.generic_resolutions[psup]

			if !resolved {
				c.generic_resolutions[psup] = sub
			} else {
				// TODO: common ancestor type
				panic("Please implement")
			}
		}

		return true

	case *UnionType:
		switch psub := sub.(type) {
		case *UnionType:
			return c.isUnionSubtype(psub, psup)
		default:
			return c.isTypeInUnion(psub, psup)
		}

	case *FunctionType:
		switch psub := sub.(type) {
		case *FunctionType:
			return c.isFunctionSubtype(psub, psup)
		default:
			return false // TODO: some resolving generics should be ok
		}

	case FunctionType:
		panic("type is not a pointer: " + sup.tekotypeToString())
	case UnionType:
		panic("type is not a pointer: " + sup.tekotypeToString())

	default:
		return c.isObjectSubtype(sub, psup)
	}
}

func (c *Checker) isObjectSubtype(sub TekoType, sup TekoType) bool {
	for name, ttype := range sup.allFields() {
		sub_ttype := getField(sub, name)

		if sub_ttype == nil {
			switch psub := sub.(type) {
			case *GenericType:
				if c.isDeclared(psub) {
					psub.addField(name, ttype)
				} else {
					panic("How should this scenario be handled?")
				}

			default:
				return false
			}
		} else if !c.isTekoSubtype(sub_ttype, ttype) {
			return false
		}
	}

	return true
}

func (c *Checker) isFunctionSubtype(fsub *FunctionType, fsup *FunctionType) bool {
	if !c.isTekoSubtype(fsub.rtype, fsup.rtype) {
		return false
	}

	if len(fsub.argdefs) != len(fsup.argdefs) {
		return false
	}

	for i, sub_argdef := range fsub.argdefs {
		sup_argdef := fsup.argdefs[i]

		if sub_argdef.name != sup_argdef.name {
			return false
		}

		// Yes, this is the right order.
		if !c.isTekoSubtype(sup_argdef.ttype, sub_argdef.ttype) {
			return false
		}

		// TODO does mutability matter? depends on value semantics
	}

	// TODO default arguments and/or all arguments are a single object?

	return true
}

func (c *Checker) isTypeInUnion(sub TekoType, sup *UnionType) bool {
	for _, ttype := range sup.types {
		if c.isTekoSubtype(sub, ttype) {
			return true
		}
	}

	return false
}

func (c *Checker) isUnionSubtype(usub *UnionType, usup *UnionType) bool {
	return false // TODO
}

func (c *Checker) isTekoEqType(t1 TekoType, t2 TekoType) bool {
	return c.isTekoSubtype(t1, t2) && c.isTekoSubtype(t2, t1)
}
