#include "nodes.h"

enum StatementType {DeclarationStmtType, AssignmentStmtType, ExpressionStmtType};

const int num_expression_types = 13;

enum ExpressionType {SimpleExpr, CallExpr, AttrExpr, SliceExpr,

                     SeqExpr, MapExpr, StructExpr, ArgExpr,

                     BinopExpr, ComparisonExpr, ConversionExpr, NotExpr,

                     IfExpr, ForExpr, WhileExpr};

struct Node {
  Tag *first_tag;
};

struct Statement : Node {
  Statement *next = 0;
  StatementType stmt_type;
};

// ------------

struct ExpressionNode : Node{
  ExpressionType expr_type;
};

struct ExpressionStmt : Statement {
  StatementType stmt_type = ExpressionStmtType;
  ExpressionNode *body = 0;
};

struct DeclarationNode : Node {
  DeclarationNode *next = 0;
  SimpleNode *label = 0;
  StructNode *argstruct = 0;
  ExpressionNode *value = 0;
};

struct DeclarationStmt : Statement {
  StatementType stmt_type = DeclarationStmtType;
  ExpressionNode *tekotype = 0;
  DeclarationNode *declist = 0;
};

struct AssignmentStmt : Statement {
  StatementType stmt_type = AssignmentStmtType;
  ExpressionNode *left = 0;
  char setter = 0;
  ExpressionNode *right = 0;
};

// ------------

enum SimpleExpressionType {LabelExpr, StringExpr, CharExpr, RealExpr, IntExpr, BoolExpr};

struct SimpleNode : ExpressionNode{
  ExpressionType expr_type = SimpleExpr;
  SimpleExpressionType data_type;
  char *val = 0; // will be cast

  SimpleNode(Tag *t) {
    first_tag = t;
    val = t->val;
    switch (t->type) {
      case LabelTag:  data_type = LabelExpr;  break;
      case StringTag: data_type = StringExpr; break;
      case CharTag:   data_type = CharExpr;   break;
      case RealTag:   data_type = RealExpr;   break;
      case IntTag:    data_type = IntExpr;    break;
      case BoolTag:   data_type = BoolExpr;   break;
      default:        throw runtime_error("Invalid simple data type: " + to_string((char) t->type));
    }
  }
};

struct CallNode : ExpressionNode {
  ExpressionType expr_type = CallExpr;
  ExpressionNode *left = 0;
  ArgNode *args = 0;
};

struct AttrNode : ExpressionNode {
  ExpressionType expr_type = AttrExpr;
  ExpressionNode *left = 0;
  string label;
};

struct SliceNode : ExpressionNode {
  ExpressionType expr_type = SliceExpr;
  ExpressionNode *left = 0;
  Brace brace;
  ExpressionNode *slice = 0;
};

// ------------

struct SeqNode : ExpressionNode {
  ExpressionType expr_type = SeqExpr;
  Brace brace;
  ExpressionNode *first = 0;
  SeqNode *next = 0;
};

struct MapNode : ExpressionNode {
  ExpressionType expr_type = MapExpr;
  ExpressionNode *key = 0;
  ExpressionNode *value = 0;
  MapNode *next = 0;
};

struct StructNode : ExpressionNode {
  ExpressionType expr_type = StructExpr;
  ExpressionNode *tekotype = 0;
  string label;
  ExpressionNode *deflt = 0;
  StructNode *next = 0;
};

struct ArgNode : ExpressionNode {
  ExpressionType expr_type = ArgExpr;
  string label;
  ExpressionNode *value = 0;
};

// ------------

struct BinopNode : ExpressionNode {
  ExpressionType expr_type = BinopExpr;
  char binop;
  ExpressionNode *left = 0;
  ExpressionNode *right = 0;
};

struct ComparisonNode : ExpressionNode {
  ExpressionType expr_type = ComparisonExpr;
  char comp;
  ExpressionNode *left = 0;
  ExpressionNode *right = 0;
};

struct ConversionNode : ExpressionNode {
  ExpressionType expr_type = ConversionExpr;
  char conv;
  ExpressionNode *left = 0;
};

struct NotNode : ExpressionNode {
  ExpressionNode *right = 0;
};

// ------------

struct IfNode : ExpressionNode {
    ExpressionNode *condition = 0;
    Statement *then_block = 0;
    IfNode *else_stmt = 0;
};

struct ForNode : ExpressionNode {
    ExpressionNode *type = 0;
    string label;
    ExpressionNode *iterator = 0;
    Statement *codeblock = 0;
};

struct WhileNode : ExpressionNode {
    ExpressionNode *condition = 0;
    Statement *codeblock = 0;
};
