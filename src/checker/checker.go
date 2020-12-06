package checker

type TypeTable struct {
  parent *TypeTable
  table map[string]TekoType
}

func (ttable *TypeTable) get(name string) TekoType {
  if val, ok := ttable.table[name]; ok {
    return val
  } else if ttable.parent == nil {
    return nil
  } else {
    return ttable.parent.get(name)
  }
}

func (ttable *TypeTable) set(name string, val TekoType) {
  if ttable.get(name) != nil {
    panic("Existing type name")
  } else {
    ttable.table[name] = val
  }
}

var stdlibTypeTable TypeTable = TypeTable{
  parent: nil,
  table: map[string]TekoType{
    "int": &IntType,
    "bool": &BoolType,
    "str": &StringType,
    "char": &CharType,
  },
}

type Checker struct {
  typeTable *TypeTable
  ttype *BasicType
}

var BaseChecker Checker = Checker{
  typeTable: &stdlibTypeTable,
  ttype: &BasicType{
    fields: map[string]TekoType{},
  },
}

func NewChecker(parent *Checker) Checker {
  c := Checker{
    typeTable: &TypeTable{
      parent: parent.typeTable,
      table: map[string]TekoType{},
    },
    ttype: &BasicType{
      fields: map[string]TekoType{},
    },
  }

  return c
}

func (c *Checker) getType() TekoType {
  return c.ttype
}

func (c *Checker) GetType() TekoType {
  return c.getType()
}

func (c *Checker) getFieldType(name string) TekoType {
  return getField(c.ttype, name)
}

func (c *Checker) declareFieldType(name string, tekotype TekoType) {
  c.ttype.setField(name, tekotype)
}

func (c *Checker) getTypeByName(name string) TekoType {
  return c.typeTable.get(name)
}
