package checker

type TekoType interface {
	tekotypeToString() string
}

type ObjectType interface {
	TekoType
	allFields() map[string]TekoType
}

func getField(ttype ObjectType, name string) TekoType {
	val, ok := ttype.allFields()[name]
	if ok {
		return val
	} else {
		return nil
	}
}

func getFieldSafe(ttype TekoType, name string) TekoType {
	switch p := ttype.(type) {
	case ObjectType:
		return getField(p, name)
	default:
		return nil
	}
}

type BasicType struct {
	name   string
	fields map[string]TekoType
}

func newBasicType(name string) *BasicType {
	return &BasicType{
		name:   name,
		fields: map[string]TekoType{},
	}
}

func tekoObjectTypeShowFields(otype ObjectType) string {
	out := "{"
	for k, v := range otype.allFields() {
		out += k + ": " + v.tekotypeToString() + ", "
	}

	return out + "}"
}

func (ttype BasicType) tekotypeToString() string {
	if ttype.name != "" {
		return ttype.name
	}

	return tekoObjectTypeShowFields(ttype)
}

func (ttype BasicType) allFields() map[string]TekoType {
	return ttype.fields
}

func (ttype *BasicType) setField(name string, tekotype TekoType) {
	if getField(ttype, name) != nil {
		panic("Field " + name + " has already been declared")
	}

	ttype.fields[name] = tekotype
}

func isTekoSubtype(sub TekoType, sup TekoType) bool {
	switch psup := sup.(type) {
	case ObjectType:
		switch psub := sub.(type) {
		case ObjectType:
			return isObjectSubtype(psub, psup)
		default:
			return false // TODO
		}

	case *UnionType:
		switch psub := sub.(type) {
		case UnionType:
			return false // TODO
		default:
			return isTypeInUnion(psub, psup)
		}

	case *FunctionType:
		switch psub := sub.(type) {
		case *FunctionType:
			return isFunctionSubtype(psub, psup)
		default:
			return false // TODO
		}

	case FunctionType:
	case UnionType:
		panic("type is not a pointer: " + sup.tekotypeToString())

	default:
		panic("Unknown kind of type: " + sup.tekotypeToString())
	}

	return false // Uhhh we panicked, Go, remember? smdh
}

func isObjectSubtype(sub ObjectType, sup ObjectType) bool {
	if sub == sup {
		return true
	}

	for name, ttype := range sup.allFields() {
		sub_ttype := getField(sub, name)
		if (sub_ttype == nil) || !isTekoEqType(ttype, sub_ttype) {
			return false
		}
	}

	return true
}

func isFunctionSubtype(fsub *FunctionType, fsup *FunctionType) bool {
	if !isTekoSubtype(fsub.rtype, fsup.rtype) {
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
		if !isTekoSubtype(sup_argdef.ttype, sub_argdef.ttype) {
			return false
		}

		// TODO does mutability matter? depends on value semantics
	}

	// TODO default arguments and/or all arguments are a single object?

	return true
}

func isTypeInUnion(sub TekoType, sup *UnionType) bool {
	for _, ttype := range sup.types {
		if isTekoSubtype(sub, ttype) {
			return true
		}
	}

	return false
}

func isTekoEqType(t1 TekoType, t2 TekoType) bool {
	return isTekoSubtype(t1, t2) && isTekoSubtype(t2, t1)
}
