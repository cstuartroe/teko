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
	case *lexparse.SimpleExpression:
		return c.evaluateSimpleType(p)

	case *lexparse.ObjectExpression:
		return c.evaluateObjectType(p)

	case *lexparse.BinopExpression:
		return c.evaluateUnionType(p)

	case *lexparse.VarExpression:
		return c.evaluateVarType(p)

	case *lexparse.SliceExpression:
		return c.evaluateSliceType(p)

	case *lexparse.MapExpression:
		return c.evaluateMapType(p)

	case *lexparse.FunctionExpression:
		return c.evaluateFunctionType(p)

	default:
		panic("Unknown type format!")
	}
}

func (c *Checker) evaluateConstantIntType(token *lexparse.Token) TekoType {
	n, err := strconv.Atoi(string(token.Value))
	if err == nil {
		return NewConstantIntType(n)
	} else {
		panic("Invalid int")
	}
}

func (c *Checker) evaluateSimpleType(expr *lexparse.SimpleExpression) TekoType {
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
		return NewConstantStringType(expr.Token().Value)

	case lexparse.BoolT:
		b := false
		if string(expr.Token().Value) == "true" {
			b = true
		}

		return ConstantBoolTypeCache[b]

	// TODO: chars and floats

	default:
		expr.Token().Raise(shared.TypeError, "Invalid type expression")
	}

	return nil
}

func (c *Checker) evaluateObjectType(expr *lexparse.ObjectExpression) *BasicType {
	out := newBasicType("") // TODO: Do we need object type name? If so, set it.

	for _, field := range expr.Fields {
		out.setField(string(field.Symbol.Value), c.evaluateType(field.Value))
	}

	return out
}

func (c *Checker) evaluateUnionType(expr *lexparse.BinopExpression) TekoType {
	if string(expr.Operation.Value) != "|" {
		expr.Operation.Raise(shared.SyntaxError, "Invalid type expression")
	}

	// TODO raise on something like `int | int`

	return c.unionTypes(
		c.evaluateType(expr.Left),
		c.evaluateType(expr.Right),
	)
}

func (c *Checker) evaluateVarType(expr *lexparse.VarExpression) *VarType {
	return newVarType(c.evaluateType(expr.Right))
}

func (c *Checker) evaluateSliceType(expr *lexparse.SliceExpression) TekoType {
	if expr.Inside != nil {
		panic("Still need to implement fixed-size arrays") // TODO
	}

	return newArrayType(c.evaluateType(expr.Left))
}

func (c *Checker) evaluateMapType(expr *lexparse.MapExpression) TekoType {
	if expr.Ktype == nil || expr.Vtype == nil {
		expr.Token().Raise(shared.SyntaxError, "Map type must specify both key and value type")
	} else if expr.HasBraces {
		expr.Token().Raise(shared.SyntaxError, "Map type cannot have contents")
	}

	return newMapType(c.evaluateType(expr.Ktype), c.evaluateType(expr.Vtype))
}

func (c *Checker) evaluateFunctionType(expr *lexparse.FunctionExpression) TekoType {
	if expr.Name != nil {
		expr.Name.Raise(shared.SyntaxError, "Function type cannot have name")
	}

	if expr.Right != nil {
		expr.Right.Token().Raise(shared.SyntaxError, "Function type cannot have body")
	}

	// TODO: take defaults and GDL into account

	argdefs := []FunctionArgDef{}

	for _, ad := range expr.Argdefs {
		argdefs = append(argdefs, FunctionArgDef{
			Name:  string(ad.Symbol.Value),
			ttype: c.evaluateType(ad.Tekotype),
		})
	}

	return &FunctionType{
		rtype:   c.evaluateType(expr.Rtype),
		argdefs: argdefs,
	}
}
