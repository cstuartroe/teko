TekoObject TekoType = TekoObject();

bool is_subtype(TekoObject sub, TekoObject sup) {
    cout << "&sub in is_subtype " << &sub << endl;
    cout << "&sup in is_subtype " << &sup << endl;
    cout << "&TekoType in is_subtype " << &TekoType << endl;

    bool sub_is_type;
    if (&sub == &TekoType) {
        sub_is_type = true;
        //} else if (is_subtype(*(sub.type), TekoType)) {
        //    sub_is_type = true;
    } else {
        sub_is_type = false;
    }

    bool sup_is_type;
    if (&sup == &TekoType) {
        sup_is_type = true;
        //} else if (is_subtype(*(sup.type), TekoType)) {
        //    sup_is_type = true;
    } else {
        sup_is_type = false;
    }

    cout << "bools finished" << endl;

    if (sub_is_type && sup_is_type) {
        unordered_map<string,TekoObject*>::const_iterator got = sub.attrs.find("parent");
        if (got == sub.attrs.end()) { return false; }
        else if (got->second == &sup) { return true; }
        else { return is_subtype(*(got->second),sup); }
    } else {
        return false;
    }
};

TekoObject::TekoObject(TekoObject* _type) {
    cout << "type in Constructor " << type << endl;;
    cout << "&TekoType in constructor " << &TekoType << endl;
    if ( is_subtype(*_type, TekoType)) {
        type = _type;
    } else {
        throw runtime_error("type snafu");
    }
}

void TekoObject::set_parent(TekoObject _parent) {
    if (is_subtype(_parent,TekoType)) {
        attrs["parent"] = &_parent;
    } else {
        throw runtime_error("Issue setting parent");
    }
}

TekoObject TekoString = TekoObject(&TekoType);
TekoObject TekoObject::TekoInt = TekoObject(&TekoType);
TekoObject TekoObject::TekoReal = TekoObject(&TekoType);
TekoObject TekoObject::TekoBool = TekoObject(&TekoType);