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
			p.grabStatement(),
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
		lexparsePanic(t.Line, t.Col+len(t.Value), 1, SyntaxError, "Unexpected EOF")
		return nil
	}
}

func (parser *Parser) advance() {
	parser.position++
}

func (parser *Parser) HasMore() bool {
	return parser.position < len(parser.tokens)
}

func (parser *Parser) expect(ttype tokenType) *Token {
	token := parser.currentToken()

	if token.TType != ttype {
		token.Raise(SyntaxError, fmt.Sprintf("Expected %s", ttype))
	}

	parser.advance()

	return token
}

// I'm attempting to make the teko parser without needing lookahead
// We'll see how sustainable that is

func (parser *Parser) optionalSemicolon() *Token {
	if parser.HasMore() && parser.currentToken().TType == SemicolonT {
		return parser.expect(SemicolonT)
	} else {
		return nil
	}
}

func (parser *Parser) grabStatement() Statement {
	if parser.currentToken().TType == TypeT {
		return parser.grabTypeStatement()
	} else {
		return parser.grabExpressionStmt()
	}
}

func (parser *Parser) grabTypeStatement() TypeStatement {
	tt := *parser.expect(TypeT)
	name := *parser.expect(SymbolT)
	parser.expect(EqualT)

	return TypeStatement{
		TypeToken:      tt,
		Name:           name,
		TypeExpression: parser.grabTypeExpression(min_prec),
		semicolon:      parser.optionalSemicolon(),
	}
}

func (parser *Parser) grabExpressionStmt() ExpressionStatement {
	return ExpressionStatement{
		Expression: parser.grabExpression(min_prec),
		semicolon:  parser.optionalSemicolon(),
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

	case DoT:
		expr = parser.grabDoExpression()

	case VarT:
		expr = parser.grabVarExpression(prec)

	case WhileT:
		expr = parser.grabWhileExpression(prec)

	default:
		expr = parser.grabSimpleExpression()
	}

	return parser.continueExpression(expr, prec)
}

func (parser *Parser) grabTypeExpression(prec int) Expression {
	if prec < add_sub_prec {
		prec = add_sub_prec
	}

	transform := parser.transform
	parser.transform = false
	out := parser.grabExpression(prec)
	parser.transform = transform
	return out
}

func (parser *Parser) grabSimpleExpression() SimpleExpression {
	t := parser.currentToken()
	if _, ok := simpleExprTokenTypes[t.TType]; ok {
		n := SimpleExpression{token: *t}
		parser.advance()
		return n
	} else {
		parser.currentToken().Raise(SyntaxError, "Illegal start to expression")
		return SimpleExpression{} // unreachable code that the compiler requires
	}
}

func (parser *Parser) continueExpression(expr Expression, prec int) Expression {
	if !parser.HasMore() {
		return expr
	}

	ttype := parser.currentToken().TType
	value := string(parser.currentToken().Value)

	var out Expression

	switch ttype {
	case ColonT:
		if !isValidDeclared(expr) {
			parser.currentToken().Raise(SyntaxError, "Illegal left-hand side to declaration")
		}

		parser.advance()

		var tekotype Expression = nil
		switch parser.currentToken().TType {
		case EqualT:
		default:
			tekotype = parser.grabTypeExpression(prec)
		}

		setter := *parser.expect(EqualT)

		out = DeclarationExpression{
			Symbol:   expr.Token(),
			Tekotype: tekotype,
			Setter:   setter,
			Right:    parser.grabExpression(prec),
		}

	case LParT:
		out = parser.makeCallExpression(expr)

	case DotT:
		out = parser.makeAttributeExpression(expr)

	case BinopT:
		op_prec := binop_precs[binops[value]]
		if prec <= op_prec {
			op := *parser.currentToken()
			parser.advance()
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
		parser.advance()

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
			panic("Other updaters not supported yet")
		}

	case EqualT:
		if prec <= setter_prec {
			setter := *parser.expect(EqualT)

			out = CallExpression{
				Receiver: AttributeExpression{
					Left: expr,
					Symbol: Token{
						Line: setter.Line,
						Col:  setter.Col,
						TType: SymbolT,
						Value: []rune("="),
					},
				},
				Args: []Expression{
					parser.grabExpression(setter_prec + 1),
				},
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
	parser.expect(LParT)

	args := []Expression{}
	kwargs := []FunctionKwarg{}
	on_kwargs := false
	cont := parser.currentToken().TType != RParT

	for cont {
		arg := parser.grabExpression(add_sub_prec) // don't want it continuing past a setter!

		if parser.currentToken().TType == EqualT {
			if string(parser.currentToken().Value) != "=" {
				parser.currentToken().Raise(SyntaxError, "Keyword argument must use =")
			}
			switch p := arg.(type) {
			case SimpleExpression:
				if p.token.TType != SymbolT {
					p.token.Raise(SyntaxError, "Left-hand side of keyword argument cannot be a value")
				}
			default:
				p.Token().Raise(SyntaxError, "Left-hand side of keyword argument cannot be a value")
			}

			parser.advance()
			on_kwargs = true

			kwargs = append(kwargs, FunctionKwarg{
				Symbol: arg.Token(),
				Value:  parser.grabExpression(add_sub_prec),
			})
		} else {
			if on_kwargs {
				parser.currentToken().Raise(SyntaxError, "All positional arguments must be before all keyword arguments")
			}

			args = append(args, arg)
		}

		if parser.currentToken().TType == CommaT {
			parser.advance()
			cont = (parser.currentToken().TType != RParT)
		} else {
			cont = false
		}
	}

	parser.expect(RParT)

	return CallExpression{
		Receiver: receiver,
		Args:     args,
		Kwargs:   kwargs,
	}
}

func (parser *Parser) makeAttributeExpression(left Expression) AttributeExpression {
	parser.expect(DotT)

	symbol := *parser.expect(SymbolT)

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
			parser.advance()
		} else {
			break
		}
	}

	parser.expect(closingType)

	return seq
}

func (parser *Parser) grabTuple() Expression {
	lpar := *parser.expect(LParT)

	seq := parser.grabSequence(RParT)
	switch len(seq) {
	case 0:
		lpar.Raise(SyntaxError, "Cannot have empty tuple")
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
	open := *parser.expect(LSquareBrT)

	return SequenceExpression{
		OpenBrace: open,
		Stype:     ArraySeqType,
		Elements:  parser.grabSequence(RSquareBrT),
	}
}

func (parser *Parser) grabSet() SequenceExpression {
	open := *parser.expect(SetT)

	parser.expect(LCurlyBrT)

	return SequenceExpression{
		OpenBrace: open,
		Stype:     SetSeqType,
		Elements:  parser.grabSequence(RCurlyBrT),
	}
}

func (parser *Parser) grabIf(prec int) IfExpression {
	if_token := *parser.expect(IfT)

	cond := parser.grabExpression(prec)

	if parser.currentToken().TType == ThenT {
		parser.advance()
	}

	then := parser.grabExpression(prec)
	var else_expr Expression = nil

	if parser.currentToken().TType == ElseT {
		parser.advance()
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
	open := *parser.expect(LCurlyBrT)

	fields := []ObjectField{}

	for parser.currentToken().TType != RCurlyBrT {
		fields = append(fields, parser.grabObjectField())

		if parser.currentToken().TType == CommaT {
			parser.advance()
		} else {
			break
		}
	}

	parser.expect(RCurlyBrT)

	return ObjectExpression{
		OpenBrace: open,
		Fields:    fields,
	}
}

func (parser *Parser) grabObjectField() ObjectField {
	symbol := *parser.expect(SymbolT)
	var value Expression

	if parser.currentToken().TType == ColonT {
		parser.advance()

		value = parser.grabExpression(min_prec)
	} else if parser.transform {
		value = SimpleExpression{symbol}
	} else {
		symbol.Raise(SyntaxError, "object property shorthand cannot be used in no-transform context")
	}

	return ObjectField{
		Symbol: symbol,
		Value:  value,
	}
}

func (parser *Parser) grabOptionalType(prec int) Expression {
	if parser.currentToken().TType == ColonT {
		parser.advance()
		return parser.grabTypeExpression(prec)
	} else {
		return nil
	}
}

func (parser *Parser) grabArgdefs() []ArgdefNode {
	argdefs := []ArgdefNode{}

	for parser.currentToken().TType != RParT {
		symbol := *parser.expect(SymbolT)

		tekotype := parser.grabOptionalType(min_prec)

		argdefs = append(argdefs, ArgdefNode{Symbol: symbol, Tekotype: tekotype})

		if parser.currentToken().TType == CommaT {
			parser.advance()
		} else {
			break
		}
	}

	return argdefs
}

func (parser *Parser) grabFunctionRight(prec int) Expression {
	if parser.currentToken().TType == LCurlyBrT {
		return DoExpression{
			DoToken: nil,
			Codeblock: parser.grabCodeblock(),
		}
	} else {
		parser.expect(ArrowT)

		return parser.grabExpression(prec)
	}
}

func (parser *Parser) grabFunctionDefinition(prec int) FunctionExpression {
	fn := *parser.expect(FnT)

	var name *Token = nil
	if parser.currentToken().TType == SymbolT {
		name = parser.currentToken()
		parser.advance()
	}

	parser.expect(LParT)

	argdefs := parser.grabArgdefs()

	parser.expect(RParT)

	rtype := parser.grabOptionalType(prec)
	right := parser.grabFunctionRight(prec)

	return FunctionExpression{
		FnToken: fn,
		Name:    name,
		Argdefs: argdefs,
		Rtype:   rtype,
		Right:   right,
	}
}

func (parser *Parser) grabCodeblock() Codeblock {
	openBr := *parser.expect(LCurlyBrT)

	statements := []Statement{}

	for parser.currentToken().TType != RCurlyBrT {
		statements = append(statements, parser.grabStatement())
	}

	parser.expect(RCurlyBrT)

	return Codeblock{
		OpenBr:     openBr,
		Statements: statements,
	}
}

func (parser *Parser) grabDoExpression() DoExpression {
	return DoExpression{
		DoToken:   parser.expect(DoT),
		Codeblock: parser.grabCodeblock(),
	}
}

func (parser *Parser) grabVarExpression(prec int) VarExpression {
	return VarExpression{
		VarToken: *parser.expect(VarT),
		Right: parser.grabExpression(prec),
	}
}

func (parser *Parser) grabWhileExpression(prec int) WhileExpression {
	return WhileExpression{
		WhileToken: *parser.expect(WhileT),
		Condition: parser.grabExpression(prec),
		Body: parser.grabExpression(prec),
	}
}
