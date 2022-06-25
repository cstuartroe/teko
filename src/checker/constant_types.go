package checker

import (
	"strconv"
)

type ConstantTypeType int

const (
	IntConstant ConstantTypeType = iota
	StringConstant
)

type ConstantType struct {
	name   string
	fields map[string]TekoType
	ctype  ConstantTypeType
}

func (t ConstantType) tekotypeToString() string {
	return t.name
}

func (t ConstantType) allFields() map[string]TekoType {
	return t.fields
}

func (t *ConstantType) setField(name string, tekotype TekoType) {
	if getField(t, name) != nil {
		panic("Field " + name + " has already been declared")
	}

	t.fields[name] = tekotype
}

func copyFields(baseType TekoType) map[string]TekoType {
	fields := map[string]TekoType{}

	for name, ttype := range baseType.allFields() {
		fields[name] = ttype
	}

	return fields
}

var constantIntTypeCache map[int]*ConstantType = map[int]*ConstantType{}

func NewConstantIntType(n int) TekoType {
	if ct, ok := constantIntTypeCache[n]; ok {
		return ct
	}

	out := &ConstantType{
		name:   strconv.Itoa(n),
		fields: copyFields(IntType),
		ctype:  IntConstant,
	}

	out.setField(strconv.Itoa(n), NullType)

	constantIntTypeCache[n] = out

	return out
}

var constantStringTypeCache map[string]*ConstantType = map[string]*ConstantType{}

func NewConstantStringType(s []rune) TekoType {
	if ct, ok := constantStringTypeCache[string(s)]; ok {
		return ct
	}

	out := &ConstantType{
		name:   "\"" + string(s) + "\"",
		fields: copyFields(StringType),
		ctype:  StringConstant,
	}

	out.setField("$"+string(s), NullType)

	constantStringTypeCache[string(s)] = out

	return out
}

func deconstantize(ttype TekoType) TekoType {
	switch p := ttype.(type) {
	case *ConstantType:
		switch p.ctype {
		case IntConstant:
			return IntType
		case StringConstant:
			return StringType
		default:
			panic("Unknown constant type: " + ttype.tekotypeToString())
		}
	default:
		return ttype
	}
}
