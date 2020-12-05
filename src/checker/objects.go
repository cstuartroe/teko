package checker

type TekoObject interface {
  getType() TekoType
  getFieldValue(name string) TekoObject
}

type SymbolTable struct {
  parent *SymbolTable
  table map[string]TekoObject
}

func (stable *SymbolTable) get(name string) TekoObject {
  if val, ok := stable.table[name]; ok {
    return val
  } else if stable.parent == nil {
    return nil
  } else {
    return stable.parent.get(name)
  }
}

func (stable *SymbolTable) set(name string, val TekoObject) {
  if stable.get(name) != nil {
    panic("Existing symbol")
  } else {
    stable.table[name] = val
  }
}

type BasicObject struct {
  ttype TekoType
  symbolTable SymbolTable
}

func (o *BasicObject) getType(name string) TekoType {
  return o.ttype
}

func (o *BasicObject) getFieldValue(name string) TekoObject {
  return o.symbolTable.get(name)
}

type Integer struct {
  value int
}

func (n Integer) getType() TekoType {
  return &IntType
}

func (n Integer) getFieldValue(name string) TekoObject {
  return nil
}

type Boolean struct {
  value bool
}

func (b Boolean) getType() TekoType {
  return &BoolType
}

func (b Boolean) getFieldValue(name string) TekoObject {
  return nil
}
