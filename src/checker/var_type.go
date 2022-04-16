package checker

type VarType struct {
	ttype TekoType
}

func (t VarType) tekotypeToString() string {
	return "var " + t.ttype.tekotypeToString()
}

func (t VarType) allFields() map[string]TekoType {
	out := map[string]TekoType{}

	for k, v := range t.ttype.allFields() {
		out[k] = v
	}

	out["="] = &FunctionType{
		rtype: t.ttype,
		argdefs: []FunctionArgDef{
			{
				Name:  "value",
				ttype: t.ttype,
			},
		},
	}

	return out
}

func newVarType(ttype TekoType) *VarType {
	switch (ttype).(type) {
	case *VarType:
		panic("VarType")
	case VarType:
		panic("VarType")
	default:
		return &VarType{deconstantize(ttype)}
	}
}

func isvar(ttype TekoType) bool {
	switch (ttype).(type) {
	case *VarType:
		return true
	default:
		return false
	}
}

func devar(ttype TekoType) TekoType {
	switch p := (ttype).(type) {
	case *VarType:
		return p.ttype
	default:
		return ttype
	}
}
