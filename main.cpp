#include "Matrix.h"
#include "ExtraMatrixFunctions.h"

int main() {
    // Enter matrices:
    Matrix A = readAndInitializeMatrix();
    std::cout << A;
    Matrix B = readAndInitializeMatrix();
    std::cout << B;

    std::cout << leftInverse(A) << std::endl;
    std::cout << leftInverse(B) << std::endl;

    return 0;
}
