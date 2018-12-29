#include <vector>
#include <variant>

using namespace std;

enum Brace { paren, curly, square, angle };
string to_string(Brace b, bool open) {
    if (open) {
        switch (b) {
            case paren:  return "("; break;
            case curly:  return "{"; break;
            case square: return "["; break;
            case angle:  return "<"; break;
        }
    } else {
        switch (b) {
            case paren:  return ")"; break;
            case curly:  return "}"; break;
            case square: return "]"; break;
            case angle:  return ">"; break;
        }
    }
}

enum BinOp { add, sub, mult, divs, exp, mod, conj, disj };
string to_string(BinOp op) {
    switch (op) {
        case add:  return "+"; break;
        case sub:  return "-"; break;
        case mult: return "*"; break;
        case divs: return "/"; break;
        case exp:  return "^"; break;
        case mod:  return "%%"; break;
        case conj: return "&&"; break;
        case disj: return "||"; break;
    }
}

enum Setter { setnormal, setadd, setsub, setmult, setdivs, setexp, setmod };
string to_string(Setter s) {
    switch (s) {
        case setnormal: return "=";  break;
        case setadd:    return "+="; break;
        case setsub:    return "-="; break;
        case setmult:   return "*="; break;
        case setdivs:   return "/="; break;
        case setexp:    return "^="; break;
        case setmod:    return "%%="; break;
    }
}

enum Comparison { eq, neq, lt, gt, leq, geq, subtype };
string to_string(Comparison c) {
    switch (c) {
        case eq:  return "=="; break;
        case neq: return "!="; break;
        case lt:  return "<";  break;
        case gt:  return ">";  break;
        case leq: return "<="; break;
        case geq: return ">="; break;
        case subtype: return "<:"; break;
    }
}

enum Conversion { toReal, toStr, toArr, toList, toSet };
string to_string(Conversion c) {
    switch (c) {
        case toReal: return ".";  break;
        case toStr:  return "$";  break;
        case toArr:  return "[]"; break;
        case toList: return "{}"; break;
        case toSet:  return "<>"; break;
    }
}

enum TagType { LabelTag, StringTag, IntTag, RealTag,
                BoolTag, IfTag, ElseTag, SemicolonTag,
                ColonTag, CommaTag, QMarkTag, BangTag,
                AttrTag, OpenTag, CloseTag, LAngleTag,
                RAngleTag, BinOpTag, SetterTag, 
                ComparisonTag, ConversionTag };

using TagValue = variant<string, int, float, bool, Brace, BinOp, Setter, Comparison, Conversion>;

struct Tag {
    TagType type;
    TagValue val;
    string to_str() {
        switch(type) {
            case LabelTag:     return "LabelTag " + get<string>(val);          break;
            case StringTag:    return "StringTag " + get<string>(val);         break;
            case IntTag:       return "IntTag " + to_string(get<int>(val));    break;
            case RealTag:      return "RealTag " + to_string(get<float>(val)); break;
            case BoolTag:      return "BoolTag " + to_string(get<bool>(val));; break;

            case IfTag:        return "IfTag";                                 break;
            case ElseTag:      return "ElseTag";                               break;
            case SemicolonTag: return "Semicolon Tag";                         break;
            case ColonTag:     return "ColonTag";                              break;
            case CommaTag:     return "CommaTag";                              break;
            case QMarkTag:     return "QMarkTag";                              break;
            case BangTag:      return "BangTag";                               break;
            case AttrTag:      return "AttrTag";                               break;

            case OpenTag:      return "OpenTag " + to_string(get<Brace>(val),true); break;
            case CloseTag:     return "CloseTag " + to_string(get<Brace>(val),false); break;
            case LAngleTag:    return "LAngleTag";                             break;
            case RAngleTag:    return "RAngleTag";                             break;
            case BinOpTag:     return "BinOpTag " + to_string(get<BinOp>(val)); break;
            case SetterTag:    return "SetterTag " + to_string(get<Setter>(val)); break;
            case ComparisonTag:return "ComparisonTag " + to_string(get<Comparison>(val)); break;
            case ConversionTag:return "ConversionTag " + to_string(get<Conversion>(val)); break;
        }
    }
};

vector<Tag> get_tags(vector<string> tokens) {
    vector<Tag> tags;

    for (int i = 0; i < tokens.size(); i++) {
        string token = tokens[i];
        Tag newtag;

        if (token == ";") {
            newtag.type = SemicolonTag; 
        } else if (token == ":") {
            newtag.type = ColonTag;
        } else if (token == ",") {
            newtag.type = CommaTag;
        } else if (token == "?") {
            newtag.type = QMarkTag;
        } else if (token == "!") {
            newtag.type = BangTag;
        }

        else if (token == "if") {
            newtag.type = IfTag;
        } else if (token == "else") {
            newtag.type = ElseTag;
        } 

        else if (token == "(") {
            newtag.type = OpenTag;
            newtag.val = paren;
        } else if (token == ")") {
            newtag.type = CloseTag;
            newtag.val = paren;
        } else if (token == "{") {
            newtag.type = OpenTag;
            newtag.val = curly;
        } else if (token == "}") {
            newtag.type = CloseTag;
            newtag.val = curly;
        } else if (token == "[") {
            newtag.type = OpenTag;
            newtag.val = square;
        } else if (token == "]") {
            newtag.type = CloseTag;
            newtag.val = square;
        } 

        else if (token == "+") {
            newtag.type = BinOpTag;
            newtag.val = add;
        } else if (token == "-") {
            newtag.type = BinOpTag;
            newtag.val = sub;
        } else if (token == "*") {
            newtag.type = BinOpTag;
            newtag.val = mult;
        } else if (token == "/") {
            newtag.type = BinOpTag;
            newtag.val = divs;
        } else if (token == "^") {
            newtag.type = BinOpTag;
            newtag.val = exp;
        } else if (token == "%%") {
            newtag.type = BinOpTag;
            newtag.val = mod;
        } else if (token == "&&") {
            newtag.type = BinOpTag;
            newtag.val = conj;
        } else if (token == "||") {
            newtag.type = BinOpTag;
            newtag.val = disj;
        }

        else if (token == "=") {
            newtag.type = SetterTag;
            newtag.val = setnormal;
        } else if (token == "+=") {
            newtag.type = SetterTag;
            newtag.val = setadd;
        } else if (token == "-=") {
            newtag.type = SetterTag;
            newtag.val = setsub;
        } else if (token == "*=") {
            newtag.type = SetterTag;
            newtag.val = setmult;
        } else if (token == "/=") {
            newtag.type = SetterTag;
            newtag.val = setdivs;
        } else if (token == "^=") {
            newtag.type = SetterTag;
            newtag.val = setexp;
        } else if (token == "%%=") {
            newtag.type = SetterTag;
            newtag.val = setmod;
        }

        else if (token == "==") {
            newtag.type = SetterTag;
            newtag.val = eq;
        } else if (token == "!=") {
            newtag.type = SetterTag;
            newtag.val = neq;
        } else if (token == "<=") {
            newtag.type = SetterTag;
            newtag.val = leq;
        } else if (token == ">=") {
            newtag.type = SetterTag;
            newtag.val = geq;
        } else if (token == "<:") {
            newtag.type = SetterTag;
            newtag.val = subtype;
        }

        // these might end up being Open/CloseTags or
        // ComparisonTags, pending semantic analysis
        else if (token == "<") {
            newtag.type = LAngleTag;
        } else if (token == ">") {
            newtag.type = RAngleTag;
        }

        else if (token == ".") {
            newtag.type = ConversionTag;
            newtag.val = toReal;
        } else if (token == "$") {
            newtag.type = ConversionTag;
            newtag.val = toStr;
        } else if (token == "[]") {
            newtag.type = ConversionTag;
            newtag.val = toArr;
        } else if (token == "{}") {
            newtag.type = ConversionTag;
            newtag.val = toList;
        } else if (token == "<>") {
            newtag.type = ConversionTag;
            newtag.val = toSet;
        }

        else if (token == "true") {
            newtag.type = BoolTag;
            newtag.val = true;
        } else if (token == "false") {
            newtag.type = BoolTag;
            newtag.val = false;
        } 

        else if (token[0] == '"') {
            newtag.type = StringTag;
            newtag.val = token;
        } else if (in_charset(token[0],nums)) {
            size_t found = token.find('.');
            if (found == string::npos) {
                newtag.type = IntTag;
                newtag.val = stoi(token);
            } else {
                newtag.type = RealTag;
                newtag.val = stof(token);
            }
        } else {
            if (i > 0 && tokens[i-1] == ".") {
                tags.pop_back(); // period . gets interpreted as a conversion unless followed by a label
                Tag prevtag;
                prevtag.type = AttrTag;
                tags.push_back(prevtag);
            } 
            newtag.type = LabelTag;
            newtag.val = token;
        }

        tags.push_back(newtag);
    }

    return tags;
}