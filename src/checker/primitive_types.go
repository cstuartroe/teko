package checker

var ToStrType *FunctionType = &FunctionType{}

var IntType *BasicType = &BasicType{
	name: "int",
	fields: map[string]TekoType{
		"add":    IntBinopType,
		"sub":    IntBinopType,
		"mult":   IntBinopType,
		"div":    IntBinopType,
		"exp":    IntBinopType,
		"mod":    IntBinopType,
		"to_str": ToStrType,
	},
}

var IntBinopType *FunctionType = &FunctionType{}

var BoolType *BasicType = &BasicType{
	name: "bool",
	fields: map[string]TekoType{
		"and": BoolBinopType,
		"or":  BoolBinopType,
	},
}

var BoolBinopType *FunctionType = &FunctionType{}

var CharType *BasicType = &BasicType{
	name:   "char",
	fields: map[string]TekoType{},
}

var VoidType *BasicType = newBasicType("")

// avoids circular initialization
func SetupFunctionTypes() {
	IntBinopType.rtype = IntType
	IntBinopType.argdefs = []FunctionArgDef{
		{
			name:  "other",
			ttype: IntType,
		},
	}

	BoolBinopType.rtype = BoolType
	BoolBinopType.argdefs = []FunctionArgDef{
		{
			name:  "other",
			ttype: BoolType,
		},
	}

	ToStrType.rtype = StringType
	ToStrType.argdefs = []FunctionArgDef{}
}
