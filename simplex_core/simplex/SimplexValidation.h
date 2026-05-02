#ifndef SIMPLEX_VALIDATION_H
#define SIMPLEX_VALIDATION_H

#include "../matrix/Matrix.h"

void canonicalFormValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS);

void simplexValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS,
                       const Matrix &constraintsRHS, const std::vector<int> &basis);

void phaseIValidation(const Matrix &constraintsLHS, const Matrix &constraintsRHS);

void twoPhaseValidation(const Matrix &objectiveFunc, const Matrix &constraintsLHS, const Matrix &constraintsRHS);

#endif
