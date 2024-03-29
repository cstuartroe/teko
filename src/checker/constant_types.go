package checker

import (
	"math/big"
)

type ConstantTypeType int

const (
	IntConstant ConstantTypeType = iota
	StringConstant
	BoolConstant
)

type ConstantType struct {
	name   string
	fields map[string]TekoType
	ctype  ConstantTypeType
}

func (t ConstantType) isDeferred() bool {
	return false
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

var constantIntTypeCache map[string]*ConstantType = map[string]*ConstantType{}

func NewConstantIntType(n *big.Int) TekoType {
	s := n.String()

	if ct, ok := constantIntTypeCache[s]; ok {
		return ct
	}

	out := &ConstantType{
		name:   s,
		fields: copyFields(IntType),
		ctype:  IntConstant,
	}

	out.setField(s, NullType)

	constantIntTypeCache[s] = out

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

func makeConstantIntType(b bool) *ConstantType {
	name := "false"
	if b {
		name = "true"
	}

	out := &ConstantType{
		name:   name,
		fields: copyFields(BoolType),
		ctype:  BoolConstant,
	}

	out.setField("?"+name, NullType)

	return out
}

var ConstantBoolTypeCache map[bool]*ConstantType = map[bool]*ConstantType{
	true:  makeConstantIntType(true),
	false: makeConstantIntType(false),
}

func deconstantize(ttype TekoType) TekoType {
	switch p := ttype.(type) {
	case *ConstantType:
		switch p.ctype {
		case IntConstant:
			return IntType
		case StringConstant:
			return StringType
		case BoolConstant:
			return BoolType
		default:
			panic("Unknown constant type: " + ttype.tekotypeToString())
		}
	default:
		return ttype
	}
}
