package checker

var ToStrType FunctionType

var IntType BasicType = BasicType{
	fields: map[string]TekoType{
		"add":    &IntBinopType,
		"sub":    &IntBinopType,
		"mult":   &IntBinopType,
		"div":    &IntBinopType,
		"exp":    &IntBinopType,
		"mod":    &IntBinopType,
		"to_str": &ToStrType,
	},
}

var IntBinopType FunctionType

var BoolType BasicType = BasicType{
	fields: map[string]TekoType{
		"and": &BoolBinopType,
		"or":  &BoolBinopType,
	},
}

var BoolBinopType FunctionType

var CharType BasicType = BasicType{
	fields: map[string]TekoType{},
}

var VoidType BasicType

// avoids circular initialization
func SetupFunctionTypes() {
	IntBinopType.rtype = &IntType
	IntBinopType.argdefs = []FunctionArgDef{
		{
			name:  "other",
			ttype: &IntType,
		},
	}

	BoolBinopType.rtype = &BoolType
	BoolBinopType.argdefs = []FunctionArgDef{
		{
			name:  "other",
			ttype: &BoolType,
		},
	}

	ToStrType.rtype = &StringType
	ToStrType.argdefs = []FunctionArgDef{}
}
