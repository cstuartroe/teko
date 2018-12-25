#include <iostream>
#include <vector>
#include <set>
#include "general.cpp"

using namespace std;

set<char> alpha {'a','b','c','d','e','f','g','h','i','j',
                'k','l','m','n','o','p','q','r','s','t',
                'u','v','w','x','y','z','A','B','C','D',
                'E','F','G','H','I','J','K','L','M','N',
                'O','P','Q','R','S','T','U','V','W','X',
                'Y','Z','_'};
set<char> nums {'0','1','2','3','4','5','6','7','8','9'};
set<char> white {' ','\t','\n'};
set<char> punct {'!','"','#','$','%','&','\'','(',')','*',
                '+',',','-','.','/',':',';','<','=','>',
                '?','@','[','\\',']','^','`','|','{','|',
                '}','~'};
set<string> punct_combos {"==","<=",">=","!=","<:","+=",
                "-=","*=","/=","%%="};

bool in_charset(char c, set<char> charset) {
    return charset.find(c) != charset.end();
}

void grab_comment(string &source) {
    if (source.substr(0,2) != "//") { 
        throw runtime_error("comment handling snafu"); 
    }
    source = source.substr(2);
    while (source.length() != 0 && source[0] != '\n') {
        source = source.substr(1);
    }
    if (source.length() != 0) {source = source.substr(1);}
}

void grab_str(string &source, vector<string> &tokens) {
    if (source[0] != '"') {
        throw runtime_error("string handling snafu");
    }
    string token = "\"";
    source = source.substr(1);
    bool stop = false;
    while (!stop) {
        if (source.length() == 0) {
            compiler_error("EOF while parsing string");
        }
        if (source.substr(0,2) == "\\\"") {
            token += "\\\"";
            source = source.substr(2);
        } else {
            stop = (source[0] == '\"');
            token += source[0];
            source = source.substr(1);
        }
    }
    tokens.push_back(token);
}

void grab_num(string &source, vector<string> &tokens) {
    string token = "";
    while (in_charset(source[0],nums)) {
        token += source[0];
        source = source.substr(1);
    }
    if (source.at(0) == '.') {
        token += source[0];
        source = source.substr(1);
        while (in_charset(source[0],nums)) {
            token += source[0];
            source = source.substr(1);
        }
    }
    tokens.push_back(token);
}

void grab_label(string &source, vector<string> &tokens) {
    string token = source.substr(0,1);
    source = source.substr(1);
    while (in_charset(source[0],alpha) || 
           in_charset(source[0],nums)) {
        token += source.at(0);
        source = source.substr(1);
    }
    tokens.push_back(token);
}

void grab_punct(string &source, vector<string> &tokens) {
    string token = "";
    while (in_charset(source[0],punct)) {
        string added = token + source.substr(0,1);
        if (added.length() > 1 && 
            punct_combos.find(added) == punct_combos.end()) {
            tokens.push_back(token);
            token = source.substr(0,1);
        } else {
            token += source[0];
        }
        source = source.substr(1);
    }
    tokens.push_back(token);
}

vector<string> tokenize(string source) {
    vector<string> tokens;
    while (source.length() > 0) {
        if (source.length() >= 2 && source.substr(0,2) == "//") {
            grab_comment(source);
        } else if (source[0] == '\"') {
            grab_str(source,tokens);
        } else if (in_charset(source[0],alpha)) {
            grab_label(source, tokens);
        } else if (in_charset(source[0],nums)) {
            grab_num(source, tokens);
        } else if (in_charset(source[0],punct)) {
            grab_punct(source, tokens);
        } else if (in_charset(source[0],white)) { 
            source = source.substr(1); 
        } else { 
            compiler_error("unknown character: " 
                + source.substr(0,1));
        }
    }
    return tokens;
}