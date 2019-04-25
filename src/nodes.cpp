#include <vector>
#include "nodes.h"

using namespace std;

const int indent_width = 4;

enum StatementType {DeclarationStmtType, AssignmentStmtType, ExpressionStmtType};

const int num_expression_types = 14;

enum ExpressionType {SimpleExpr, CallExpr, AttrExpr, SliceExpr,

                     SeqExpr, MapExpr, StructExpr, ArgExpr,

                     PrefixExpr, InfixExpr, SuffixExpr,

                     IfExpr, ForExpr, WhileExpr};

struct Node {
  Tag *first_tag;

  virtual string to_str(int indent) { return ""; }
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

  string to_str(int indent) {
    string s = "expr:\n";
    s += string(indent*indent_width, ' ');
    s += body->to_str(indent) + ";\n";
    return s;
  }
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

  string to_str(int indent) {
    string s;
    switch (data_type) {
      case LabelExpr:  s = *((string*) val); break;
      case StringExpr: s = "\"" + teko_escape(*((string*) val)) + "\""; break;
      case CharExpr:   s = "\'" + teko_escape(to_string(val[0])) + "\'"; break;
      case IntExpr:    s = to_string(*((int*) val)); break;
      case RealExpr:   s = to_string(*((float*) val)); break;
      case BoolExpr:   s = *((bool*) val) ? "true" : "false"; break;
    }
    return s;
  }
};

struct AttrNode : ExpressionNode {
  ExpressionType expr_type = AttrExpr;
  ExpressionNode *left = 0;
  string label;

  string to_str(int indent) {
    return left->to_str(indent) + "." + label;
  }
};

// ------------

struct MapNode : ExpressionNode {
  ExpressionType expr_type = MapExpr;
  ExpressionNode *key = 0;
  ExpressionNode *value = 0;

  string to_str(int indent) {
    return key->to_str(indent) + ": " + value->to_str(indent);
  }
};

struct SeqNode : ExpressionNode {
  ExpressionType expr_type = SeqExpr;
  Brace brace;
  vector<ExpressionNode*> elems;

  string to_str(int indent) {
    string s = to_string(brace, true);
    for (ExpressionNode *expr : elems) {
      s += expr->to_str(indent) + ", ";
    }
    s.pop_back();
    s.pop_back();
    s += to_string(brace, false);
    return s;
  }
};

struct SliceNode : ExpressionNode {
  ExpressionType expr_type = SliceExpr;
  ExpressionNode *left = 0;
  SeqNode *slice = 0;

  string to_str(int indent) {
    return left->to_str(indent) + slice->to_str(indent);
  }
};

struct StructElemNode : Node {
  ExpressionNode *tekotype = 0;
  string label = "";
  ExpressionNode *deflt = 0;
};

struct StructNode : ExpressionNode {
  ExpressionType expr_type = StructExpr;
  vector<StructElemNode*> elems;
};

struct ArgNode : Node {
  string label = "";
  ExpressionNode *value = 0;

  string to_str(int indent) {
    if (label == "") {
      return value->to_str(indent);
    } else {
      return label + " = " + value->to_str(indent);
    }
  }
};

struct ArgsNode : ExpressionNode {
  ExpressionType expr_type = ArgExpr;
  vector<ArgNode*> args;

  string to_str(int indent) {
    string s = "(";
    for (ArgNode *arg: args) {
      s += arg->to_str(indent) + ", ";
    }
    s.pop_back();
    s.pop_back();
    s += ")";
    return s;
  }
};

struct CallNode : ExpressionNode {
  ExpressionType expr_type = CallExpr;
  ExpressionNode *left = 0;
  ArgsNode *args = 0;

  string to_str(int indent) {
    return left->to_str(indent) + args->to_str(indent);
  }
};

// ------------

struct PrefixNode : ExpressionNode {
  ExpressionType expr_type = PrefixExpr;
  char prefix;
  ExpressionNode *right = 0;
};

struct InfixNode : ExpressionNode {
  ExpressionType expr_type = InfixExpr;
  ExpressionNode *left = 0;
  char infix;
  ExpressionNode *right = 0;

  string to_str(int indent) {
    string s = "(" + left->to_str(indent);
    s += " " + infixes[infix] + " ";
    s += right->to_str(indent) + ")";
    return s;
  }
};

struct SuffixNode : ExpressionNode {
  ExpressionType expr_type = SuffixExpr;
  ExpressionNode *left = 0;
  char suffix;

  string to_str(int indent) {
    string s = left->to_str(indent);
    s += suffixes[suffix];
    return s;
  }
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
