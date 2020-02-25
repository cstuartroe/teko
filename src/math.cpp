struct TekoInt : TekoObject {
    int *val;

    TekoInt(int *_val) : TekoObject(TekoIntType) {
        val = _val;
    }
};