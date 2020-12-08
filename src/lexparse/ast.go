package lexparse

import (
  "strings"
  "fmt"
)

type Node interface {
  ntype() string
  children() []Node
  Token() Token
  child_strings(indent int) []string
}

func node_to_str(node Node, indent int) string {
  space := strings.Repeat(" ", indent*INDENT_AMOUNT)
  out := fmt.Sprintf("%s%s\n", space, node.ntype())
  for _, s := range node.child_strings(indent+1) {
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

func (c *Codeblock) ntype() string {
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

func (s ExpressionStatement) ntype() string {
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

func (s ExpressionStatement) statementNode() { }

//---

type SimpleExpression struct {
  token Token
}

func (e SimpleExpression) ntype() string {
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

func (e SimpleExpression) expressionNode() { }

//---

type DeclarationExpression struct {
  Tekotype Expression
  Declareds []Declared
}

func (e DeclarationExpression) ntype() string {
  return "DeclarationExpression"
}

func (e DeclarationExpression) children() []Node {
  out := []Node{}
  for _, d := range e.Declareds {
    out = append(out, d)
  }
  return out
}

func (e DeclarationExpression) Token() Token {
  return e.Tekotype.Token()
}

func (e DeclarationExpression) child_strings(indent int) []string {
  out := []string{
    node_to_str(e.Tekotype, indent),
  }
  for _, d := range e.Declareds {
    out = append(out, node_to_str(d, indent))
  }
  return out
}

func (e DeclarationExpression) expressionNode() { }

//---

type Declared struct {
  Symbol Token
  Setter Token
  Right Expression
}

func (d Declared) ntype() string {
  return "Declared"
}

func (d Declared) children() []Node {
  return []Node{d.Right}
}

func (d Declared) child_strings(indent int) []string {
  return []string{
    d.Symbol.to_indented_str(indent),
    d.Setter.to_indented_str(indent),
    node_to_str(d.Right, indent),
  }
}

func (d Declared) Token() Token {
  return d.Symbol
}

//---

type CallExpression struct {
  Receiver Expression
  Args []Expression
  Kwargs []FunctionKwarg
  ResolvedArgs map[string]Expression
}

func (e CallExpression) ntype() string {
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

func (e CallExpression) expressionNode() { }

//---

type FunctionKwarg struct {
  Symbol Token
  Value Expression
}

func (k FunctionKwarg) ntype() string {
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