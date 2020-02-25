struct TekoPrintFunction : TekoFunction {
    TekoPrintFunction(TekoFunctionType *_type) : TekoFunction(_type, vector<Statement*>()) { }

    void execute(Interpreter i8er, ArgsNode *args);
};

TekoStruct *TekoPrintStruct = new TekoStruct(vector<string>{"s"}, vector<TekoStructField>{TekoStructField(TekoStringType)});
TekoFunctionType *TekoPrintType = new TekoFunctionType(TekoVoidType, TekoPrintStruct);
TekoPrintFunction *TekoPrint = new TekoPrintFunction(TekoPrintType);


