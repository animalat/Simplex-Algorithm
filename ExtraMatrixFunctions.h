#ifndef EXTRA_MATRIX_FUNC_H
#define EXTRA_MATRIX_FUNC_H

#include "Matrix.h"

Matrix identityMatrix(int size);

Matrix rowSwapMatrix(int row1, int row2, int size);

// row1 <- row1 + row2
Matrix rowAddMatrix(int row1, int row2, int size, double factor = 1.0);

Matrix rowMultiplyMatrix(int row, int size, double factor);

Matrix leftInverse(const Matrix &matrix);

#endif
