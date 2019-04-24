#include <fstream>

using namespace std;

char alpha_chars[54] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_";
char num_chars[11]   = "0123456789";
char hex_chars[17]   = "0123456789abcdef";
char white_chars[4]  = " \t\n";
char punct_chars[32] = "!\"#$%&\'()*+,-./:;<=>?@[\\]^`{|}~";

// givien an ascii char representing a digit, returns char with integer value
char ctoi(char c) {
	return c - 48;
}

char itoc(char i) {
	return i + 48;
}

bool in_charset(char c, char *charset) {
    while (*charset != 0) {
        if (c == *charset) { return true; }
        charset++;
    }
    return false;
}

bool in_stringset(string s, string *stringset, int size) {
    for (int i = 0; i < size; i++) {
        if (s == stringset[i]) { return true ;}
    }
    return false;
}

// accepts a hexidecimal digit and returns corresponding value
char xtoi(char c) {
    for (int i = 0; i < 16; i++) {
        if (hex_chars[i] == c) { return i; }
    }
    throw runtime_error("uh oh, how'd we get here? 1");
}

char string_index(string s, string *stringlist, int size) {
    for (char i = 0; i < size; i++) {
        if (s == stringlist[i]) { return i ;}
    }
    throw runtime_error("bad call to string_index");
    // return 0;
}

string teko_escape(string s) {
  string out = "";
  for (int i = 0; i < s.length(); i++) {
    switch (s[i]) {
      case '\n': out += "\\n"; break;
      case '\t': out += "\\t"; break;
      default: out += s[i]; break;
    }
  }
  return out;
}
