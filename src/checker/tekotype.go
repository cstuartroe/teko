package checker

type TekoType interface {
	tekotypeToString() string
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
	name   string
	fields map[string]TekoType
}

func newBasicType(name string) *BasicType {
	return &BasicType{
		name:   name,
		fields: map[string]TekoType{},
	}
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
