package lexparse

import (
  "fmt"
)

var simpleExprTokenTypes map[tokenType]bool = map[tokenType]bool {
  SymbolT: true,
  StringT: true,
  CharT: true,
  IntT: true,
  FloatT: true,
  BoolT: true,
  NullT: true,
}

type Parser struct {
  tokens []Token
  position int
}

func (parser *Parser) LexFile(filename string) {
  parser.tokens = LexFile(filename)
}

func (parser *Parser) currentToken() *Token {
  if parser.HasMore() {
    return &(parser.tokens[parser.position])
  } else {
    t := parser.tokens[parser.position-1]
    lexparsePanic(t.Line, t.Col + len(t.Value), 1, "Unexpected EOF")
    return nil
  }
}

func (parser *Parser) Advance() {
  parser.position++
}

func (parser *Parser) HasMore() bool {
  return parser.position < len(parser.tokens)
}

func (parser *Parser) Expect(ttype tokenType) {
  if parser.currentToken().TType != ttype {
    tokenPanic(*parser.currentToken(), fmt.Sprintf("Expected %s", ttype))
  }
}

// I'm attempting to make the teko parser without needing lookahead
// We'll see how sustainable that is

func (parser *Parser) GrabStatement() Statement {
  return parser.grabExpressionStmt()
}

func (parser *Parser) grabExpressionStmt() ExpressionStatement {
  return ExpressionStatement{
    expression: parser.grabExpression(add_sub_prec),
  }
}

type precedence int

const (
  add_sub_prec precedence = iota
  mult_div_prec
  exp_prec
  max_prec
)

func (parser *Parser) grabExpression(prec precedence) Expression {
  expr := parser.grabSimpleExpression()

  return parser.continueExpression(expr, prec)
}

func (parser *Parser) grabSimpleExpression() SimpleExpression {
  t := parser.currentToken()
  if _, ok := simpleExprTokenTypes[t.TType]; ok {
    n := SimpleExpression{Token: *t}
    parser.Advance()
    return n
  } else {
    tokenPanic(*parser.currentToken(), "Illegal start to simple expression")
    return SimpleExpression{} // unreachable code that the compiler requires
  }
}

func (parser *Parser) continueExpression(expr Expression, prec precedence) Expression {
  if parser.currentToken().TType == SymbolT && prec < max_prec {
    return DeclarationExpression{
      tekotype: expr,
      declareds: parser.grabDeclaredChain(),
    }
  }

  return expr
}

func (parser *Parser) grabDeclaredChain() []Declared {
  chain := []Declared{
    parser.grabDeclared(),
  }

  cont := true
  for cont && parser.HasMore() && parser.currentToken().TType == CommaT {
    parser.Advance()
    if parser.HasMore() && parser.currentToken().TType == SymbolT {
      chain = append(chain, parser.grabDeclared())
    } else {
      cont = false
    }
  }

  return chain
}

func (parser *Parser) grabDeclared() Declared {
  parser.Expect(SymbolT)
  symbol := parser.currentToken()
  parser.Advance()

  // TODO: function argdef

  parser.Expect(SetterT);
  setter := parser.currentToken()
  parser.Advance()

  right := parser.grabExpression(add_sub_prec)

  return Declared{
    symbol: *symbol,
    setter: *setter,
    right: right,
  }
}
