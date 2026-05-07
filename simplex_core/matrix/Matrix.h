#ifndef MATRIX_H
#define MATRIX_H

#include <iostream>
#include <vector>

class Matrix {
    public: 
        Matrix(int rows, int cols);
        Matrix() = default;

        inline double operator()(int row, int col) const {
            return entries_[row * cols_ + col];
        }

        inline double &operator()(int row, int col) {
            return entries_[row * cols_ + col];
        }

        double at(int row, int col) const;
        double &at(int row, int col);

        int getRows() const;
        int getCols() const;
        
        friend std::ostream &operator<<(std::ostream &os, const Matrix &matrix);
        friend std::istream &operator>>(std::istream &is, Matrix &matrix);
        
        Matrix operator*(const Matrix &rhs) const;
        Matrix operator+(const Matrix &rhs) const;
        Matrix operator-(const Matrix &rhs) const;

        void swapRows(int row1, int row2);
        void addRows(int row1, int row2, double factor = 1.0);
        void scaleRow(int row, double factor);
    private:
        int rows_, cols_;
        std::vector<double> entries_;

        void validateNegativeRowsOrCols() const;
        void validateExceedingRowsOrCols(int row, int col) const;
        void additionSubtractionCheck(const Matrix &rhs) const;
        void multiplicationCheck(const Matrix &rhs) const;

        void addRowsScalar(int row1, int row2, double factor, int start);
        #if defined(__AVX2__) && defined(__FMA__)
        void addRowsSIMD(int row1, int row2, double factor);
        #endif
};

void printMatrixBasic(std::ostream &os, const Matrix &matrix);
Matrix readAndInitializeMatrix();
Matrix readAndInitializeMatrixQuiet();

#endif
