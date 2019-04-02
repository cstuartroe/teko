#include <stdio.h>
#include <fstream>
#include "general.cpp"
#include "tokenize.cpp"

using namespace std;

int main(int argc, char *argv[]) {   
    string filename;
    if (argc == 1) {
        printf("Please give a filename.");
        exit (EXIT_FAILURE);
    } else if (argc > 2) {
        printf("Usage: teko <file.to>");
        exit (EXIT_FAILURE);
    }

    filename = argv[1];

    Tokenizer toker = Tokenizer();

    ifstream t;
    t.open(filename);
    string line;
    while(!t.eof()){
        getline(t,line);
        toker.digest_line(line);
    }
    t.close();

    Token *curr = toker.start;

    while (curr->next != 0) {
        printf("token: %s %d\n", curr->s.c_str(), curr->type);
        curr = curr->next;
    }
}