package checker

type TekoType interface {
	tekotypeToString() string
	allFields() map[string]TekoType
	isDeferred() bool
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
	name     string
	fields   map[string]TekoType
	deferred bool
}

func newBasicType(name string) *BasicType {
	return &BasicType{
		name:   name,
		fields: map[string]TekoType{},
	}
}

func (ttype BasicType) isDeferred() bool {
	return ttype.deferred
}

func tekoObjectTypeShowFields(otype TekoType) string {
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

func (c *Checker) greatestCommonAncestor(t1 TekoType, t2 TekoType) TekoType {
	if c.isTekoEqType(t1, t2) {
		return t1
	}

	fields := map[string]TekoType{}

	for name, ttype1 := range t1.allFields() {
		ttype2 := getField(t2, name)

		if ttype2 != nil {
			fields[name] = c.unionTypes(ttype1, ttype2)
		}
	}

	return &BasicType{
		name:   "",
		fields: fields,
	}
}
