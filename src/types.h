#include <unordered_map>
#include "nodes.cpp"

struct TekoTypeField;
struct TekoStructField;
struct TekoObject;
struct TekoType;
struct TekoStruct;
struct TekoFunctionType;
struct TekoFunction;

TekoType *TekoTypeType, *TekoObjectType, *TekoVoidType, *TekoBoolType, *TekoIntType,
         *TekoRealType, *TekoCharType, *TekoStringType, *TekoNamespaceType, *TekoStructType, *TekoFunctionTypeType;

bool is_subtype(TekoType *sub, TekoType *sup);
bool is_instance(TekoObject *o, TekoType *t);

struct Interpreter;