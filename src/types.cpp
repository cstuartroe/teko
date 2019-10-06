#include "nodes.h"
#include "types.h"

struct Variable {
    TekoType *type;
    TekoObject *value;
    bool is_mutable;

    Variable(TekoType *_type, bool _is_mutable) {
        type = _type;
        value = 0;
        is_mutable = _is_mutable;
    }

    TekoObject *get() {
        return value;
    }

    void set(TekoObject *o) {
        if (!is_instance(o, type)) {
            throw runtime_error("Value has wrong type!");
        } else if (!is_mutable && (value != 0)) {
            throw runtime_error("Constant already set!");
        } else {
            value = o;
        }
    }
};

struct TekoObject {
    TekoType* type;
    unordered_map<string, Variable*> attrs;
    char* val;
    string name;

    TekoObject(TekoType *_type) {
        declare("type", TekoTypeType, true);
        set_type(_type);
    }

    TekoObject() {
    }

    void declare(string attr_name, TekoType *t, bool constant) {
        attrs[attr_name] = new Variable(t, constant);
    }

    TekoObject *get(string attr_name) {
        return attrs[attr_name]->get();
    }

    void set(string attr_name, TekoObject* value) {
        if (value->name == ""){
            value->name = to_str() + "." + attr_name;
        }
        attrs[attr_name]->set(value);
    }

    void set_type(TekoType *_type) {
        type = _type;
        set("type", (TekoObject*) type);
    }

    string to_str() {
        return name;
    }
};

struct TekoType : TekoObject {
    TekoType *parent;

    TekoType(TekoType *_parent) : TekoObject(TekoTypeType) {
        parent = _parent;
    }

    TekoType() {
    }
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

struct TekoFunctionType : TekoType {
    TekoType *rtype;
    vector<TekoType*> argtypes;
    vector<string> argnames;

    TekoFunctionType(TekoType *_rtype, vector<TekoType*> _argtypes, vector<string> _argnames) : TekoType(TekoTypeType) {
        rtype = _rtype;
        if (argtypes.size() != argnames.size()) {
            throw runtime_error("Function type mismatch");
        }
        argtypes = _argtypes;
        argnames = _argnames;
    }
};

struct TekoFunction : TekoObject {
    vector<Statement*> stmts;

    TekoFunction(TekoFunctionType *_type, vector<Statement*> _stmts) : TekoObject(_type) {
        stmts = _stmts;
    }

    void execute(Interpreter i8er, ArgsNode *args);
};

struct TekoPrintFunction : TekoFunction {
    TekoPrintFunction() : TekoFunction(0, vector<Statement*>()) { }

    void execute(Interpreter i8er, ArgsNode *args);
};

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

        TekoVoid = new TekoType(TekoObjectType);

        TekoBool = new TekoType(TekoObjectType);
        TekoInt = new TekoType(TekoObjectType);
        TekoReal = new TekoType(TekoObjectType);
        TekoChar = new TekoType(TekoObjectType);
        TekoString = new TekoType(TekoObjectType);

        TekoModule = new TekoType(TekoObjectType);

        TekoPrintType = new TekoFunctionType(TekoVoid, vector<TekoType*>{TekoString}, vector<string>{"s"});
        TekoPrint = new TekoPrintFunction();
        TekoPrint->type = TekoPrintType;
    }

    void check(Statement *stmt) {
        printf("%s", stmt->to_str(0).c_str());
    }
};