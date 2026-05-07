CXX := g++
CXXFLAGSUNOPTIMIZED := -g -fno-omit-frame-pointer -mavx2 -mfma -O0 -std=c++23 -Wall -Wextra -I./matrix -I./simplex -I./common
CXXFLAGS := -mavx2 -mfma -O3 -std=c++23 -Wall -Wextra -I./matrix -I./simplex -I./common

SRCS := main.cpp \
		common/Utils.cpp \
        matrix/Matrix.cpp \
        matrix/ExtraMatrixFunctions.cpp \
        simplex/Simplex.cpp \
		simplex/SimplexValidation.cpp

OBJS := $(SRCS:.cpp=.o)
TARGET := simplex_solver

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CXX) $(CXXFLAGS) -o $@ $^

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

clean:
	rm -f $(OBJS) $(TARGET)

.PHONY: all clean
