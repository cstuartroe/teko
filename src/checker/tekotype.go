package checker

type TekoType interface {
  allFields() map[string]TekoType
}

func getField(ttype TekoType, name string) TekoType {
  val, ok := ttype.allFields()[name]
  if ok {
    return val
  } else {
    return nil
  }
}

type BasicType struct {
  fields map[string]TekoType
}

func (ttype *BasicType) allFields() map[string]TekoType {
  return ttype.fields
}

func (ttype *BasicType) setField(name string, tekotype TekoType) {
  if getField(ttype, name) != nil {
    panic("Field " + name + " has already been declared")
  }

  ttype.fields[name] = tekotype
}

type FunctionArgDef struct {
  name string
  mutable bool
  byref bool
  ttype TekoType
}

type FunctionType struct {
  rtype TekoType
  argdefs []FunctionArgDef
}

func (ttype *FunctionType) allFields() map[string]TekoType {
  return map[string]TekoType{}
}

// Constants

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

// type checking

func isTekoSubtype(sub TekoType, sup TekoType) bool {
  for name, ttype := range sup.allFields() {
    sub_ttype := getField(sub, name)
    if (sub_ttype == nil) || !isTekoEqType(ttype, sub_ttype) {
      return false
    }
  }

  return true
}

func isTekoEqType(t1 TekoType, t2 TekoType) bool {
  if t1 == t2 {
    return true
  }

  return isTekoSubtype(t1, t2) && isTekoSubtype(t2, t1)
}
