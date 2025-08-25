#include "matrix/Matrix.h"
#include "matrix/ExtraMatrixFunctions.h"
#include "simplex/Simplex.h"

constexpr int oneArg = 2;
constexpr int humanReadableOutput = 1;

// use the --humanReadable flag for human-readable output
int main(int argc, char *argv[]) {
    if (argc > oneArg) {
        throw std::invalid_argument("Too many arguments, expected up to 1 argument");
    }

    std::string progArgument;
    if (argc == oneArg) {
        progArgument = argv[humanReadableOutput];
    }

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

    const std::string humanReadableFlag = "--humanReadable";
    if (progArgument == humanReadableFlag) {
        std::cout << res.solution << getResultTypeString(res.type) << res.certificate;
    } else {
        printMatrixBasic(std::cout, res.solution);
        std::cout << std::endl;
        std::cout << getResultTypeString(res.type) << std::endl;
        printMatrixBasic(std::cout, res.certificate);
        std::cout << std::endl;
    }

    return 0;
}
