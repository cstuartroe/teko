#include <stdio.h>
#include "general.cpp"
#include "parser.cpp"
#include "types.cpp"

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
    p.parse();

    TekoTypeChecker ttc = TekoTypeChecker();
    ttc.check(p.first_stmt);
}
