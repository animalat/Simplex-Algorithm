#include "Simplex.h"
#include "../matrix/ExtraMatrixFunctions.h"
#include <stdexcept>
#include <limits>

Matrix getSubMatrix(const Matrix &matrix, const std::vector<int> &basis) {
    if (basis.size() > static_cast<size_t>(std::numeric_limits<int>::max())) {
        throw std::invalid_argument("Basis size is greater than an integer");
    }
    
    const int basisSize = static_cast<int>(basis.size());
    Matrix subMatrix(matrix.getRows(), basisSize);
    for (int x = 0; x < basisSize; ++x) {
        for (int y = 0; y < matrix.getRows(); ++y) {
            subMatrix(y, x) = matrix(y, basis[x]);
        }
    }

    return subMatrix;
}

void canonicalForm(Matrix &objectiveFunc, double &constantTerm, Matrix &constraintsLHS, Matrix &constraintsRHS, const std::vector<int> &basis) {
    // Note that this function needs constraint x >= 0, for Ax = b
    if (objectiveFunc.getRows() != 1) {
        throw std::invalid_argument("Objective function not a row vector");
    } else if (objectiveFunc.getCols() != constraintsLHS.getCols()) {
        throw std::invalid_argument("Objective function and constraint matrix differ in column size");
    }
    
    Matrix basisInverse = leftInverse(getSubMatrix(constraintsLHS, basis));
    
    // Change objective function
    Matrix inverseTranspose = transpose(basisInverse);
    // y^T = ((A_B^{-T}) * c_B)^T
    Matrix yT = transpose(inverseTranspose * transpose(getSubMatrix(objectiveFunc, basis)));
    objectiveFunc = objectiveFunc - (yT * constraintsLHS);
    constantTerm += (yT * constraintsRHS)(0, 0);

    // Change constraints
    constraintsLHS = basisInverse * constraintsLHS;
    constraintsRHS = basisInverse * constraintsRHS;

    return;
}

int simplex(Matrix &objectiveFunc, double constantTerm, Matrix &constraintsLHS, Matrix &constraintsRHS, std::vector<int> &basis) {
    canonicalForm(objectiveFunc, constantTerm, constraintsLHS, constraintsRHS, basis);
    
}