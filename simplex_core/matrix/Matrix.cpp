#include "Matrix.h"
#include "ExtraMatrixFunctions.h"
#include "../common/Constants.h"
#include <stdexcept>
#include <algorithm>
#include <cmath>
#include <iomanip>
#include <immintrin.h>

void Matrix::validateNegativeRowsOrCols() const {
    if constexpr (DO_VALIDATION) {
        if (rows_ < 0 || cols_ < 0) {
            throw std::invalid_argument("Invalid matrix size");
        }
    }
}

// Constructor
Matrix::Matrix(int rows, int cols) : rows_{rows}, cols_{cols}, entries_(rows * cols) {
    validateNegativeRowsOrCols();
}

void Matrix::validateExceedingRowsOrCols(int row, int col) const {
    if constexpr (DO_VALIDATION) {
        if (row < 0 || row >= rows_) {
            throw std::out_of_range("Matrix row out of bounds");
        } else if (col < 0 || col >= cols_) {
            throw std::out_of_range("Matrix column out of bounds");
        }
    }
}

// Accessors and mutators
double Matrix::at(int row, int col) const {
    validateExceedingRowsOrCols(row, col);

    return entries_[row * cols_ + col];
}

double &Matrix::at(int row, int col) {
    validateExceedingRowsOrCols(row, col);
    
    return entries_[row * cols_ + col];
}

int Matrix::getRows() const {
    return rows_;
}

int Matrix::getCols() const {
    return cols_;
}

// Print matrix
void printTopBottom(int numCols, char symbolStart, char symbolEnd, const std::vector<int> &columnSpacing) {
    std::cout << symbolStart;
    for (int i = 0; i < numCols; ++i) {
        std::cout << " ";
        for (int j = 0; j < columnSpacing[i]; ++j) {
            std::cout << " ";
        }
    }
    std::cout << " " << symbolEnd << std::endl;
    return;
}

int getLengthOfNum(float num) {
    // get length with a log
    int width = (std::abs(num) < 1) ? 1 : static_cast<int>(std::log10(std::abs(num))) + 1;
    if (num < 0.0) {
        // negative sign
        ++width;
    }

    // +1 for decimal point
    return width + 1 + DECIMAL_PRECISION;
}

std::ostream &operator<<(std::ostream &os, const Matrix &matrix) {
    // get column spacing for each column
    std::vector<int> columnSpacing(matrix.getCols(), 0);
    for (int y = 0; y < matrix.getRows(); ++y) {
        for (int x = 0; x < matrix.getCols(); ++x) {
            columnSpacing[x] = std::max(columnSpacing[x], getLengthOfNum(matrix(y, x)));
        }
    }
    
    printTopBottom(matrix.getCols(), '_', '_', columnSpacing);

    for (int y = 0; y < matrix.getRows(); ++y) {
        os << "|";
        for (int x = 0; x < matrix.getCols(); ++x) {
            os << " " << std::setw(columnSpacing[x]) 
                      << std::fixed << std::setprecision(2) 
                      << matrix(y, x);
        }
        os << " |" << std::endl;
    }

    printTopBottom(matrix.getCols(), '-', '-', columnSpacing);
    return os;
}

// Read matrix
std::istream &operator>>(std::istream &is, Matrix &matrix) {
    for (int y = 0; y < matrix.getRows(); ++y) {
        for (int x = 0; x < matrix.getCols(); ++x) {
            double curEntry;
            if (!(is >> curEntry)) {
                throw std::runtime_error("Invalid entry");
            }
            matrix(y, x) = curEntry;
        }
    }
    return is;
}

void printMatrixBasic(std::ostream &os, const Matrix &matrix) {
    for (int y = 0; y < matrix.getRows(); ++y) {
        for (int x = 0; x < matrix.getCols(); ++x) {
            os << matrix(y, x) << " ";
        }
    }
    return;
}

Matrix readAndInitializeMatrix() {
    int numRows, numCols;

    std::cout << "Please enter the number of rows for the matrix" << std::endl;
    if (!(std::cin >> numRows)) {
        throw std::runtime_error("Invalid input for numRows");
    }
    std::cout << "Please enter the number of columns for the matrix" << std::endl;
    if (!(std::cin >> numCols)) {
        throw std::runtime_error("Invalid input for numCols");
    }
    
    Matrix result(numRows, numCols);
    std::cout << "Please type your rows (left to right):" << std::endl;
    std::cin >> result;

    return result;
}

Matrix readAndInitializeMatrixQuiet() {
    int numRows, numCols;

    if (!(std::cin >> numRows)) {
        throw std::runtime_error("Invalid input for numRows");
    }
    if (!(std::cin >> numCols)) {
        throw std::runtime_error("Invalid input for numCols");
    }
    
    Matrix result(numRows, numCols);
    std::cin >> result;

    return result;
}

void Matrix::multiplicationCheck(const Matrix &rhs) const {
    if constexpr (DO_VALIDATION) {
        if (this->getCols() != rhs.getRows()) {
            throw std::out_of_range("Invalid matrix dimensions");
        }
    }
}

// Operations
Matrix Matrix::operator*(const Matrix &rhs) const {
    multiplicationCheck(rhs);
    
    Matrix result(this->getRows(), rhs.getCols());
    for (int i = 0; i < getRows(); ++i) {
        for (int j = 0; j < getCols(); ++j) {
            const double temp = (*this)(i, j);
            for (int k = 0; k < rhs.getCols(); ++k) {
                result(i, k) += temp * rhs(j, k);
            }
        }
    }

    return result;
}


void Matrix::additionSubtractionCheck(const Matrix &rhs) const {
    if constexpr (DO_VALIDATION) {
        if (this->getCols() != rhs.getCols() || this->getRows() != rhs.getRows()) {
            throw std::out_of_range("Invalid matrix dimensions");
        }
    }
}


Matrix Matrix::operator+(const Matrix &rhs) const {
    additionSubtractionCheck(rhs);

    Matrix result(this->getRows(), this->getCols());
    for (int i = 0; i < result.getRows(); ++i) {
        for (int j = 0; j < result.getCols(); ++j) {
            result(i, j) = (*this)(i, j) + rhs(i, j);
        }
    }

    return result;
}

Matrix Matrix::operator-(const Matrix &rhs) const {
    additionSubtractionCheck(rhs);

    Matrix result(this->getRows(), this->getCols());
    for (int i = 0; i < result.getRows(); ++i) {
        for (int j = 0; j < result.getCols(); ++j) {
            result(i, j) = (*this)(i, j) - rhs(i, j);
        }
    }

    return result;
}

void Matrix::swapRows(int row1, int row2) {
    for (int i = 0; i < getCols(); ++i) {
        std::swap((*this)(row1, i), (*this)(row2, i));
    }
}

void Matrix::addRowsScalar(int row1, int row2, double factor, int start) {
    for (int i = start; i < getCols(); ++i) {
        (*this)(row1, i) += factor * (*this)(row2, i);
    }
}

#if defined(__AVX2__) && defined(__FMA__)
void Matrix::addRowsSIMD(int row1, int row2, double factor) {
    __m256d avxFactor = _mm256_set1_pd(factor);
    for (int i = 0; i <= getCols() - 4; i += 4) {
        __m256d avxRow1 = _mm256_loadu_pd(&entries_.data()[row1 * cols_ + i]);
        __m256d avxRow2 = _mm256_loadu_pd(&entries_.data()[row2 * cols_ + i]);

        __m256d avxResult = _mm256_fmadd_pd(avxRow2, avxFactor, avxRow1);

        _mm256_storeu_pd(&entries_.data()[row1 * cols_ + i], avxResult);
    }
}
#endif

void Matrix::addRows(int row1, int row2, double factor) {
#if defined(__AVX2__) && defined(__FMA__)
    addRowsSIMD(row1, row2, factor);
    addRowsScalar(row1, row2, factor, getCols() - getCols() % 4);
#else
    addRowsScalar(row1, row2, factor, 0);
#endif
}

void Matrix::scaleRow(int row, double factor) {
    for (int i = 0; i < getCols(); ++i) {
        (*this)(row, i) *= factor;
    }
}
