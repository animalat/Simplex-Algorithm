# Compiler and flags
CXX := g++
CXXFLAGS := -std=c++17 -Wall -Wextra -I./matrix -I./simplex

# Source files
SRCS := main.cpp \
        matrix/Matrix.cpp \
        matrix/ExtraMatrixFunctions.cpp \
        simplex/Simplex.cpp

# Object files
OBJS := $(SRCS:.cpp=.o)

# Target executable
TARGET := my_program

# Default target
all: $(TARGET)

# Link
$(TARGET): $(OBJS)
	$(CXX) $(CXXFLAGS) -o $@ $^

# Compile .cpp -> .o
%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

clean:
	rm -f $(OBJS) $(TARGET)

.PHONY: all clean
