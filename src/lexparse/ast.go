package lexparse

import (
  "strings"
  "fmt"
)

type Node interface {
  ntype() string
  children() []Node
  token() Token
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

type ExpressionStatement struct {
  expression Expression
}

func (s ExpressionStatement) ntype() string {
  return "ExpressionStatement"
}

func (s ExpressionStatement) children() []Node {
  return []Node{s.expression}
}

func (s ExpressionStatement) token() Token {
  return s.expression.token()
}

func (s ExpressionStatement) child_strings(indent int) []string {
  return []string{
    node_to_str(s.expression, indent),
  }
}

func (s ExpressionStatement) statementNode() { }

//---

type SimpleExpression struct {
  Token Token
}

func (e SimpleExpression) ntype() string {
  return "SimpleExpression"
}

func (e SimpleExpression) children() []Node {
  return []Node{}
}

func (e SimpleExpression) token() Token {
  return e.Token
}

func (e SimpleExpression) child_strings(indent int) []string {
  return []string{
    e.Token.to_indented_str(indent),
  }
}

func (e SimpleExpression) expressionNode() { }

//---

type DeclarationExpression struct {
  tekotype Expression
  declareds []Declared
}

func (e DeclarationExpression) ntype() string {
  return "DeclarationExpression"
}

func (e DeclarationExpression) children() []Node {
  out := []Node{}
  for _, d := range e.declareds {
    out = append(out, d)
  }
  return out
}

func (e DeclarationExpression) token() Token {
  return e.tekotype.token()
}

func (e DeclarationExpression) child_strings(indent int) []string {
  out := []string{
    node_to_str(e.tekotype, indent),
  }
  for _, d := range e.declareds {
    out = append(out, node_to_str(d, indent))
  }
  return out
}

func (e DeclarationExpression) expressionNode() { }

//---

type Declared struct {
  symbol Token
  setter Token
  right Expression
}

func (d Declared) ntype() string {
  return "Declared"
}

func (d Declared) children() []Node {
  return []Node{d.right}
}

func (d Declared) child_strings(indent int) []string {
  return []string{
    d.symbol.to_indented_str(indent),
    d.setter.to_indented_str(indent),
    node_to_str(d.right, indent),
  }
}

func (d Declared) token() Token {
  return d.symbol
}
