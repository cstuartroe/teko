#include <unordered_map>

struct Variable;
struct TekoObject;
struct TekoType;

TekoType *TekoTypeType, *TekoObjectType, *TekoBool, *TekoInt, *TekoReal, *TekoChar;

bool is_subtype(TekoType *sub, TekoType *sup);
bool is_instance(TekoObject *o, TekoType *t);