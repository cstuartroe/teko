#include <unordered_map>

class TekoObject {
public:
    TekoObject* type;
    string to_str();
    unordered_map<string,TekoObject*> attrs;
    variant<string, int, float, bool> val;

    //static TekoObject TekoType, TekoString, TekoInt, TekoReal, TekoBool;

    void set_parent(TekoObject _parent);

    TekoObject get(string attr_label) {
        if (attr_label == "type") {
            return *type;
        //} else if (attr_label == "__str__") {
        //    to = TekoObject(&TekoString);
        //    to.attrs["__str__"] = to_str();
        //    return to;
        } else {
            return *(attrs[attr_label]);
        }
    }

    TekoObject(TekoObject* _type);

    TekoObject() {
    }
};

bool is_subtype(TekoObject sub, TekoObject super);

TekoObject TekoType, TekoString, TekoInt, TekoReal, TekoBool;