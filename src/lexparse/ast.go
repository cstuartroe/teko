package lexparse

type Node interface {
	Token() *Token
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	Semicolon() *Token
}

//---

type Codeblock struct {
	OpenBr     *Token
	Statements []Statement
}

func (c Codeblock) Token() *Token {
	return c.OpenBr
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

//---

type SimpleExpression struct {
	token *Token
}

func (e SimpleExpression) Token() *Token {
	return e.token
}

func (e SimpleExpression) expressionNode() {}

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

type FunctionKwarg struct {
	Symbol *Token
	Value  Expression
}

func (k FunctionKwarg) Token() *Token {
	return k.Symbol
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

//---

type AttributeExpression struct {
	Left   Expression
	Symbol *Token
}

func (e AttributeExpression) Token() *Token {
	return e.Symbol
}

func (e AttributeExpression) expressionNode() {}

//---

type TupleExpression struct {
	Elements []Expression
	LPar     *Token
}

func (e TupleExpression) Token() *Token {
	return e.LPar
}

func (e TupleExpression) expressionNode() {}

//--

type SuffixExpression struct {
	Left   Expression
	Suffix *Token
}

func (e SuffixExpression) Token() *Token {
	return e.Suffix
}

func (e SuffixExpression) expressionNode() {}

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

//---

type seqType string

const ArraySeqType seqType = "array"
const SetSeqType seqType = "set"

type SequenceExpression struct {
	OpenBrace *Token
	Stype     seqType
	Elements  []Expression
}

func (e SequenceExpression) Token() *Token {
	return e.OpenBrace
}

func (e SequenceExpression) expressionNode() {}

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

//---

type VarExpression struct {
	VarToken *Token
	Right    Expression
}

func (e VarExpression) Token() *Token {
	return e.VarToken
}

func (e VarExpression) expressionNode() {}

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

//---

type ScopeExpression struct {
	ScopeToken *Token
	Codeblock  *Codeblock
}

func (e ScopeExpression) Token() *Token {
	return e.ScopeToken
}

func (e ScopeExpression) expressionNode() {}

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
	Inside Expression
}

func (e SliceExpression) Token() *Token {
	return e.Left.Token()
}

func (e SliceExpression) expressionNode() {}

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
