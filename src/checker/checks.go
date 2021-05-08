package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
)

func (c *Checker) checkStatement(stmt lexparse.Statement) {
	switch p := stmt.(type) {
	case lexparse.ExpressionStatement:
		c.checkExpressionStatement(p)
	default:
		lexparse.TokenPanic(stmt.Token(), "Unknown statement type")
	}
}

func (c *Checker) checkExpressionStatement(stmt lexparse.ExpressionStatement) {
	c.checkExpression(stmt.Expression, nil)
}

func (c *Checker) checkExpression(expr lexparse.Expression, expectedType TekoType) TekoType {
	var ttype TekoType

	switch p := expr.(type) {

	case lexparse.SimpleExpression:
		ttype = c.checkSimpleExpression(p)

	case lexparse.DeclarationExpression:
		ttype = c.checkDeclaration(p)

	case lexparse.CallExpression:
		ttype = c.checkCallExpression(p)

	case lexparse.AttributeExpression:
		ttype = c.checkAttributeExpression(p)

	case lexparse.IfExpression:
		ttype = c.checkIfExpression(p)

	case lexparse.SequenceExpression:
		ttype = c.checkSequenceExpression(p, expectedType)

	case lexparse.ObjectExpression:
		ttype = c.checkObjectExpression(p)

	default:
		lexparse.TokenPanic(expr.Token(), "Cannot typecheck expression type: "+expr.Ntype())
	}

	if ttype == nil {
		lexparse.TokenPanic(expr.Token(), "Evaluated to nil type")
	}
	if (expectedType != nil) && !isTekoEqType(ttype, expectedType) {
		lexparse.TokenPanic(expr.Token(), "Incorrect type")
	}
	return ttype
}

func (c *Checker) checkSimpleExpression(expr lexparse.SimpleExpression) TekoType {
	t := expr.Token()
	switch t.TType {
	case lexparse.SymbolT:
		ttype := c.getFieldType(string(t.Value))
		if ttype == nil {
			lexparse.TokenPanic(t, "Undefined variable")
			return nil
		} else {
			return ttype
		}

	case lexparse.IntT:
		return &IntType
	case lexparse.BoolT:
		return &BoolType
	case lexparse.StringT:
		return &StringType
	case lexparse.CharT:
		return &CharType
	default:
		panic("Unknown simple expression type")
	}
}

func (c *Checker) checkDeclaration(decl lexparse.DeclarationExpression) TekoType {
	var tekotype TekoType = c.evaluateType(decl.Tekotype)

	for _, declared := range decl.Declareds {
		c.declare(declared, tekotype)
	}

	return tekotype // TODO: decide what type declarations return and return it
}

func (c *Checker) evaluateType(expr lexparse.Expression) TekoType {
	switch p := expr.(type) {
	case lexparse.SimpleExpression:
		return c.evaluateNamedType(p)
	default:
		panic("Unknown type format!")
	}
}

func (c *Checker) evaluateNamedType(expr lexparse.SimpleExpression) TekoType {
	if expr.Token().TType != lexparse.SymbolT {
		lexparse.TokenPanic(expr.Token(), "Invalid type expression")
		return nil
	}

	return c.getTypeByName(string(expr.Token().Value))
}

func (c *Checker) declare(declared lexparse.Declared, tekotype TekoType) TekoType {
	// TODO: function types

	c.checkExpression(declared.Right, tekotype)

	name := string(declared.Symbol.Value)
	if c.getFieldType(name) == nil {
		c.declareFieldType(name, tekotype)
	} else {
		lexparse.TokenPanic(declared.Symbol, "Name has already been declared")
	}

	return tekotype
}

func (c *Checker) checkCallExpression(expr lexparse.CallExpression) TekoType {
	receiver_tekotype := c.checkExpression(expr.Receiver, nil)

	switch ftype := receiver_tekotype.(type) {
	case *FunctionType:
		args_by_name := ResolveArgs(ftype.argdefs, expr)

		for _, argdef := range ftype.argdefs {
			arg, ok := args_by_name[argdef.name]
			if !ok {
				lexparse.TokenPanic(
					arg.Token(),
					"Argument was not passed: "+argdef.name,
				)
			}

			c.checkExpression(arg, argdef.ttype)
		}

		return ftype.rtype

	case FunctionType:
		panic("Use *FunctionType instead of FunctionType")

	default:
		lexparse.TokenPanic(expr.Token(), "Expression does not have a function type")
		return nil
	}
}

func (c *Checker) checkAttributeExpression(expr lexparse.AttributeExpression) TekoType {
	left_tekotype := c.checkExpression(expr.Left, nil)

	tekotype := getField(left_tekotype, string(expr.Symbol.Value))
	if tekotype != nil {
		return tekotype
	} else {
		lexparse.TokenPanic(expr.Symbol, "No such field: "+string(expr.Symbol.Value))
		return nil
	}
}

func (c *Checker) checkIfExpression(expr lexparse.IfExpression) TekoType {
	c.checkExpression(expr.Condition, &BoolType)

	then_tekotype := c.checkExpression(expr.Then, nil)
	else_tekotype := c.checkExpression(expr.Else, nil)

	if then_tekotype != else_tekotype {
		lexparse.TokenPanic(expr.Else.Token(), "Then and else blocks have mismatching types")
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
			lexparse.TokenPanic(expr.Token(), "With no expected type, sequence cannot be empty")
		} else {
			etype = c.checkExpression(expr.Elements[0], nil)

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

func (c *Checker) checkObjectExpression(expr lexparse.ObjectExpression) TekoType {
	fields := map[string]TekoType{}

	for _, of := range expr.Fields {
		if _, ok := fields[string(of.Symbol.Value)]; ok {
			lexparse.TokenPanic(of.Symbol, "Duplicate member")
		}

		fields[string(of.Symbol.Value)] = c.checkExpression(of.Value, nil)
	}

	return BasicType{fields}
}
