#include <stdio.h>
#include <fstream>
#include "tagger.cpp"
#include "nodes.cpp"

enum Precedence {NoPrec, Compare, BoolOp, AddSub, MultDiv, Exp, TopPrec};

Precedence getPrec(string infix) {
  if      (infix == "+")  { return AddSub; }
  else if (infix == "-")  { return AddSub; }
  else if (infix == "*")  { return MultDiv; }
  else if (infix == "/")  { return MultDiv; }
  else if (infix == "^")  { return Exp; }
  else if (infix == "%")  { return MultDiv; }
  else if (infix == "&")  { return BoolOp; }
  else if (infix == "|")  { return BoolOp; }
  else if (infix == "in") { return BoolOp; }
  else if (infix == "to") { return BoolOp; }
  else if (infix == "==") { return Compare; }
  else if (infix == "!=") { return Compare; }
  else if (infix == "<")  { return Compare; }
  else if (infix == ">")  { return Compare; }
  else if (infix == "<=") { return Compare; }
  else if (infix == ">=") { return Compare; }
  else if (infix == "<:") { return Compare; }
  else { throw runtime_error("Invalid infix: " + infix); }
}

struct TekoParser {
  Tag* first_tag = 0;
  Tag* curr_tag = 0;
  vector<string> lines;
  Statement *first_stmt = 0;
  Statement *curr_stmt = 0;

  TekoParser(string filename) {
    Tokenizer toker = Tokenizer();

    ifstream t;
    t.open(filename);
    string line;
    while(!t.eof()){
        getline(t,line);
        lines.push_back(line);
        toker.digest_line(line);
    }
    t.close();

    Token *curr_token = toker.start;

    while (curr_token->next != 0) {
        Tag *new_tag = from_token(*curr_token);
        if (curr_tag != 0) {
          curr_tag->next = new_tag;
        } else {
          first_tag = new_tag;
        }
        curr_tag = new_tag;
        curr_token = curr_token->next;
    }
  }

  void TekoParserException(string message, Tag t) {
      printf("%s\n", string(12,'-').c_str());
      printf("Teko Parsing Error: %s\n", message.c_str());
      printf("at line %d, column %d\n", t.line_number, t.col);
      printf("    %s\n", lines[t.line_number-1].c_str());
      printf("%s\n", (string(t.col+3,' ') + string((int) t.s.size(), '~')).c_str());
      exit (EXIT_FAILURE);
  }

  void printout() {
    curr_tag = first_tag;
    while (curr_tag != 0) {
      printf("%s\n", curr_tag->to_str().c_str());
      curr_tag = curr_tag->next;
    }
  }

  void parse() {
    curr_tag = first_tag;
    while (curr_tag != 0) {
      Statement *new_stmt = grab_statement();
      if (first_stmt == 0) {
        first_stmt = new_stmt;
      } else {
        curr_stmt->next = new_stmt;
      }
      curr_stmt = new_stmt;
      printf("%s", curr_stmt->to_str(0).c_str());
    }
  }

  Tag *advance() {
    Tag *prev = curr_tag;
    if (curr_tag != 0) {
      curr_tag = curr_tag->next;
    }
    return prev;
  }

  Statement *grab_statement() {
    ExpressionNode *expr1 = grab_expression(NoPrec);
    Statement *stmt;
    if (curr_tag->type == LabelTag) {
      DeclarationStmt *decl_stmt = new DeclarationStmt();
      decl_stmt->first_tag = expr1->first_tag;
      decl_stmt->tekotype = expr1;
      decl_stmt->declist = grab_declaration();
      stmt = decl_stmt;
    } else {
      ExpressionStmt *expr_stmt = new ExpressionStmt();
      expr_stmt->first_tag = expr1->first_tag;
      expr_stmt->body = expr1;
      stmt = expr_stmt;
    }

    if (curr_tag->type == SemicolonTag) {
      advance();
      return stmt;
    } else {
      TekoParserException("Expected semicolon", *curr_tag);
    }
  }

  DeclarationNode *grab_declaration() {
    return new DeclarationNode();
  }

  ExpressionNode *grab_expression(Precedence prec) {
    ExpressionNode *expr;
    Tag *expr_first_tag = curr_tag;

    switch (curr_tag->type) {
      case LabelTag:  expr = new SimpleNode(curr_tag); advance(); break;
      case StringTag: expr = new SimpleNode(curr_tag); advance(); break;
      case CharTag:   expr = new SimpleNode(curr_tag); advance(); break;
      case IntTag:    expr = new SimpleNode(curr_tag); advance(); break;
      case RealTag:   expr = new SimpleNode(curr_tag); advance(); break;
      case BoolTag:   expr = new SimpleNode(curr_tag); advance(); break;
      case OpenTag:   {
        if (*((Brace*)curr_tag->val) == paren) {
          expr = grab_struct();
          if (expr == 0) {
            expr = grab_args();
          }
        } else {
          expr = grab_sequence();
        }
      } break;
      case PrefixTag: {
        PrefixNode *prefix_expr = new PrefixNode();
        prefix_expr->first_tag = curr_tag;
        prefix_expr->prefix = *curr_tag->val;
        prefix_expr->right = grab_expression(TopPrec);
        expr = prefix_expr;
      } break;
      default: TekoParserException("Illegal start to expression: " + curr_tag->s, *curr_tag);
    }
    expr->first_tag = expr_first_tag;

    expr = continue_expression(prec, expr);
    return expr;
  }

  ExpressionNode *continue_expression(Precedence prec, ExpressionNode *left_expr) {
    if (curr_tag->type == AttrTag) {
      if (curr_tag->next->type == LabelTag) {
        AttrNode *attr_expr = new AttrNode();
        attr_expr->first_tag = curr_tag;
        attr_expr->left = left_expr;
        advance();
        attr_expr->label = curr_tag->s;
        advance();
        return continue_expression(prec, attr_expr);
      } else {
        curr_tag->type = SuffixTag;
        char *suf = new char(string_index(".", suffixes, num_suffixes));
        curr_tag->val = suf;
      }
    }

    switch (curr_tag->type) {
      case InfixTag: {
        Precedence new_prec = getPrec(infixes[*curr_tag->val]);
        if (new_prec <= prec) {
          return left_expr;
        } else {
          InfixNode *infix_expr = new InfixNode();
          infix_expr->first_tag = curr_tag;
          infix_expr->left = left_expr;
          infix_expr->infix = *curr_tag->val;
          advance();
          infix_expr->right = grab_expression(new_prec);
          return continue_expression(prec, infix_expr);
        }
      } break;

      case OpenTag: {
        Brace b = *((Brace *) curr_tag->val);
        switch (b) {
          case paren: {
            CallNode *call_expr = new CallNode();
            call_expr->first_tag = left_expr->first_tag;
            call_expr->left = left_expr;
            call_expr->args = grab_args();
            return continue_expression(prec, call_expr);
          } break;

          default: {
            SliceNode *slice_expr = new SliceNode();
            slice_expr->first_tag = curr_tag;
            slice_expr->left = left_expr;
            slice_expr->slice = grab_sequence();
            return continue_expression(prec, slice_expr);
          } break;
        }
      } break;

      case SuffixTag: {
        SuffixNode *suff_expr = new SuffixNode;
        suff_expr->first_tag = curr_tag;
        suff_expr->left = left_expr;
        suff_expr->suffix = *curr_tag->val;
        advance();
        return continue_expression(prec, suff_expr);
      }

      default: return left_expr;
    }
  }

  SeqNode *grab_sequence() {
    SeqNode *seq_expr = new SeqNode();
    if (curr_tag->type != OpenTag) { throw runtime_error("erroneously called grab_sequence"); }
    seq_expr->first_tag = curr_tag;
    seq_expr->brace = *((Brace*) curr_tag->val);
    advance();

    while (curr_tag->type != CloseTag) {
      ExpressionNode *elem = grab_expression(NoPrec);
      if (curr_tag->type == ColonTag) {
        MapNode *map_expr = new MapNode();
        map_expr->first_tag = elem->first_tag;
        map_expr->key = elem;
        advance();
        map_expr->value = grab_expression(NoPrec);
        elem = map_expr;
      }
      seq_expr->elems.push_back(elem);

      if (curr_tag->type == CommaTag) {
        advance();
      } else if (curr_tag->type != CloseTag) {
        TekoParserException("Expected a comma", *curr_tag);
      }
    }

    if (*((Brace*) curr_tag->val) != seq_expr->brace) {
      TekoParserException("Mismatched brace", *curr_tag);
    } else {
      advance();
    }

    return seq_expr;
  }

  ArgsNode *grab_args() {
    ArgsNode *args_expr = new ArgsNode();
    if (curr_tag->type != OpenTag || *((Brace*) curr_tag->val) != paren) { throw runtime_error("erroneously called grab_args"); }
    args_expr->first_tag = curr_tag;
    advance();

    while (curr_tag->type != CloseTag) {
      ArgNode *arg = new ArgNode();
      ExpressionNode* expr = grab_expression(NoPrec);
      if (expr->expr_type == SimpleExpr && ((SimpleNode*) expr)->data_type == LabelExpr && curr_tag->type == SetterTag && setters[*curr_tag->val] == "=") {
        arg->label = ((SimpleNode*) expr)->to_str(0);
        advance();
        expr = grab_expression(NoPrec);
      }
      arg->value = expr;
      args_expr->args.push_back(arg);

      if (curr_tag->type == CommaTag) {
        advance();
      } else if (curr_tag->type != CloseTag) {
        TekoParserException("Expected a comma", *curr_tag);
      }
    }

    if (*((Brace*) curr_tag->val) != paren) {
      TekoParserException("Mismatched brace", *curr_tag);
    } else {
      advance();
    }

    return args_expr;
  }

  StructNode *grab_struct() {
    Tag* orig_tag = curr_tag;

    StructNode *struct_expr = new StructNode();
    if (curr_tag->type != OpenTag || *((Brace*) curr_tag->val) != paren) { throw runtime_error("erroneously called grab_args"); }
    struct_expr->first_tag = curr_tag;
    advance();

    ExpressionNode* expr = grab_expression(NoPrec);

    if (curr_tag->type != LabelTag) {
      curr_tag = orig_tag;
      return 0;
    } else {
      curr_tag = orig_tag;
      advance();
    }

    while (curr_tag->type != CloseTag) {
      StructElemNode *elem = new StructElemNode();
      elem->tekotype = grab_expression(NoPrec);
      if (curr_tag->type != LabelTag) {
        TekoParserException("Expected label", *curr_tag);
      } else {
        elem->label = curr_tag->s;
        advance();
      }

      if (curr_tag->type == QMarkTag) {
        advance();
        elem->deflt = grab_expression(NoPrec);
      }

      struct_expr->elems.push_back(elem);

      if (curr_tag->type == CommaTag) {
        advance();
      } else if (curr_tag->type != CloseTag) {
        TekoParserException("Expected a comma", *curr_tag);
      }
    }

    if (*((Brace*) curr_tag->val) != paren) {
      TekoParserException("Mismatched brace", *curr_tag);
    } else {
      advance();
    }

    return struct_expr;
  }
};
