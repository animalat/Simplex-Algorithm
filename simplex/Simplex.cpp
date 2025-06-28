#include "Simplex.h"
#include "../matrix/ExtraMatrixFunctions.h"
#include "../common/Constants.h"
#include "../common/Utils.h"
#include <unordered_set>
#include <algorithm>
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

void canonicalForm(Matrix &objectiveFunc, double &constantTerm, 
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
    return;
}

// Phase II of the algorithm
void simplex(Matrix objectiveFunc, double constantTerm, 
             Matrix constraintsLHS, Matrix constraintsRHS, 
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

    // Store some original data for obtaining certificates
    const Matrix origObjective = objectiveFunc;
    const Matrix origConstraintLHS = constraintsLHS;

    Matrix currentSolution(constraintsLHS.getCols(), 1);
    while (true) {
        canonicalForm(objectiveFunc, constantTerm, 
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

                // get certificate (c_B * A_B^{-1}), notice that our objective is a row vector in this.
                Matrix certificate = getSubMatrix(origObjective, basis) * 
                                     leftInverse(getSubMatrix(origConstraintLHS, basis));
                result.certificate = certificate;

                result.solution = currentSolution;
                return;
            }
        }
        
        // leaving variable is curMinIndex
        double curMinValue = std::numeric_limits<double>::infinity();
        constexpr int UNBOUNDED_INDEX = -1;
        int curMinIndex = UNBOUNDED_INDEX;
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

        if (curMinIndex == UNBOUNDED_INDEX) {
            // Unbounded case
            Matrix certificateUnbounded(constraintsLHS.getCols(), 1);
            certificateUnbounded(enteringVariableCol, 0) = 1;           // This was the entering variable 
                                                                        // (before we knew the LP is unbounded)
            // Create the certificate: -t*A_{enteringVariableCol}
            for (int i = 0; i < convertToInt(basis.size()); ++i) {
                double entry = constraintsLHS(i, enteringVariableCol);
                if (std::abs(entry) < EPSILON) {
                    entry = 0.0;
                }
                // negate (since the entries are negative), i.e., -t*A_{enteringVariableCol}
                certificateUnbounded(basis[i], 0) = -entry;
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
            throw std::invalid_argument("Invalid basis (tried removing non-existent element)");
        }

        // Insert enteringVariableCol into basis in sorted order
        auto insertIt = std::lower_bound(basis.begin(), basis.end(), enteringVariableCol);
        basis.insert(insertIt, enteringVariableCol);
    }
}

// Phase I of the algorithm (determine a possibly feasible basis)
std::vector<int> phaseI(const Matrix &constraintsLHS, const Matrix &constraintsRHS) {
    if (constraintsRHS.getRows() <= 0 || constraintsRHS.getCols() != 1) {
        throw std::invalid_argument("Invalid constraintsRHS dimensions");
    } else if (constraintsLHS.getRows() <= 0 || constraintsLHS.getCols() <= 0) {
        throw std::invalid_argument("Invalid constraintsLHS dimensions");
    } else if (constraintsLHS.getRows() != constraintsRHS.getRows()) {
        throw std::invalid_argument("constraintsLHS and constraintsRHS differ in height");
    }
    
    // auxiliaryLHS will be augmented matrix in form [A | I]
    Matrix auxiliaryLHS(constraintsLHS.getRows(), constraintsLHS.getCols() + constraintsLHS.getRows());
    Matrix auxiliaryRHS = constraintsRHS;

    // Populate auxiliaryLHS with values from constraintsLHS
    for (int x = 0; x < constraintsLHS.getCols(); ++x) {
        for (int y = 0; y < constraintsLHS.getRows(); ++y) {
            auxiliaryLHS(y, x) = constraintsLHS(y, x);
        }
    }

    // Multiply rows s.t. auxiliaryRHS >= \vector{0}
    for (int i = 0; i < auxiliaryRHS.getRows(); ++i) {
        if (auxiliaryRHS(i, 0) < EPSILON) {
            auxiliaryLHS.scaleRow(i, -1.0);
            auxiliaryRHS.scaleRow(i, -1.0);
        }
    }

    // Augment with I:
    for (int i = 0; i < constraintsLHS.getRows(); ++i) {
        auxiliaryLHS(i, constraintsLHS.getCols() + i) = 1.0;
    }

    // Populate basis with values m, m + 1, ..., m + n - 1
    // Where m = constraintsLHS.getCols(), n = constraintsLHS.getRows()
    std::vector<int> basis;
    for (int i = 0; i < constraintsLHS.getRows(); ++i) {
        basis.emplace_back(constraintsLHS.getCols() + i);
    }

    // Populate the auxiliary objective, \vector{-1}x' 
    // where x' is the auxiliary variables
    Matrix auxiliaryObjective(1, auxiliaryLHS.getCols());
    for (int i = constraintsLHS.getCols(); i < auxiliaryLHS.getCols(); ++i) {
        auxiliaryObjective(0, i) = -1.0;
    }

    // Run the simplex algorithm on our auxiliary LP to 
    // determine a feasible solution or infeasibility.
    LPResult result;
    simplex(auxiliaryObjective, 0.0, auxiliaryLHS, auxiliaryRHS, basis, result);



    // Note that the basis may not be feasible for the original LP, 
    // we must determine this in the caller.
    return basis;
}
