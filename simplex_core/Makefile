CXX := g++
CXXFLAGSUNOPTIMIZED := -g -fno-omit-frame-pointer -O0 -std=c++23 -Wall -Wextra -I./matrix -I./simplex -I./common
CXXFLAGS := -g -fno-omit-frame-pointer -O3 -std=c++23 -Wall -Wextra -I./matrix -I./simplex -I./common

SRCS := main.cpp \
		common/Utils.cpp \
        matrix/Matrix.cpp \
        matrix/ExtraMatrixFunctions.cpp \
        simplex/Simplex.cpp

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
