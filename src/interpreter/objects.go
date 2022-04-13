package interpreter

type TekoObject interface {
	getFieldValue(name string) *TekoObject
}

// This function is necessitated by dumb Go pointer semantics
func tp(t TekoObject) *TekoObject {
	return &t
}

type SymbolTable struct {
	parent *SymbolTable
	table  map[string]*TekoObject
}

func newSymbolTable(parent *SymbolTable) SymbolTable {
	return SymbolTable{
		parent: parent,
		table:  map[string]*TekoObject{},
	}
}

func (stable *SymbolTable) get(name string) *TekoObject {
	if val, ok := stable.table[name]; ok {
		return val
	} else if stable.parent == nil {
		return nil
	} else {
		return stable.parent.get(name)
	}
}

func (stable *SymbolTable) set(name string, val *TekoObject) {
	if stable.get(name) != nil {
		panic("Existing symbol")
	} else {
		stable.table[name] = val
	}
}

func (stable *SymbolTable) cached_get(name string, f func() *TekoObject) *TekoObject {
	if attr := stable.get(name); attr != nil {
		return attr
	} else {
		v := f()
		stable.set(name, v)
		return v
	}
}

type BasicObject struct {
	symbolTable SymbolTable
}

func (o BasicObject) getFieldValue(name string) *TekoObject {
	return o.symbolTable.get(name)
}

type Integer struct {
	value       int // TODO: arbitrary-precision integer
	symbolTable SymbolTable
}

type Boolean struct {
	value bool
}

var True *TekoObject = tp(Boolean{
	value: true,
})

var False *TekoObject = tp(Boolean{
	value: false,
})

func (b Boolean) getFieldValue(name string) *TekoObject {
	switch name {
	default:
		panic("Attributes haven't been implemented for bools yet")
	}
}

var Null *TekoObject = tp(BasicObject{})
