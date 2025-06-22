#include "matrix/Matrix.h"
#include "matrix/ExtraMatrixFunctions.h"
#include "simplex/Simplex.h"

int main() {
    // Enter matrices:
    Matrix A = readAndInitializeMatrix();
    std::cout << A;
    Matrix B = readAndInitializeMatrix();
    std::cout << B;
    Matrix C = readAndInitializeMatrix();
    std::cout << C;
    double z;
    std::cin >> z;

    canonicalForm(C, z, A, B, std::vector{0, 1, 3});
    std::cout << C << z << '\n' << A << B << std::endl;

    return 0;
}
