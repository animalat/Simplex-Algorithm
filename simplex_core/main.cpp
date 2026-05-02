#include "matrix/Matrix.h"
#include "matrix/ExtraMatrixFunctions.h"
#include "simplex/Simplex.h"
#include <chrono>

constexpr int maxNumArgs = 3;

// use the --humanReadable flag for human-readable output
// use the --timed flag to display time taken for Simplex
int main(int argc, char *argv[]) {
    if (argc > maxNumArgs) {
        throw std::invalid_argument("Too many arguments, expected up to 2 arguments");
    }

    bool isHumanReadable = false;
    bool isTimed = false;
    for (int i = 1; i < argc; ++i) {
        const std::string humanReadableFlag = "--humanReadable";
        const std::string timedFlag = "--timed";

        isHumanReadable = isHumanReadable || argv[i] == humanReadableFlag;
        isTimed = isTimed || argv[i] == timedFlag;
    }

    // Enter matrices:
    Matrix constraintsLHS;
    Matrix constraintsRHS;
    Matrix objectiveFunc;
    double constantTerm;

    std::vector<int> basis{0, 2, 4, 5};
    LPResult res = {LPResultType::INFEASIBLE, Matrix(0, 0), Matrix(0, 0)};
    
    if (isHumanReadable) {
        constraintsLHS = readAndInitializeMatrix();
        constraintsRHS = readAndInitializeMatrix();
        objectiveFunc = readAndInitializeMatrix();
        std::cout << "Please enter the constant for the objective function." << std::endl;
        std::cin >> constantTerm;
    } else {
        constraintsLHS = readAndInitializeMatrixQuiet();
        constraintsRHS = readAndInitializeMatrixQuiet();
        objectiveFunc = readAndInitializeMatrixQuiet();
        std::cin >> constantTerm;
    }

    auto start = std::chrono::steady_clock::now();
    twoPhase(objectiveFunc, constantTerm, constraintsLHS, constraintsRHS, res);
    auto end = std::chrono::steady_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::nanoseconds>(end - start);

    if (isTimed) {
        std::cout << "Time: " << duration.count() << std::endl;
    }

    if (isHumanReadable) {
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
