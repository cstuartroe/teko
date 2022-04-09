package checker

type UnionType struct {
	types []TekoType
}

func (ut UnionType) tekotypeToString() string {
	out := ""

	for i, tt := range ut.types {
		if i > 0 {
			out += " | "
		}
		out += tt.tekotypeToString()
	}

	return out
}

func (ut UnionType) allFields() map[string]TekoType {
	return map[string]TekoType{} // TODO should actually be all fields shared by types
}

func (c *Checker) unionTypes(t1 TekoType, t2 TekoType) TekoType {
	if c.isTekoEqType(t1, t2) {
		return t1
	}

	types := []TekoType{}

	switch p1 := t1.(type) {
	case UnionType:
		types = append(types, p1.types...)
	default:
		types = append(types, t1)
	}

	switch p2 := t2.(type) {
	case UnionType:
		types = append(types, p2.types...)
	default:
		types = append(types, t2)
	}

	return &UnionType{types}
}
