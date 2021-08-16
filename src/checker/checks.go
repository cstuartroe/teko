package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

func (c *Checker) checkStatement(stmt lexparse.Statement) {
	switch p := stmt.(type) {
	case lexparse.ExpressionStatement:
		c.checkExpressionStatement(p)

	case lexparse.TypeStatement:
		c.checkTypeStatement(p)

	default:
		stmt.Token().Raise(lexparse.NotImplementedError, "Unknown statement type")
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
	var ttype TekoType

	switch p := expr.(type) {

	case lexparse.SimpleExpression:
		ttype = c.checkSimpleExpression(p, expectedType)

	case lexparse.DeclarationExpression:
		ttype = c.checkDeclaration(p, expectedType)

	case lexparse.UpdateExpression:
		ttype = c.checkUpdate(p)

	case lexparse.CallExpression:
		ttype = c.checkCallExpression(p)

	case lexparse.AttributeExpression:
		ttype = c.checkAttributeExpression(p)

	case lexparse.IfExpression:
		ttype = c.checkIfExpression(p, expectedType)

	case lexparse.SequenceExpression:
		ttype = c.checkSequenceExpression(p, expectedType)

	case lexparse.ObjectExpression:
		ttype = c.checkObjectExpression(p, expectedType)

	case lexparse.FunctionExpression:
		// TODO use expected type in function definitions
		ttype = c.checkFunctionDefinition(p)

	case lexparse.DoExpression:
		ttype = c.checkDoExpression(p, expectedType)

	default:
		expr.Token().Raise(lexparse.NotImplementedError, "Cannot typecheck expression type")
	}

	if ttype == nil {
		expr.Token().Raise(lexparse.TypeError, "Evaluated to nil type")
	}
	if (expectedType != nil) && !isTekoSubtype(ttype, expectedType) {
		expr.Token().Raise(lexparse.TypeError, "Actual type does not fulfill expected type")
	}
	if expectedType == nil {
		return ttype
	} else {
		return expectedType
	}
}

func (c *Checker) checkSimpleExpression(expr lexparse.SimpleExpression, expectedType TekoType) TekoType {
	t := expr.Token()
	switch t.TType {
	case lexparse.SymbolT:
		ttype := c.getFieldType(string(t.Value))
		if ttype == nil {
			t.Raise(lexparse.NameError, "Undefined variable")
			return nil
		} else {
			return ttype
		}

	// TODO bool, char, and float constant types

	case lexparse.IntT:
		if expectedType != nil && isTekoEqType(expectedType, IntType) {
			return IntType
		} else {
			return c.evaluateConstantIntType(t)
		}

	case lexparse.BoolT:
		return BoolType

	case lexparse.StringT:
		if expectedType != nil && isTekoEqType(expectedType, StringType) {
			return StringType
		} else {
			return newConstantStringType(t.Value)
		}

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
	// TODO: function types

	evaluated_ttype := c.checkExpression(right, tekotype)

	c.declareFieldType(symbol, evaluated_ttype)

	return evaluated_ttype
}

func validUpdated(updated lexparse.Expression) bool {
	switch p := updated.(type) {
	case lexparse.SimpleExpression:
		return p.Token().TType == lexparse.SymbolT
	case lexparse.AttributeExpression:
		return validUpdated(p.Left)
	default:
		return false
	}
}

func (c *Checker) checkUpdate(update lexparse.UpdateExpression) TekoType {
	if !validUpdated(update.Updated) {
		update.Updated.Token().Raise(lexparse.SyntaxError, "Invalid left hand side")
	}

	ttype := c.checkExpression(update.Updated, nil)
	c.checkExpression(update.Right, ttype)
	return ttype
}

func (c *Checker) checkCallExpression(expr lexparse.CallExpression) TekoType {
	receiver_tekotype := c.checkExpression(expr.Receiver, nil)

	switch ftype := receiver_tekotype.(type) {
	case *FunctionType:
		args_by_name := ResolveArgs(ftype.argnames(), expr)

		for _, argdef := range ftype.argdefs {
			arg, ok := args_by_name[argdef.name]
			if !ok {
				arg.Token().Raise(lexparse.ArgumentError, "Argument was not passed: "+argdef.name)
			}

			c.checkExpression(arg, argdef.ttype)
		}

		return ftype.rtype

	case FunctionType:
		panic("Use *FunctionType instead of FunctionType")

	default:
		expr.Token().Raise(lexparse.TypeError, "Expression does not have a function type")
		return nil
	}
}

func (c *Checker) checkAttributeExpression(expr lexparse.AttributeExpression) TekoType {
	left_tekotype := c.checkExpression(expr.Left, nil)

	tekotype := getFieldSafe(left_tekotype, string(expr.Symbol.Value))
	if tekotype != nil {
		return tekotype
	} else {
		expr.Symbol.Raise(lexparse.NameError, "No such field: "+string(expr.Symbol.Value))
		return nil
	}
}

func (c *Checker) checkIfExpression(expr lexparse.IfExpression, expectedType TekoType) TekoType {
	c.checkExpression(expr.Condition, BoolType)

	then_tekotype := c.checkExpression(expr.Then, expectedType)
	else_tekotype := c.checkExpression(expr.Else, expectedType)

	if then_tekotype != else_tekotype {
		expr.Else.Token().Raise(lexparse.TypeError, "Then and else blocks have mismatching types")
	}

	return then_tekotype
}

func deconstantize(ttype TekoType) TekoType {
	switch p := ttype.(type) {
	case *ConstantType:
		switch p.ctype {
		case IntConstant:
			return IntType
		case StringConstant:
			return StringType
		default:
			panic("Unknown constant type")
		}
	default:
		return ttype
	}
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
			expr.Token().Raise(lexparse.TypeError, "With no expected type, sequence cannot be empty")
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

func (c *Checker) checkObjectExpression(expr lexparse.ObjectExpression, expectedType TekoType) TekoType {
	var expectedObjType ObjectType = nil

	switch p := expectedType.(type) {
	case ObjectType:
		expectedObjType = p
	case nil:
		expectedType = VoidType
	default:
		expr.Token().Raise(lexparse.TypeError, "Did not expect object type")
	}

	fields := map[string]TekoType{}

	for _, of := range expr.Fields {
		field_name := string(of.Symbol.Value)

		if _, ok := fields[field_name]; ok {
			of.Symbol.Raise(lexparse.NameError, "Duplicate member")
		}

		fields[field_name] = c.checkExpression(of.Value, getField(expectedObjType, field_name))
	}

	return &BasicType{fields}
}

func (c *Checker) checkFunctionDefinition(expr lexparse.FunctionExpression) TekoType {
	blockChecker := NewChecker(c)

	argdefs := []FunctionArgDef{}

	for _, ad := range expr.Argdefs {
		if ad.Tekotype == nil {
			ad.Symbol.Raise(lexparse.NotImplementedError, "Type inference for function arguments not yet implemented")
		}

		ttype := c.evaluateType(ad.Tekotype)
		blockChecker.declareFieldType(ad.Symbol, ttype)

		// TODO: get argdefs from blockChecker
		argdefs = append(argdefs, FunctionArgDef{
			name:  string(ad.Symbol.Value),
			ttype: ttype,
		})
	}

	var rtype TekoType = nil
	if expr.Rtype != nil {
		rtype = c.evaluateType(expr.Rtype)
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
	var out TekoType = VoidType
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
