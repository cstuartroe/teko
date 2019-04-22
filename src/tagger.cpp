#include "tokenize.cpp"

using namespace std;

const int num_braces = 4;

char braces[9] = "(){}[]<>";

enum Brace { paren, curly, square, angle };

Brace brace_fromc(char c) {
    switch(c) {
        case '(': return paren;
        case ')': return paren;
        case '{': return curly;
        case '}': return curly;
        case '[': return square;
        case ']': return square;
        case '<': return angle;
        case '>': return angle;
    }
}

bool open_fromc(char c) {
  switch(c) {
      case '(': return true;
      case ')': return false;
      case '{': return true;
      case '}': return false;
      case '[': return true;
      case ']': return false;
      case '<': return true;
      case '>': return false;
  }
}

string to_string(Brace b, bool open) {
    if (open) {
        switch (b) {
            case paren:  return "(";
            case curly:  return "{";
            case square: return "[";
            case angle:  return "<";
        }
    } else {
        switch (b) {
            case paren:  return ")";
            case curly:  return "}";
            case square: return "]";
            case angle:  return ">";
        }
    }
}

// ------------

const int num_binops = 9;

string binops[num_binops+1] = {"+", "-", "*", "/", "^", "%", "&", "|", "in"};

enum BinOp { add, sub, mult, divs, exp, mod, conj, disj, in};

// ------------

const int num_setters = 8;

string setters[num_setters+1] = {"=", "+=", "-=", "*=", "/=", "/=", "^=", "%=", "->"};

enum Setter { setnormal, setadd, setsub, setmult, setdivs, setexp, setmod, routine};

// ------------

const int num_comparisons = 7;

string comparisons[num_comparisons+1] = {"==", "!=", "<", ">", "<=", ">=", "<:"};

enum Comparison { eq, neq, lt, gt, leq, geq, subtype};

// ------------

const int num_conversions = 5;

string conversions[num_conversions+1] = {".", "$", "[]", "{}", "<>"};

enum Conversion { toReal, toStr, toArr, toList, toSet};

// ------------

const int num_tagtypes = 29;

enum TagType { LabelTag, StringTag, IntTag, RealTag, BoolTag, CharTag,

               IfTag, ElseTag, ForTag, WhileTag,

               SemicolonTag, ColonTag, CommaTag, QMarkTag, BangTag, AttrTag,

               OpenTag, CloseTag,

               BinOpTag, SetterTag, ComparisonTag, ConversionTag,

               VarTag, VisibilityTag, VartypeTag, AnnotationTag, CommandTag };

// ------------

TagType punct_tagtype(string s) {
    if      (s == ";") { return SemicolonTag; }
    else if (s == ":") { return ColonTag; }
    else if (s == ",") { return CommaTag; }
    else if (s == "?") { return QMarkTag; }
    else if (s == "!") { return BangTag; }

    else if (s.length() == 1 && in_charset(s[0], braces)) {
      if (open_fromc(s[0])) {
        return OpenTag;
      } else {
        return CloseTag;
      }
    }

    else if (in_stringset(s, binops, num_binops)) {
      return BinOpTag;
    }

    else if (in_stringset(s, setters, num_setters)) {
      return SetterTag;
    }

    else if (in_stringset(s, comparisons, num_comparisons)) {
      return ComparisonTag;
    }

    else if (in_stringset(s, conversions, num_conversions)) {
      return ConversionTag;
    }
}

struct Tag {
    Tag *next;
    TagType type;
    char *val; // will be cast to a different type of pointer depending on TagType
    string s;
    int line_number, col;

    Tag() {}

    string to_str() {
        switch(type) {
            case LabelTag:     return "LabelTag "  + (*((string*) val));
            case StringTag:    return "StringTag \"" + teko_escape(*((string*) val)) + "\"";
            case IntTag:       return "IntTag "    + to_string(*((int*)    val));
            case RealTag:      return "RealTag "   + to_string(*((float*)  val));
            case BoolTag:      return "BoolTag "   + (*((bool*)   val)) ? "true" : "false";
            case CharTag:      { string out = "CharTag '"; out += *val; out += "'"; return out; }

            case IfTag:        return "IfTag";
            case ElseTag:      return "ElseTag";
            case ForTag:       return "ForTag";
            case WhileTag:     return "WhileTag";

            case SemicolonTag: return "Semicolon Tag";
            case ColonTag:     return "ColonTag";
            case CommaTag:     return "CommaTag";
            case QMarkTag:     return "QMarkTag";
            case BangTag:      return "BangTag";
            case AttrTag:      return "AttrTag";

            case OpenTag:      return "OpenTag "       + to_string(*((Brace*) val), true);
            case CloseTag:     return "CloseTag "      + to_string(*((Brace*) val), false);

            case BinOpTag:     return "BinOpTag "      + binops[*val];
            case SetterTag:    return "SetterTag "     + setters[*val];
            case ComparisonTag:return "ComparisonTag " + comparisons[*val];
            case ConversionTag:return "ConversionTag " + conversions[*val];

            case VarTag:       return "VarTag";
            case VisibilityTag: return "VisibilityTag";
            case VartypeTag:   return "VartypeTag";
            case AnnotationTag:return "AnnotationTag";
            case CommandTag:   return "CommandTag";
        }
    }
};

Tag *from_token(Token token) {
    Tag *tag = new Tag();
    tag->line_number = token.line_number;
    tag->col = token.col;
    tag->s = token.s;

    switch (token.type) {
        case LABEL_T: {
            tag->type = LabelTag;
            string* sp = new string(token.s);
            tag->val = (char*) sp;
            break;
        }

        case NUM_T: {
            tag->type = IntTag;
            int* ip = new int(stoi(token.s));
            tag->val = (char*) ip;
            break;
        }

        case STRING_T: {
            tag->type = StringTag;
            string* sp = new string(token.s);
            tag->val = (char*) sp;
            break;
        }

        case PUNCT_T: {
            tag->type = punct_tagtype(token.s);
            char index;
            switch(tag->type) {
                case OpenTag:  index = brace_fromc(token.s[0]); break;
                case CloseTag: index = brace_fromc(token.s[0]); break;

                case BinOpTag:      index = string_index(token.s, binops,      num_binops);      break;
                case SetterTag:     index = string_index(token.s, setters,     num_setters);     break;
                case ComparisonTag: index = string_index(token.s, comparisons, num_comparisons); break;
                case ConversionTag: index = string_index(token.s, conversions, num_conversions); break;

                default: index = 0; break; // for ColonTag, etc.
            }
            char* heap_index = new char(index);
            tag->val = heap_index;
            break;
        }

        case CHAR_T: {
            tag->type = CharTag;
            char *cp = new char(token.s[0]);
            tag->val = cp;
            break;
        }
    }

    return tag;
}
