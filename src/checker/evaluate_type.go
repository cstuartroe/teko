package checker

import (
	"github.com/cstuartroe/teko/src/lexparse"
	"strconv"
)

func (c *Checker) evaluateType(expr lexparse.Expression) TekoType {
	if expr == nil {
		return nil
	}

	switch p := (expr).(type) {
	case lexparse.SimpleExpression:
		return c.evaluateSimpleType(p)

	case lexparse.ObjectExpression:
		return c.evaluateObjectType(p)

	case lexparse.BinopExpression:
		return c.evaluateUnionType(p)

	default:
		panic("Unknown type format!")
	}
}

func (c *Checker) evaluateConstantIntType(token lexparse.Token) TekoType {
	n, err := strconv.Atoi(string(token.Value))
	if err == nil {
		return newConstantIntType(n)
	} else {
		panic("Invalid int")
	}
}

func (c *Checker) evaluateSimpleType(expr lexparse.SimpleExpression) TekoType {
	switch expr.Token().TType {

	case lexparse.SymbolT:
		ttype := c.getTypeByName(string(expr.Token().Value))
		if ttype == nil {
			expr.Token().Raise(lexparse.NameError, "No type called "+string(expr.Token().Value))
		} else {
			return ttype
		}

	case lexparse.IntT:
		return c.evaluateConstantIntType(expr.Token())

	case lexparse.StringT:
		return newConstantStringType(expr.Token().Value)

	// TODO: bool, chars and floats

	default:
		expr.Token().Raise(lexparse.TypeError, "Invalid type expression")
	}

	return nil
}

func (c *Checker) evaluateObjectType(expr lexparse.ObjectExpression) ObjectType {
	out := newBasicType("")

	for _, field := range expr.Fields {
		out.setField(string(field.Symbol.Value), c.evaluateType(field.Value))
	}

	return out
}

func (c *Checker) evaluateUnionType(expr lexparse.BinopExpression) *UnionType {
	if string(expr.Operation.Value) != "|" {
		expr.Operation.Raise(lexparse.SyntaxError, "Invalid type expression")
	}

	return unionTypes(
		c.evaluateType(expr.Left),
		c.evaluateType(expr.Right),
	)
}
