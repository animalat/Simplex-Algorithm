#include "Simplex.h"
#include "../matrix/ExtraMatrixFunctions.h"
#include <stdexcept>

Matrix getSubMatrix(const Matrix &matrix, const std::vector<int> &basis) {
    Matrix subMatrix(basis.size(), basis.size());
    for (const int &x : basis) {
        for (int y = 0; y < matrix.getRows(); ++y) {
            subMatrix(y, x) = matrix(y, x);
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
    Matrix inverseTranpose = transpose(basisInverse);
    // dualTranspose = y^T = ((A_B^{-T}) * c_B)^T
    Matrix dualTranspose = transpose(inverseTranpose * getSubMatrix(objectiveFunc, basis));
    objectiveFunc = objectiveFunc - (dualTranspose * constraintsLHS);
    constantTerm += (dualTranspose * constraintsRHS)(0, 0);

    // Change constraints
    constraintsLHS = basisInverse * constraintsLHS;
    constraintsRHS = basisInverse * constraintsRHS;

    return;
}
