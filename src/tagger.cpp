#include <vector>
#include <set>

using namespace std;

enum TagType { LabelType, StringType, IntType, RealType,
				BoolType, IfType, ElseType, SemicolonType,
				ColonType, CommaType, QMarkType, AttrType, 
				LParType, RParType, LBracType, RBracType, 
				BinOpType, SetterType, ComparisonType, 
				ConversionType };

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

class AttrTag : public Tag {
public:
	AttrTag() { 
		type = AttrType; 
		repr = "AttrTag";
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

class LBracTag : public Tag {
public:
	LBracTag() {
		type = LBracType;
		repr = "LBracTag";
	}
};

class RBracTag : public Tag {
public:
	RBracTag() {
		type = RBracType;
		repr = "RBracTag";
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

vector<Tag> get_tags(vector<string> tokens) {
	vector<Tag> tags;
	for (int i = 0; i < tokens.size(); i++) {
		string token = tokens[i];
		if (token == ";") {
			tags.push_back(SemicolonTag()); 
		} else if (token == ":") {
			tags.push_back(ColonTag());
		} else if (token == ",") {
			tags.push_back(CommaTag());
		} else if (token == "?") {
			tags.push_back(QMarkTag());
		}

		else if (token == "(") {
			tags.push_back(LParTag());
		} else if (token == ")") {
			tags.push_back(RParTag());
		} else if (token == "{") {
			tags.push_back(LBracTag());
		} else if (token == "}") {
			tags.push_back(RBracTag());
		} 

		else if (token == "if") {
			tags.push_back(IfTag());
		} else if (token == "else") {
			tags.push_back(ElseTag());
		} 

		else if (token == "+") {
			tags.push_back(BinOpTag(add));
		} else if (token == "-") {
			tags.push_back(BinOpTag(sub));
		} else if (token == "*") {
			tags.push_back(BinOpTag(mult));
		} else if (token == "/") {
			tags.push_back(BinOpTag(divs));
		} else if (token == "^") {
			tags.push_back(BinOpTag(exp));
		} else if (token == "%%") {
			tags.push_back(BinOpTag(mod));
		} else if (token == "&&") {
			tags.push_back(BinOpTag(conj));
		} else if (token == "||") {
			tags.push_back(BinOpTag(disj));
		}

		else if (token == "=") {
			tags.push_back(SetterTag(normal));
		} else if (token == "+=") {
			tags.push_back(SetterTag(setadd));
		} else if (token == "-=") {
			tags.push_back(SetterTag(setsub));
		} else if (token == "*=") {
			tags.push_back(SetterTag(setmult));
		} else if (token == "/=") {
			tags.push_back(SetterTag(setdivs));
		} else if (token == "^=") {
			tags.push_back(SetterTag(setexp));
		} else if (token == "%%=") {
			tags.push_back(SetterTag(setmod));
		}

		else if (token == "==") {
			tags.push_back(ComparisonTag(eq));
		} else if (token == "!=") {
			tags.push_back(ComparisonTag(neq));
		} else if (token == "<") {
			tags.push_back(ComparisonTag(lt));
		} else if (token == ">") {
			tags.push_back(ComparisonTag(gt));
		} else if (token == "<=") {
			tags.push_back(ComparisonTag(leq));
		} else if (token == ">=") {
			tags.push_back(ComparisonTag(geq));
		} else if (token == "<:") {
			tags.push_back(ComparisonTag(subtype));
		}

		else if (token == ".") {
			tags.push_back(ConversionTag(toReal));
		} else if (token == "$") {
			tags.push_back(ConversionTag(toStr));
		} else if (token == "[]") {
			tags.push_back(ConversionTag(toArr));
		} else if (token == "{}") {
			tags.push_back(ConversionTag(toList));
		} else if (token == "<>") {
			tags.push_back(ConversionTag(toSet));
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
			if (i > 0 && tokens[i-1] == ".") {
				tags.pop_back(); // period . gets interpreted as a conversion unless followed by a label
				tags.push_back(AttrTag());
			} 
			tags.push_back(LabelTag(token));
		}
	}
	return tags;
}