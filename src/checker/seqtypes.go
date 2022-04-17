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

type sequence_variant int

const (
	mutable_sequence   sequence_variant = iota // retrieved elements are mutable, and (for arrays) an append method is present
	hashable_sequence                          // elements are not overwritable, and so the object has a hash_method
	sequence_supertype                         // the fewest assumptions are made - elements are immutable but there is no hash method
)

func (c *Checker) makeMapFields(keyType TekoType, valueType TekoType, variant sequence_variant) map[string]TekoType {
	at_rtype := valueType
	if variant == mutable_sequence {
		at_rtype = newVarType(valueType)
	}

	at_t := atType(keyType, at_rtype)
	in_t := includesType(valueType)

	out := map[string]TekoType{
		"at":       at_t,
		"size":     IntType,
		"includes": in_t,
		"to_str": &FunctionType{
			rtype:   StringType,
			argdefs: []FunctionArgDef{},
		},
	}

	if variant == hashable_sequence && c.isTekoSubtype(valueType, Hashable) {
		out["hash"] = HashType
	}

	return out
}

func (c *Checker) newMapType(keyType TekoType, valueType TekoType, variant sequence_variant) *MapType {
	return &MapType{
		ktype:  keyType,
		vtype:  valueType,
		fields: c.makeMapFields(keyType, valueType, variant),
	}
}

// Arrays, strings

type ArrayType struct {
	etype   TekoType
	variant sequence_variant
	fields  map[string]TekoType
}

func (t ArrayType) tekotypeToString() string {
	if t.etype == CharType && t.variant == hashable_sequence {
		return "string"
	}

	out := t.etype.tekotypeToString()

	if t.variant == mutable_sequence {
		out = "(var " + out + ")"
	}

	out += "[]"

	if t.variant == hashable_sequence {
		out += "!"
	}

	return out
}

func (t ArrayType) allFields() map[string]TekoType {
	return t.fields
}

func (c *Checker) newArrayType(etype TekoType, variant sequence_variant) *ArrayType {
	atype := &ArrayType{
		etype:   etype,
		variant: variant,
	}

	c.setArrayFields(atype)

	return atype
}

func (c *Checker) makeArrayAddType(atype *ArrayType) *FunctionType {
	rtype := atype
	// if atype.variant != mutable_sequence {
	// 	rtype = c.newArrayType(atype.etype, mutable_sequence)
	// }

	argtype := atype
	if atype.variant != sequence_supertype {
		argtype = c.newArrayType(atype.etype, sequence_supertype)
	}

	return &FunctionType{
		rtype: rtype,
		argdefs: []FunctionArgDef{
			{
				Name:  "other",
				ttype: argtype,
			},
		},
	}
}

func (c *Checker) setArrayFields(atype *ArrayType) {
	atype.fields = c.makeMapFields(IntType, atype.etype, atype.variant)

	atype.fields["add"] = c.makeArrayAddType(atype)

	if atype.variant == mutable_sequence {
		atype.fields["append"] = &FunctionType{
			rtype: NullType,
			argdefs: []FunctionArgDef{
				{
					Name:  "element",
					ttype: atype.etype,
				},
			},
		}
	}
}

var StringType *ArrayType = &ArrayType{
	etype:   CharType,
	variant: hashable_sequence,
}

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
	emptyChecker.setArrayFields(StringType)
}

var mapGenericA *GenericType = newGenericType("")
var mapGenericB *GenericType = newGenericType("")

var arrayMapType *FunctionType = &FunctionType{
	rtype: emptyChecker.newArrayType(mapGenericB, mutable_sequence),
	argdefs: []FunctionArgDef{
		{
			Name: "f",
			ttype: &FunctionType{
				rtype: mapGenericB,
				argdefs: []FunctionArgDef{
					{
						Name:  "e",
						ttype: mapGenericA,
					},
				},
			},
		},
		{
			Name:  "l",
			ttype: emptyChecker.newArrayType(mapGenericA, sequence_supertype),
		},
	},
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
