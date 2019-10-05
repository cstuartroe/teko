#include "nodes.h"
#include "types.h"

struct Variable {
    TekoType *type;
    TekoObject *value;
    bool constant;

    Variable(TekoType *_type, bool _constant) {
        type = _type;
        value = 0;
        constant = _constant;
    }

    TekoObject *get() {
        return value;
    }

    void set(TekoObject *o) {
        if (!is_instance(o, type)) {
            throw runtime_error("Value has wrong type!");
        } else if (constant && (value != 0)) {
            throw runtime_error("Constant already set!");
        } else {
            value = o;
        }
    }
};

struct TekoObject {
    TekoType* type;
    string to_str();
    unordered_map<string, Variable*> attrs;
    char* val;

    void declare(string attr_name, TekoType *t, bool constant) {
        attrs[attr_name] = new Variable(t, constant);
    }

    TekoObject *get(string attr_name) {
        return attrs[attr_name]->get();
    }

    void set(string attr_name, TekoObject* value) {
        attrs[attr_name]->set(value);
    }

    void set_type(TekoType *_type) {
        type = _type;
        set("type", (TekoObject*) type);
    }

    TekoObject(TekoType *_type) {
        declare("type", TekoTypeType, true);
        set_type(_type);
    }

    TekoObject() {}
};

struct TekoType : TekoObject {
    TekoType *parent;

    TekoType(TekoType *_parent) : TekoObject(TekoTypeType) {
        parent = _parent;
    }

    TekoType() {}
};

bool is_subtype(TekoType *sub, TekoType *sup) {
    if (sub == sup) {
        return true;
    } else if (sub == TekoObjectType) {
        return false;
    } else {
        return is_subtype(sub->parent, sup);
    }
};

bool is_instance(TekoObject *o, TekoType *t) {
    return is_subtype(o->type, t);
}





struct TekoTypeChecker {
    TekoTypeChecker() {
        typesetup();
    }

    void typesetup() {
        TekoObjectType = new TekoType();
        TekoTypeType = new TekoType();

        TekoObjectType->type = TekoTypeType;
        TekoObjectType->parent = 0;
        TekoTypeType->type = TekoTypeType;
        TekoTypeType->parent = TekoObjectType;

        TekoBool = new TekoType(TekoObjectType);
        TekoInt = new TekoType(TekoObjectType);
        TekoReal = new TekoType(TekoObjectType);
        TekoChar = new TekoType(TekoObjectType);
    }

    void check(Statement *stmt) {
        printf("%s", stmt->to_str(0).c_str());
        if (stmt->next != 0) {
            check(stmt->next);
        }
    }
};