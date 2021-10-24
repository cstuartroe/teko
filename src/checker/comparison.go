package checker

func isTekoSubtype(sub TekoType, sup TekoType) bool {
	if sub == sup {
		return true
	}

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
		case *UnionType:
			return isUnionSubtype(psub, psup)
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

func isUnionSubtype(usub *UnionType, usup *UnionType) bool {
	return false
}

func isTekoEqType(t1 TekoType, t2 TekoType) bool {
	return isTekoSubtype(t1, t2) && isTekoSubtype(t2, t1)
}
