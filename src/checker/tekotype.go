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

// avoids circular initialization
func SetupFunctionTypes() {
  IntBinopType.rtype = &IntType
  IntBinopType.argdefs = append(
    IntBinopType.argdefs,
    FunctionArgDef{
      name: "other",
      ttype: &IntType,
    },
  )

  BoolBinopType.rtype = &BoolType
  BoolBinopType.argdefs = append(
    BoolBinopType.argdefs,
    FunctionArgDef{
      name: "other",
      ttype: &BoolType,
    },
  )
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
