#include <vector>
#include <set>

using namespace std;

enum TagType { LabelType, StringType, IntType, RealType,
                BoolType, IfType, ElseType, SemicolonType,
                ColonType, CommaType, QMarkType, BangType,
                AttrType, OpenType, CloseType, LAngleType,
                RAngleType, BinOpType, SetterType, 
                ComparisonType, ConversionType };

class Tag {
protected:
    TagType type;
    string repr;
public:
    TagType getType() { return type; }
    string to_str() { return repr; }
};

class LabelTag : public Tag {
    string label;
public:
    LabelTag(string _label) { 
        type = LabelType; 
        label = _label;
        repr = "LabelTag " + label;
    }
};

class StringTag : public Tag {
    string val;
public:
    StringTag(string _val) { 
        type = StringType; 
        val = _val;
        repr = "StringTag " + val;
    }
};

class IntTag : public Tag {
    int val;
public:
    IntTag(int _val) { 
        type = IntType; 
        val = _val;
        repr = "IntTag " + std::to_string(val);
    }
};

class RealTag : public Tag {
    float val;
public:
    RealTag(float _val) { 
        type = RealType;
        val = _val;
        repr = "RealTag " + std::to_string(val);
    }
};

class BoolTag : public Tag {
    bool val;
public:
    BoolTag(bool _val) { 
        type = BoolType; 
        val = _val;
        repr = "BoolTag " + std::to_string(val);
    }
};

class IfTag : public Tag {
public:
    IfTag() { 
        type = IfType; 
        repr = "IfTag";
    }
};

class ElseTag : public Tag {
public:
    ElseTag() {    
        type = ElseType; 
        repr = "ElseTag";
    }
};

class SemicolonTag : public Tag {
public:
    SemicolonTag() { 
        type = SemicolonType; 
        repr = "SemicolonTag";
    }
};

class ColonTag : public Tag {
public:
    ColonTag() { 
        type = ColonType; 
        repr = "ColonTag";
    }
};

class CommaTag : public Tag {
public:
    CommaTag() { 
        type = CommaType; 
        repr = "CommaTag";
    }
};

class QMarkTag : public Tag {
public:
    QMarkTag() { 
        type = QMarkType; 
        repr = "QMarkTag";
    }
};

class BangTag : public Tag {
public:
    BangTag() { 
        type = BangType; 
        repr = "BangTag";
    }
};

class AttrTag : public Tag {
public:
    AttrTag() { 
        type = AttrType; 
        repr = "AttrTag";
    }
};

enum Brace { paren, curly, square, angle };

class OpenTag : public Tag {
    Brace b;
public:
    OpenTag(Brace _b) {
        type = OpenType;
        b = _b;
        string bstr;
        switch (b) {
            case paren:  bstr = "("; break;
            case curly:  bstr = "{"; break;
            case square: bstr = "["; break;
            case angle:  bstr = "<"; break;
        }
        repr = "OpenTag " + bstr;
    }
    Brace getBrace() {return b;}
};

class CloseTag : public Tag {
    Brace b;
public:
    CloseTag(Brace _b) {
        type = CloseType;
        b = _b;
        string bstr;
        switch (b) {
            case paren:  bstr = ")"; break;
            case curly:  bstr = "}"; break;
            case square: bstr = "]"; break;
            case angle:  bstr = ">"; break;
        }
        repr = "CloseTag " + bstr;
    }
    Brace getBrace() {return b;}
};

// < and > need their own tag type because
// they might be braces or comparisons
// semantics is needed to disambiguate

class LAngleTag : public Tag {
public:
    LAngleTag() {
        type = LAngleType;
        repr = "LAngleTag";
    }
};

class RAngleTag : public Tag {
public:
    RAngleTag() {
        type = RAngleType;
        repr = "RAngleTag";
    }
};

enum BinOp { add, sub, mult, divs, exp, mod, conj, disj };

class BinOpTag : public Tag {
    BinOp op;
public:
    BinOpTag(BinOp _op) {
        type = BinOpType;
        op = _op;
        string opstr;
        switch (op) {
            case add:  opstr = "+"; break;
            case sub:  opstr = "-"; break;
            case mult: opstr = "*"; break;
            case divs: opstr = "/"; break;
            case exp:  opstr = "^"; break;
            case mod:  opstr = "%%"; break;
            case conj: opstr = "&&"; break;
            case disj: opstr = "||"; break;
        }
        repr = "BinOpTag " + opstr;
    }
};

enum Setter { normal, setadd, setsub, setmult, setdivs, setexp, setmod };

class SetterTag : public Tag {
    Setter s;
public:
    SetterTag(Setter _s) {
        type = SetterType;
        s = _s;
        string sstr;
        switch (s) {
            case normal:  sstr = "=";  break;
            case setadd:  sstr = "+="; break;
            case setsub:  sstr = "-="; break;
            case setmult: sstr = "*="; break;
            case setdivs: sstr = "/="; break;
            case setexp:  sstr = "^="; break;
            case setmod:  sstr = "%%="; break;
        }
        repr = "SetterTag " + sstr;
    }
};

enum Comparison { eq, neq, lt, gt, leq, geq, subtype };

class ComparisonTag : public Tag {
    Comparison c;
public:
    ComparisonTag(Comparison _c) {
        type = ComparisonType;
        c = _c;
        string cstr;
        switch (c) {
            case eq:  cstr = "=="; break;
            case neq: cstr = "!="; break;
            case lt:  cstr = "<";  break;
            case gt:  cstr = ">";  break;
            case leq: cstr = "<="; break;
            case geq: cstr = ">="; break;
            case subtype: cstr = "<:"; break;
        }
        repr = "ComparisonTag " + cstr;
    }
};

enum Conversion { toReal, toStr, toArr, toList, toSet };

class ConversionTag : public Tag {
    Conversion c;
public:
    ConversionTag(Conversion _c) {
        type = ConversionType;
        c = _c;
        string cstr;
        switch (c) {
            case toReal: cstr = ".";  break;
            case toStr:  cstr = "$";  break;
            case toArr:  cstr = "[]"; break;
            case toList: cstr = "{}"; break;
            case toSet:  cstr = "<>"; break;
        }
        repr = "ConversionTag " + cstr;
    }
};

vector<Tag*> get_tags(vector<string> tokens) {
    vector<Tag*> tags;
    for (int i = 0; i < tokens.size(); i++) {
        string token = tokens[i];
        if (token == ";") {
            SemicolonTag st = SemicolonTag();
            tags.push_back(&st); 
        } else if (token == ":") {
            ColonTag ct = ColonTag();
            tags.push_back(&ct);
        } else if (token == ",") {
            CommaTag ct = CommaTag();
            tags.push_back(&ct);
        } else if (token == "?") {
            QMarkTag qt = QMarkTag();
            tags.push_back(&qt);
        } else if (token == "!") {
            BangTag bt = BangTag();
            tags.push_back(&bt);
        }

        else if (token == "(") {
            OpenTag ot = OpenTag(paren);
            tags.push_back(&ot);
        } else if (token == ")") {
            CloseTag ct = CloseTag(paren);
            tags.push_back(&ct);
        } else if (token == "{") {
            OpenTag ot = OpenTag(curly);
            tags.push_back(&ot);
        } else if (token == "}") {
            CloseTag ct = CloseTag(curly);
            tags.push_back(&ct);
        } else if (token == "[") {
            OpenTag ot = OpenTag(square);
            tags.push_back(&ot);
        } else if (token == "]") {
            CloseTag ct = CloseTag(square);
            tags.push_back(&ct);
        } 

        else if (token == "if") {
            IfTag it = IfTag();
            tags.push_back(&it);
        } else if (token == "else") {
            ElseTag et = ElseTag();
            tags.push_back(&et);
        } 

        else if (token == "+") {
            BinOpTag bot = BinOpTag(add);
            tags.push_back(&bot);
        } else if (token == "-") {
            BinOpTag bot = BinOpTag(sub);
            tags.push_back(&bot);
        } else if (token == "*") {
            BinOpTag bot = BinOpTag(mult);
            tags.push_back(&bot);
        } else if (token == "/") {
            BinOpTag bot = BinOpTag(divs);
            tags.push_back(&bot);
        } else if (token == "^") {
            BinOpTag bot = BinOpTag(exp);
            tags.push_back(&bot);
        } else if (token == "%%") {
            BinOpTag bot = BinOpTag(mod);
            tags.push_back(&bot);
        } else if (token == "&&") {
            BinOpTag bot = BinOpTag(conj);
            tags.push_back(&bot);
        } else if (token == "||") {
            BinOpTag bot = BinOpTag(disj);
            tags.push_back(&bot);
        }

        else if (token == "=") {
            SetterTag st = SetterTag(normal);
            tags.push_back(&st);
        } else if (token == "+=") {
            SetterTag st = SetterTag(setadd);
            tags.push_back(&st);
        } else if (token == "-=") {
            SetterTag st = SetterTag(setsub);
            tags.push_back(&st);
        } else if (token == "*=") {
            SetterTag st = SetterTag(setmult);
            tags.push_back(&st);
        } else if (token == "/=") {
            SetterTag st = SetterTag(setdivs);
            tags.push_back(&st);
        } else if (token == "^=") {
            SetterTag st = SetterTag(setexp);
            tags.push_back(&st);
        } else if (token == "%%=") {
            SetterTag st = SetterTag(setmod);
            tags.push_back(&st);
        }

        else if (token == "==") {
            ComparisonTag ct = ComparisonTag(eq);
            tags.push_back(&ct);
        } else if (token == "!=") {
            ComparisonTag ct = ComparisonTag(neq);
            tags.push_back(&ct);
        } else if (token == "<=") {
            ComparisonTag ct = ComparisonTag(leq);
            tags.push_back(&ct);
        } else if (token == ">=") {
            ComparisonTag ct = ComparisonTag(geq);
            tags.push_back(&ct);
        } else if (token == "<:") {
            ComparisonTag ct = ComparisonTag(subtype);
            tags.push_back(&ct);
        }

        // these might end up being Open/CloseTags or
        // ComparisonTags, pending semantic analysis
        else if (token == "<") {
            LAngleTag lt = LAngleTag();
            tags.push_back(&lt);
        } else if (token == ">") {
            RAngleTag rt = RAngleTag();
            tags.push_back(&rt);
        }

        else if (token == ".") {
            ConversionTag ct = ConversionTag(toReal);
            tags.push_back(&ct);
        } else if (token == "$") {
            ConversionTag ct = ConversionTag(toStr);
            tags.push_back(&ct);
        } else if (token == "[]") {
            ConversionTag ct = ConversionTag(toArr);
            tags.push_back(&ct);
        } else if (token == "{}") {
            ConversionTag ct = ConversionTag(toList);
            tags.push_back(&ct);
        } else if (token == "<>") {
            ConversionTag ct = ConversionTag(toSet);
            tags.push_back(&ct);
        }

        else if (token == "true") {
            BoolTag bt = BoolTag(true);
            tags.push_back(&bt);
        } else if (token == "false") {
            BoolTag bt = BoolTag(false);
            tags.push_back(&bt);
        } 

        else if (token[0] == '"') {
            StringTag st = StringTag(token);
            tags.push_back(&st);
        } else if (in_charset(token[0],nums)) {
            size_t found = token.find('.');
            if (found == string::npos) {
                IntTag it = IntTag(stoi(token));
                tags.push_back(&it);
            } else {
                RealTag rt = RealTag(stof(token));
                tags.push_back(&rt);
            }
        } else {
            if (i > 0 && tokens[i-1] == ".") {
                tags.pop_back(); // period . gets interpreted as a conversion unless followed by a label
                AttrTag at = AttrTag();
                tags.push_back(&at);
            }
            LabelTag lt = LabelTag(token);
            tags.push_back(&lt);
        }
    }
    return tags;
}