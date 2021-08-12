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
	fields map[string]TekoType
	ctype  ConstantTypeType
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

func newConstantIntType(n int) TekoType {
	out := ConstantType{
		fields: copyFields(IntType),
		ctype:  IntConstant,
	}

	out.setField(strconv.Itoa(n), VoidType)

	return &out
}

func newConstantStringType(s []rune) TekoType {
	out := ConstantType{
		fields: copyFields(IntType),
		ctype:  StringConstant,
	}

	out.setField("$"+string(s), VoidType)

	return &out
}
