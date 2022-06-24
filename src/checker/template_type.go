package checker

type TemplateType struct {
	template            TekoType
	declared_generics   []*GenericType
	generic_resolutions map[*GenericType]TekoType
}

func (t TemplateType) tekotypeToString() string {
	out := t.template.tekotypeToString()

	out += "["

	for i, g := range t.declared_generics {
		if i > 0 {
			out += ", "
		}

		out += g.tekotypeToString()
	}

	return out + "]"
}

func (t TemplateType) allFields() map[string]TekoType {
	if t.generic_resolutions == nil {
		panic("Unresolved template")
	} else {
		return degenericize(t.template, t.generic_resolutions).allFields()
	}
}

func (t *TemplateType) resolve(generic_resolutions map[*GenericType]TekoType) *TemplateType {
	if t.generic_resolutions != nil {
		panic("Already resolved")
	}

	return &TemplateType{
		template:            t.template,
		declared_generics:   t.declared_generics,
		generic_resolutions: generic_resolutions,
	}
}

func (t *TemplateType) resolveByList(generic_resolutions_slice ...TekoType) *TemplateType {
	if len(generic_resolutions_slice) != len(t.declared_generics) {
		panic("Wrong number of resolutions")
	}

	generic_resolutions := map[*GenericType]TekoType{}

	for i, g := range generic_resolutions_slice {
		generic_resolutions[t.declared_generics[i]] = g
	}

	return t.resolve(generic_resolutions)
}

func (t *TemplateType) degenericize(generic_resolutions map[*GenericType]TekoType) *TemplateType {
	if t.template == nil {
		panic("wuh")
	}

	new_generic_resolutions := map[*GenericType]TekoType{}
	var updated bool = false

	for g, r := range t.generic_resolutions {
		new_generic_resolutions[g] = r

		switch p := r.(type) {
		case *GenericType:
			if res, ok := generic_resolutions[p]; ok {
				new_generic_resolutions[g] = res
				updated = true
			}
		}
	}

	if updated {
		return &TemplateType{
			template:            t.template,
			declared_generics:   t.declared_generics,
			generic_resolutions: new_generic_resolutions,
		}
	} else {
		return t
	}
}
