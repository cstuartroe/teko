#include <stdio.h>
#include "general.cpp"
#include "parser.cpp"
#include "exec.cpp"

using namespace std;

int main(int argc, char *argv[]) {
    string filename;
    if (argc == 1) {
        printf("Please give a filename.\n");
        exit (EXIT_FAILURE);
    } else if (argc > 2) {
        printf("Usage: teko <file.to>\n");
        exit (EXIT_FAILURE);
    }

    filename = argv[1];

    TekoParser p = TekoParser(filename);
    printf("o\n");
    TekoParser *pp = &p;
    Interpreter i8er = Interpreter(pp);

    p.parse();
    i8er.execute(p.first_stmt);
}
