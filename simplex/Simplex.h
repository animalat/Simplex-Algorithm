#ifndef SIMPLEX_H
#define SIMPLEX_H

#include "../matrix/Matrix.h"
#include <vector>

enum class LPResultType { OPTIMAL, UNBOUNDED, INFEASIBLE };

struct LPResult {
    LPResultType type;
    // solution is the optimal solution if OPTIMAL, 
    // feasible (r) if UNBOUNDED, 
    // and empty if INFEASIBLE.
    Matrix solution;
    // certificate is y if OPTIMAL (s.t. (c - y^TA) <= 0),
    // d (s.t. Ad = 0) if UNBOUNDED
    // y s.t. y^TA >= 0 but y^Tb < 0 if INFEASIBLE.
    Matrix certificate;
};

/**
 * Transforms the LP into canonical form (suitable for the simplex algorithm).
 * 
 * Assumes LP is in standard equality form (SEF):
 *     maximize    c^Tx
 *     subject to  Ax = b
 *                 x >= 0
 *
 * After this transformation, the basic variables will correspond to the provided basis,
 * and the constraint matrix will be reduced accordingly. The objective function will also
 * be rewritten to account for the current basis.
 *
 * @param objectiveFunc   Row vector (1 x n). objective coefficients c^T
 *                        Gets modified to the reduced objective function.
 * 
 * @param constantTerm    Scalar (double). constant term from the objective function z
 *                        Gets updated during transformation (constantTerm + y^Tb).
 * 
 * @param constraintsLHS  Matrix (m x n). LHS matrix A in constraints Ax = b
 *                        Gets overwritten to A_B^{-1}A after canonicalization.
 *
 * @param constraintsRHS  Column vector (m x 1). RHS vector b in Ax = b
 *                        Gets overwritten to A_B^{-1}b.
 * 
 * @param basis           Vector of indices (size m), column indices in A that form the initial basis.
 *                        These must point to linearly independent columns in A.
 * 
 * @return Matrix         Returns y = ((A_B^{-T}) * c_B) for use as a certificate
 */
Matrix canonicalForm(Matrix &objectiveFunc, double &constantTerm, 
                   Matrix &constraintsLHS, Matrix &constraintsRHS, 
                   const std::vector<int> &basis);

/**
 * Runs the simplex algorithm to solve a linear program in standard form:
 *     maximize    c^Tx
 *     subject to  Ax = b
 *                 x >= 0
 *
 * This assumes a feasible basis is already provided.
 * Further, the LP must be in standard equality form (SEF) with x >= 0
 *
 * @param objectiveFunc   Row vector (1 x n), objective coefficients c^T.
 *                        Gets modified during pivot steps to the reduced form.
 * 
 * @param constantTerm    Scalar (double), constant term z in the objective.
 *                        Updated as the algorithm proceeds.
 * 
 * @param constraintsLHS  Matrix (m x n), constraint matrix A.
 *                        Gets transformed as pivot operations proceed.
 *
 * @param constraintsRHS  Column vector (m x 1), right-hand side vector b.
 *                        Updated with each pivot.
 * 
 * @param basis           Vector of indices (size m), current basis column indices.
 *                        Gets updated in-place as the algorithm pivots.
 */
void simplex(Matrix &objectiveFunc, double constantTerm, 
            Matrix &constraintsLHS, Matrix &constraintsRHS, 
            std::vector<int> &basis, LPResult &result);

#endif