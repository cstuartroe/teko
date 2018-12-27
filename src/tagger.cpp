#include <vector>

using namespace std;

enum TagType { LabelType, StringType, IntType, RealType,
                BoolType, IfType, ElseType, SemicolonType,
                ColonType, CommaType, QMarkType, BangType,
                AttrType, OpenType, CloseType, LAngleType,
                RAngleType, BinOpType, SetterType, 
                ComparisonType, ConversionType };

enum Brace { paren, curly, square, angle };
enum BinOp { add, sub, mult, divs, exp, mod, conj, disj };
enum Setter { normal, setadd, setsub, setmult, setdivs, setexp, setmod };
enum Comparison { eq, neq, lt, gt, leq, geq, subtype };
enum Conversion { toReal, toStr, toArr, toList, toSet };

class Tag {
    TagType type;
    string repr;

    string label;
    string str_val;
    int int_val;
    float real_val;
    bool bool_val;

    Brace brace;
    BinOp op;
    Setter setter;
    Comparison comp;
    Conversion conv;

    Tag(TagType _type) {
        type = _type;
    }
public:
    TagType getType() { return type; }

    string to_str() { return repr; }

    Brace getBrace() {
        if (type == OpenType || type == CloseType) { return brace; }
        else { throw runtime_error(repr + " has no attribute brace"); }
    }

    static Tag LabelTag(string label) {
        Tag t = Tag(LabelType);
        t.label = label;
        t.repr = "LabelTag " + label;
        return t;
    }
    static Tag StringTag(string val) {
        Tag t = Tag(StringType);
        t.str_val = val;
        t.repr = "StringTag " + val;
        return t;
    }
    static Tag IntTag(int val) {
        Tag t = Tag(IntType);
        t.int_val = val;
        t.repr = "IntTag " + std::to_string(val);
        return t;
    }
    static Tag RealTag(float val) {
        Tag t = Tag(RealType);
        t.real_val = val;
        t.repr = "RealTag " + std::to_string(val);
        return t;
    }
    static Tag BoolTag(bool val) {
        Tag t = Tag(BoolType);
        t.bool_val = val;
        t.repr = "BoollTag " + std::to_string(val);
        return t;
    }
    static Tag IfTag() {
        Tag t = Tag(IfType);
        t.repr = "IfTag";
        return t;
    }
    static Tag ElseTag() {
        Tag t = Tag(ElseType);
        t.repr = "ElseTag";
        return t;
    }
    static Tag SemicolonTag() {
        Tag t = Tag(SemicolonType);
        t.repr = "SemicolonTag";
        return t;
    }
    static Tag ColonTag() {
        Tag t = Tag(ElseType);
        t.repr = "ColonTag";
        return t;
    }
    static Tag CommaTag() {
        Tag t = Tag(CommaType);
        t.repr = "CommaTag";
        return t;
    }
    static Tag QMarkTag() {
        Tag t = Tag(QMarkType);
        t.repr = "QMarkTag";
        return t;
    }
    static Tag BangTag() {
        Tag t = Tag(BangType);
        t.repr = "BangTag";
        return t;
    }
    static Tag AttrTag() {
        Tag t = Tag(AttrType);
        t.repr = "AttrTag";
        return t;
    }
    static Tag OpenTag(Brace b) {
        Tag t = Tag(OpenType);
        t.brace = b;
        string bstr;
        switch (b) {
            case paren:  bstr = "("; break;
            case curly:  bstr = "{"; break;
            case square: bstr = "["; break;
            case angle:  bstr = "<"; break;
        }
        t.repr = "OpenTag " + bstr;
        return t;
    }
    static Tag CloseTag(Brace b) {
        Tag t = Tag(CloseType);
        t.brace = b;
        string bstr;
        switch (b) {
            case paren:  bstr = ")"; break;
            case curly:  bstr = "}"; break;
            case square: bstr = "]"; break;
            case angle:  bstr = ">"; break;
        }
        t.repr = "CloseTag " + bstr;
        return t;
    }
    static Tag LAngleTag() {
        Tag t = Tag(LAngleType);
        t.repr = "LAngleTag";
        return t;
    }
    static Tag RAngleTag() {
        Tag t = Tag(RAngleType);
        t.repr = "RAngleTag";
        return t;
    }
    static Tag BinOpTag(BinOp op) {
        Tag t = Tag(BinOpType);
        t.op = op;
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
        t.repr = "BinOpTag " + opstr;
        return t;
    }
    static Tag SetterTag(Setter s) {
        Tag t = Tag(SetterType);
        t.setter = s;
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
        t.repr = "SetterTag " + sstr;
        return t;
    }
    static Tag ComparisonTag(Comparison c) {
        Tag t = Tag(ComparisonType);
        t.comp = c;
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
        t.repr = "ComparisonTag " + cstr;
        return t;
    }
    static Tag ConversionTag(Conversion c) {
        Tag t = Tag(ConversionType);
        t.conv = c;
        string cstr;
        switch (c) {
            case toReal: cstr = ".";  break;
            case toStr:  cstr = "$";  break;
            case toArr:  cstr = "[]"; break;
            case toList: cstr = "{}"; break;
            case toSet:  cstr = "<>"; break;
        }
        t.repr = "ConversionTag " + cstr;
        return t;
    }
};

vector<Tag> get_tags(vector<string> tokens) {
    vector<Tag> tags;
    for (int i = 0; i < tokens.size(); i++) {
        string token = tokens[i];
        if (token == ";") {
            tags.push_back(Tag::SemicolonTag()); 
        } else if (token == ":") {
            tags.push_back(Tag::ColonTag());
        } else if (token == ",") {
            tags.push_back(Tag::CommaTag());
        } else if (token == "?") {
            tags.push_back(Tag::QMarkTag());
        } else if (token == "!") {
            tags.push_back(Tag::BangTag());
        }

        else if (token == "(") {
            tags.push_back(Tag::OpenTag(paren));
        } else if (token == ")") {
            tags.push_back(Tag::CloseTag(paren));
        } else if (token == "{") {
            tags.push_back(Tag::OpenTag(curly));
        } else if (token == "}") {
            tags.push_back(Tag::CloseTag(curly));
        } else if (token == "[") {
            tags.push_back(Tag::OpenTag(square));
        } else if (token == "]") {
            tags.push_back(Tag::CloseTag(square));
        } 

        else if (token == "if") {
            tags.push_back(Tag::IfTag());
        } else if (token == "else") {
            tags.push_back(Tag::ElseTag());
        } 

        else if (token == "+") {
            tags.push_back(Tag::BinOpTag(add));
        } else if (token == "-") {
            tags.push_back(Tag::BinOpTag(sub));
        } else if (token == "*") {
            tags.push_back(Tag::BinOpTag(mult));
        } else if (token == "/") {
            tags.push_back(Tag::BinOpTag(divs));
        } else if (token == "^") {
            tags.push_back(Tag::BinOpTag(exp));
        } else if (token == "%%") {
            tags.push_back(Tag::BinOpTag(mod));
        } else if (token == "&&") {
            tags.push_back(Tag::BinOpTag(conj));
        } else if (token == "||") {
            tags.push_back(Tag::BinOpTag(disj));
        }

        else if (token == "=") {
            tags.push_back(Tag::SetterTag(normal));
        } else if (token == "+=") {
            tags.push_back(Tag::SetterTag(setadd));
        } else if (token == "-=") {
            tags.push_back(Tag::SetterTag(setsub));
        } else if (token == "*=") {
            tags.push_back(Tag::SetterTag(setmult));
        } else if (token == "/=") {
            tags.push_back(Tag::SetterTag(setdivs));
        } else if (token == "^=") {
            tags.push_back(Tag::SetterTag(setexp));
        } else if (token == "%%=") {
            tags.push_back(Tag::SetterTag(setmod));
        }

        else if (token == "==") {
            tags.push_back(Tag::ComparisonTag(eq));
        } else if (token == "!=") {
            tags.push_back(Tag::ComparisonTag(neq));
        } else if (token == "<=") {
            tags.push_back(Tag::ComparisonTag(leq));
        } else if (token == ">=") {
            tags.push_back(Tag::ComparisonTag(geq));
        } else if (token == "<:") {
            tags.push_back(Tag::ComparisonTag(subtype));
        }

        // these might end up being Open/CloseTags or
        // ComparisonTags, pending semantic analysis
        else if (token == "<") {
            tags.push_back(Tag::LAngleTag());
        } else if (token == ">") {
            tags.push_back(Tag::RAngleTag());
        }

        else if (token == ".") {
            tags.push_back(Tag::ConversionTag(toReal));
        } else if (token == "$") {
            tags.push_back(Tag::ConversionTag(toStr));
        } else if (token == "[]") {
            tags.push_back(Tag::ConversionTag(toArr));
        } else if (token == "{}") {
            tags.push_back(Tag::ConversionTag(toList));
        } else if (token == "<>") {
            tags.push_back(Tag::ConversionTag(toSet));
        }

        else if (token == "true") {
            tags.push_back(Tag::BoolTag(true));
        } else if (token == "false") {
            tags.push_back(Tag::BoolTag(false));
        } 

        else if (token[0] == '"') {
            tags.push_back(Tag::StringTag(token));
        } else if (in_charset(token[0],nums)) {
            size_t found = token.find('.');
            if (found == string::npos) {
                tags.push_back(Tag::IntTag(stoi(token)));
            } else {
                tags.push_back(Tag::RealTag(stof(token)));
            }
        } else {
            if (i > 0 && tokens[i-1] == ".") {
                tags.pop_back(); // period . gets interpreted as a conversion unless followed by a label
                tags.push_back(Tag::AttrTag());
            } 
            tags.push_back(Tag::LabelTag(token));
        }
    }
    return tags;
}