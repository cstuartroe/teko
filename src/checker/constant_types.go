package checker

var IntType BasicType = BasicType{
  fields: map[string]TekoType{
    "add": &IntBinopType,
    "sub": &IntBinopType,
    "mult": &IntBinopType,
    "div": &IntBinopType,
    "exp": &IntBinopType,
    "mod": &IntBinopType,
  },
}

var IntBinopType FunctionType

var BoolType BasicType = BasicType{
  fields: map[string]TekoType{
    "and": &BoolBinopType,
    "or": &BoolBinopType,
  },
}

var BoolBinopType FunctionType

func atType(keyType TekoType, valueType TekoType) FunctionType {
  return FunctionType{
    rtype: valueType,
    argdefs: []FunctionArgDef{
      FunctionArgDef{
        name: "key",
        ttype: keyType,
      },
    },
  }
}

func mapType(keyType TekoType, valueType TekoType) BasicType {
  t := atType(keyType, valueType)
  return BasicType{
    fields: map[string]TekoType{
      "at": &t,
      "size": &IntType,
    },
  }
}

func arrayType(memberType TekoType) BasicType {
  return mapType(&IntType, memberType)
}

var CharType BasicType = BasicType{
  fields: map[string]TekoType{},
}

var StringType BasicType = arrayType(&CharType)

var PrintType FunctionType = FunctionType{
  rtype: nil,
  argdefs: []FunctionArgDef{
    FunctionArgDef{
      name: "s",
      ttype: &StringType,
    },
  },
}

// avoids circular initialization
func SetupFunctionTypes() {
  IntBinopType.rtype = &IntType
  IntBinopType.argdefs =  []FunctionArgDef{
    FunctionArgDef{
      name: "other",
      ttype: &IntType,
    },
  }

  BoolBinopType.rtype = &BoolType
  BoolBinopType.argdefs = []FunctionArgDef{
    FunctionArgDef{
      name: "other",
      ttype: &BoolType,
    },
  }
}