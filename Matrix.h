#ifndef MATRIX_H
#define MATRIX_H

#include <iostream>
#include <vector>

class Matrix {
    public: 
        Matrix(int rows, int cols) : rows_{rows}, cols_{cols}, 
                                     entries_(rows, std::vector<double>(cols, 0.0)) {}

        double operator()(int row, int col) const;
        double &operator()(int row, int col);
        int getRows() const;
        int getCols() const;
        
        friend std::ostream& operator<<(std::ostream& os, const Matrix& matrix);
        friend std::istream& operator>>(std::istream& is, Matrix& matrix);
        
        Matrix operator*(const Matrix &rhs) const;
    private:
        int rows_, cols_;
        std::vector<std::vector<double>> entries_;
};

Matrix readAndInitializeMatrix();

#endif