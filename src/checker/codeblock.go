package checker

import (
  "github.com/cstuartroe/teko/src/lexparse"
)

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
  },
}

type Codeblock struct {
  statements []lexparse.Node
  typeTable *TypeTable
  ttype *BasicType
  symbolTable *SymbolTable
  parser lexparse.Parser
}

func (c *Codeblock) GetStatements() []lexparse.Node {
  return c.statements
}

func (c *Codeblock) getType() TekoType {
  return c.ttype
}

func (c *Codeblock) getFieldValue(name string) TekoObject {
  return c.symbolTable.get(name)
}

func (c *Codeblock) startFile(filename string) {
  c.parser.LexFile(filename)
}

func (c *Codeblock) hasMore() bool {
  return c.parser.HasMore()
}

func (c *Codeblock) grabStatement() {
  c.statements = append(
    c.statements,
    c.parser.GrabStatement(),
  )
  c.parser.Expect(lexparse.SemicolonT); c.parser.Advance()
}

func NewCodeblock(parent *Codeblock) Codeblock {
  c := Codeblock{
    statements: []lexparse.Node{},
    typeTable: &TypeTable{
      parent: parent.typeTable,
      table: map[string]TekoType{},
    },
    ttype: &BasicType{
      fields: map[string]TekoType{},
    },
    symbolTable: &SymbolTable{
      parent: parent.symbolTable,
      table: map[string]TekoObject{},
    },
  }

  c.parser = lexparse.Parser{}

  return c
}

var BaseCodeblock Codeblock = Codeblock{
  statements: []lexparse.Node{},
  typeTable: &stdlibTypeTable,
  ttype: &BasicType{
    fields: map[string]TekoType{},
  },
}
