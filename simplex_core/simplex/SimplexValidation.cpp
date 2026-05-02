#include "SimplexValidation.h"
#include "../common/Constants.h"
#include "../common/Utils.h"

void canonicalFormValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS) {
    if constexpr (DO_VALIDATION) {
        if (objectiveFunc.getRows() != 1) {
            throw std::invalid_argument("Objective function not a row vector");
        } else if (objectiveFunc.getCols() != constraintsLHS.getCols()) {
            throw std::invalid_argument("Objective function and constraint matrix differ in column size");
        }
    }
}

void simplexValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS,
                       const Matrix &constraintsRHS, const std::vector<int> &basis) {
    if constexpr (DO_VALIDATION) {
        if (objectiveFunc.getCols() <= 0 || objectiveFunc.getRows() <= 0) {
            throw std::invalid_argument("Invalid objectiveFunc size");
        } else if (constraintsLHS.getCols() <= 0 || constraintsLHS.getRows() <= 0) {
            throw std::invalid_argument("Invalid constraintsLHS size");
        } else if (constraintsRHS.getCols() <= 0 || constraintsRHS.getRows() <= 0) {
            throw std::invalid_argument("Invalid constraintsRHS size");
        } else if (constraintsLHS.getRows() != constraintsRHS.getRows()) {
            throw std::invalid_argument("constraintsLHS must be same height constraintsRHS");
        } else if (convertToInt(basis.size()) != constraintsLHS.getRows()) {
            throw std::invalid_argument("constraintsLHS height does not match basis size");
        }
    }
}

void phaseIValidation(const Matrix &constraintsLHS, const Matrix &constraintsRHS) {
    if constexpr (DO_VALIDATION) {
        if (constraintsRHS.getRows() <= 0 || constraintsRHS.getCols() != 1) {
            throw std::invalid_argument("Invalid constraintsRHS dimensions");
        } else if (constraintsLHS.getRows() <= 0 || constraintsLHS.getCols() <= 0) {
            throw std::invalid_argument("Invalid constraintsLHS dimensions");
        } else if (constraintsLHS.getRows() != constraintsRHS.getRows()) {
            throw std::invalid_argument("constraintsLHS and constraintsRHS differ in height");
        }
    }
}

void twoPhaseValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS, const Matrix &constraintsRHS) {
    if constexpr (DO_VALIDATION) {
        if (objectiveFunc.getCols() <= 0 || objectiveFunc.getRows() <= 0) {
            throw std::invalid_argument("Invalid objectiveFunc size");
        } else if (constraintsLHS.getCols() <= 0 || constraintsLHS.getRows() <= 0) {
            throw std::invalid_argument("Invalid constraintsLHS size");
        } else if (constraintsRHS.getCols() <= 0 || constraintsRHS.getRows() <= 0) {
            throw std::invalid_argument("Invalid constraintsRHS size");
        } else if (constraintsLHS.getRows() != constraintsRHS.getRows()) {
            throw std::invalid_argument("constraintsLHS must be same height constraintsRHS");
        }
    }
}
