package checker

type DeferredType struct {
	ttype TekoType
}

func (t DeferredType) tekotypeToString() string {
	return t.ttype.tekotypeToString()
}

func (t DeferredType) allFields() map[string]TekoType {
	return t.ttype.allFields()
}

func undefer(t TekoType) TekoType {
	switch p := t.(type) {
	case *DeferredType:
		return p.ttype
	default:
		return t
	}
}
