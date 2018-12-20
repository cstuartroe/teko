#include <iostream>
#include <vector>
#include <fstream>

using namespace std;

vector<string> tokenize(string content) {
	vector<string> tokens;
	for ( int i = 0 ; i < content.length(); i++) {
   		cout << content[i];
	}
	return tokens;
}

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
	while(t){
		getline(t, line);
		buffer += line + " ";
	}
	t.close();

	vector<string> tokens = tokenize(buffer);
}