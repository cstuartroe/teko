package checker

import (
	"strconv"
)

func newConstantType(baseType TekoType, newFieldName string) TekoType {
	out := newBasicType()

	for name, ttype := range baseType.allFields() {
		out.fields[name] = ttype
	}

	out.setField(newFieldName, VoidType)

	return out
}

func newConstantIntType(n int) TekoType {
	return newConstantType(IntType, strconv.Itoa(n))
}

func newConstantStringType(s []rune) TekoType {
	return newConstantType(StringType, "$"+string(s))
}
