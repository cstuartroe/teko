package lexparse

type Node interface {
	Token() Token
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
	OpenBr     Token
	Statements []Statement
}

func (c Codeblock) Token() Token {
	return c.OpenBr
}

//---

type ExpressionStatement struct {
	Expression Expression
	semicolon  *Token
}

func (s ExpressionStatement) Token() Token {
	return s.Expression.Token()
}

func (s ExpressionStatement) Semicolon() *Token {
	return s.semicolon
}

//---

type TypeStatement struct {
	TypeToken      Token
	Name           Token
	TypeExpression Expression
	semicolon      *Token
}

func (s TypeStatement) Token() Token {
	return s.TypeToken
}

func (s TypeStatement) Semicolon() *Token {
	return s.semicolon
}

//---

type SimpleExpression struct {
	token Token
}

func (e SimpleExpression) Token() Token {
	return e.token
}

func (e SimpleExpression) expressionNode() {}

//---

type DeclarationExpression struct {
	Symbol   Token
	Tekotype Expression
	Setter   Token
	Right    Expression
}

func (e DeclarationExpression) Token() Token {
	return e.Symbol
}

func (e DeclarationExpression) expressionNode() {}

//---

type UpdateExpression struct {
	Updated Expression
	Setter  Token
	Right   Expression
}

func (n UpdateExpression) Token() Token {
	return n.Setter
}

func (n UpdateExpression) expressionNode() {}

//---

type CallExpression struct {
	Receiver Expression
	Args     []Expression
	Kwargs   []FunctionKwarg
}

func (e CallExpression) Token() Token {
	return e.Receiver.Token()
}

func (e CallExpression) expressionNode() {}

type FunctionKwarg struct {
	Symbol Token
	Value  Expression
}

func (k FunctionKwarg) Token() Token {
	return k.Symbol
}

//---

type BinopExpression struct {
	Left      Expression
	Operation Token
	Right     Expression
}

func (e BinopExpression) Token() Token {
	return e.Operation
}

func (e BinopExpression) expressionNode() {}

//---

type AttributeExpression struct {
	Left   Expression
	Symbol Token
}

func (e AttributeExpression) Token() Token {
	return e.Symbol
}

func (e AttributeExpression) expressionNode() {}

//---

type TupleExpression struct {
	Elements []Expression
	LPar     Token
}

func (e TupleExpression) Token() Token {
	return e.LPar
}

func (e TupleExpression) expressionNode() {}

//--

type SuffixExpression struct {
	Left   Expression
	Suffix Token
}

func (e SuffixExpression) Token() Token {
	return e.Suffix
}

func (e SuffixExpression) expressionNode() {}

//---

type IfExpression struct {
	If        Token
	Condition Expression
	Then      Expression
	Else      Expression
}

func (e IfExpression) Token() Token {
	return e.If
}

func (e IfExpression) expressionNode() {}

//---

type seqType string

const ArraySeqType seqType = "array"
const SetSeqType seqType = "set"

type SequenceExpression struct {
	OpenBrace Token
	Stype     seqType
	Elements  []Expression
}

func (e SequenceExpression) Token() Token {
	return e.OpenBrace
}

func (e SequenceExpression) expressionNode() {}

//---

type ObjectField struct {
	Symbol Token
	Value  Expression
}

func (of ObjectField) Token() Token {
	return of.Symbol
}

type ObjectExpression struct {
	OpenBrace Token
	Fields    []ObjectField
}

func (e ObjectExpression) Token() Token {
	return e.OpenBrace
}

func (e ObjectExpression) expressionNode() {}

//---

type ArgdefNode struct {
	Symbol   Token
	Tekotype Expression
}

func (a ArgdefNode) Token() Token {
	return a.Symbol
}

type FunctionExpression struct {
	FnToken Token
	Name    *Token
	Argdefs []ArgdefNode
	Rtype   Expression
	Right   Expression
}

func (e FunctionExpression) Token() Token {
	return e.FnToken
}

func (e FunctionExpression) expressionNode() {}

//---

type DoExpression struct {
	DoToken   *Token // commonly omitted
	Codeblock Codeblock
}

func (e DoExpression) Token() Token {
	if e.DoToken != nil {
		return *e.DoToken
	} else {
		return e.Codeblock.Token()
	}
}

func (e DoExpression) expressionNode() {}
