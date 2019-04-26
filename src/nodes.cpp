#include <vector>
#include "nodes.h"

using namespace std;

const int indent_width = 4;

enum StatementType {DeclarationStmtType, NamespaceStmtType, ExpressionStmtType};

const int num_expression_types = 14;

enum ExpressionType {SimpleExpr, CallExpr, AttrExpr, SliceExpr,

                     SeqExpr, MapExpr, StructExpr, ArgExpr,

                     PrefixExpr, InfixExpr, AssignmentExpr, SuffixExpr,

                     CodeblockExpr, IfExpr, ForExpr, WhileExpr};

struct Node {
    Tag *first_tag;

    virtual string to_str(int indent) {
        return "";
    }
};

struct Statement : Node {
    Statement *next = 0;
    StatementType stmt_type;
};

// ------------

struct ExpressionNode : Node {
    ExpressionType expr_type;
};

struct ExpressionStmt : Statement {
    StatementType stmt_type = ExpressionStmtType;
    ExpressionNode *body = 0;

    string to_str(int indent) {
        string s = string(indent*indent_width, ' ');
        s += body->to_str(indent) + ";\n";
        return s;
    }
};

// ------------

enum SimpleExpressionType {LabelExpr, StringExpr, CharExpr, RealExpr, IntExpr,
                           BoolExpr, BitsExpr, BytesExpr};

struct SimpleNode : ExpressionNode {
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
        case BitsTag:   data_type = BitsExpr;   break;
        case BytesTag:  data_type = BytesExpr;  break;
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
        case BitsExpr:   s = *((string*) val); break;
        case BytesExpr:  s = *((string*) val); break;
        }
        return s;
    }
};

struct AttrNode : ExpressionNode {
    ExpressionType expr_type = AttrExpr;
    ExpressionNode *left = 0;
    string label;

    string to_str(int indent) {
        return "(" + left->to_str(indent) + "." + label + ")";
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

    string to_str(int indent) {
        string s = tekotype->to_str(indent);
        s += " " + label;
        if (deflt != 0) {
            s += " ? " + deflt->to_str(indent);
        }
        return s;
    }
};

struct StructNode : ExpressionNode {
    ExpressionType expr_type = StructExpr;
    vector<StructElemNode*> elems;

    string to_str(int indent) {
        string s = "(";
        for (StructElemNode *elem: elems) {
            s += elem->to_str(indent) + ", ";
        }
        if (elems.size() > 0) {
            s.pop_back();
            s.pop_back();
        }
        s += ")";
        return s;
    }
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
        if (args.size() > 0) {
            s.pop_back();
            s.pop_back();
        }
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

struct AssignmentNode : ExpressionNode {
    ExpressionType expr_type = AssignmentExpr;
    ExpressionNode *left = 0;
    char setter = 0;
    ExpressionNode *right = 0;

    string to_str(int indent) {
        string s = "(" + left->to_str(indent);
        s += " " + setters[setter] + " ";
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

struct Codeblock : ExpressionNode {
    ExpressionType expr_type = CodeblockExpr;
    vector<Statement*> stmts;

    string to_str(int indent) {
        string s = "{\n";
        for (Statement* stmt : stmts) {
            s += stmt->to_str(indent+1);
        }
        s += string(indent*indent_width, ' ') + "}";
        return s;
    }
};

struct IfNode : ExpressionNode {
    ExpressionType expr_type = IfExpr;
    ExpressionNode *condition = 0;
    Codeblock *then_block = 0;
    Codeblock *else_block = 0;

    string to_str(int indent) {
        string s = "if ";
        s += condition->to_str(indent) + " ";
        s += then_block->to_str(indent);
        if (else_block != 0) {
            s += else_block->to_str(indent);
        }
        return s;
    }
};

struct ForNode : ExpressionNode {
    ExpressionType expr_type = ForExpr;
    ExpressionNode *tekotype = 0;
    string label;
    ExpressionNode *iterator = 0;
    Codeblock *codeblock = 0;

    string to_str(int indent) {
        string s = "for (";
        s += (tekotype == 0) ? "" : (tekotype->to_str(indent) + " ");
        s += label;
        s += " in " + iterator->to_str(indent);
        s += ") " + codeblock->to_str(indent);
        return s;
    }
};

struct WhileNode : ExpressionNode {
    ExpressionType expr_type = WhileExpr;
    ExpressionNode *condition = 0;
    Codeblock *codeblock = 0;

    string to_str(int indent) {
        string s = "while (" + condition->to_str(indent) + ") ";
        s += codeblock->to_str(indent);
        return s;
    }
};

// ------------

struct DeclarationNode : Node {
    string label;
    StructNode *argstruct = 0;
    char setter;
    ExpressionNode *value = 0;

    string to_str(int indent) {
        string s = label;
        if (argstruct != 0) {
            s += argstruct->to_str(indent);
        }

        if (value != 0) {
            s += " " + setters[setter] + " " + value->to_str(indent);
        }
        return s;
    }
};

struct AnnotationNode : Node {
    vector<ExpressionNode*> params;

    string to_str(int indent) {
        string s = "";
        for (ExpressionNode *param : params) {
            s += " " + param->to_str(indent) + ",";
        }
        if (params.size() > 0) {
            s.pop_back();
        }
        return s;
    }
};

struct DeclarationStmt : Statement {
    StatementType stmt_type = DeclarationStmtType;
    AnnotationNode *annots[num_annotations] = {0};
    bool vts[num_vartypes] = {false};
    ExpressionNode *tekotype = 0;
    vector<DeclarationNode*> declist;

    string to_str(int indent) {
        string s = string(indent*indent_width, ' ');;
        for (int i = 0; i < num_annotations; i++) {
            if (annots[i] != 0) {
                s += annotations[i] + annots[i]->to_str(indent) + "\n";
            }
        }

        bool has_vt = false;
        for (int i = 0; i < num_vartypes; i++) {
            if (vts[i]) {
                has_vt = true;
                s += vartypes[i] + " ";
            }
        }

        if (tekotype != 0) {
            s += tekotype->to_str(indent) + " ";
        } else if (!has_vt) {
            s += "let ";
        }
        for (DeclarationNode *decl_expr: declist) {
            s += decl_expr->to_str(indent) + ", ";
        }
        s.pop_back();
        s.pop_back();
        s += ";\n";
        return s;
    }
};

struct NamespaceStmt : Statement {
    StatementType stmt_type = NamespaceStmtType;
    string label;
    Codeblock *codeblock;

    string to_str(int indent) {
        string s = string(indent*indent_width, ' ');
        s += "namespace " + label + " ";
        s += codeblock->to_str(indent);
        s += ";\n";
        return s;
    }
};
