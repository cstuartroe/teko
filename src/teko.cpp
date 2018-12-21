#include <iostream>
#include <vector>
#include <fstream>
#include "tokenize.cpp"

using namespace std;

int main(int argc, char *argv[])
{
    string filename;
    if (argc == 1) {
        cout << "Please give a filename." << endl;
        exit (EXIT_FAILURE);
    } else if (argc > 2) {
        cout << "Usage: teko <file.to>" << endl;
        exit (EXIT_FAILURE);
    }

    filename = argv[1];

    ifstream t;
    t.open(filename);
    string buffer;
    string line;
    while(true){
        if (t.eof()) {break;}
        getline(t, line);
        buffer += line + "\n";
    }
    t.close();

    vector<string> tokens = tokenize(buffer);
    for(int i = 0; i < tokens.size(); i++) {
        cout << tokens[i] << endl;
    }
}