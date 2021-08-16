package checker

type TekoType interface {
	tekotype()
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
	fields map[string]TekoType
}

func newBasicType() *BasicType {
	return &BasicType{map[string]TekoType{}}
}

func (ttype BasicType) tekotype() {}

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

	case UnionType:
		switch psub := sub.(type) {
		case UnionType:
			return false // TODO
		default:
			return isTypeInUnion(psub, psup)
		}

	case FunctionType:
		return false // TODO

	default:
		panic("Unknown kind of type")
	}
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

func isTypeInUnion(sub TekoType, sup UnionType) bool {
	for _, ttype := range sup.types {
		if isTekoSubtype(sub, ttype) {
			return true
		}
	}

	return false
}

func isTekoEqType(t1 TekoType, t2 TekoType) bool {
	if t1 == t2 {
		return true
	}

	return isTekoSubtype(t1, t2) && isTekoSubtype(t2, t1)
}
