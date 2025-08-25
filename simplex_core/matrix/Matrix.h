#ifndef MATRIX_H
#define MATRIX_H

#include <iostream>
#include <vector>

class Matrix {
    public: 
        Matrix(int rows, int cols);
        Matrix() = default;

        double operator()(int row, int col) const;
        double &operator()(int row, int col);
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
        std::vector<std::vector<double>> entries_;
};

void printMatrixBasic(std::ostream &os, const Matrix &matrix);
Matrix readAndInitializeMatrix();

#endif
