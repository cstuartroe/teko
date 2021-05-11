package interpreter

type TekoObject interface {
	getFieldValue(name string) TekoObject
}

type SymbolTable struct {
	parent *SymbolTable
	table  map[string]TekoObject
}

func newSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		parent: parent,
		table:  map[string]TekoObject{},
	}
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

func cached_get(t *SymbolTable, name string, f func() TekoObject) TekoObject {
	if attr := t.get(name); attr != nil {
		return attr
	} else {
		v := f()
		t.set(name, v)
		return v
	}
}

type BasicObject struct {
	symbolTable *SymbolTable
}

func (o BasicObject) getFieldValue(name string) TekoObject {
	return o.symbolTable.get(name)
}

type Integer struct {
	value       int // TODO: arbitrary-precision integer
	symbolTable *SymbolTable
}

type Boolean struct {
	value bool
}

var True *Boolean = &Boolean{
	value: true,
}

var False *Boolean = &Boolean{
	value: false,
}

func (b Boolean) getFieldValue(name string) TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for bools yet")
	}
}

type String struct {
	value []rune
}

func (s String) getFieldValue(name string) TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for strings yet")
	}
}
