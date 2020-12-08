package checker

type TekoType interface {
	allFields() map[string]TekoType
}

func getField(ttype TekoType, name string) TekoType {
	val, ok := ttype.allFields()[name]
	if ok {
		return val
	} else {
		return nil
	}
}

type BasicType struct {
	fields map[string]TekoType
}

func (ttype *BasicType) allFields() map[string]TekoType {
	return ttype.fields
}

func (ttype *BasicType) setField(name string, tekotype TekoType) {
	if getField(ttype, name) != nil {
		panic("Field " + name + " has already been declared")
	}

	ttype.fields[name] = tekotype
}

type FunctionArgDef struct {
	name    string
	mutable bool
	byref   bool
	ttype   TekoType
}

type FunctionType struct {
	rtype   TekoType
	argdefs []FunctionArgDef
}

func (ttype FunctionType) allFields() map[string]TekoType {
	return map[string]TekoType{}
}

// type checking

func isTekoSubtype(sub TekoType, sup TekoType) bool {
	for name, ttype := range sup.allFields() {
		sub_ttype := getField(sub, name)
		if (sub_ttype == nil) || !isTekoEqType(ttype, sub_ttype) {
			return false
		}
	}

	return true
}

func isTekoEqType(t1 TekoType, t2 TekoType) bool {
	if t1 == t2 {
		return true
	}

	return isTekoSubtype(t1, t2) && isTekoSubtype(t2, t1)
}
