package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

func (c Checker) CheckTree(codeblock lexparse.Codeblock) {
	for _, stmt := range codeblock.Statements {
		c.checkStatement(stmt)
	}
}

func (c *Checker) checkStatement(stmt lexparse.Statement) {
	switch p := stmt.(type) {
	case lexparse.ExpressionStatement:
		c.checkExpressionStatement(p)

	case lexparse.TypeStatement:
		c.checkTypeStatement(p)

	default:
		stmt.Token().Raise(shared.NotImplementedError, "Unknown statement type")
	}
}

func (c *Checker) checkExpressionStatement(stmt lexparse.ExpressionStatement) {
	c.checkExpression(stmt.Expression, nil)
}

func (c *Checker) checkTypeStatement(stmt lexparse.TypeStatement) {
	c.declareNamedType(
		stmt.Name,
		c.evaluateType(stmt.TypeExpression),
	)
}

func (c *Checker) checkExpression(expr lexparse.Expression, expectedType TekoType) TekoType {
	return devar(c.checkExpressionAllowingVar(expr, expectedType))
}

func (c *Checker) checkExpressionAllowingVar(expr lexparse.Expression, expectedType TekoType) TekoType {
	var ttype TekoType

	switch p := expr.(type) {

	case lexparse.SimpleExpression:
		ttype = c.checkSimpleExpression(p)

	case lexparse.DeclarationExpression:
		ttype = c.checkDeclaration(p, expectedType)

	case lexparse.CallExpression:
		ttype = c.checkCallExpression(p)

	case lexparse.AttributeExpression:
		ttype = c.checkAttributeExpression(p)

	case lexparse.IfExpression:
		ttype = c.checkIfExpression(p, expectedType)

	case lexparse.SequenceExpression:
		ttype = c.checkSequenceExpression(p, expectedType)

	case lexparse.MapExpression:
		ttype = c.checkMapExpression(p, expectedType)

	case lexparse.ObjectExpression:
		ttype = c.checkObjectExpression(p, expectedType)

	case lexparse.FunctionExpression:
		// TODO use expected type in function definitions
		ttype = c.checkFunctionDefinition(p, expectedType)

	case lexparse.DoExpression:
		ttype = c.checkDoExpression(p, expectedType)

	case lexparse.VarExpression:
		expr.Token().Raise(shared.SyntaxError, "Illegal start to expression")

	case lexparse.WhileExpression:
		ttype = c.checkWhileExpression(p)

	case lexparse.ScopeExpression:
		ttype = c.checkScopeExpression(p)

	case lexparse.ComparisonExpression:
		ttype = c.checkComparisonExpression(p)

	default:
		expr.Token().Raise(shared.NotImplementedError, "Cannot typecheck expression type")
	}

	if ttype == nil {
		expr.Token().Raise(shared.TypeError, "Evaluated to nil type")
	}
	if (expectedType != nil) && !c.isTekoSubtype(ttype, expectedType) {
		expr.Token().Raise(shared.TypeError, "Actual type "+ttype.tekotypeToString()+" does not fulfill expected type "+expectedType.tekotypeToString())
	}

	if expectedType == nil {
		return ttype
	} else {
		return degenericize(expectedType, c.generic_resolutions)
	}
}

func (c *Checker) checkSimpleExpression(expr lexparse.SimpleExpression) TekoType {
	t := expr.Token()
	switch t.TType {
	case lexparse.SymbolT:
		ttype := c.getFieldType(string(t.Value))
		if ttype == nil {
			t.Raise(shared.NameError, "Undefined variable")
			return nil
		} else {
			return ttype
		}

	// TODO bool, char, and float constant types

	case lexparse.IntT:
		return c.evaluateConstantIntType(t)

	case lexparse.BoolT:
		return BoolType

	case lexparse.StringT:
		return newConstantStringType(t.Value)

	case lexparse.CharT:
		return CharType

	default:
		panic("Unknown simple expression type")
	}
}

func (c *Checker) checkDeclaration(decl lexparse.DeclarationExpression, expectedType TekoType) TekoType {
	var tekotype TekoType = c.evaluateType(decl.Tekotype)

	if tekotype == nil {
		tekotype = expectedType
	}

	return c.declare(decl.Symbol, decl.Right, tekotype)
}

func (c *Checker) declare(symbol lexparse.Token, right lexparse.Expression, tekotype TekoType) TekoType {
	var output_type TekoType

	switch p := (tekotype).(type) {
	case *VarType:
		c.checkExpression(right, p.ttype)
		output_type = tekotype
	default:
		output_type = c.checkExpression(right, tekotype)
	}

	c.declareFieldType(symbol, output_type)

	return output_type
}

func (c *Checker) checkCallExpression(expr lexparse.CallExpression) TekoType {
	call_checker := NewChecker(c)

	receiver_tekotype := call_checker.checkExpression(expr.Receiver, nil)

	switch ftype := receiver_tekotype.(type) {
	case *FunctionType:
		args_by_name := ResolveArgs(ftype.argdefs, expr)

		for _, argdef := range ftype.argdefs {
			arg, ok := args_by_name[argdef.Name]
			if !ok {
				expr.Token().Raise(shared.ArgumentError, "Argument was not passed: "+argdef.Name)
			}

			call_checker.checkExpression(arg, argdef.ttype)
		}

		return degenericize(ftype.rtype, call_checker.generic_resolutions)

	case FunctionType:
		panic("Use *FunctionType instead of FunctionType")

	default:
		expr.Token().Raise(shared.TypeError, "Expression does not have a function type")
		return nil
	}
}

func (c *Checker) checkAttributeExpression(expr lexparse.AttributeExpression) TekoType {
	var left_tekotype TekoType

	if string(expr.Symbol.Value) == "=" {
		left_tekotype = c.checkExpressionAllowingVar(expr.Left, nil)
	} else {
		left_tekotype = c.checkExpression(expr.Left, nil)
	}

	name := string(expr.Symbol.Value)
	tekotype := getField(left_tekotype, name)

	if tekotype != nil {
		return tekotype
	} else {
		switch p := left_tekotype.(type) {
		case *GenericType:
			if !c.isDeclared(p) {
				panic("Generic must be declared")
			}

			out := newGenericType("")

			p.addField(name, out)
			c.declareGeneric(out)

			return out

		default:
			expr.Symbol.Raise(shared.NameError, "No such field: "+string(expr.Symbol.Value)+" on "+left_tekotype.tekotypeToString())
			return nil
		}
	}
}

func (c *Checker) checkIfExpression(expr lexparse.IfExpression, expectedType TekoType) TekoType {
	c.checkExpression(expr.Condition, BoolType)

	then_tekotype := c.checkExpression(expr.Then, expectedType)
	else_tekotype := c.checkExpression(expr.Else, expectedType)

	if then_tekotype != else_tekotype {
		expr.Else.Token().Raise(shared.TypeError, "Then and else blocks have mismatching types")
	}

	return then_tekotype
}

func (c *Checker) checkSequenceExpression(expr lexparse.SequenceExpression, expectedType TekoType) TekoType {
	var etype TekoType
	var seqtype TekoType = expectedType

	switch p := expectedType.(type) {
	case *ArrayType:
		if expr.Stype == lexparse.ArraySeqType {
			etype = p.etype
		} else {
			return nil
		}
	case *SetType:
		if expr.Stype == lexparse.SetSeqType {
			etype = p.etype
		} else {
			return nil
		}
	case nil:
		if len(expr.Elements) == 0 {
			// TODO see whether it's possible to defer inference
			expr.Token().Raise(shared.TypeError, "With no expected type, sequence cannot be empty")
		} else {
			etype = deconstantize(c.checkExpression(expr.Elements[0], nil))

			if expr.Stype == lexparse.ArraySeqType {
				seqtype = newArrayType(etype)
			} else if expr.Stype == lexparse.SetSeqType {
				seqtype = newSetType(etype)
			} else {
				panic("Unknown sequence type: " + expr.Stype)
			}
		}
	default:
		return nil
	}

	for _, element := range expr.Elements {
		c.checkExpression(element, etype)
	}

	return seqtype
}

func (c *Checker) checkMapExpression(expr lexparse.MapExpression, expectedType TekoType) TekoType {
	var ktype, vtype TekoType = nil, nil

	switch p := expectedType.(type) {
	case nil:
		if expr.Ktype != nil {
			ktype = c.evaluateType(expr.Ktype)
			vtype = c.evaluateType(expr.Vtype)
		}

	case *MapType:
		if expr.Ktype != nil {
			stated_ktype := c.evaluateType(expr.Ktype)

			if !c.isTekoEqType(stated_ktype, p.ktype) {
				expr.Ktype.Token().Raise(shared.TypeError, "Expected key type "+p.ktype.tekotypeToString()+" not "+ktype.tekotypeToString())
			}

			stated_vtype := c.evaluateType(expr.Vtype)

			if !c.isTekoEqType(stated_vtype, p.vtype) {
				expr.Vtype.Token().Raise(shared.TypeError, "Expected value type "+p.vtype.tekotypeToString()+" not "+vtype.tekotypeToString())
			}
		}

		// Want ktype and vtype to come from expectedType just in case one of them is equivalent to but
		// structurally different from stated_ktype/stated_vtype
		ktype, vtype = p.ktype, p.vtype
	}

	if (ktype == nil) != (vtype == nil) {
		panic("?")
	}

	i := 0

	if ktype == nil {
		if len(expr.KVPairs) == 0 {
			expr.MapToken.Raise(shared.TypeError, "Map without stated or expected type cannot be empty")
		}

		ktype = deconstantize(c.checkExpression(expr.KVPairs[0].Key, nil))
		vtype = c.checkExpression(expr.KVPairs[0].Value, nil)
		i += 1
	}

	if !c.isTekoSubtype(ktype, Hashable) {
		expr.Token().Raise(shared.TypeError, "Unhashable type: "+ktype.tekotypeToString())
	}

	for i < len(expr.KVPairs) {
		c.checkExpression(expr.KVPairs[i].Key, ktype)
		c.checkExpression(expr.KVPairs[i].Value, vtype)
		i += 1
	}

	return newMapType(ktype, vtype)
}

func (c *Checker) checkObjectExpression(expr lexparse.ObjectExpression, expectedType TekoType) TekoType {
	fields := map[string]TekoType{}

	for _, of := range expr.Fields {
		field_name := string(of.Symbol.Value)

		if _, ok := fields[field_name]; ok {
			of.Symbol.Raise(shared.NameError, "Duplicate member")
		}

		var expectedFieldType TekoType = nil
		if expectedType != nil && !isGeneric(expectedType) {
			expectedFieldType = getField(expectedType, field_name)
			if expectedFieldType == nil {
				of.Symbol.Raise(shared.TypeError, "Object literal cannot have unreachable field")
			}
		}

		switch p := of.Value.(type) {
		case lexparse.VarExpression:
			fields[field_name] = newVarType(c.checkExpression(p.Right, devar(expectedFieldType)))
		default:
			fields[field_name] = c.checkExpression(of.Value, expectedFieldType)
		}
	}

	return &BasicType{
		name:   "",
		fields: fields,
	}
}

func (c *Checker) checkFunctionDefinition(expr lexparse.FunctionExpression, expectedType TekoType) TekoType {
	var expectedFType *FunctionType = nil
	switch fp := expectedType.(type) {
	case *FunctionType:
		expectedFType = fp
	}

	blockChecker := NewChecker(c)
	blockChecker.declared_generics = map[*GenericType]bool{}

	for _, gd := range expr.GDL.Declarations {
		g := newGenericType(string(gd.Name.Value))

		blockChecker.declareGeneric(g)
		blockChecker.declareNamedType(gd.Name, g)
	}

	argdefs := []FunctionArgDef{}

	for i, ad := range expr.Argdefs {
		var ttype TekoType

		if ad.Tekotype != nil {
			ttype = blockChecker.evaluateType(ad.Tekotype)
			if isvar(ttype) {
				ad.Tekotype.Token().Raise(shared.TypeError, "Function arguments cannot be mutable. Complain to Conor if you hate this fact.")
			}
		} else if expectedFType != nil && len(expectedFType.argdefs) > i {
			ttype = expectedFType.argdefs[i].ttype
		} else {
			g := newGenericType("")
			blockChecker.declareGeneric(g)
			ttype = g
		}

		blockChecker.declareFieldType(ad.Symbol, ttype)

		// TODO: get argdefs from blockChecker
		argdefs = append(argdefs, FunctionArgDef{
			Name:    string(ad.Symbol.Value),
			ttype:   ttype,
			Default: ad.Default,
		})
	}

	var rtype TekoType = nil
	if expr.Rtype != nil {
		rtype = blockChecker.evaluateType(expr.Rtype)
	} else if expectedFType != nil {
		rtype = expectedFType.rtype
	}

	ftype := &FunctionType{
		rtype:   blockChecker.checkExpression(expr.Right, rtype),
		argdefs: argdefs,
	}

	if expr.Name != nil {
		c.declareFieldType(*expr.Name, ftype)
	}

	return ftype
}

func (c *Checker) checkDoExpression(expr lexparse.DoExpression, expectedType TekoType) TekoType {
	var out TekoType = NullType
	blockChecker := NewChecker(c)

	for i, stmt := range expr.Codeblock.Statements {
		if (i == len(expr.Codeblock.Statements)-1) && stmt.Semicolon() == nil {
			switch p := stmt.(type) {
			case lexparse.ExpressionStatement:
				out = blockChecker.checkExpression(p.Expression, expectedType)
			default:
				blockChecker.checkStatement(stmt)
			}
		} else {
			blockChecker.checkStatement(stmt)
		}
	}

	return out
}

func (c *Checker) checkWhileExpression(expr lexparse.WhileExpression) TekoType {
	c.checkExpression(expr.Condition, BoolType)

	return newArrayType(c.checkExpression(expr.Body, nil))
}

func (c *Checker) checkScopeExpression(expr lexparse.ScopeExpression) TekoType {
	blockChecker := NewChecker(c)

	for _, stmt := range expr.Codeblock.Statements {
		blockChecker.checkStatement(stmt)
	}

	return &BasicType{
		fields: blockChecker.ctype.fields,
	}
}

func (c *Checker) checkComparisonExpression(expr lexparse.ComparisonExpression) TekoType {
	rtype := c.checkCallExpression(lexparse.ComparisonCallExpression(expr))

	if !c.isTekoEqType(rtype, IntType) {
		expr.Token().Raise(shared.TypeError, "compare method must return an int")
	}

	return BoolType
}
