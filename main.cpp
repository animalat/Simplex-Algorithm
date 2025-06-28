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

    std::vector<int> basis{0, 2, 4, 5};
    LPResult res = {LPResultType::INFEASIBLE, Matrix(0, 0), Matrix(0, 0)};
    
    twoPhase(C, z, A, B, res);
    std::cout << res.certificate << res.solution << (res.type == LPResultType::INFEASIBLE);

    return 0;
}
