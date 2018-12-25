#include <iostream>

using namespace std;

void compiler_error(string message) {
    cout << "Teko Compiler Error: " << message << endl;
    exit (EXIT_FAILURE);
}