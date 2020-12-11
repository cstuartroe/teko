package interpreter

type TekoObject interface {
	getFieldValue(name string) TekoObject
}

type SymbolTable struct {
	parent *SymbolTable
	table  map[string]TekoObject
}

func newSymbolTable(parent *SymbolTable) SymbolTable {
	return SymbolTable{
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

type BasicObject struct {
	symbolTable SymbolTable
}

func (o *BasicObject) getFieldValue(name string) TekoObject {
	return o.symbolTable.get(name)
}

type Integer struct {
	value       int // TODO: arbitrary-precision integer
	symbolTable SymbolTable
}

type Boolean struct {
	value bool
}

func (b Boolean) getFieldValue(name string) TekoObject {
	return nil
}

type String struct {
	value []rune
}

func (s String) getFieldValue(name string) TekoObject {
	return nil
}
