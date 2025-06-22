#ifndef SIMPLEX_H
#define SIMPLEX_H

#include "../matrix/Matrix.h"
#include <vector>

void canonicalForm(Matrix &objectiveFunc, double &constantTerm, Matrix &constraintsLHS, Matrix &constraintsRHS, const std::vector<int> &basis);

int simplex(Matrix &objectiveFunc, double constantTerm, Matrix &constraintsLHS, Matrix &constraintsRHS, std::vector<int> &basis);

#endif