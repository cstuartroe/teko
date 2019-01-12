#include <iostream>
#include <vector>
#include <fstream>
#include "tokenize.cpp"
#include "tagger.cpp"
#include "parser.cpp"
#include "types.cpp"

using namespace std;

int main(int argc, char *argv[]) {   
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
    vector<Line> lines = get_lines(tags);
    for (int i = 0; i < lines.size(); i++) {
        vector<Tag> line_tags = lines[i].tags;
        for (int j = 0; j < line_tags.size(); j++) {
            cout << line_tags[j].to_str() << endl;
        }
        cout << endl;
    }
}