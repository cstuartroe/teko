#include <iostream>
#include <vector>
#include <fstream>
#include "tokenize.cpp"
#include "tagger.cpp"

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
    vector<Tag> tags = get_tags(tokens);
    for(int i = 0; i < tags.size(); i++) {
        cout << tags[i].to_str() << endl;
    }
}