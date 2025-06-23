#include "Utils.h"
#include <iostream>
#include <stdexcept>
#include <limits>

int convertToInt(size_t n) {
    if (n > static_cast<size_t>(std::numeric_limits<int>::max())) {
        throw std::invalid_argument("Basis size is greater than an integer");
    }

    return static_cast<int>(n);
}
