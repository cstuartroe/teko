#include <string>

char alpha[54] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_";
char num[11] = "0123456789";
char hex[17] = "0123456789abcdef";
char white[4] = " \t\n";
char punct[32] = "!\"#$%&\'()*+,-./:;<=>?@[\\]^`{|}~";

const int num_punct_combos = 14;
std::string punct_combos[num_punct_combos] = {"==","<=",">=","!=","<:","+=","-=",
                                              "*=","/=","^=","%%=","{}","[]","<>"};

// accepts a hexidecimal digit and returns corresponding value
char xtoi(char c) {
    for (int i = 0; i < 16; i++) {
        if (hex[i] == c) { return i; }
    }
    throw std::runtime_error("uh oh, how'd we get here? 1");
}

enum STATE {LABEL_T, NUM_T, STRING_T, PUNCT_T, CHAR_T, 
             LINE_COMMENT, BLOCK_COMMENT, BLANK};

struct Token {
    Token () {
        next = 0;
    }
    std::string s;
    int line_number, col;
    Token *next;
    STATE type;
};

bool in_charset(char c, char *charset) {
    while (*charset != 0) {
        if (c == *charset) { return true; }
        charset++;
    }
    return false;
}

bool in_stringset(std::string s, std::string *stringset, char size) {
    for (int i = 0; i < size; i++) {
        if (s == stringset[i]) { return true ;}
    }
    return false;
}

struct Tokenizer {
    Token *start;
    Token *curr;
    std::string buff;
    int line_number = 0, col;
    STATE state;
    bool decimal;

    Tokenizer () {
        start = new Token();
        curr = start;
        buff = "";
        state = BLANK;
        decimal = false;
    }

    void newtoken() {
        curr->type = state;
        curr->next = new Token;
        curr = curr->next;
        state = BLANK;
    }

    // currently, parse_char returns \x00 if the current buff is not yet a valid escape sequence
    // maybe it should be a more obscure character like \x14, or maybe I should make a response struct with a bool field
    char parse_char() {
        int len = buff.length();
        char out;

        if (len == 0 || len > 4) { // no escape sequence should be more than 4 characters
            throw std::runtime_error("uh oh, how'd we get here? 2");
        }

        if (buff[0] != '\\') { // a non-escaped character
            if (len != 1) {
                throw std::runtime_error("uh oh, how'd we get here? 3");
            }
            out = buff[0];
            buff = "";
            return out;
        }

        if (len == 1) { // a single backslash - not an escape sequence yet!
            return 0;
        } 

        switch (buff[1]) {
            case '\\': out = '\\'; break;
            case '\"': out = '\"'; break;
            case '\'': out = '\''; break;
            case 'n':  out = '\n'; break;
            case 't':  out = '\t'; break;
            case 'x':  out = parse_hex(); break;
            default:   throw std::runtime_error("Invalid escape sequence");
        }

        if (out != 0) { buff = ""; }
        return out;
    }

    char parse_hex() {
        if (buff.length() < 4) {
            return 0;
        }

        if (buff[0] != '\\' || buff[1] != 'x' || !in_charset(buff[2],hex) || !in_charset(buff[3],hex)) {
            throw std::runtime_error("Invalid escape sequence");
        }

        return xtoi(buff[2])*16 + xtoi(buff[3]);
    }

    void digest_line(std::string line) {
        line_number++;
        for (int i = 0; i < line.length(); i++) {
            col = i+1;
            digest(line[i]);
        }
        digest('\n');
    }

    void digest(char c) {
        switch (state) {
            case BLANK:    digest_blank(c);  break;
            case LABEL_T:  digest_label(c);  break;
            case NUM_T:    digest_num(c);    break;
            case STRING_T: digest_string(c); break;
            case PUNCT_T:  digest_punct(c);  break;
            case CHAR_T:   digest_char(c);   break;
            case LINE_COMMENT:  digest_line_comment(c);  break;
            case BLOCK_COMMENT: digest_block_comment(c); break;
        }
    }

    void digest_blank(char c) {
        if (in_charset(c, white)) { return; }

        curr->line_number = line_number;
        curr->col = col;

        if (c == '\"') {
            state = STRING_T;
            return;
        } else if (c == '\'') {
            state = CHAR_T;
            return;
        } else if (in_charset(c, alpha)) {
            state = LABEL_T;
        } else if (in_charset(c, num)) {
            state = NUM_T;
        } else if (in_charset(c, punct)) {
            state = PUNCT_T;
        } else  { 
            throw std::runtime_error("Disallowed character");
        }

        digest(c);
    }

    void digest_label(char c) {
        if (in_charset(c, alpha) || in_charset(c, num))  {
            curr->s += c;
        } else {
            newtoken();
            digest_blank(c);
        }
    }

    void digest_num(char c) {
        if (in_charset(c, num)) {
            curr->s += c;
        } else if (c == '.' && !decimal) {
            decimal = true;
            curr->s += c;
        } else {
            newtoken();
            digest_blank(c);            
        }
    }

    void digest_punct(char c) {
        if (in_charset(c,num) && buff == ".") {
            state = NUM_T;
            decimal = true;
            digest(c);
            return;
        }

        if (!in_charset(c, punct)) {
            curr->s = buff.substr(0, buff.length());
            buff = "";
            newtoken();
            digest_blank(c);
            return;
        }

        buff += c;

        if (buff == "/*") {
            buff = "";
            state = BLOCK_COMMENT;
            return;
        }

        if (buff == "//") {
            buff = "";
            state = LINE_COMMENT;
            return;
        }

        if (buff.length() > 1 && !in_stringset(buff, punct_combos, num_punct_combos)) {
            curr->s = buff.substr(0, buff.length() - 1);
            buff = "";
            newtoken();
            digest_blank(c);
        }
    }

    void digest_string(char c) {
        buff += c;
        if (buff == "\"") {
            buff = "";
            newtoken(); 
            return;          
        }
        char parsed = parse_char();
        if (parsed != 0) {
            curr->s += parsed;
        }
    }

    void digest_char(char c) {
        buff += c;
        if (buff == "\'") {
            if (curr->s.length() == 0){
                throw std::runtime_error("Empty character is prohibited");
            }
            buff = "";
            newtoken();
        } else {
            if (curr->s.length() != 0) {
                throw std::runtime_error("Strings cannot use single quotes");
            }
            char parsed = parse_char();
            if (parsed != 0) {
                curr->s += parsed;
            }
        }
    }

    void digest_line_comment(char c) {
        if (col == 1) {
            state = BLANK;
            digest_blank(c);
        }
    }

    void digest_block_comment(char c) {
        buff += c;
        if (buff == "*/") {
            state = BLANK;
            buff = "";
        } else if (buff != "*") {
            buff = "";
        }
    }
};