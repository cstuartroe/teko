#include <vector>

enum LineType { declaration, assignment, expression, 
                if_statement, while_block, for_block };

struct Line {
    vector<Tag> tags;
    LineType type;
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
    l.type = expression;
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