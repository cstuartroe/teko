#include <vector>
#include <variant>

enum LineType { DeclarationLine, AssignmentLine, ExpressionLine, 
                IfStatementLine, WhileBlockLine, ForBlockLine };

class Line;

struct CodeBlock {
    vector<Line> lines;
};

struct Assignment {
};

struct Declaration {
    TekoObject type;
    vector<Assignment> assignments;
};

struct Expression {
};

struct IfStatement {
    Expression condition;
    CodeBlock then_block;
    IfStatement* else_stmt;
};

struct WhileBlock {
    Expression condition;
    CodeBlock loop;
};

struct ForBlock {
    TekoObject type;
    TekoObject label;
    TekoObject iterable;
    CodeBlock loop;
};

struct Line {
    LineType type;
    variant<Declaration,Assignment,Expression,IfStatement,
        WhileBlock,ForBlock> content;
    vector<Tag> tags;
};

vector<Tag> match_braces(vector<Tag> &tags, int &location) {
    vector<Tag> out;
    Brace b;
    if (tags[location].type == OpenTag) {
        b = get<Brace>(tags[location].val);
        out.push_back(tags[location]);
        location++;
    } else { throw runtime_error("brace matching issue"); }

    while (location < tags.size() && tags[location].type != CloseTag) {
        if (tags[location].type == OpenTag) {
            vector<Tag> brace_section = match_braces(tags,location);
            out.insert(out.end(),brace_section.begin(),brace_section.end()) ;
        } else {
            out.push_back(tags[location]);
            location++;
        }
    }

    if (location == tags.size()) {
        compiler_error("EOF during brace section");
    } else {
        if (get<Brace>(tags[location].val) != b) {
            compiler_error("Mismatched braces");
        } else {
            out.push_back(tags[location]);
            location++;
        }
    }

    return out;
}

LineType determine_line_type(vector<Tag> line_tags) {
    if (line_tags[0].type == IfTag) {
        return IfStatementLine;
    } else if (line_tags[0].type == ForTag) {
        return ForBlockLine;
    } else if (line_tags[0].type == WhileTag) {
        return WhileBlockLine;
    } else if (line_tags[0].type == LetTag) {
        return DeclarationLine;
    } else if (line_tags.size() >= 2) {
        if (line_tags[0].type == LabelTag && line_tags[1].type == LabelTag) {
            return DeclarationLine;
        } else if (line_tags[0].type == LabelTag && line_tags[1].type == SetterTag) {
            return AssignmentLine;
        }
    }
    return ExpressionLine;
}

void grab_line(vector<Tag> &tags, vector<Line> &lines, int &location) {
    vector<Tag> line_tags;
    while (location < tags.size() && tags[location].type != SemicolonTag ) {
        if (tags[location].type == OpenTag) {
            vector<Tag> brace_section = match_braces(tags,location);
            line_tags.insert(line_tags.end(),brace_section.begin(),brace_section.end()) ;
        } else {
            line_tags.push_back(tags[location]);
            location++;
        }
    }

    if (location == tags.size()) {
        compiler_error("EOF in middle of line");
    } else {
        location++;
    }

    Line l;
    l.tags = line_tags;
    l.type = determine_line_type(line_tags);
    lines.push_back(l);
}

vector<Line> get_lines(vector<Tag> &tags) {
    vector<Line> lines;
    int location = 0;
    while(tags.size() > location) {
        grab_line(tags, lines, location);
    }
    return lines;
}