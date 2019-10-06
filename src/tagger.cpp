#include "tokenize.cpp"

using namespace std;

void TekoTokenizerException(string message, Token t) {
    printf("%s\n", string(12,'-').c_str());
    printf("Teko Tokenizing Error: %s\n", message.c_str());
    printf("at line %d, column %d\n", t.line_number, t.col);
    printf("%s\n", t.s.c_str());
    exit (EXIT_FAILURE);
}

// ------------

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
    throw runtime_error("illegal brace: " + to_string((char) b));
}

// ------------

const int num_prefixes = 2;

string prefixes[num_prefixes] = {"!", "~"};

// ------------

const int num_infixes = 17;

string infixes[num_infixes] = {"+", "-", "*", "/", "^", "%", "&", "|", "in", "to",
                               "==", "!=", "<", ">", "<=", ">=", "<:"};
// ------------

const int num_setters = 9;

string setters[num_setters] = {"=", "+=", "-=", "*=", "/=", "^=", "%=", "=&", "->"};

// ------------

const int num_suffixes = 7;

string suffixes[num_suffixes] = {".", "$", "#", "[]", "{}", "++", "--"};

// ------------

const int num_vartypes = 2;

string vartypes[num_vartypes] = {"async", "var"};

// ------------

const int num_annotations = 6;

string annotations[num_annotations] = {"@IO", "@hangs", "@updates", "@sees",
                                       "@modifies", "@throws"};

bool annotations_params[num_annotations] = {false, false, false, true, true, true};

// ------------

const int num_tagtypes = 27;

enum TagType { LabelTag, StringTag, IntTag, RealTag, BoolTag, CharTag,

               BitsTag, BytesTag,

               IfTag, ElseTag, ForTag, WhileTag,

               SemicolonTag, ColonTag, CommaTag, QMarkTag, AttrTag,

               OpenTag, CloseTag,

               PrefixTag, InfixTag, SetterTag, SuffixTag,

               LetTag, VisibilityTag, VartypeTag, AnnotationTag,

               CommandTag, NamespaceTag };

// ------------

TagType punct_tagtype(Token t) {
    string s = t.s;

    if      (s == ";") { return SemicolonTag; }
    else if (s == ":") { return ColonTag; }
    else if (s == ".") { return AttrTag; }
    else if (s == ",") { return CommaTag; }
    else if (s == "?") { return QMarkTag; }

    else if (s.length() == 1 && in_charset(s[0], braces)) {
        if (open_fromc(s[0])) {
            return OpenTag;
        } else {
            return CloseTag;
        }
    }

    else if (in_stringset(s, prefixes, num_prefixes)) {
        return PrefixTag;
    }

    else if (in_stringset(s, infixes, num_infixes)) {
        return InfixTag;
    }

    else if (in_stringset(s, setters, num_setters)) {
        return SetterTag;
    }

    else if (in_stringset(s, suffixes, num_suffixes)) {
        return SuffixTag;
    }

    else {
        TekoTokenizerException("Unknown token: " + s, t);
    }
}

struct Tag {
    Tag *next = 0;
    TagType type;
    char *val = 0; // will be cast to a different type of pointer depending on TagType
    string s;
    int line_number, col;

    Tag() {
    }

    string to_str() {
        switch(type) {
            case LabelTag:      return "LabelTag "  + *((string*) val);
            case StringTag:     return "StringTag \"" + teko_escape(*((string*) val)) + "\"";
            case IntTag:        return "IntTag "    + to_string(*((int*)    val));
            case RealTag:       return "RealTag "   + to_string(*((float*)  val));
            case BoolTag:       return "BoolTag "   + (*((bool*)   val)) ? "true" : "false";
            case CharTag:       { string out = "CharTag '"; out += *val; out += "'"; return out; }
            case BitsTag:       return "BitsTag "   + *((string*) val);
            case BytesTag:      return "BytesTag "  + *((string*) val);

            case IfTag:         return "IfTag";
            case ElseTag:       return "ElseTag";
            case ForTag:        return "ForTag";
            case WhileTag:      return "WhileTag";

            case SemicolonTag:  return "Semicolon Tag";
            case ColonTag:      return "ColonTag";
            case CommaTag:      return "CommaTag";
            case QMarkTag:      return "QMarkTag";
            case AttrTag:       return "AttrTag";

            case OpenTag:       return "OpenTag "   + to_string(*((Brace*) val), true);
            case CloseTag:      return "CloseTag "  + to_string(*((Brace*) val), false);

            case PrefixTag:     return "PrefixTag " + prefixes[*val];
            case InfixTag:      return "InfixTag "  + infixes[*val];
            case SetterTag:     return "SetterTag " + setters[*val];
            case SuffixTag:     return "SuffixTag " + suffixes[*val];

            case LetTag:        return "LetTag";
            case VisibilityTag: return "VisibilityTag";
            case VartypeTag:    return "VartypeTag " + vartypes[*val];
            case AnnotationTag: return "AnnotationTag " + annotations[*val];
            case CommandTag:    return "CommandTag";
            case NamespaceTag:  return "NamespaceTag";
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
            if      (token.s == "if")        { tag->type = IfTag; }
            else if (token.s == "else")      { tag->type = ElseTag; }
            else if (token.s == "for")       { tag->type = ForTag; }
            else if (token.s == "while")     { tag->type = WhileTag; }
            else if (token.s == "let")       { tag->type = LetTag; }
            else if (token.s == "namespace") { tag->type = NamespaceTag; }

            else if (token.s[0] == '@') {
                if (in_stringset(token.s, annotations, num_annotations)) {
                    tag->type = AnnotationTag;
                    tag->val = new char(string_index(token.s, annotations, num_annotations));
                } else {
                    TekoTokenizerException("Illegal token", token);
                }
            }

            else if (in_stringset(token.s, infixes, num_infixes)) {
                tag->type = InfixTag;
                tag->val = new char(string_index(token.s, infixes, num_infixes));
            }

            else if (in_stringset(token.s, vartypes, num_vartypes)) {
                tag->type = VartypeTag;
                tag->val = new char(string_index(token.s, vartypes, num_vartypes));
            }

            else {
                tag->type = LabelTag;
                string* sp = new string(token.s);
                tag->val = (char*) sp;
            }
            break;
        }

        case NUM_T: {
            if (token.s.find(".") == string::npos) {
                tag->type = IntTag;
                int* ip = new int(stoi(token.s));
                tag->val = (char*) ip;
            } else {
                tag->type = RealTag;
                float* fp = new float(stof(token.s));
                tag->val = (char*) fp;
            }
            break;
        }

        case STRING_T: {
            tag->type = StringTag;
            string* sp = new string(token.s);
            tag->val = (char*) sp;
            break;
        }

        case PUNCT_T: {
            tag->type = punct_tagtype(token);
            switch(tag->type) {
                case OpenTag:
                case CloseTag: {
                    Brace *bp = new Brace;
                    *bp = brace_fromc(token.s[0]);
                    tag->val = (char *) bp;
                } break;

                case PrefixTag: tag->val = new char(string_index(token.s, prefixes, num_prefixes)); break;
                case InfixTag:  tag->val = new char(string_index(token.s, infixes,  num_infixes));  break;
                case SetterTag: tag->val = new char(string_index(token.s, setters,  num_setters));  break;
                case SuffixTag: tag->val = new char(string_index(token.s, suffixes, num_suffixes)); break;

                default: tag->val = 0; break; // for ColonTag, etc.
            }
            break;
        }

        case CHAR_T: {
            tag->type = CharTag;
            char *cp = new char(token.s[0]);
            tag->val = cp;
            break;
        }

        case BITS_T: {
            tag->type = BitsTag;
            string* bp = new string(token.s);
            tag->val = (char*) bp;
        }

        case BYTES_T: {
            tag->type = BytesTag;
            string* Bp = new string(token.s);
            tag->val = (char*) Bp;
        }
    }
    return tag;
}
