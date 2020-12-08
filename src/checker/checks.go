package checker

import (
	"fmt"
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
	c.checkExpression(stmt.Expression)
}

func (c *Checker) checkExpression(expr lexparse.Expression) TekoType {
	switch p := expr.(type) {
	case lexparse.SimpleExpression:
		return c.checkSimpleExpression(p)
	case lexparse.DeclarationExpression:
		return c.checkDeclaration(p)
	case lexparse.CallExpression:
		return c.checkCallExpression(p)
	default:
		panic("Unknown expression type")
	}
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

	return nil // TODO: decide what type declarations return and return it
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

	right_tekotype := c.checkExpression(declared.Right)

	if !isTekoSubtype(right_tekotype, tekotype) {
		lexparse.TokenPanic(declared.Right.Token(), "Incorrect type")
		return nil
	}

	name := string(declared.Symbol.Value)
	if c.getFieldType(name) == nil {
		c.declareFieldType(name, tekotype)
	} else {
		lexparse.TokenPanic(declared.Symbol, "Name has already been declared")
	}

	return tekotype
}

func getArgTypeByName(argdefs []FunctionArgDef, name string) TekoType {
	for _, argdef := range argdefs {
		if argdef.name == name {
			return argdef.ttype
		}
	}
	return nil
}

func (c *Checker) checkCallExpression(expr lexparse.CallExpression) TekoType {
	receiver_tekotype := c.checkExpression(expr.Receiver)

	switch ftype := receiver_tekotype.(type) {
	case *FunctionType:
		args_by_name := map[string]lexparse.Expression{}

		if len(expr.Args) > len(ftype.argdefs) {
			lexparse.TokenPanic(
				expr.Args[len(ftype.argdefs)].Token(),
				fmt.Sprintf("Too many arguments (%d expected, %d given)", len(ftype.argdefs), len(expr.Args)),
			)
		}

		for i, arg := range expr.Args {
			args_by_name[ftype.argdefs[i].name] = arg
		}

		for _, kwarg := range expr.Kwargs {
			name := string(kwarg.Symbol.Value)

			if _, ok := args_by_name[name]; ok {
				lexparse.TokenPanic(
					kwarg.Token(),
					fmt.Sprintf("Argument already passed: %s", name),
				)
			}

			argtype := getArgTypeByName(ftype.argdefs, name)
			if argtype != nil {
				args_by_name[name] = kwarg.Value
			} else {
				lexparse.TokenPanic(
					kwarg.Symbol,
					fmt.Sprintf("Function doesn't take argument %s", name),
				)
			}
		}

		for _, argdef := range ftype.argdefs {
			arg, ok := args_by_name[argdef.name]
			if !ok {
				lexparse.TokenPanic(
					arg.Token(),
					"Argument was not passed: "+argdef.name,
				)
			}

			if !isTekoEqType(c.checkExpression(arg), argdef.ttype) {
				lexparse.TokenPanic(
					arg.Token(),
					"Incorrect argument type",
				)
			}
		}

		return ftype.rtype

	default:
		lexparse.TokenPanic(expr.Token(), "Expression does not have a function type")
		return nil
	}
}
