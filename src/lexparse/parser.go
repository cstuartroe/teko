package lexparse

import (
	"fmt"

	"github.com/cstuartroe/teko/src/shared"
)

var simpleExprTokenTypes map[tokenType]bool = map[tokenType]bool{
	SymbolT: true,
	StringT: true,
	CharT:   true,
	IntT:    true,
	FloatT:  true,
	BoolT:   true,
}

type Parser struct {
	Tokens    []Token
	position  int
	Codeblock *Codeblock
}

func (parser *Parser) ParseFile(filename string, transform bool) {
	lexer := FromFile(filename)
	tokens := lexer.Lex()
	parser.Parse(tokens, transform)
}

func (parser *Parser) Parse(tokens []Token, transform bool) {
	parser.Tokens = tokens
	parser.Codeblock = &Codeblock{}
	parser.skipComments()

	for parser.HasMore() {
		parser.Codeblock.Statements = append(
			parser.Codeblock.Statements,
			parser.grabStatement(),
		)
	}

	if transform {
		parser.Codeblock = parser.Codeblock.Transform()
	}
}

func (parser *Parser) currentToken() *Token {
	if parser.HasMore() {
		return &(parser.Tokens[parser.position])
	} else {
		t := parser.Tokens[parser.position-1]
		raiseTekoError(t.Line, t.Col+len(t.Value), shared.SyntaxError, "Unexpected EOF")
		return nil
	}
}

func (parser *Parser) advance() {
	parser.position++
	parser.skipComments()
}

func (parser *Parser) skipComments() {
	if !parser.HasMore() {
		return
	}

	t := parser.currentToken()

	if t.TType == LineCommentT || t.TType == BlockCommentT {
		parser.advance()
	}
}

func (parser *Parser) HasMore() bool {
	return parser.position < len(parser.Tokens)
}

func (parser *Parser) expect(ttype tokenType) *Token {
	token := parser.currentToken()

	if token.TType != ttype {
		token.Raise(shared.SyntaxError, fmt.Sprintf("Expected %s", ttype))
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

func (parser *Parser) grabTypeStatement() *TypeStatement {
	tt := parser.expect(TypeT)
	name := parser.expect(SymbolT)
	parser.expect(EqualT)

	return &TypeStatement{
		TypeToken:      tt,
		Name:           name,
		TypeExpression: parser.grabTypeExpression(min_prec),
		semicolon:      parser.optionalSemicolon(),
	}
}

func (parser *Parser) grabExpressionStmt() *ExpressionStatement {
	return &ExpressionStatement{
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

	case MapT:
		expr = parser.grabMap()

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

	case ForT:
		expr = parser.grabForLoop(prec)

	case ScopeT:
		expr = parser.grabScopeExpression()

	case SwitchT:
		expr = parser.grabSwitchExpression()

	default:
		expr = parser.grabSimpleExpression()
	}

	return parser.continueExpression(expr, prec)
}

func (parser *Parser) grabTypeExpression(prec int) Expression {
	if prec < add_sub_prec {
		prec = add_sub_prec
	}

	return parser.grabExpression(prec)
}

func (parser *Parser) grabSimpleExpression() *SimpleExpression {
	t := parser.currentToken()
	if _, ok := simpleExprTokenTypes[t.TType]; ok {
		n := &SimpleExpression{token: t}
		parser.advance()
		return n
	} else {
		parser.currentToken().Raise(shared.SyntaxError, "Illegal start to expression")
		return nil // unreachable code that the compiler requires
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
		if prec <= setter_prec {
			if !isValidDeclared(expr) {
				parser.currentToken().Raise(shared.SyntaxError, "Illegal left-hand side to declaration")
			}

			parser.advance()

			var tekotype Expression = nil
			switch parser.currentToken().TType {
			case EqualT:
			default:
				tekotype = parser.grabTypeExpression(prec)
			}

			setter := parser.expect(EqualT)

			out = &DeclarationExpression{
				Symbol:   expr.Token(),
				Tekotype: tekotype,
				Setter:   setter,
				Right:    parser.grabExpression(prec),
			}
		}

	case LParT:
		out = parser.makeCallExpression(expr)

	case DotT:
		out = parser.makeAttributeExpression(expr)

	case BinopT:
		op_prec := binop_precs[binops[value]]
		if prec <= op_prec {
			op := parser.currentToken()
			parser.advance()
			right := parser.grabExpression(op_prec + 1)

			out = &BinopExpression{
				Left:      expr,
				Operation: op,
				Right:     right,
			}
		}

	case SuffixT:
		suffix := parser.currentToken()
		parser.advance()

		out = &SuffixExpression{
			Left:   expr,
			Suffix: suffix,
		}

	case UpdaterT:
		if prec <= setter_prec {
			panic("Other updaters not supported yet")
		}

	case EqualT:
		if prec <= setter_prec {
			setter := *parser.expect(EqualT)

			out = &CallExpression{
				Receiver: &AttributeExpression{
					Left:   expr,
					Symbol: fakeToken(&setter, SymbolT, "="),
				},
				Args: []Expression{
					parser.grabExpression(setter_prec + 1),
				},
			}
		}

	case ComparisonT:
		if prec <= comparison_prec {
			return &ComparisonExpression{
				Left:       expr,
				Comparator: parser.expect(ComparisonT),
				Right:      parser.grabExpression(comparison_prec + 1),
			}
		}

	case PipeT:
		if prec <= setter_prec {
			pipe := parser.expect(PipeT)
			function := parser.grabExpression(setter_prec + 1)

			out = &PipeExpression{
				PipeToken: pipe,
				Arg:       expr,
				Function:  function,
			}
		}

	case LSquareBrT:
		br := parser.expect(LSquareBrT)

		var inside Expression = nil
		if parser.currentToken().TType != RSquareBrT {
			inside = parser.grabExpression(min_prec)
		}

		parser.expect(RSquareBrT)

		out = &SliceExpression{
			Left:   expr,
			OpenBr: br,
			Inside: inside,
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
	case *SimpleExpression:
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

func (parser *Parser) makeCallExpression(receiver Expression) *CallExpression {
	parser.expect(LParT)

	args := []Expression{}
	kwargs := []*FunctionKwarg{}
	on_kwargs := false
	cont := parser.currentToken().TType != RParT

	for cont {
		arg := parser.grabExpression(add_sub_prec)

		if parser.currentToken().TType == EqualT {
			switch p := arg.(type) {
			case *SimpleExpression:
				if p.token.TType != SymbolT {
					p.token.Raise(shared.SyntaxError, "Left-hand side of keyword argument cannot be a value")
				}
			default:
				p.Token().Raise(shared.SyntaxError, "Left-hand side of keyword argument cannot be a value")
			}

			parser.advance()
			on_kwargs = true

			kwargs = append(kwargs, &FunctionKwarg{
				Symbol: arg.Token(),
				Value:  parser.grabExpression(add_sub_prec),
			})
		} else {
			if on_kwargs {
				parser.currentToken().Raise(shared.SyntaxError, "All positional arguments must be before all keyword arguments")
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

	return &CallExpression{
		Receiver: receiver,
		Args:     args,
		Kwargs:   kwargs,
	}
}

func (parser *Parser) makeAttributeExpression(left Expression) *AttributeExpression {
	parser.expect(DotT)

	return &AttributeExpression{
		Left:   left,
		Symbol: parser.expect(SymbolT),
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
	lpar := parser.expect(LParT)

	seq := parser.grabSequence(RParT)
	switch len(seq) {
	case 0:
		lpar.Raise(shared.SyntaxError, "Cannot have empty tuple")
		return nil
	case 1:
		return seq[0]
	default:
		return &TupleExpression{
			Elements: seq,
			LPar:     lpar,
		}
	}
}

func (parser *Parser) grabArray() *SequenceExpression {
	open := parser.expect(LSquareBrT)

	return &SequenceExpression{
		OpenBrace: open,
		Stype:     ArraySeqType,
		Elements:  parser.grabSequence(RSquareBrT),
	}
}

func (parser *Parser) grabSet() *SequenceExpression {
	open := parser.expect(SetT)

	// TODO: Typehints for sets would be nice

	parser.expect(LCurlyBrT)

	return &SequenceExpression{
		OpenBrace: open,
		Stype:     SetSeqType,
		Elements:  parser.grabSequence(RCurlyBrT),
	}
}

func (parser *Parser) grabMap() *MapExpression {
	map_token := parser.expect(MapT)

	var ktype, vtype Expression = nil, nil
	if parser.currentToken().TType == LSquareBrT {
		parser.advance()
		ktype = parser.grabTypeExpression(min_prec)
		parser.expect(CommaT)
		vtype = parser.grabTypeExpression(min_prec)
		parser.expect(RSquareBrT)
	}

	kvpairs := []*KVPair{}
	has_braces := false

	if parser.currentToken().TType == LCurlyBrT {
		parser.advance()
		has_braces = true

		for parser.currentToken().TType != RCurlyBrT {
			key := parser.grabExpression(setter_prec + 1)
			parser.expect(ColonT)
			value := parser.grabExpression(min_prec)

			kvpairs = append(kvpairs, &KVPair{Key: key, Value: value})

			if parser.currentToken().TType == CommaT {
				parser.advance()
			} else {
				break
			}
		}

		parser.expect(RCurlyBrT)
	}

	return &MapExpression{
		MapToken:  map_token,
		Ktype:     ktype,
		Vtype:     vtype,
		HasBraces: has_braces,
		KVPairs:   kvpairs,
	}
}

func (parser *Parser) grabControlBlock(prec int) Expression {
	if parser.currentToken().TType == LCurlyBrT {
		return &DoExpression{
			DoToken:   nil,
			Codeblock: parser.grabCodeblock(),
		}
	} else {
		if parser.currentToken().TType == DoT {
			parser.currentToken().Raise(shared.SyntaxError, "Do keyword not needed in control block")
		}

		return parser.grabExpression(prec)
	}
}

func (parser *Parser) grabIf(prec int) *IfExpression {
	if_token := parser.expect(IfT)

	cond := parser.grabExpression(prec)

	if parser.currentToken().TType == ThenT {
		parser.advance()
	}

	then := parser.grabControlBlock(prec)
	var else_expr Expression = nil

	if parser.currentToken().TType == ElseT {
		parser.advance()
		else_expr = parser.grabControlBlock(prec)
	}

	return &IfExpression{
		If:        if_token,
		Condition: cond,
		Then:      then,
		Else:      else_expr,
	}
}

func (parser *Parser) grabObject() *ObjectExpression {
	open := parser.expect(LCurlyBrT)

	fields := []*ObjectField{}

	for parser.currentToken().TType != RCurlyBrT {
		fields = append(fields, parser.grabObjectField())

		if parser.currentToken().TType == CommaT {
			parser.advance()
		} else {
			break
		}
	}

	parser.expect(RCurlyBrT)

	return &ObjectExpression{
		OpenBrace: open,
		Fields:    fields,
	}
}

func (parser *Parser) grabObjectField() *ObjectField {
	symbol := parser.expect(SymbolT)
	var value Expression

	if parser.currentToken().TType == ColonT {
		parser.advance()

		value = parser.grabExpression(min_prec)
	} else {
		value = &SimpleExpression{symbol}
	}

	return &ObjectField{
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

func (parser *Parser) grabArgdefs() []*ArgdefNode {
	argdefs := []*ArgdefNode{}

	for parser.currentToken().TType != RParT {
		symbol := parser.expect(SymbolT)

		tekotype := parser.grabOptionalType(min_prec)

		var dft Expression = nil
		if parser.currentToken().TType == EqualT {
			parser.advance()
			dft = parser.grabExpression(max_prec)
		}

		argdefs = append(argdefs, &ArgdefNode{
			Symbol:   symbol,
			Tekotype: tekotype,
			Default:  dft,
		})

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
		return &DoExpression{
			DoToken:   nil,
			Codeblock: parser.grabCodeblock(),
		}
	} else {
		parser.expect(ArrowT)

		return parser.grabExpression(prec)
	}
}

func (parser *Parser) grabGenericDeclarationList() *GenericDeclarationList {
	out := &GenericDeclarationList{
		OpenBrace: parser.expect(LSquareBrT),
	}

	for parser.currentToken().TType != RSquareBrT {
		name := parser.expect(SymbolT)

		var supertype Expression = nil
		if parser.currentToken().TType == SubtypeT {
			parser.advance()
			supertype = parser.grabExpression(add_sub_prec)
		}

		out.Declarations = append(out.Declarations, GenericDeclaration{Name: name, Supertype: supertype})

		if parser.currentToken().TType == CommaT {
			parser.advance()
		} else {
			break
		}
	}

	parser.expect(RSquareBrT)

	return out
}

func (parser *Parser) grabFunctionDefinition(prec int) *FunctionExpression {
	fn := parser.expect(FnT)

	var name *Token = nil
	if parser.currentToken().TType == SymbolT {
		name = parser.expect(SymbolT)
	}

	var gdl *GenericDeclarationList
	if parser.currentToken().TType == LSquareBrT {
		gdl = parser.grabGenericDeclarationList()
	}

	parser.expect(LParT)

	argdefs := parser.grabArgdefs()

	parser.expect(RParT)

	rtype := parser.grabOptionalType(prec)
	right := parser.grabFunctionRight(prec)

	return &FunctionExpression{
		FnToken: fn,
		Name:    name,
		GDL:     gdl,
		Argdefs: argdefs,
		Rtype:   rtype,
		Right:   right,
	}
}

func (parser *Parser) grabCodeblock() *Codeblock {
	openBr := parser.expect(LCurlyBrT)

	statements := []Statement{}

	for parser.currentToken().TType != RCurlyBrT {
		statements = append(statements, parser.grabStatement())
	}

	parser.expect(RCurlyBrT)

	return &Codeblock{
		OpenBr:     openBr,
		Statements: statements,
	}
}

func (parser *Parser) grabDoExpression() *DoExpression {
	return &DoExpression{
		DoToken:   parser.expect(DoT),
		Codeblock: parser.grabCodeblock(),
	}
}

func (parser *Parser) grabVarExpression(prec int) *VarExpression {
	return &VarExpression{
		VarToken: parser.expect(VarT),
		Right:    parser.grabExpression(prec),
	}
}

func (parser *Parser) grabWhileExpression(prec int) *WhileExpression {
	return &WhileExpression{
		WhileToken: parser.expect(WhileT),
		Condition:  parser.grabExpression(prec),
		Body:       parser.grabControlBlock(prec),
	}
}

func (parser *Parser) grabForLoop(prec int) Expression {
	for_t := parser.expect(ForT)
	iterand := parser.expect(SymbolT)
	ttype := parser.grabOptionalType(min_prec)
	parser.expect(InT)
	iterator := parser.grabExpression(prec)
	body := parser.grabControlBlock(prec)

	return &ForExpression{
		ForToken: for_t,
		Iterand:  iterand,
		Tekotype: ttype,
		Iterator: iterator,
		Body:     body,
	}
}

func (parser *Parser) grabScopeExpression() *ScopeExpression {
	return &ScopeExpression{
		ScopeToken: parser.expect(ScopeT),
		Codeblock:  parser.grabCodeblock(),
	}
}

func (parser *Parser) grabCaseBlock() *CaseBlock {
	case_token := parser.expect(CaseT)
	ttype := parser.grabTypeExpression(min_prec)

	if parser.currentToken().TType == ColonT {
		parser.advance()
	}

	body := parser.grabControlBlock(min_prec)

	return &CaseBlock{
		Case:  case_token,
		TType: ttype,
		Body:  body,
	}
}

func (parser *Parser) grabSwitchExpression() *SwitchExpression {
	Switch := parser.expect(SwitchT)
	Symbol := parser.expect(SymbolT)
	parser.expect(LCurlyBrT)

	Cases := []*CaseBlock{}

	for parser.currentToken().TType == CaseT {
		Cases = append(Cases, parser.grabCaseBlock())
	}

	var Default Expression = nil
	if parser.currentToken().TType == DefaultT {
		parser.advance()
		Default = parser.grabDoExpression()
	}

	parser.expect(RCurlyBrT)

	return &SwitchExpression{
		Switch:  Switch,
		Symbol:  Symbol,
		Cases:   Cases,
		Default: Default,
	}
}

func ComparisonCallExpression(expr *ComparisonExpression) *CallExpression {
	return &CallExpression{
		Receiver: &AttributeExpression{
			Left:   expr.Left,
			Symbol: fakeToken(expr.Comparator, expr.Comparator.TType, "compare"),
		},
		Args: []Expression{
			expr.Right,
		},
		Kwargs: []*FunctionKwarg{},
	}
}
