package lexparse

import (
	"fmt"
	"strings"
)

type Node interface {
	Ntype() string
	children() []Node
	Token() Token
	child_strings(indent int) []string
}

func node_to_str(node Node, indent int) string {
	space := strings.Repeat(" ", indent*INDENT_AMOUNT)
	out := fmt.Sprintf("%s%s\n", space, node.Ntype())
	for _, s := range node.child_strings(indent + 1) {
		out += s
	}
	return out
}

func PrintNode(node Node) {
	fmt.Print(node_to_str(node, 0))
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

//---

type Codeblock struct {
	statements []Statement
}

func (c *Codeblock) Ntype() string {
	return "Codeblock"
}

func (c *Codeblock) children() []Node {
	out := []Node{}
	for _, stmt := range c.statements {
		out = append(out, stmt)
	}
	return out
}

func (c *Codeblock) Token() Token {
	return Token{}
}

func (c *Codeblock) child_strings(indent int) []string {
	out := []string{}
	for _, stmt := range c.statements {
		out = append(out, node_to_str(stmt, indent))
	}
	return out
}

func (c *Codeblock) GetStatements() []Statement {
	return c.statements
}

//---

type ExpressionStatement struct {
	Expression Expression
}

func (s ExpressionStatement) Ntype() string {
	return "ExpressionStatement"
}

func (s ExpressionStatement) children() []Node {
	return []Node{s.Expression}
}

func (s ExpressionStatement) Token() Token {
	return s.Expression.Token()
}

func (s ExpressionStatement) child_strings(indent int) []string {
	return []string{
		node_to_str(s.Expression, indent),
	}
}

func (s ExpressionStatement) statementNode() {}

//---

type SimpleExpression struct {
	token Token
}

func (e SimpleExpression) Ntype() string {
	return "SimpleExpression"
}

func (e SimpleExpression) children() []Node {
	return []Node{}
}

func (e SimpleExpression) Token() Token {
	return e.token
}

func (e SimpleExpression) child_strings(indent int) []string {
	return []string{
		e.token.to_indented_str(indent),
	}
}

func (e SimpleExpression) expressionNode() {}

//---

type DeclarationExpression struct {
	Symbol Token
	Tekotype  *Expression
	Setter Token
	Right  Expression
}

func (e DeclarationExpression) Ntype() string {
	return "DeclarationExpression"
}

func (e DeclarationExpression) children() []Node {
	return []Node{*e.Tekotype, e.Right}
}

func (e DeclarationExpression) Token() Token {
	return e.Symbol
}

func (e DeclarationExpression) child_strings(indent int) []string {
	return []string{
		e.Symbol.to_indented_str(indent),
		node_to_str(*e.Tekotype, indent),
		e.Setter.to_indented_str(indent),
		node_to_str(e.Right, indent),
	}
}

func (e DeclarationExpression) expressionNode() {}

//---

type UpdateExpression struct {
	Updated Expression
	Setter  Token
	Right   Expression
}

func (n UpdateExpression) Ntype() string {
	return "UpdateNode"
}

func (n UpdateExpression) children() []Node {
	return []Node{n.Updated, n.Right}
}

func (n UpdateExpression) Token() Token {
	return n.Setter
}

func (n UpdateExpression) child_strings(indent int) []string {
	return []string{
		node_to_str(n.Updated, indent),
		n.Setter.to_indented_str(indent),
		node_to_str(n.Right, indent),
	}
}

func (n UpdateExpression) expressionNode() {}

//---

type CallExpression struct {
	Receiver Expression
	Args     []Expression
	Kwargs   []FunctionKwarg
}

func (e CallExpression) Ntype() string {
	return "CallExpression"
}

func (e CallExpression) children() []Node {
	out := []Node{e.Receiver}
	for _, arg := range e.Args {
		out = append(out, arg)
	}
	for _, kwarg := range e.Kwargs {
		out = append(out, kwarg)
	}
	return out
}

func (e CallExpression) child_strings(indent int) []string {
	out := []string{}
	for _, arg := range e.children() {
		out = append(out, node_to_str(arg, indent))
	}
	return out
}

func (e CallExpression) Token() Token {
	return e.Receiver.Token()
}

func (e CallExpression) expressionNode() {}

type FunctionKwarg struct {
	Symbol Token
	Value  Expression
}

func (k FunctionKwarg) Ntype() string {
	return "FunctionKwarg"
}

func (k FunctionKwarg) children() []Node {
	return []Node{k.Value}
}

func (k FunctionKwarg) child_strings(indent int) []string {
	return []string{
		k.Symbol.to_indented_str(indent),
		node_to_str(k.Value, indent),
	}
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

func (e BinopExpression) Ntype() string {
	return "BinopExpression"
}

func (e BinopExpression) children() []Node {
	return []Node{e.Left, e.Right}
}

func (e BinopExpression) child_strings(indent int) []string {
	return []string{
		node_to_str(e.Left, indent),
		e.Operation.to_indented_str(indent),
		node_to_str(e.Right, indent),
	}
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

func (e AttributeExpression) Ntype() string {
	return "AttributeExpression"
}

func (e AttributeExpression) children() []Node {
	return []Node{e.Left}
}

func (e AttributeExpression) child_strings(indent int) []string {
	return []string{
		node_to_str(e.Left, indent),
		e.Symbol.to_indented_str(indent),
	}
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

func (e TupleExpression) Ntype() string {
	return "TupleExpression"
}

func (e TupleExpression) children() []Node {
	out := []Node{}
	for _, expr := range e.Elements {
		out = append(out, expr)
	}
	return out
}

func (e TupleExpression) child_strings(indent int) []string {
	out := []string{}
	for _, expr := range e.Elements {
		out = append(out, node_to_str(expr, indent))
	}
	return out
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

func (e SuffixExpression) Ntype() string {
	return "SuffixExpression"
}

func (e SuffixExpression) children() []Node {
	return []Node{e.Left}
}

func (e SuffixExpression) child_strings(indent int) []string {
	return []string{
		node_to_str(e.Left, indent),
		e.Suffix.to_indented_str(indent),
	}
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

func (e IfExpression) Ntype() string {
	return "IfExpression"
}

func (e IfExpression) children() []Node {
	return []Node{
		e.Condition,
		e.Then,
		e.Else,
	}
}

func (e IfExpression) child_strings(indent int) []string {
	return []string{
		node_to_str(e.Condition, indent),
		node_to_str(e.Then, indent),
		node_to_str(e.Else, indent),
	}
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

func (e SequenceExpression) Ntype() string {
	return "SequenceExpression (" + string(e.Stype) + ")"
}

func (e SequenceExpression) children() []Node {
	out := []Node{}
	for _, expr := range e.Elements {
		out = append(out, expr)
	}
	return out
}

func (e SequenceExpression) child_strings(indent int) []string {
	out := []string{}
	for _, expr := range e.Elements {
		out = append(out, node_to_str(expr, indent))
	}
	return out
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

func (of ObjectField) Ntype() string {
	return "ObjectField"
}

func (of ObjectField) children() []Node {
	return []Node{of.Value}
}

func (of ObjectField) child_strings(indent int) []string {
	return []string{
		of.Symbol.to_indented_str(indent),
		node_to_str(of.Value, indent),
	}
}

func (of ObjectField) Token() Token {
	return of.Symbol
}

type ObjectExpression struct {
	OpenBrace Token
	Fields    []ObjectField
}

func (e ObjectExpression) Ntype() string {
	return "ObjectExpression"
}

func (e ObjectExpression) children() []Node {
	out := []Node{}
	for _, of := range e.Fields {
		out = append(out, of)
	}
	return out
}

func (e ObjectExpression) child_strings(indent int) []string {
	out := []string{}
	for _, of := range e.Fields {
		out = append(out, node_to_str(of, indent))
	}
	return out
}

func (e ObjectExpression) Token() Token {
	return e.OpenBrace
}

func (e ObjectExpression) expressionNode() {}
