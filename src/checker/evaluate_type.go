package checker

import (
	"strconv"

	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
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

	case lexparse.VarExpression:
		return c.evaluateVarType(p)

	case lexparse.SliceExpression:
		return c.evaluateSliceType(p)

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
			expr.Token().Raise(shared.NameError, "No type called "+string(expr.Token().Value))
		} else {
			return ttype
		}

	case lexparse.IntT:
		return c.evaluateConstantIntType(expr.Token())

	case lexparse.StringT:
		return newConstantStringType(expr.Token().Value)

	// TODO: bool, chars and floats

	default:
		expr.Token().Raise(shared.TypeError, "Invalid type expression")
	}

	return nil
}

func (c *Checker) evaluateObjectType(expr lexparse.ObjectExpression) *BasicType {
	out := newBasicType("") // TODO: Do we need object type name? If so, set it.

	for _, field := range expr.Fields {
		out.setField(string(field.Symbol.Value), c.evaluateType(field.Value))
	}

	return out
}

func (c *Checker) evaluateUnionType(expr lexparse.BinopExpression) TekoType {
	if string(expr.Operation.Value) != "|" {
		expr.Operation.Raise(shared.SyntaxError, "Invalid type expression")
	}

	// TODO raise on something like `int | int`

	return c.unionTypes(
		c.evaluateType(expr.Left),
		c.evaluateType(expr.Right),
	)
}

func (c *Checker) evaluateVarType(expr lexparse.VarExpression) *VarType {
	return newVarType(c.evaluateType(expr.Right))
}

func (c *Checker) evaluateSliceType(expr lexparse.SliceExpression) TekoType {
	if expr.Inside != nil {
		panic("Still need to implement fixed-size arrays") // TODO
	}

	return newArrayType(c.evaluateType(expr.Left))
}
