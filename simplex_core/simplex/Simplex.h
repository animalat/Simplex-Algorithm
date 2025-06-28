#ifndef SIMPLEX_H
#define SIMPLEX_H

#include "../matrix/Matrix.h"
#include <vector>
#include <stdexcept>

// Represents the corresponding possible outcomes for a linear program
// (Fundamental Theorem of Linear Programming)
enum class LPResultType { OPTIMAL, UNBOUNDED, INFEASIBLE };

// LPResult is used to return the outcome
struct LPResult {
    LPResultType type;
    // solution is the optimal solution if OPTIMAL, 
    // feasible (r) if UNBOUNDED, 
    // and empty if INFEASIBLE.
    Matrix solution;
    // certificate is y if OPTIMAL (s.t. (c - y^TA) <= 0),
    // d (s.t. Ad = 0) if UNBOUNDED
    // y s.t. y^TA >= 0 but y^Tb < 0 if INFEASIBLE.
    // y = c_B^T * A_B^{-1}
    Matrix certificate;
};

struct PhaseIResult {
    bool feasible;
    std::vector<int> basis;
    // certificate is y if INFEASIBLE
    // otherwise, not meaningful.
    // y = c_B^T * A_B^{-1}
    Matrix certificate;
};

/**
 * Transforms the LP into canonical form (suitable for the simplex algorithm).
 * 
 * Assumes LP is in standard equality form (SEF):
 *     maximize    c^Tx + z
 *     subject to  Ax = b
 *                 x >= 0
 *
 * After this transformation, the basic variables will correspond to the provided basis,
 * and the constraint matrix will be reduced accordingly. The objective function will also
 * be rewritten to account for the current basis.
 * 
 * NOTE: Throws exceptions for invalid input.
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
 */
void canonicalForm(Matrix &objectiveFunc, double &constantTerm, 
                   Matrix &constraintsLHS, Matrix &constraintsRHS, 
                   const std::vector<int> &basis);

/**
 * Runs the simplex algorithm to solve a linear program in standard form:
 *     maximize    c^Tx + z
 *     subject to  Ax = b
 *                 x >= 0
 *
 * This assumes a feasible basis is already provided.
 * Further, the LP must be in standard equality form (SEF) with x >= 0
 *
 * NOTE: Throws exceptions for invalid input.
 * 
 * @param objectiveFunc   Row vector (1 x n), objective coefficients c^T.
 * 
 * @param constantTerm    Scalar (double), constant term z in the objective.
 * 
 * @param constraintsLHS  Matrix (m x n), constraint matrix A.
 *
 * @param constraintsRHS  Column vector (m x 1), right-hand side vector b.
 * 
 * @param basis           Vector of indices (size m), current basis column indices.
 *                        This function updates basis to the final basis used.
 * 
 * @param result          Returns an LPResult with the solution type,
 *                        solution (optimal if such exists, otherwise feasible)
 *                        and the corresponding certificate.
 * 
 *                        In the optimal case: solution is optimal, certificate verifies this.
 * 
 *                        In the unbounded case: solution is feasible, certificate, call it r,
 *                        has property Ar = vector{0}, cr > 0, meaning, for feasible solution x,
 *                        x + tr, (t >= 0), is a feasible solution with arbitrarily large
 *                        objective value for arbitrarily large t.
 */
void simplex(Matrix objectiveFunc, double constantTerm, 
             Matrix constraintsLHS, Matrix constraintsRHS, 
             std::vector<int> &basis, LPResult &result);

/**
 * Runs Phase I algorithm on LP, returning a PhaseIResult.
 *
 * @param constraintsLHS  The LHS of the constraints (A in Ax=b)
 * 
 * @param constraintsRHS  The RHS of the constraints (b in Ax=b)
 * 
 * @return PhaseIResult   Returns a PhaseIResult with solution type (true if feasible, false otherwise)
 *                        If feasible, field basis will store the correct basis. Empty otherwise.
 *                        If infeasible, field certificate contains a certificate of infeasibility, meaningless otherwise.
 */
PhaseIResult phaseI(const Matrix &constraintsLHS, const Matrix &constraintsRHS);

/**
 * Runs the 2-Phase algorithm on a passed LP
 * 
 * @param objectiveFunc   The multiplied row vector in the objective function (c^T in c^T * x + z)
 * 
 * @param constantTerm    The constant term in the objective function (z in c^T * x + z)
 * 
 * @param constraintsLHS  The LHS of the constraints (A in Ax=b)
 * 
 * @param constraintsRHS  The RHS of the constraints (b in Ax=b)
 * 
 * @param result          Returns the result of the LP (type, solution, certificate)
 *                        If infeasible, note that solution is meaningless.
 */
void twoPhase(const Matrix &objectiveFunc, double constantTerm, 
              const Matrix &constraintsLHS, const Matrix &constraintsRHS,
              LPResult &result);

#endif
