package lexparse

import (
  "strings"
  "fmt"
)

type Node interface {
  ntype() string
  children() []Node
  token() Token
}

func PrintNode(node Node, indent int) {
  space := strings.Repeat(" ", indent*2)
  fmt.Printf("%s%s\n", space, node.ntype())
  fmt.Printf("%s%s\n", space, node.token().to_str())
  for _, child := range node.children() {
    PrintNode(child, indent+1)
  }
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

func (d Declared) token() Token {
  return d.symbol
}
