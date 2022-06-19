package checker

// Common function types

func atType(keyType TekoType, valueType TekoType) *FunctionType {
	return &FunctionType{
		rtype: valueType,
		argdefs: []FunctionArgDef{
			{
				Name:  "key",
				ttype: keyType,
			},
		},
	}
}

func includesType(etype TekoType) *FunctionType {
	return &FunctionType{
		rtype: BoolType,
		argdefs: []FunctionArgDef{
			{
				Name:  "element",
				ttype: etype,
			},
		},
	}
}

func setAddType(etype TekoType) *FunctionType {
	return &FunctionType{
		rtype: NullType,
		argdefs: []FunctionArgDef{
			{
				Name:  "element",
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

func (t SetType) tekotypeToString() string {
	return t.etype.tekotypeToString() + "{}"
}

func (t SetType) allFields() map[string]TekoType {
	return t.fields
}

func makeSetFields(etype TekoType) map[string]TekoType {
	in_t := includesType(etype)
	add_t := setAddType(etype)

	return map[string]TekoType{
		"add":      add_t,
		"size":     IntType,
		"includes": in_t,
	}
}

func newSetType(etype TekoType) *SetType {
	etype = deconstantize(etype)

	stype := &SetType{
		etype:  etype,
		fields: makeSetFields(etype),
	}

	setOpType := &FunctionType{
		rtype: stype,
		argdefs: []FunctionArgDef{
			{
				Name:  "other",
				ttype: stype,
			},
		},
	}

	stype.fields["and"] = setOpType
	stype.fields["or"] = setOpType
	stype.fields["sub"] = setOpType

	return stype
}

// Maps

type MapType struct {
	ktype  TekoType
	vtype  TekoType
	fields map[string]TekoType
}

func (t MapType) tekotypeToString() string {
	return t.vtype.tekotypeToString() + "[" + t.ktype.tekotypeToString() + "]"
}

func (t MapType) allFields() map[string]TekoType {
	return t.fields
}

func makeMapFields(keyType TekoType, valueType TekoType) map[string]TekoType {
	at_t := atType(keyType, valueType)
	in_t := includesType(valueType)

	return map[string]TekoType{
		"at":       at_t,
		"size":     IntType,
		"includes": in_t,
		"to_str": &FunctionType{
			rtype:   StringType,
			argdefs: []FunctionArgDef{},
		},
	}
}

func newMapType(keyType TekoType, valueType TekoType) *MapType {
	return &MapType{
		ktype:  keyType,
		vtype:  valueType,
		fields: makeMapFields(keyType, valueType),
	}
}

var HashType *FunctionType = &FunctionType{
	rtype:   IntType,
	argdefs: []FunctionArgDef{},
}

var Hashable TekoType = &BasicType{
	name: "Hashable",
	fields: map[string]TekoType{
		"hash": HashType,
	},
}

// Arrays, strings

type ArrayType struct {
	etype  TekoType
	fields map[string]TekoType
}

func (t ArrayType) tekotypeToString() string {
	return t.etype.tekotypeToString() + "[]"
}

func (t ArrayType) allFields() map[string]TekoType {
	return t.fields
}

func _newArrayTypeWithoutForEach(etype TekoType) *ArrayType {
	etype = deconstantize(etype)

	atype := &ArrayType{
		etype:  etype,
		fields: makeMapFields(IntType, etype),
	}

	atype.fields["add"] = makeBinopType(atype)

	return atype
}

var mapGeneric *GenericType = newGenericType("")

func forEachType(etype TekoType, rtype TekoType) *FunctionType {
	g := newGenericType("")

	return &FunctionType{
		rtype: rtype,
		argdefs: []FunctionArgDef{
			{
				Name: "f",
				ttype: &FunctionType{
					rtype: g,
					argdefs: []FunctionArgDef{
						{
							Name: "e",
						},
					},
				},
			},
		},
	}
}

func newArrayType(etype TekoType) *ArrayType {
	atype := _newArrayTypeWithoutForEach(etype)

	atype.fields["forEach"] = forEachType(etype, atype)

	return atype
}

var StringType *BasicType = &BasicType{
	name: "string",
}

var ToStrType *FunctionType = &FunctionType{}

var PrintType *FunctionType = &FunctionType{
	rtype: NullType,
	argdefs: []FunctionArgDef{
		{
			Name:  "s",
			ttype: StringType,
		},
	},
}

func SetupStringTypes() {
	ToStrType.rtype = StringType
	ToStrType.argdefs = []FunctionArgDef{}

	StringType.fields = makeMapFields(IntType, CharType)
	StringType.fields["add"] = makeBinopType(StringType)
	StringType.fields["hash"] = HashType
	StringType.fields["forEach"] = forEachType(CharType, StringType)
}

func (c *Checker) generalEtype(ttype TekoType) TekoType {
	if !c.isTekoSubtype(getField(ttype, "size"), IntType) {
		return nil
	}

	at_type := getField(ttype, "at")

	switch p := at_type.(type) {
	case *FunctionType:
		if len(p.argdefs) != 1 {
			return nil
		} else if !c.isTekoSubtype(IntType, p.argdefs[0].ttype) {
			return nil
		}

		return p.rtype
	default:
		return nil
	}
}
