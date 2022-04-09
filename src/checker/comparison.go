package checker

import "fmt"

var primitives map[*BasicType]bool = map[*BasicType]bool{
	IntType:  true,
	BoolType: true,
	CharType: true,
	VoidType: true,
}

func (c *Checker) isTekoSubtype(sub TekoType, sup TekoType) bool {
	return c.isTekoSubtypeWithAncestry(newAncestry(sub), newAncestry(sup))
}

type Ancestors map[TekoType]int

// For structural type cycle detection
type TekoTypeWithAncestry struct {
	ttype     TekoType
	ancestors Ancestors
}

var ancestor_constants map[int]TekoType = map[int]TekoType{}

func getAncestorConstant(i int) TekoType {
	if t, ok := ancestor_constants[i]; ok {
		return t
	} else {
		t = &BasicType{
			fields: map[string]TekoType{
				(fmt.Sprintf("#%d", i)): VoidType,
			},
		}
		ancestor_constants[i] = t
		return t
	}
}

func newAncestry(ttype TekoType) TekoTypeWithAncestry {
	return TekoTypeWithAncestry{
		ancestors: map[TekoType]int{},
		ttype:     ttype,
	}
}

func (a TekoTypeWithAncestry) next() Ancestors {
	var newAncestors Ancestors = map[TekoType]int{}

	for anc, i := range a.ancestors {
		newAncestors[anc] = i + 1
	}

	newAncestors[a.ttype] = 1

	return newAncestors
}

func (a TekoTypeWithAncestry) getField(name string) TekoTypeWithAncestry {
	field := getField(a.ttype, name)
	newAncestors := a.next()

	if i, ok := newAncestors[field]; ok {
		field = getAncestorConstant(i)
	}

	return TekoTypeWithAncestry{
		ancestors: newAncestors,
		ttype:     field,
	}
}

func (a TekoTypeWithAncestry) degenericize() TekoTypeWithAncestry {
	switch p := a.ttype.(type) {
	case *GenericType:
		return TekoTypeWithAncestry{
			ttype:     p.ttype,
			ancestors: a.ancestors,
		}
	default:
		panic("Not a generic type")
	}
}

func (c *Checker) isTekoSubtypeWithAncestry(sub TekoTypeWithAncestry, sup TekoTypeWithAncestry) bool {
	if sub.ttype == sup.ttype {
		return true
	}

	switch psub := sub.ttype.(type) {
	case *GenericType:
		if c.isDeclared(psub) && psub.ttype == nil {
			psub.ttype = sup.ttype
			return true
		}
	}

	switch psup := sup.ttype.(type) {
	case *GenericType:
		if !(psup.ttype == nil || c.isTekoSubtypeWithAncestry(sub, sup.degenericize())) {
			return false
		}

		if !c.isDeclared(psup) {
			_, resolved := c.generic_resolutions[psup]

			// fmt.Printf("%p\n", c)
			// fmt.Println(c.generic_resolutions)
			// fmt.Println(psup)
			// fmt.Println(psup.tekotypeToString())
			// fmt.Printf("%p\n", sub)
			// fmt.Println(sub.tekotypeToString())

			if !resolved {
				c.generic_resolutions[psup] = sub.ttype
			} else {
				// TODO: common ancestor type
				panic("Please implement")
			}

			// fmt.Println(c.generic_resolutions)
			// fmt.Println("")
		}

		return true

	case *UnionType:
		switch psub := sub.ttype.(type) {
		case *UnionType:
			return c.isUnionSubtype(psub, sub.ancestors, psup, sup.ancestors)
		default:
			return c.isTypeInUnion(sub, psup, sup.ancestors)
		}

	case *FunctionType:
		switch sub.ttype.(type) {
		case *FunctionType:
			return c.isFunctionSubtype(sub, sup)
		default:
			return false // TODO: some resolving generics should be ok
		}

	case FunctionType:
		panic("type is not a pointer: " + psup.tekotypeToString())
	case UnionType:
		panic("type is not a pointer: " + psup.tekotypeToString())

	default:
		return c.isObjectSubtype(sub, sup)
	}
}

func (c *Checker) isObjectSubtype(sub TekoTypeWithAncestry, sup TekoTypeWithAncestry) bool {
	for name := range sup.ttype.allFields() {
		sup_ta := sup.getField(name)
		sub_ta := sub.getField(name)

		if sub_ta.ttype == nil {
			switch psub := sub.ttype.(type) {
			case *GenericType:
				if c.isDeclared(psub) {
					psub.addField(name, sup_ta.ttype)
				} else {
					panic("How should this scenario be handled?")
				}

			default:
				return false
			}
		} else if !c.isTekoSubtypeWithAncestry(sub_ta, sup_ta) {
			return false
		}
	}

	return true
}

func (a TekoTypeWithAncestry) getRtype() TekoTypeWithAncestry {
	switch p := a.ttype.(type) {
	case *FunctionType:
		return TekoTypeWithAncestry{
			ttype:     p.rtype,
			ancestors: a.next(),
		}
	default:
		panic("Not a function type")
	}
}

type ArgDefWithAncestry struct {
	name string
	ta   TekoTypeWithAncestry
}

func (a TekoTypeWithAncestry) getArgdefs() []ArgDefWithAncestry {
	switch p := a.ttype.(type) {
	case *FunctionType:
		out := []ArgDefWithAncestry{}
		next := a.next()

		for _, ad := range p.argdefs {
			out = append(out, ArgDefWithAncestry{
				name: ad.name,
				ta: TekoTypeWithAncestry{
					ttype:     ad.ttype,
					ancestors: next,
				},
			})
		}

		return out

	default:
		panic("Not a function type")
	}
}

func (c *Checker) isFunctionSubtype(sub TekoTypeWithAncestry, sup TekoTypeWithAncestry) bool {
	if !c.isTekoSubtypeWithAncestry(sub.getRtype(), sup.getRtype()) {
		return false
	}

	sub_argdefs := sub.getArgdefs()
	sup_argdefs := sup.getArgdefs()

	if len(sub_argdefs) != len(sup_argdefs) {
		return false
	}

	for i, sub_argdef := range sub_argdefs {
		sup_argdef := sup_argdefs[i]

		if sub_argdef.name != sup_argdef.name {
			return false
		}

		// Yes, this is the right order.
		if !c.isTekoSubtypeWithAncestry(sup_argdef.ta, sub_argdef.ta) {
			return false
		}

		// TODO does mutability matter? depends on value semantics
	}

	// TODO default arguments and/or all arguments are a single object?

	return true
}

func (c *Checker) isTypeInUnion(sub TekoTypeWithAncestry, sup *UnionType, asup Ancestors) bool {
	for _, ttype := range sup.types {
		if c.isTekoSubtypeWithAncestry(sub, TekoTypeWithAncestry{ttype: ttype, ancestors: asup}) {
			return true
		}
	}

	return false
}

func (c *Checker) isUnionSubtype(usub *UnionType, asub Ancestors, usup *UnionType, asup Ancestors) bool {
	return false // TODO
}

func (c *Checker) isTekoEqType(t1 TekoType, t2 TekoType) bool {
	return c.isTekoSubtype(t1, t2) && c.isTekoSubtype(t2, t1)
}
