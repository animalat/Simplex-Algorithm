#include "ExtraMatrixFunctions.h"
#include <algorithm>
#include <cmath>

constexpr double EPSILON = 1e-9;

Matrix identityMatrix(int size) {
    if (size < 0) {
        throw std::invalid_argument("Invalid matrix size");
    }

    Matrix result(size, size);
    for (int i = 0; i < size; ++i) {
        result(i, i) = 1;
    }
    return result;
}

Matrix rowSwapMatrix(int row1, int row2, int size) {
    if (size <= 0)
        throw std::invalid_argument("Invalid size");
    if (row1 < 0 || row1 >= size || row2 < 0 || row2 >= size)
        throw std::out_of_range("Row index out of bounds");

    Matrix result = identityMatrix(size);
    result(row1, row1) = 0;
    result(row1, row2) = 1;

    result(row2, row2) = 0;
    result(row2, row1) = 1;

    return result;
}

Matrix rowAddMatrix(int row1, int row2, int size, double factor) {
    if (size <= 0)
        throw std::invalid_argument("Invalid size");
    if (row1 < 0 || row1 >= size || row2 < 0 || row2 >= size)
        throw std::out_of_range("Row index out of bounds");

    Matrix result = identityMatrix(size);
    result(row1, row2) = factor;
    return result;
}

Matrix rowMultiplyMatrix(int row, int size, double factor) {
    if (size <= 0)
        throw std::invalid_argument("Invalid size");
    if (row < 0 || row >= size)
        throw std::out_of_range("Row index out of bounds");

    Matrix result = identityMatrix(size);
    result(row, row) *= factor;
    return result;
}


Matrix leftInverse(const Matrix &matrix) {
    if (matrix.getCols() != matrix.getRows()) {
        throw std::domain_error("Matrix has no inverse");
    }

    Matrix matrixCopy = matrix;
    Matrix result = identityMatrix(matrixCopy.getCols());

    for (int x = 0; x < matrixCopy.getCols(); ++x) {
        // make sure matrixCopy(x, x) is non-zero
        if (std::abs(matrixCopy(x, x)) < EPSILON) {
            for (int y = x + 1; y < matrixCopy.getRows(); ++y) {
                if (std::abs(matrixCopy(y, x)) >= EPSILON) {
                    result.swapRows(y, x);
                    matrixCopy.swapRows(y, x);
                    break;
                }

                if (y == (matrixCopy.getRows() - 1)) {
                    throw std::domain_error("Matrix has no inverse");
                }
            }
        }

        // normalize the pivot row 
        result.multRow(x, 1.0/matrixCopy(x, x));
        matrixCopy.multRow(x, 1.0/matrixCopy(x, x));

        // row elimination
        for (int y = 0; y < matrixCopy.getRows(); ++y) {
            if (y == x) {
                continue;
            }

            result.addRows(y, x, -matrixCopy(y, x));
            matrixCopy.addRows(y, x, -matrixCopy(y, x));
        }
    }

    return result;
}
