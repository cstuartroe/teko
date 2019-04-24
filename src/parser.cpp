#include <stdio.h>
#include <fstream>

enum LineType { DeclarationLine, AssignmentLine, ExpressionLine,
                IfStatementLine, WhileBlockLine, ForBlockLine };

struct Parser {
  Tag* first_tag;
  Tag* curr_tag;
  Statement *first_stmt;
  Statement *curr_stmt;

  Parser(string filename) {
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
    curr_tag = 0;

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
};
