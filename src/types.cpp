#include <unordered_map>

class TekoObject;

bool is_subtype(TekoObject *sub, TekoObject *super);

class TekoObject {
public:
	TekoObject* type;
	string to_str();
	static bool is_subtype(TekoObject sub, TekoObject super);
	unordered_map<string,TekoObject*> attrs;
	variant<string, int, float, bool> val;

	static TekoObject TekoType, TekoString, TekoInt, TekoReal, TekoBool;

	void set_parent(TekoObject _parent) {
		if (is_subtype(_parent,TekoObject::TekoType)) {
			attrs["parent"] = &_parent;
		} else {
			throw runtime_error("Issue setting parent");
		}
	}

	TekoObject get(string attr_label) {
		if (attr_label == "type") {
			return type;
		//} else if (attr_label == "__str__") {
		//	to = TekoObject(&TekoString);
		//	to.attrs["__str__"] = to_str();
		//	return to;
		} else {
			return *(attrs[attr_label]);
		}
	}

	TekoObject(TekoObject* _type) {
		if (is_subtype(*_type, TekoObject::TekoType)) {
			type = _type;
		} else {
            throw runtime_error("type snafu");
		}
	}
};

bool TekoObject::is_subtype(TekoObject sub, TekoObject sup) {
	bool sub_is_type = (&sub == &TekoObject::TekoType) || is_subtype(sub.type, TekoObject::TekoType);
	bool sup_is_type = (&sup == &TekoObject::TekoType) || is_subtype(sup.type, TekoObject::TekoType);
	if (sub_is_type && sup_is_type) {
		unordered_map<string,TekoObject*>::const_iterator got = sub.attrs.find("parent");
		if (got == sub.attrs.end()) { return false; }
		else if (got->second == &sup) { return true; }
		else { return is_subtype(*(got->second),sup); }
	} else {
		return false;
	}
};

TekoObject TekoObject::TekoType = TekoObject(NULL);

TekoObject TekoString = TekoObject(&TekoObject::TekoType);
TekoObject TekoInt = TekoObject(&TekoObject::TekoType);
TekoObject TekoReal = TekoObject(&TekoObject::TekoType);
TekoObject TekoBool = TekoObject(&TekoObject::TekoType);