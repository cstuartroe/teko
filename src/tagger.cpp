#include <vector>
#include <set>

using namespace std;

enum TagType { LabelType, StringType, IntType, RealType,
				BoolType, IfType, ElseType, SemicolonType,
				LParType, RParType, BinOpType, SetterType,
				ComparisonType };

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

class LParTag : public Tag {
public:
	LParTag() {
		type = LParType;
		repr = "LParTag";
	}
};

class RParTag : public Tag {
public:
	RParTag() {
		type = RParType;
		repr = "RParTag";
	}
};

enum BinOp { add, sub, mult, divs, exp, mod };

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
		}
		repr = "BinOpTag " + opstr;
	}
};

enum Setter { normal };

class SetterTag : public Tag {
	Setter s;
public:
	SetterTag(Setter _s) {
		type = SetterType;
		s = _s;
		string sstr;
		switch (s) {
			case normal: sstr = "="; break;
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
			case eq: cstr = "=="; break;
		}
		repr = "ComparisonTag " + cstr;
	}
};

vector<Tag> get_tags(vector<string> tokens) {
	vector<Tag> tags;
	for (int i = 0; i < tokens.size(); i++) {
		string token = tokens[i];
		if (token == ";") {
			tags.push_back(SemicolonTag()); 
		} else if (token == "(") {
			tags.push_back(LParTag());
		} else if (token == ")") {
			tags.push_back(RParTag());
		} else if (token == "if") {
			tags.push_back(IfTag());
		} else if (token == "else") {
			tags.push_back(ElseTag());
		} 

		else if (token == "+") {
			tags.push_back(BinOpTag(add));
		}

		else if (token == "=") {
			tags.push_back(SetterTag(normal));
		} 

		else if (token == "==") {
			tags.push_back(ComparisonTag(eq));
		}

		else if (token == "true") {
			tags.push_back(BoolTag(true));
		} else if (token == "false") {
			tags.push_back(BoolTag(false));
		} 

		else if (token[0] == '"') {
			tags.push_back(StringTag(token));
		} else if (in_charset(token[0],nums)) {
			size_t found = token.find('.');
			if (found == string::npos) {
				tags.push_back(IntTag(stoi(token)));
			} else {
				tags.push_back(RealTag(stof(token)));
			}
		} else {
			tags.push_back(LabelTag(token));
		}
	}
	return tags;
}