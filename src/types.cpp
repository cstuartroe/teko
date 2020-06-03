#include "nodes.h"
#include "types.h"

struct TekoTypeField {
    TekoType *type = 0;
    bool is_mutable = false;

    TekoTypeField(TekoType *_type, bool _is_mutable) {
        type = _type;
        is_mutable = _is_mutable;
    }

    TekoTypeField() {
    }
};

struct TekoStructField {
    TekoType *type;
    TekoObject *deflt = 0;
    bool byref = false;

    TekoStructField(TekoType *_type) {
        type = _type;
    }
};

struct TekoObject {
    TekoType* type;
    unordered_map<string, TekoObject*> attrs;
    string name = "";

    TekoObject(TekoType *_type) {
        setType(_type);
    }

    TekoObject() {
    }

    void setType(TekoType *_type) {
        type = _type;
        attrs["type"] = (TekoObject*) type;
    }

    void set(string attr_name, TekoObject* value);

    TekoTypeField getField(string attr_name);

    string to_str() {
        return name;
    }
};

struct TekoType : TekoObject {
    TekoType *parent;
    unordered_map<string, TekoTypeField> fields;

    TekoType(TekoType *_parent) : TekoObject(TekoTypeType) {
        setParent(_parent);
    }

    TekoType() {
    }

    void setParent(TekoType *_parent) {
        parent = _parent;
        attrs["parent"] = parent;
    }

    void setField(string name, TekoType *type, bool is_mutable) {
        if (getField(name) != 0) {
            throw TekoCompilerException("Name already a field for type: " + name);
        }
        TekoTypeField field = TekoTypeField(type, is_mutable);
        fields[name] = field;
//      TODO:  attrs["fields"] = TekoMap(fields);
    }

    TekoTypeField *getField(string name) {
        if (fields.count(name) == 1) {
            return &fields[name];
        } else if (parent == 0) {
            return 0;
        } else {
            return parent->getField(name);
        }
    }
};

void TekoObject::set(string attr_name, TekoObject* value) {
    if (value->name == ""){
        value->name = to_str() + "." + attr_name;
    }
    TekoTypeField *field = type->getField(attr_name);
    if (field == 0) {
        throw TekoCompilerException("Object of type " + type->name + " has no field " + attr_name);
    } else if (!is_instance(value, field->type)) {
        throw TekoCompilerException("Value has wrong type!");
    } else if (!field->is_mutable && attrs[attr_name] != 0) {
        throw TekoCompilerException("Constant already set!");
    }
    attrs[attr_name] = value;
}

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

struct TekoStruct : TekoType {
    vector<string> argnames;
    vector<TekoStructField> argtypes;

    TekoStruct (vector<string> _argnames, vector<TekoStructField> _argtypes) : TekoType(TekoStructType) {
        if (_argtypes.size() != _argnames.size()) {
            throw runtime_error("Function type mismatch");
        }
        argnames = _argnames;
        argtypes = _argtypes;
    }

};

struct TekoFunctionType : TekoType {
    TekoType *rtype;
    TekoStruct *args;

    TekoFunctionType(TekoType *_rtype, TekoStruct *_args) : TekoType(TekoTypeType) {
        rtype = _rtype;
        args = _args;
    }
};

struct TekoFunction : TekoObject {
    vector<Statement*> stmts;

    TekoFunction(TekoFunctionType *_type, vector<Statement*> _stmts) : TekoObject(_type) {
        stmts = _stmts;
    }

    void execute(Interpreter i8er, ArgsNode *args);
};

struct TekoNamespace : TekoObject {
    unordered_map<string, TekoTypeField> fields;
    TekoNamespace *outer;

    TekoNamespace(TekoNamespace *_outer) : TekoObject(TekoNamespaceType) {
        outer = _outer;
    }

    TekoNamespace *getFieldOwner(string attr_name) {
        // TODO: attributes of type?
        if (fields.count(attr_name) == 1) {
            return this;
        } else if (outer == 0) {
            return 0;
        } else {
            return outer->getFieldOwner(attr_name);
        }
    }

    void declareField(string attr_name, TekoType *type, bool is_mutable) {
        if (getFieldOwner(attr_name) != 0 || type->getField(attr_name) != 0) {
            throw TekoCompilerException("Variable name already declared: " + attr_name);
        }
        fields[attr_name] = TekoTypeField(type, is_mutable);
    }

    void set(string attr_name, TekoObject* value) {
        TekoNamespace *owner = getFieldOwner(attr_name);
        if (owner == 0) {
            throw TekoCompilerException("Namespace has no variable " + attr_name);
        }

        if (value->name == ""){
            value->name = owner->to_str() + "." + attr_name;
        }
        TekoTypeField field = owner->fields[attr_name];
        if (!is_instance(value, field.type)) {
            throw TekoCompilerException("Value has wrong type!");
        } else if (!field.is_mutable && owner->attrs[attr_name] != 0) {
            throw TekoCompilerException("Constant already set!");
        }
        owner->attrs[attr_name] = value;
    }

    TekoObject *get(string attr_name) {
        TekoNamespace *owner = getFieldOwner(attr_name);
        if (owner == 0) {
            throw runtime_error("What happened here?");
        }
        return owner->attrs[attr_name];
    }
};

struct TekoStandardNS : TekoNamespace {
    TekoStandardNS() : TekoNamespace (0) {
        name = "";

        TekoObjectType = new TekoType();
        TekoTypeType = new TekoType();

        TekoObjectType->setType(TekoTypeType);
        TekoObjectType->setParent(0);
        TekoObjectType->setField("type", TekoTypeType, false);

        TekoTypeType->setType(TekoTypeType);
        TekoTypeType->setParent(TekoObjectType);
        TekoTypeType->setField("parent", TekoTypeType, false);

        TekoVoidType = new TekoType(TekoObjectType);
        TekoStructType = new TekoType(TekoObjectType);

        TekoFunctionTypeType = new TekoType(TekoObjectType);
        TekoFunctionTypeType->setField("rtype", TekoTypeType, false);
        TekoFunctionTypeType->setField("args",  TekoStructType, false);

        TekoBoolType = new TekoType(TekoObjectType);
        TekoIntType = new TekoType(TekoObjectType);
        TekoRealType = new TekoType(TekoObjectType);
        TekoCharType = new TekoType(TekoObjectType);
        TekoStringType = new TekoType(TekoObjectType);

        TekoNamespaceType = new TekoType(TekoObjectType);
        type = TekoNamespaceType;

        TekoStruct *to_str_struct = new TekoStruct(vector<string>(), vector<TekoStructField>());
        TekoFunctionType *to_str_type = new TekoFunctionType(TekoStringType, to_str_struct);

        TekoObjectType->setField("to_str", to_str_type, false);

//        TekoMapType field_map_type = new TekoMapType(TekoStringType, TekoTypeType)
//        TekoTypeType->setField("fields", , false);

        declareField("obj",  TekoTypeType, false);
        fields["type"] = TekoTypeField(TekoTypeType, false);
        declareField("void", TekoTypeType, false);
        declareField("bool", TekoTypeType, false);
        declareField("int",  TekoTypeType, false);
        declareField("real", TekoTypeType, false);
        declareField("char", TekoTypeType, false);
        declareField("str",  TekoTypeType, false);
        declareField("namespace",  TekoTypeType, false);

        set("obj",  TekoObjectType);
        set("type", TekoTypeType);
        set("void", TekoVoidType);
        set("bool", TekoBoolType);
        set("int",  TekoIntType);
        set("real", TekoRealType);
        set("char", TekoCharType);
        set("str",  TekoStringType);
        set("namespace", TekoNamespaceType);
    }

    void check(Statement *stmt) {
        printf("%s", stmt->to_str(0).c_str());
    }
};

