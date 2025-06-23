#include "Simplex.h"
#include "../matrix/ExtraMatrixFunctions.h"
#include "../common/Constants.h"
#include "../common/Utils.h"
#include <algorithm>
#include <stdexcept>
#include <limits>

Matrix getSubMatrix(const Matrix &matrix, const std::vector<int> &basis) {
    Matrix subMatrix(matrix.getRows(), convertToInt(basis.size()));
    for (int x = 0; x < convertToInt(basis.size()); ++x) {
        for (int y = 0; y < matrix.getRows(); ++y) {
            subMatrix(y, x) = matrix(y, basis[x]);
        }
    }

    return subMatrix;
}

Matrix canonicalForm(Matrix &objectiveFunc, double &constantTerm, 
                     Matrix &constraintsLHS, Matrix &constraintsRHS, 
                     const std::vector<int> &basis) {
    // Note that this function needs constraint x >= 0, for Ax = b
    if (objectiveFunc.getRows() != 1) {
        throw std::invalid_argument("Objective function not a row vector");
    } else if (objectiveFunc.getCols() != constraintsLHS.getCols()) {
        throw std::invalid_argument("Objective function and constraint matrix differ in column size");
    }
    
    Matrix basisInverse = leftInverse(getSubMatrix(constraintsLHS, basis));
    
    // Change objective function
    Matrix inverseTranspose = transpose(basisInverse);
    // y = (A_B^{-T}) * c_B
    Matrix y = inverseTranspose * transpose(getSubMatrix(objectiveFunc, basis));
    Matrix yT = transpose(y);
    objectiveFunc = objectiveFunc - (yT * constraintsLHS);
    constantTerm += (yT * constraintsRHS)(0, 0);

    // Change constraints
    constraintsLHS = basisInverse * constraintsLHS;
    constraintsRHS = basisInverse * constraintsRHS;

    // Return certificate
    return y;
}

void simplex(Matrix &objectiveFunc, double constantTerm, 
             Matrix &constraintsLHS, Matrix &constraintsRHS, 
             std::vector<int> &basis, LPResult &result) {
    if (objectiveFunc.getCols() <= 0 || objectiveFunc.getRows() <= 0) {
        throw std::invalid_argument("Invalid objectiveFunc size");
    } else if (constraintsLHS.getCols() <= 0 || constraintsLHS.getRows() <= 0) {
        throw std::invalid_argument("Invalid constraintsLHS size");
    } else if (constraintsRHS.getCols() <= 0 || constraintsRHS.getRows() <= 0) {
        throw std::invalid_argument("Invalid constraintsRHS size");
    } else if (constraintsLHS.getRows() != constraintsRHS.getRows()) {
        throw std::invalid_argument("constraintsLHS must be same height constraintsRHS");
    }

    Matrix certificate(constraintsLHS.getRows(), 1);
    Matrix currentSolution(constraintsLHS.getCols(), 1);
    while (true) {
        certificate = canonicalForm(objectiveFunc, constantTerm, 
                                    constraintsLHS, constraintsRHS, 
                                    basis);
        // Set new basic feasible solution
        currentSolution = Matrix(constraintsLHS.getCols(), 1);
        for (int i = 0; i < convertToInt(basis.size()); ++i) {
            currentSolution(basis[i], 0) = constraintsRHS(i, 0);
        }

        // Find first positive element in objective function (Bland's rule)
        int enteringVariableCol = -1;
        for (int i = 0; i < objectiveFunc.getCols(); ++i) {
            if (objectiveFunc(0, i) > EPSILON) {
                enteringVariableCol = i;
                break;
            }

            if (i == objectiveFunc.getCols() - 1) {
                // We've found an optimal solution
                result.type = LPResultType::OPTIMAL;
                result.certificate = certificate;
                result.solution = currentSolution;
                return;
            }
        }
        
        // leaving variable is curMinIndex
        double curMinValue = std::numeric_limits<double>::infinity();
        int curMinIndex = -1;
        for (int i = 0; i < constraintsLHS.getRows(); ++i) {
            if (constraintsLHS(i, enteringVariableCol) < EPSILON) {
                continue;
            }

            const double currentValue = constraintsRHS(i, 0) / constraintsLHS(i, enteringVariableCol);
            if (currentValue < (curMinValue - EPSILON)) {
                curMinIndex = i;
                curMinValue = currentValue;
            } else if (std::abs(currentValue - curMinValue) < EPSILON && i < curMinIndex) {
                // i < curMinIndex is the tiebreaker rule (Bland's rule)
                curMinIndex = i;
                curMinValue = currentValue;
            }
        }

        if (curMinIndex == -1) {
            // unbounded
            Matrix certificateUnbounded(constraintsLHS.getCols(), 1);
            certificateUnbounded(enteringVariableCol, 0) = 1;
            for (const int &k : basis) {
                double entry = constraintsLHS(k, enteringVariableCol);
                if (std::abs(entry) < EPSILON) {
                    entry = 0.0;
                }
                certificateUnbounded(k, 0) = entry;
            }
            result.type = LPResultType::UNBOUNDED;
            result.certificate = certificateUnbounded;
            result.solution = currentSolution;
            return;
        }

        // Remove curMinIndex from basis
        auto removeIt = std::find(basis.begin(), basis.end(), basis[curMinIndex]);
        if (removeIt != basis.end()) {
            basis.erase(removeIt);
        } else {
            throw std::invalid_argument("this shouldn't have happened");
        }

        // Insert enteringVariableCol into basis in sorted order
        auto insertIt = std::lower_bound(basis.begin(), basis.end(), enteringVariableCol);
        basis.insert(insertIt, enteringVariableCol);
    }
}