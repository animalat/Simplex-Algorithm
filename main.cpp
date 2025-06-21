#include "Matrix.h"

int main() {
    // Enter matrices:
    Matrix A = readAndInitializeMatrix();
    std::cout << A;
    Matrix B = readAndInitializeMatrix();
    std::cout << B;

    std::cout << (A*B);

    return 0;
}
