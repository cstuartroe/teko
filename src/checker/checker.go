package checker

import (
    "github.com/cstuartroe/teko/src/lexparse"
)

func (c *Codeblock) checkNextStatement() {
  stmt := c.parser.GrabStatement()
  c.statements = append(c.statements, stmt)
  // lexparse.PrintNode(stmt)

  switch p := stmt.(type) {
  case lexparse.ExpressionStatement: c.checkExpressionStatement(p)
  default: lexparse.TokenPanic(stmt.Token(), "Unknown statement type")
  }
}

func (c *Codeblock) checkExpressionStatement(stmt lexparse.ExpressionStatement) {
  c.checkExpression(stmt.Expression)
}

func (c *Codeblock) checkExpression(expr lexparse.Expression) TekoType {
  switch p := expr.(type) {
  case lexparse.SimpleExpression: return c.checkSimpleExpression(p)
  case lexparse.DeclarationExpression: return c.checkDeclaration(p)
  default: panic("Unknown expression type")
  }
}

func (c *Codeblock) checkSimpleExpression(expr lexparse.SimpleExpression) TekoType {
  t := expr.Token()
  switch t.TType {
  case lexparse.SymbolT:
    ttype := c.getFieldType(string(t.Value))
    if ttype == nil {
      lexparse.TokenPanic(t, "Undefined variable"); return nil
    } else {
      return ttype
    }

  case lexparse.IntT: return &IntType
  case lexparse.BoolT: return &BoolType
  default: panic("Unknown simple expression type")
  }
}

func (c *Codeblock) checkDeclaration(decl lexparse.DeclarationExpression) TekoType {
  var tekotype TekoType = c.evaluateType(decl.Tekotype)

  for _, declared := range decl.Declareds {
    c.declare(declared, tekotype)
  }

  return nil // TODO: decide what type declarations return and return it
}

func (c *Codeblock) evaluateType(expr lexparse.Expression) TekoType {
  switch p := expr.(type) {
  case lexparse.SimpleExpression: return c.evaluateNamedType(p)
  default: panic("Unknown type format!")
  }
}

func (c *Codeblock) evaluateNamedType(expr lexparse.SimpleExpression) TekoType {
  if expr.Token().TType != lexparse.SymbolT {
    lexparse.TokenPanic(expr.Token(), "Invalid type expression")
    return nil
  }

  return c.getTypeByName(string(expr.Token().Value))
}

func (c *Codeblock) declare(declared lexparse.Declared, tekotype TekoType) TekoType {
  // TODO: function types

  right_tekotype := c.checkExpression(declared.Right)

  if !isTekoEqType(tekotype, right_tekotype) {
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
