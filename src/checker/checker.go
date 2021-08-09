package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

type TypeTable struct {
	parent *TypeTable
	table  map[string]TekoType
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

var stdlibTypeTable *TypeTable = &TypeTable{
	parent: nil,
	table: map[string]TekoType{
		"int":  IntType,
		"bool": BoolType,
		"str":  StringType,
		"char": CharType,
	},
}

type CheckerType struct {
	fields map[string]TekoType
	parent *CheckerType
}

func (ctype *CheckerType) allFields() map[string]TekoType {
	var out map[string]TekoType

	if ctype.parent == nil {
		out = map[string]TekoType{}
	} else {
		out = ctype.parent.allFields()
	}

	for name, ttype := range ctype.fields {
		if _, ok := out[name]; ok {
			panic("Field somehow got declared twice: " + name)
		}

		out[name] = ttype
	}

	return out
}

func (ctype *CheckerType) setField(name string, tekotype TekoType) {
	if getField(ctype, name) != nil {
		panic("Field " + name + " has already been declared")
	}

	ctype.fields[name] = tekotype
}

var BaseCheckerTypeFields map[string]TekoType = map[string]TekoType{
	"print": PrintType,
}

var baseCheckerType *CheckerType = &CheckerType{
	fields: BaseCheckerTypeFields,
	parent: nil,
}

type Checker struct {
	typeTable *TypeTable
	ctype     *CheckerType
}

var BaseChecker *Checker = &Checker{
	typeTable: stdlibTypeTable,
	ctype:     baseCheckerType,
}

func NewChecker(parent *Checker) Checker {
	c := Checker{
		typeTable: &TypeTable{
			parent: parent.typeTable,
			table:  map[string]TekoType{},
		},
		ctype: &CheckerType{
			fields: map[string]TekoType{},
			parent: parent.ctype,
		},
	}

	return c
}

func (c *Checker) getType() TekoType {
	return c.ctype
}

func (c *Checker) GetType() TekoType {
	return c.getType()
}

func (c *Checker) getFieldType(name string) TekoType {
	return getField(c.ctype, name)
}

func (c *Checker) declareFieldType(token lexparse.Token, tekotype TekoType) {
	name := string(token.Value)

	if c.getFieldType(name) != nil {
		token.Raise(lexparse.NameError, "Field has already been declared: " + name)
	}

	c.ctype.setField(name, tekotype)
}

func (c *Checker) getTypeByName(name string) TekoType {
	return c.typeTable.get(name)
}
