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
				ttype: devar(etype),
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
				ttype: devar(etype),
			},
		},
	}
}

// Sets

type SetType struct {
	etype  TekoType
	fields map[string]TekoType
}

func (t SetType) isDeferred() bool {
	return false
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

func (t MapType) isDeferred() bool {
	return false
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
	etype TekoType
	ttype *TemplateType
}

func (t ArrayType) isDeferred() bool {
	return false
}

func (t ArrayType) tekotypeToString() string {
	return t.etype.tekotypeToString() + "[]"
}

func (t ArrayType) allFields() map[string]TekoType {
	return t.ttype.allFields()
}

var sequenceGeneric *GenericType = newGenericType("T")
var sequenceResolutionGeneric *GenericType = newGenericType("U")

var ArrayTemplate *TemplateType = &TemplateType{
	declared_generics: []*GenericType{sequenceGeneric},
	template:          &BasicType{},
}

var MutableArrayTemplate *TemplateType = &TemplateType{
	declared_generics: []*GenericType{sequenceGeneric},
	template:          &BasicType{},
}

var arrayTypeCache map[TekoType]*ArrayType = map[TekoType]*ArrayType{}

func newArrayType(etype TekoType) *ArrayType {
	etype = deconstantize(etype)

	if at, ok := arrayTypeCache[etype]; ok {
		return at
	}

	var ttype *TemplateType

	switch p := etype.(type) {
	case *VarType:
		ttype = MutableArrayTemplate.resolveByList(p.ttype)
	default:
		ttype = ArrayTemplate.resolveByList(etype)
	}

	at := &ArrayType{
		etype: etype,
		ttype: ttype,
	}

	arrayTypeCache[etype] = at

	return at
}

func constructArrayTemplates() {
	this_type := newArrayType(sequenceGeneric)
	other_type := newArrayType(sequenceResolutionGeneric)

	fields := makeMapFields(IntType, sequenceGeneric)
	mutable_fields := makeMapFields(IntType, newVarType(sequenceGeneric))

	addType := &FunctionType{
		rtype: this_type,
		argdefs: []FunctionArgDef{
			{
				Name:  "o",
				ttype: this_type,
			},
		},
	}

	fields["add"] = addType
	mutable_fields["add"] = addType

	forEachType := &FunctionType{
		rtype: other_type,
		argdefs: []FunctionArgDef{
			{
				Name: "f",
				ttype: &FunctionType{
					rtype: sequenceResolutionGeneric,
					argdefs: []FunctionArgDef{
						{
							ttype: sequenceGeneric,
						},
					},
				},
			},
		},
	}

	fields["forEach"] = forEachType
	mutable_fields["forEach"] = forEachType

	mutable_fields["push"] = &FunctionType{
		rtype: NullType,
		argdefs: []FunctionArgDef{
			{
				Name:  "element",
				ttype: sequenceGeneric,
			},
		},
	}

	ArrayTemplate.template.(*BasicType).fields = fields
	MutableArrayTemplate.template.(*BasicType).fields = mutable_fields
}

func SetupSequenceTypes() {
	constructArrayTemplates()
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
}
