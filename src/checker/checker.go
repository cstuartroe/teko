package checker

import (
	"fmt"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

type TypeTable struct {
	parent            *TypeTable
	table             map[string]TekoType
	declared_generics map[*GenericType]bool
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

type CheckerType struct {
	fields map[string]TekoType
	parent *CheckerType
}

func (ctype CheckerType) isDeferred() bool {
	return false
}

func (ctype CheckerType) tekotypeToString() string {
	return tekoObjectTypeShowFields(ctype)
}

func (ctype CheckerType) allFields() map[string]TekoType {
	var out map[string]TekoType

	if ctype.parent == nil {
		out = map[string]TekoType{}
	} else {
		out = ctype.parent.allFields()
	}

	for name, ttype := range ctype.fields {
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

type Checker struct {
	typeTable           *TypeTable
	ctype               *CheckerType
	generic_resolutions map[*GenericType]TekoType
	deferral_allowed    bool
}

func NewChecker(parent *Checker) Checker {
	c := Checker{
		typeTable: &TypeTable{
			parent:            parent.typeTable,
			table:             map[string]TekoType{},
			declared_generics: map[*GenericType]bool{},
		},
		ctype: &CheckerType{
			fields: map[string]TekoType{},
			parent: parent.ctype,
		},
		generic_resolutions: map[*GenericType]TekoType{},
	}

	return c
}

func (c *Checker) GetType() TekoType {
	return c.ctype
}

func (c *Checker) getFieldType(name string) TekoType {
	return getField(c.ctype, name)
}

func (c *Checker) declareFieldType(token *lexparse.Token, tekotype TekoType) {
	name := string(token.Value)

	if c.getFieldType(name) != nil {
		token.Raise(shared.NameError, fmt.Sprintf("Variable called %s has already been declared", name))
	}

	c.ctype.setField(name, tekotype)
}

func (c *Checker) getTypeByName(name string) TekoType {
	return c.typeTable.get(name)
}

func (c *Checker) declareNamedType(token *lexparse.Token, tekotype TekoType) {
	name := string(token.Value)

	if c.typeTable.get(name) != nil {
		token.Raise(shared.NameError, "Type has already been declared: "+name)
	}

	c.typeTable.set(name, tekotype)
}
