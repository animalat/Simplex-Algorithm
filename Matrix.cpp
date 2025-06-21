#include "Matrix.h"
#include <stdexcept>
#include <algorithm>
#include <cmath>
#include <iomanip>

const int DECIMAL_PRECISION = 2;

// Accessors and mutators
double Matrix::operator()(int row, int col) const {
    if (row < 0 || row >= rows_) {
        throw std::out_of_range("Matrix indices out of bounds");
    } else if (col < 0 || col >= cols_) {
        throw std::out_of_range("Matrix indices out of bounds");
    }

    return entries_[row][col];
}

double &Matrix::operator()(int row, int col) {
    if (row < 0 || row >= rows_) {
        std::cout << "Invalid index";
    } else if (col < 0 || col >= cols_) {
        std::cout << "Invalid index";
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
    if (num < 0.0f) {
        // negative sign
        ++width;
    }

    // +1 for decimal point
    return width + 1 + DECIMAL_PRECISION;
}

void Matrix::printMatrix() const {
    // get column spacing for each column
    std::vector<int> columnSpacing(this->getCols(), 0);
    for (int y = 0; y < this->getRows(); ++y) {
        for (int x = 0; x < this->getCols(); ++x) {
            columnSpacing[x] = std::max(columnSpacing[x], getLengthOfNum((*this)(y, x)));
        }
    }
    
    printTopBottom(this->getCols(), '_', '_', columnSpacing);

    for (int y = 0; y < this->getRows(); ++y) {
        std::cout << "|";
        for (int x = 0; x < this->getCols(); ++x) {
            std::cout << " " << std::setw(columnSpacing[x]) 
                      << std::fixed << std::setprecision(2) 
                      << (*this)(y, x);
        }
        std::cout << " |" << std::endl;
    }

    printTopBottom(this->getCols(), '-', '-', columnSpacing);
    return;
}

void Matrix::readMatrix() {
    for (int y = 0; y < this->getRows(); ++y) {
        std::cout << "Row " << (y + 1) << ": " << std::endl;
        for (int x = 0; x < this->getCols(); ++x) {
            double curEntry;
            if (!(std::cin >> curEntry)) {
                throw std::runtime_error("Invalid entry");
            }
            (*this)(y, x) = curEntry;
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

    Matrix result = Matrix(numRows, numCols);
    result.readMatrix();

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

