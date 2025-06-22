#include "Matrix.h"
#include "ExtraMatrixFunctions.h"
#include <stdexcept>
#include <algorithm>
#include <cmath>
#include <iomanip>

constexpr int DECIMAL_PRECISION = 2;

// Constructor
Matrix::Matrix(int rows, int cols) : rows_{rows}, cols_{cols}, 
                                entries_(rows, std::vector<double>(cols, 0.0)) {
    if (rows < 0 || cols < 0) {
        throw std::invalid_argument("Invalid matrix size");
    }
}

// Accessors and mutators
double Matrix::operator()(int row, int col) const {
    if (row < 0 || row >= rows_) {
        throw std::out_of_range("Matrix row out of bounds");
    } else if (col < 0 || col >= cols_) {
        throw std::out_of_range("Matrix column out of bounds");
    }

    return entries_[row][col];
}

double &Matrix::operator()(int row, int col) {
    if (row < 0 || row >= rows_) {
        throw std::out_of_range("Matrix row out of bounds");
    } else if (col < 0 || col >= cols_) {
        throw std::out_of_range("Matrix column out of bounds");
    }
    
    return entries_[row][col];
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

std::ostream& operator<<(std::ostream& os, const Matrix& matrix) {
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

std::istream& operator>>(std::istream& is, Matrix& matrix) {
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

// Operations
Matrix Matrix::operator*(const Matrix &rhs) const {
    if (this->getCols() != rhs.getRows()) {
        throw std::out_of_range("Invalid matrix dimensions");
    }
    
    Matrix result(this->getRows(), rhs.getCols());
    
    for (int i = 0; i < result.getRows(); ++i) {
        for (int j = 0; j < result.getCols(); ++j) {
            for (int k = 0; k < rhs.getRows(); ++k) {
                result(i, j) += (*this)(i, k) * rhs(k, j);
            }
        }
    }

    return result;
}


Matrix Matrix::operator+(const Matrix &rhs) const {
    if (this->getCols() != rhs.getCols() || this->getRows() != rhs.getRows()) {
        throw std::out_of_range("Invalid matrix dimensions");
    }

    Matrix result(this->getRows(), this->getCols());
    for (int i = 0; i < result.getRows(); ++i) {
        for (int j = 0; j < result.getCols(); ++j) {
            result(i, j) = (*this)(i, j) + rhs(i, j);
        }
    }

    return result;
}

Matrix Matrix::operator-(const Matrix &rhs) const {
    if (this->getCols() != rhs.getCols() || this->getRows() != rhs.getRows()) {
        throw std::out_of_range("Invalid matrix dimensions");
    }

    Matrix result(this->getRows(), this->getCols());
    for (int i = 0; i < result.getRows(); ++i) {
        for (int j = 0; j < result.getCols(); ++j) {
            result(i, j) = (*this)(i, j) - rhs(i, j);
        }
    }

    return result;
}

// NOTE: These row operations are implemented by multiplying
//       with elementary matrices for clarity and conceptual simplicity.
//       This is NOT optimized for performance and results in higher time complexity
//       when used in matrix inversion or elimination algorithms.
//       Suitable for educational/demo purposes only.
//       It is recommended to directly implement these functions if faster time is desired.
void Matrix::swapRows(int row1, int row2) {
    (*this) = rowSwapMatrix(row1, row2, this->getRows()) * (*this);
}

void Matrix::addRows(int row1, int row2, double factor) {
    (*this) = rowAddMatrix(row1, row2, this->getRows(), factor) * (*this);
}

void Matrix::scaleRow(int row, double factor) {
    (*this) = rowMultiplyMatrix(row, this->getRows(), factor) * (*this);
}
