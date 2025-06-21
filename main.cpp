#include "Matrix.h"

int main() {
    // Enter matrices:
    Matrix A = readAndInitializeMatrix();
    A.printMatrix();
    Matrix B = readAndInitializeMatrix();
    B.printMatrix();

    (A*B).printMatrix();
    
    return 0;
}
