package lexparse

import (
	"fmt"
)

var simpleExprTokenTypes map[tokenType]bool = map[tokenType]bool{
	SymbolT: true,
	StringT: true,
	CharT:   true,
	IntT:    true,
	FloatT:  true,
	BoolT:   true,
}

func ParseFile(filename string, transform bool) Codeblock {
	p := Parser{transform: transform}
	p.LexFile(filename)

	c := Codeblock{}
	for p.HasMore() {
		c.statements = append(
			c.statements,
			p.GrabStatement(),
		)
	}

	return c
}

type Parser struct {
	tokens    []Token
	position  int
	transform bool
}

func (parser *Parser) LexFile(filename string) {
	parser.tokens = LexFile(filename)
}

func (parser *Parser) currentToken() *Token {
	if parser.HasMore() {
		return &(parser.tokens[parser.position])
	} else {
		t := parser.tokens[parser.position-1]
		lexparsePanic(t.Line, t.Col+len(t.Value), 1, "Unexpected EOF")
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
		TokenPanic(*parser.currentToken(), fmt.Sprintf("Expected %s", ttype))
	}
}

// I'm attempting to make the teko parser without needing lookahead
// We'll see how sustainable that is

func (parser *Parser) GrabStatement() Statement {
	stmt := parser.grabExpressionStmt()
	parser.Expect(SemicolonT)
	parser.Advance()
	return stmt
}

func (parser *Parser) grabExpressionStmt() ExpressionStatement {
	return ExpressionStatement{
		Expression: parser.grabExpression(min_prec),
	}
}

func (parser *Parser) grabExpression(prec int) Expression {
	var expr Expression

	switch parser.currentToken().TType {
	case LParT:
		expr = parser.grabTuple()
	default:
		expr = parser.grabSimpleExpression()
	}

	return parser.continueExpression(expr, prec)
}

func (parser *Parser) grabSimpleExpression() SimpleExpression {
	t := parser.currentToken()
	if _, ok := simpleExprTokenTypes[t.TType]; ok {
		n := SimpleExpression{token: *t}
		parser.Advance()
		return n
	} else {
		TokenPanic(*parser.currentToken(), "Illegal start to simple expression")
		return SimpleExpression{} // unreachable code that the compiler requires
	}
}

func (parser *Parser) continueExpression(expr Expression, prec int) Expression {
	ttype := parser.currentToken().TType
	value := string(parser.currentToken().Value)

	var out Expression

	switch ttype {
	case SymbolT:
		if prec <= min_prec {
			out = DeclarationExpression{
				Tekotype:  expr,
				Declareds: parser.grabDeclaredChain(),
			}
		}

	case LParT:
		out = parser.makeCallExpression(expr)

	case BinopT:
		op_prec := binop_precs[binops[value]]
		if prec <= op_prec {
			op := *parser.currentToken()
			parser.Advance()
			right := parser.grabExpression(op_prec + 1)

			if parser.transform {
				out = CallExpression{
					Receiver: AttributeExpression{
						Left: expr,
						Symbol: Token{
							Line:  op.Line,
							Col:   op.Col,
							TType: SymbolT,
							Value: []rune(binops[value]),
						},
					},
					Args: []Expression{right},
				}

			} else {
				out = BinopExpression{
					Left:      expr,
					Operation: op,
					Right:     right,
				}
			}
		}
	}

	if out != nil {
		return parser.continueExpression(out, prec)
	} else {
		return expr
	}
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

	parser.Expect(SetterT)
	setter := parser.currentToken()
	parser.Advance()

	right := parser.grabExpression(min_prec)

	return Declared{
		Symbol: *symbol,
		Setter: *setter,
		Right:  right,
	}
}

func (parser *Parser) makeCallExpression(receiver Expression) CallExpression {
	parser.Expect(LParT)
	parser.Advance()

	args := []Expression{}
	kwargs := []FunctionKwarg{}
	on_kwargs := false
	cont := true

	for cont {
		arg := parser.grabExpression(add_sub_prec) // don't want it continuing past a setter!

		if parser.currentToken().TType == SetterT {
			if string(parser.currentToken().Value) != "=" {
				TokenPanic(*parser.currentToken(), "Keyword argument must use =")
			}
			switch p := arg.(type) {
			case SimpleExpression:
				if p.token.TType != SymbolT {
					TokenPanic(p.token, "Left-hand side of keyword argument cannot be a value")
				}
			default:
				TokenPanic(p.Token(), "Left-hand side of keyword argument cannot be a value")
			}

			parser.Advance()
			on_kwargs = true

			kwargs = append(kwargs, FunctionKwarg{
				Symbol: arg.Token(),
				Value:  parser.grabExpression(add_sub_prec),
			})
		} else {
			if on_kwargs {
				TokenPanic(*parser.currentToken(), "All positional arguments must be before all keyword arguments")
			}

			args = append(args, arg)
		}

		if parser.currentToken().TType == CommaT {
			parser.Advance()
			cont = (parser.currentToken().TType != RParT)
		} else {
			cont = false
		}
	}

	parser.Expect(RParT)
	parser.Advance()

	return CallExpression{
		Receiver: receiver,
		Args:     args,
		Kwargs:   kwargs,
	}
}

func (parser *Parser) grabSequence(closingType tokenType) []Expression {
	seq := []Expression{}

	for parser.currentToken().TType != closingType {
		seq = append(seq, parser.grabExpression(min_prec))

		if parser.currentToken().TType == CommaT {
			parser.Advance()
		} else {
			break
		}
	}

	parser.Expect(closingType)
	parser.Advance()

	return seq
}

func (parser *Parser) grabTuple() Expression {
	parser.Expect(LParT)
	lpar := *parser.currentToken()
	parser.Advance()

	seq := parser.grabSequence(RParT)
	switch len(seq) {
	case 0:
		TokenPanic(lpar, "Cannot have empty tuple")
		return nil
	case 1:
		return seq[0]
	default:
		return TupleExpression{
			Elements: seq,
			LPar:     lpar,
		}
	}
}
