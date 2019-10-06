#include <unordered_map>

struct Variable;
struct TekoObject;
struct TekoType;
struct TekoFunctionType;
struct TekoFunction;

TekoType *TekoTypeType, *TekoObjectType, *TekoVoid, *TekoBool, *TekoInt, *TekoReal, *TekoChar, *TekoString, *TekoModule;

TekoFunctionType *TekoPrintType;

struct TekoPrintFunction;

TekoPrintFunction *TekoPrint;

bool is_subtype(TekoType *sub, TekoType *sup);
bool is_instance(TekoObject *o, TekoType *t);

struct Interpreter;