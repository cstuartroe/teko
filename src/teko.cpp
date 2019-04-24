#include <stdio.h>
#include "general.cpp"
#include "tagger.cpp"
#include "nodes.cpp"
#include "parser.cpp"

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

    Parser p = Parser(filename);
    p.printout();
}
