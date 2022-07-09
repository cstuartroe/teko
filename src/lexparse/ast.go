package lexparse

import "github.com/cstuartroe/teko/src/shared"

func fakeToken(t *Token, ttype tokenType, value string) *Token {
	return &Token{
		Line:  t.Line,
		Col:   t.Col,
		TType: ttype,
		Value: []rune(value),
	}
}

//---

type Node interface {
	Token() *Token
}

type Expression interface {
	Node
	expressionNode()
	Transform() Expression
}

func transformNullable(e Expression) Expression {
	if e == nil {
		return nil
	} else {
		return e.Transform()
	}
}

func transformEach(nodes []Expression) []Expression {
	out := []Expression{}

	for _, n := range nodes {
		out = append(out, n.Transform())
	}

	return out
}

type Statement interface {
	Node
	Semicolon() *Token
	Transform() Statement
}

//---

type Codeblock struct {
	OpenBr     *Token
	Statements []Statement
}

func (c Codeblock) Token() *Token {
	return c.OpenBr
}

func (c Codeblock) Transform() *Codeblock {
	statements := []Statement{}

	for _, s := range c.Statements {
		statements = append(statements, s.Transform())
	}

	return &Codeblock{
		OpenBr:     c.OpenBr,
		Statements: statements,
	}
}

//---

type ExpressionStatement struct {
	Expression Expression
	semicolon  *Token
}

func (s ExpressionStatement) Token() *Token {
	return s.Expression.Token()
}

func (s ExpressionStatement) Semicolon() *Token {
	return s.semicolon
}

func (s ExpressionStatement) Transform() Statement {
	return &ExpressionStatement{
		Expression: s.Expression.Transform(),
		semicolon:  s.semicolon,
	}
}

//---

type TypeStatement struct {
	TypeToken      *Token
	Name           *Token
	TypeExpression Expression
	semicolon      *Token
}

func (s TypeStatement) Token() *Token {
	return s.TypeToken
}

func (s TypeStatement) Semicolon() *Token {
	return s.semicolon
}

func (s *TypeStatement) Transform() Statement {
	return s
}

//---

type SimpleExpression struct {
	token *Token
}

func (e SimpleExpression) Token() *Token {
	return e.token
}

func (e SimpleExpression) expressionNode() {}

func (e *SimpleExpression) Transform() Expression {
	return e
}

//---

type DeclarationExpression struct {
	Symbol   *Token
	Tekotype Expression
	Setter   *Token
	Right    Expression
}

func (e DeclarationExpression) Token() *Token {
	return e.Symbol
}

func (e DeclarationExpression) expressionNode() {}

func (e DeclarationExpression) Transform() Expression {
	return &DeclarationExpression{
		Symbol:   e.Symbol,
		Tekotype: e.Tekotype,
		Setter:   e.Setter,
		Right:    e.Right.Transform(),
	}
}

//---

type CallExpression struct {
	Receiver Expression
	Args     []Expression
	Kwargs   []*FunctionKwarg
}

func (e CallExpression) Token() *Token {
	return e.Receiver.Token()
}

func (e CallExpression) expressionNode() {}

func (e CallExpression) Transform() Expression {
	kwargs := []*FunctionKwarg{}

	for _, kw := range e.Kwargs {
		kwargs = append(kwargs, kw.Transform())
	}

	return &CallExpression{
		Receiver: e.Receiver.Transform(),
		Args:     transformEach(e.Args),
		Kwargs:   kwargs,
	}
}

type FunctionKwarg struct {
	Symbol *Token
	Value  Expression
}

func (k FunctionKwarg) Token() *Token {
	return k.Symbol
}

func (k FunctionKwarg) Transform() *FunctionKwarg {
	return &FunctionKwarg{
		Symbol: k.Symbol,
		Value:  k.Value.Transform(),
	}
}

//---

type BinopExpression struct {
	Left      Expression
	Operation *Token
	Right     Expression
}

func (e BinopExpression) Token() *Token {
	return e.Operation
}

func (e BinopExpression) expressionNode() {}

func (e BinopExpression) Transform() Expression {
	return &CallExpression{
		Receiver: &AttributeExpression{
			Left:   e.Left.Transform(),
			Symbol: fakeToken(e.Operation, SymbolT, binops[string(e.Operation.Value)]),
		},
		Args: []Expression{e.Right.Transform()},
	}
}

//---

type ComparisonExpression struct {
	Left       Expression
	Comparator *Token
	Right      Expression
}

func (e ComparisonExpression) Token() *Token {
	return e.Comparator
}

func (e ComparisonExpression) expressionNode() {}

func (e ComparisonExpression) Transform() Expression {
	return &ComparisonExpression{
		Left:       e.Left.Transform(),
		Comparator: e.Comparator,
		Right:      e.Right.Transform(),
	}
}

//---

type AttributeExpression struct {
	Left   Expression
	Symbol *Token
}

func (e AttributeExpression) Token() *Token {
	return e.Symbol
}

func (e AttributeExpression) expressionNode() {}

func (e AttributeExpression) Transform() Expression {
	return &AttributeExpression{
		Left:   e.Left.Transform(),
		Symbol: e.Token(),
	}
}

//---

type TupleExpression struct {
	Elements []Expression
	LPar     *Token
}

func (e TupleExpression) Token() *Token {
	return e.LPar
}

func (e TupleExpression) expressionNode() {}

func (e TupleExpression) Transform() Expression {
	return &TupleExpression{
		Elements: transformEach(e.Elements),
		LPar:     e.LPar,
	}
}

//--

type SuffixExpression struct {
	Left   Expression
	Suffix *Token
}

func (e SuffixExpression) Token() *Token {
	return e.Suffix
}

func (e SuffixExpression) expressionNode() {}

func (e SuffixExpression) Transform() Expression {
	return &CallExpression{
		Receiver: &AttributeExpression{
			Left:   e.Left.Transform(),
			Symbol: fakeToken(e.Suffix, SymbolT, suffixes[string(e.Suffix.Value)]),
		},
	}
}

//---

type IfExpression struct {
	If        *Token
	Condition Expression
	Then      Expression
	Else      Expression
}

func (e IfExpression) Token() *Token {
	return e.If
}

func (e IfExpression) expressionNode() {}

func (e IfExpression) Transform() Expression {
	return &IfExpression{
		If:        e.If,
		Condition: e.Condition.Transform(),
		Then:      e.Then.Transform(),
		Else:      transformNullable(e.Else),
	}
}

//---

type seqType string

const ArraySeqType seqType = "array"
const SetSeqType seqType = "set"

type SequenceExpression struct {
	OpenBrace *Token
	Var       *Token
	Stype     seqType
	Elements  []Expression
}

func (e SequenceExpression) Token() *Token {
	return e.OpenBrace
}

func (e SequenceExpression) expressionNode() {}

func (e SequenceExpression) Transform() Expression {
	return &SequenceExpression{
		OpenBrace: e.OpenBrace,
		Var:       e.Var,
		Stype:     e.Stype,
		Elements:  transformEach(e.Elements),
	}
}

//---

type KVPair struct {
	Key   Expression
	Value Expression
}

type MapExpression struct {
	MapToken  *Token
	Ktype     Expression
	Vtype     Expression
	HasBraces bool
	KVPairs   []*KVPair
}

func (e MapExpression) Token() *Token {
	return e.MapToken
}

func (e MapExpression) expressionNode() {}

func (e MapExpression) Transform() Expression {
	kvpairs := []*KVPair{}

	for _, p := range e.KVPairs {
		kvpairs = append(kvpairs, &KVPair{
			Key:   p.Key.Transform(),
			Value: p.Value.Transform(),
		})
	}

	return &MapExpression{
		MapToken:  e.MapToken,
		Ktype:     e.Ktype,
		Vtype:     e.Vtype,
		HasBraces: e.HasBraces,
		KVPairs:   kvpairs,
	}
}

//---

type ObjectField struct {
	Symbol *Token
	Value  Expression
}

func (of ObjectField) Token() *Token {
	return of.Symbol
}

type ObjectExpression struct {
	OpenBrace *Token
	Fields    []*ObjectField
}

func (e ObjectExpression) Token() *Token {
	return e.OpenBrace
}

func (e ObjectExpression) expressionNode() {}

func (e ObjectExpression) Transform() Expression {
	fields := []*ObjectField{}

	for _, f := range e.Fields {
		var value Expression = f.Value

		// if value == nil {
		// 	value = &SimpleExpression{f.Symbol}
		// }

		fields = append(fields, &ObjectField{
			Symbol: f.Symbol,
			Value:  value.Transform(),
		})
	}

	return &ObjectExpression{
		OpenBrace: e.OpenBrace,
		Fields:    fields,
	}
}

//---

type ArgdefNode struct {
	Symbol   *Token
	Tekotype Expression
	Default  Expression
}

func (a ArgdefNode) Token() *Token {
	return a.Symbol
}

type FunctionExpression struct {
	FnToken *Token
	Name    *Token
	GDL     *GenericDeclarationList
	Argdefs []*ArgdefNode
	Rtype   Expression
	Right   Expression
}

func (e FunctionExpression) Token() *Token {
	return e.FnToken
}

func (e FunctionExpression) expressionNode() {}

func (e FunctionExpression) Transform() Expression {
	if e.Right == nil {
		e.Token().Raise(shared.SyntaxError, "Expected ->")
	}

	argdefs := []*ArgdefNode{}

	for _, ad := range e.Argdefs {
		argdefs = append(argdefs, &ArgdefNode{
			Symbol:   ad.Symbol,
			Tekotype: ad.Tekotype,
			Default:  transformNullable(ad.Default),
		})
	}

	return &FunctionExpression{
		FnToken: e.FnToken,
		Name:    e.Name,
		GDL:     e.GDL,
		Argdefs: argdefs,
		Rtype:   e.Rtype,
		Right:   e.Right.Transform(),
	}
}

//---

type DoExpression struct {
	DoToken   *Token // commonly omitted
	Codeblock *Codeblock
}

func (e DoExpression) Token() *Token {
	if e.DoToken != nil {
		return e.DoToken
	} else {
		return e.Codeblock.Token()
	}
}

func (e DoExpression) expressionNode() {}

func (e DoExpression) Transform() Expression {
	return &DoExpression{
		DoToken:   e.DoToken,
		Codeblock: e.Codeblock.Transform(),
	}
}

//---

type VarExpression struct {
	VarToken *Token
	Right    Expression
}

func (e VarExpression) Token() *Token {
	return e.VarToken
}

func (e VarExpression) expressionNode() {}

func (e VarExpression) Transform() Expression {
	return &VarExpression{
		VarToken: e.VarToken,
		Right:    e.Right.Transform(),
	}
}

//---

type WhileExpression struct {
	WhileToken *Token
	Condition  Expression
	Body       Expression
}

func (e WhileExpression) Token() *Token {
	return e.WhileToken
}

func (e WhileExpression) expressionNode() {}

func (e WhileExpression) Transform() Expression {
	return &WhileExpression{
		WhileToken: e.WhileToken,
		Condition:  e.Condition.Transform(),
		Body:       e.Body.Transform(),
	}
}

//---

type ForExpression struct {
	ForToken *Token
	Iterand  *Token
	Tekotype Expression
	Iterator Expression
	Body     Expression
}

func (e ForExpression) Token() *Token {
	return e.ForToken
}

func (e ForExpression) expressionNode() {}

func (e ForExpression) Transform() Expression {
	return &CallExpression{
		Receiver: &AttributeExpression{
			Left:   e.Iterator.Transform(),
			Symbol: fakeToken(e.ForToken, SymbolT, "forEach"),
		},
		Args: []Expression{
			&FunctionExpression{
				FnToken: e.ForToken,
				Argdefs: []*ArgdefNode{
					{
						Symbol:   e.Iterand,
						Tekotype: e.Tekotype,
					},
				},
				Right: e.Body.Transform(),
			},
		},
	}
}

//---

type ScopeExpression struct {
	ScopeToken *Token
	Codeblock  *Codeblock
}

func (e ScopeExpression) Token() *Token {
	return e.ScopeToken
}

func (e ScopeExpression) expressionNode() {}

func (e ScopeExpression) Transform() Expression {
	return &ScopeExpression{
		ScopeToken: e.ScopeToken,
		Codeblock:  e.Codeblock.Transform(),
	}
}

//---

type PipeExpression struct {
	PipeToken *Token
	Arg       Expression
	Function  Expression
}

func (e PipeExpression) Token() *Token {
	return e.PipeToken
}

func (e PipeExpression) expressionNode() {}

func (e PipeExpression) Transform() Expression {
	return &CallExpression{
		Receiver: e.Function.Transform(),
		Args:     []Expression{e.Arg.Transform()},
	}
}

//---

type GenericDeclaration struct {
	Name      *Token
	Supertype Expression
}

func (gd GenericDeclaration) Token() *Token {
	return gd.Name
}

type GenericDeclarationList struct {
	OpenBrace    *Token
	Declarations []GenericDeclaration
}

func (gd GenericDeclarationList) Token() *Token {
	return gd.OpenBrace
}

//---

type SliceExpression struct {
	Left   Expression
	OpenBr *Token
	Inside Expression
}

func (e SliceExpression) Token() *Token {
	return e.Left.Token()
}

func (e SliceExpression) expressionNode() {}

func (e SliceExpression) Transform() Expression {
	if e.Inside == nil {
		e.OpenBr.Raise(shared.SyntaxError, "Empty slice")
	}

	return &CallExpression{
		Receiver: &AttributeExpression{
			Left:   e.Left,
			Symbol: fakeToken(e.OpenBr, SymbolT, "at"),
		},
		Args: []Expression{e.Inside},
	}
}

//---

type Any interface{}

type CaseBlock struct {
	Case     *Token
	TType    Expression
	GateType Any
	Body     Expression
}

func (e CaseBlock) Token() *Token {
	return e.Case
}

type SwitchExpression struct {
	Switch  *Token
	Symbol  *Token
	Cases   []*CaseBlock
	Default Expression
}

func (e SwitchExpression) Token() *Token {
	return e.Switch
}

func (e SwitchExpression) expressionNode() {}

func (e SwitchExpression) Transform() Expression {
	cases := []*CaseBlock{}

	for _, c := range e.Cases {
		cases = append(cases, &CaseBlock{
			Case:     c.Case,
			TType:    c.TType,
			GateType: c.GateType,
			Body:     c.Body.Transform(),
		})
	}

	return &SwitchExpression{
		Switch:  e.Switch,
		Symbol:  e.Symbol,
		Cases:   cases,
		Default: transformNullable(e.Default),
	}
}

//---

type SetterExpression struct {
	Left   Expression
	Setter *Token
	Right  Expression
}

func (e SetterExpression) Token() *Token {
	return e.Setter
}

func (e SetterExpression) expressionNode() {}

func (e SetterExpression) Transform() Expression {
	if e.Setter.TType == EqualT {
		return &CallExpression{
			Receiver: &AttributeExpression{
				Left:   e.Left.Transform(),
				Symbol: fakeToken(e.Setter, SymbolT, "="),
			},
			Args: []Expression{
				e.Right.Transform(),
			},
		}
	} else {
		be := &BinopExpression{
			Left:      e.Left,
			Operation: fakeToken(e.Setter, BinopT, updaters[string(e.Setter.Value)]),
			Right:     e.Right,
		}

		se := &SetterExpression{
			Left:   e.Left,
			Setter: fakeToken(e.Setter, EqualT, "="),
			Right:  be,
		}

		return se.Transform()
	}
}
