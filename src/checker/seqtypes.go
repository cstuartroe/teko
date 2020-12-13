package checker

// Common function types

func atType(keyType TekoType, valueType TekoType) FunctionType {
	return FunctionType{
		rtype: valueType,
		argdefs: []FunctionArgDef{
			{
				name:  "key",
				ttype: keyType,
			},
		},
	}
}

func includesType(etype TekoType) FunctionType {
	return FunctionType{
		rtype: &BoolType,
		argdefs: []FunctionArgDef{
			{
				name:  "element",
				ttype: etype,
			},
		},
	}
}

func setAddType(etype TekoType) FunctionType {
	return FunctionType{
		rtype: &VoidType,
		argdefs: []FunctionArgDef{
			{
				name:  "element",
				ttype: etype,
			},
		},
	}
}

// Sets

type SetType struct {
	etype  TekoType
	fields map[string]TekoType
}

func (t SetType) allFields() map[string]TekoType {
	return t.fields
}

func makeSetFields(etype TekoType) map[string]TekoType {
	in_t := includesType(etype)
	add_t := setAddType(etype)

	return map[string]TekoType{
		"add":      &add_t,
		"size":     &IntType,
		"includes": &in_t,
	}
}

func newSetType(etype TekoType) SetType {
	return SetType{
		etype:  etype,
		fields: makeSetFields(etype),
	}
}

// Maps

type MapType struct {
	ktype  TekoType
	vtype  Tekotype
	fields map[string]TekoType
}

func (t MapType) allFields() map[string]TekoType {
	return t.fields
}

func makeMapFields(keyType TekoType, valueType TekoType) map[string]TekoType {
	at_t := atType(keyType, valueType)
	in_t := includesType(valueType)

	return map[string]TekoType{
		"at":       &at_t,
		"size":     &IntType,
		"includes": &in_t,
	}
}

func newMapType(keyType TekoType, valueType TekoType) MapType {
	return MapType{
		ktype:  keyType,
		vtype:  valueType,
		fields: makeMapFields(keyType, valueType),
	}
}

// Arrays, strings

type ArrayType struct {
	etype  TekoType
	fields map[string]TekoType
}

func (t ArrayType) allFields() map[string]TekoType {
	return t.fields
}

func makeArrayFields(etype TekoType) map[string]TekoType {
	fields := makeMapFields(&IntType, memberType)
	fields["add"] = FunctionType{
		rtype: &fields,
		argdefs: []FunctionArgDef{
			name:  "other",
			ttype: &fields,
		},
	}
	return fields
}

func newArrayType(etype TekoType) BasicType {
	return ArrayType{
		etype:  etype,
		fields: makeArrayFields(etype),
	}
}

var StringType BasicType = newArrayType(&CharType)

var PrintType FunctionType = FunctionType{
	rtype: nil,
	argdefs: []FunctionArgDef{
		{
			name:  "s",
			ttype: &StringType,
		},
	},
}
