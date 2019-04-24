#include <stdio.h>
#include <fstream>
#include "tagger.cpp"
#include "nodes.cpp"

void TekoParserException(string message, Tag t) {
    printf("Teko Parsing Error: %s\n", message.c_str());
    printf("at line %d, column %d\n", t.line_number, t.col);
    exit (EXIT_FAILURE);
}

struct TekoParser {
  Tag* first_tag = 0;
  Tag* curr_tag = 0;
  Statement *first_stmt = 0;
  Statement *curr_stmt = 0;

  TekoParser(string filename) {
    Tokenizer toker = Tokenizer();

    ifstream t;
    t.open(filename);
    string line;
    while(!t.eof()){
        getline(t,line);
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
    ExpressionNode *expr1 = grab_expression();
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

  ExpressionNode *grab_expression() {
    ExpressionNode *expr;
    Tag *expr_first_tag = curr_tag;

    switch (curr_tag->type) {
      case LabelTag: expr = new SimpleNode(curr_tag); advance(); break;
      case StringTag:expr = new SimpleNode(curr_tag); advance(); break;
      case CharTag:  expr = new SimpleNode(curr_tag); advance(); break;
      case IntTag:   expr = new SimpleNode(curr_tag); advance(); break;
      case RealTag:  expr = new SimpleNode(curr_tag); advance(); break;
      case BoolTag:  expr = new SimpleNode(curr_tag); advance(); break;
      default: TekoParserException("Illegal start to expression: " + curr_tag->s, *curr_tag);
    }

    expr->first_tag = expr_first_tag;
    return expr;
  }

  DeclarationNode *grab_declaration() {
    return new DeclarationNode();
  }
};
