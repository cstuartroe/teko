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
		c.Statements = append(
			c.Statements,
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

	case IfT:
		expr = parser.grabIf(prec)

	case LSquareBrT:
		expr = parser.grabArray()

	case SetT:
		expr = parser.grabSet()

	case LCurlyBrT:
		expr = parser.grabObject()

	case FnT:
		expr = parser.grabFunctionDefinition(prec)

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
	case ColonT:
		if !isValidDeclared(expr) {
			TokenPanic(*parser.currentToken(), "Illegal left-hand side to declaration")
		}

		parser.Advance()

		var tekotype Expression = nil
		switch (parser.currentToken().TType) {
		case DefinerT:
		default:
			tekotype = parser.grabExpression(prec)
		}

		setter := *parser.currentToken()
		parser.Expect(DefinerT); parser.Advance()

		out = DeclarationExpression{
			Symbol: expr.Token(),
			Tekotype: tekotype,
			Setter: setter,
			Right: parser.grabExpression(prec),
		}

	case LParT:
		out = parser.makeCallExpression(expr)

	case DotT:
		out = parser.makeAttributeExpression(expr)

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

	case SuffixT:
		suffix := *parser.currentToken()
		parser.Advance()

		if parser.transform {
			out = CallExpression{
				Receiver: AttributeExpression{
					Left: expr,
					Symbol: Token{
						Line:  suffix.Line,
						Col:   suffix.Col,
						TType: SymbolT,
						Value: []rune(suffixes[value]),
					},
				},
			}
		} else {
			out = SuffixExpression{
				Left:   expr,
				Suffix: suffix,
			}
		}

	case UpdaterT:
		if prec <= setter_prec {
			updater := *parser.currentToken()
			parser.Advance()

			if value == "<-" {
				out = UpdateExpression{
					Updated: expr,
					Setter:  updater,
					Right:   parser.grabExpression(setter_prec + 1),
				}

			} else {
				panic("Other updaters not supported yet")
			}
		}
	}

	if out != nil {
		return parser.continueExpression(out, prec)
	} else {
		return expr
	}
}

func isValidDeclared(declared Expression) bool {
	switch (declared).(type) {
	case SimpleExpression:
		switch declared.Token().TType {
		case SymbolT:
			return true
		default:
			return false
		}
	default:
		return false
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

		if parser.currentToken().TType == DefinerT {
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

func (parser *Parser) makeAttributeExpression(left Expression) AttributeExpression {
	parser.Expect(DotT)
	parser.Advance()

	parser.Expect(SymbolT)
	symbol := *parser.currentToken()
	parser.Advance()

	return AttributeExpression{
		Left:   left,
		Symbol: symbol,
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

func (parser *Parser) grabArray() SequenceExpression {
	parser.Expect(LSquareBrT)
	open := *parser.currentToken()
	parser.Advance()

	return SequenceExpression{
		OpenBrace: open,
		Stype:     ArraySeqType,
		Elements:  parser.grabSequence(RSquareBrT),
	}
}

func (parser *Parser) grabSet() SequenceExpression {
	parser.Expect(SetT)
	open := *parser.currentToken()
	parser.Advance()

	parser.Expect(LCurlyBrT)
	parser.Advance()

	return SequenceExpression{
		OpenBrace: open,
		Stype:     SetSeqType,
		Elements:  parser.grabSequence(RCurlyBrT),
	}
}

func (parser *Parser) grabIf(prec int) IfExpression {
	parser.Expect(IfT)
	if_token := *parser.currentToken()
	parser.Advance()

	cond := parser.grabExpression(prec)

	if parser.currentToken().TType == ThenT {
		parser.Advance()
	}

	then := parser.grabExpression(prec)
	var else_expr Expression = nil

	if parser.currentToken().TType == ElseT {
		parser.Advance()
		else_expr = parser.grabExpression(prec)
	}

	return IfExpression{
		If:        if_token,
		Condition: cond,
		Then:      then,
		Else:      else_expr,
	}
}

func (parser *Parser) grabObject() ObjectExpression {
	parser.Expect(LCurlyBrT)
	open := *parser.currentToken()
	parser.Advance()

	fields := []ObjectField{}

	for parser.currentToken().TType != RCurlyBrT {
		fields = append(fields, parser.grabObjectField())

		if parser.currentToken().TType == CommaT {
			parser.Advance()
		} else {
			break
		}
	}

	parser.Expect(RCurlyBrT)
	parser.Advance()

	return ObjectExpression{
		OpenBrace: open,
		Fields:    fields,
	}
}

func (parser *Parser) grabObjectField() ObjectField {
	parser.Expect(SymbolT)
	symbol := *parser.currentToken()
	parser.Advance()

	parser.Expect(ColonT)
	parser.Advance()

	return ObjectField{
		Symbol: symbol,
		Value:  parser.grabExpression(min_prec),
	}
}

func (parser *Parser) grabOptionalType(prec int) Expression {
	if parser.currentToken().TType == ColonT {
		parser.Advance()
		return parser.grabExpression(prec)
	} else {
		return nil
	}
}

func (parser *Parser) grabArgdefs() []ArgdefNode {
	argdefs := []ArgdefNode{}

	for parser.currentToken().TType != RParT {
		parser.Expect(SymbolT)
		symbol := *parser.currentToken()
		parser.Advance()

		tekotype := parser.grabOptionalType(min_prec)

		argdefs = append(argdefs, ArgdefNode{Symbol: symbol, Tekotype: tekotype})

		if parser.currentToken().TType == CommaT {
			parser.Advance()
		} else {
			break
		}
	}

	return argdefs
}

func (parser *Parser) grabFunctionRight(prec int) Expression {
	// TODO: do-block syntactic sugar (also for if statements)

	parser.Expect(ArrowT)
	parser.Advance()

	return parser.grabExpression(prec)
}

func (parser *Parser) grabFunctionDefinition(prec int) FunctionExpression {
	parser.Expect(FnT)
	fn := *parser.currentToken()
	parser.Advance()

	var name *Token = nil
	if parser.currentToken().TType == SymbolT {
		name = parser.currentToken()
		parser.Advance()
	}

	parser.Expect(LParT)
	parser.Advance()

	argdefs := parser.grabArgdefs()

	parser.Expect(RParT)
	parser.Advance()

	rtype := parser.grabOptionalType(prec)
	right := parser.grabFunctionRight(prec)

	return FunctionExpression{
		FnToken: fn,
		Name: name,
		Argdefs: argdefs,
		Rtype:   rtype,
		Right:   right,
	}
}
