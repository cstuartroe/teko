package lexparse

import (
  "fmt"
  "strings"
)

type NodeType int

const (
  ModuleN NodeType = iota

  DeclarationS
  ExpressionS
  TypeDefnS

  ArrayTypeN
  SetTypeN
  MapTypeN
  FunctionTypeN

  DeclaredChainN
  DeclaredN
  DeclaredLeftN
  FunctionArgdefN

  SimpleExprN
)

var simpleExprTokenTypes map[tokenType]bool = map[tokenType]bool {
  SymbolT: true,
  StringT: true,
  CharT: true,
  IntT: true,
  FloatT: true,
  BoolT: true,
  NullT: true,
}

// How the Node struct is used:
//
// token is the symbol token/int token/etc for SimpleExprN,
// and the binop/comparison/setter/prefix/suffix token for those
// nodetypes. Otherwise it's usually just the first token of
// a node and not important aside for parse error printouts.
//
// left and right are obviously used in branching, but some nodetypes
// only have one descendant, e.g. ExpressionS, PrefixN, SuffixN.
// For node types with one descendant, only left is used.
// SimpleExprN is the only node type that uses neither left or right.

type Node struct {
  ntype NodeType
  token *Token
  left *Node
  right *Node
}

func (node Node) Printout(indent int) {
  space := strings.Repeat(" ", indent*2)
  fmt.Printf("%s%d\n", space, node.ntype)
  fmt.Printf("%s%s\n", space, node.token.to_str())
  if node.left != nil {
    node.left.Printout(indent+1)
  }
  if node.right != nil {
    node.right.Printout(indent+1)
  }
}

type Parser struct {
  tokens []Token
  position int
  isTypeName func(name string) bool
}

func (parser *Parser) LexFile(filename string) {
  parser.tokens = LexFile(filename)
}

func NewParser(isTypeName func(string) bool) Parser {
  return Parser{
    tokens: []Token{},
    position: 0,
    isTypeName: isTypeName,
  }
}

func (parser *Parser) next() *Token {
  if parser.HasMore() {
    return &(parser.tokens[parser.position])
  } else {
    t := parser.tokens[parser.position-1]
    fmt.Println("ruh roh")
    fmt.Println(t.to_str())
    lexparsePanic(t.Line, t.Col + len(t.Value), 1, "Unexpected EOF")
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
  if parser.next().TType != ttype {
    tokenPanic(*parser.next(), fmt.Sprintf("Expected %s", ttype))
  }
}

// I'm attempting to make the teko parser without needing lookahead
// We'll see how sustainable that is

func (parser *Parser) GrabStatement() Node {
  if parser.isTypeName(string(parser.next().Value)) {
    return parser.grabDeclaration()
  } else {
    return parser.grabExpressionStmt()
  }
}

func (parser *Parser) grabDeclaration() Node {
  return Node {
    ntype: DeclarationS,
    token: parser.next(),
    left: parser.grabType(),
    right: parser.grabDeclaredChain(),
  }
}

func (parser *Parser) grabType() *Node {
  if parser.next().TType != SymbolT {
    tokenPanic(*parser.next(), "Illegal start to type")
  }

  typenode := Node{ntype: SimpleExprN, token: parser.next()}
  parser.Advance()

  // TODO: arrays, sets, maps, functions

  return &typenode
}

func (parser *Parser) grabDeclaredChain() *Node {
  left := parser.grabDeclared()
  var right *Node

  if parser.HasMore() && (parser.next().TType == CommaT) {
    parser.Advance()
    right = parser.grabDeclaredChain()
  } else {
    right = nil
  }

  return &Node {
    ntype: DeclaredChainN,
    token: left.token,
    left: left,
    right: right,
  }
}

func (parser *Parser) grabDeclared() *Node {
  if parser.next().TType != SymbolT {
    tokenPanic(*parser.next(), "Must give name of declared variable next")
  }

  left := parser.grabSimpleExpression()

  // TODO: function argdef

  parser.Expect(SetterT); setter_t := parser.next(); parser.Advance()

  right := parser.grabExpression()

  return &Node{
    ntype: DeclaredN,
    token: setter_t,
    left: left,
    right: right,
  }
}

func (parser *Parser) grabExpressionStmt() Node {
  return Node{
    ntype: ExpressionS,
    token: parser.next(),
    left: parser.grabExpression(),
  }
}

func (parser *Parser) grabExpression() *Node {
  return parser.grabSimpleExpression()
}

func (parser *Parser) grabSimpleExpression() *Node {
  if _, ok := simpleExprTokenTypes[parser.next().TType]; ok {
    n := Node{
      ntype: SimpleExprN,
      token: parser.next(),
    }
    parser.Advance()
    return &n
  } else {
    tokenPanic(*parser.next(), "Illegal start to simple expression")
    return nil // unreachable code that the compiler requires
  }
}
