import sys


def validate_args():
    if len(sys.argv) != 2:
        print(f'Usage: python {sys.argv[0]} <arg1>')
        sys.exit(1)

    n = 0
    try:
        n = int(sys.argv[1])
        if n <= 0:
            raise ValueError
    except:
        print(f'<arg1> must be a positive integer')
        sys.exit(1)
    return n


def print_matrix(m):
    ROWS = len(m)
    COLS = len(m[0])

    print(ROWS)
    print(COLS)
    for i in range(ROWS):
        print(*(f'{m[i][j]}' for j in range(COLS)))


def main():
    n = validate_args()

    constraints_lhs = [[0] * n for _ in range(n)]
    for i in range(n):
        constraints_lhs[i][i] = 1
    for i in range(1, n):
        for j in range(n - i):
            constraints_lhs[j + i][j] = 2 ** (i + 1)
    print_matrix(constraints_lhs)

    constraints_rhs = [[5 ** i] for i in range(1, n + 1)]
    print_matrix(constraints_rhs)

    BASE = 10
    objective_func = [[BASE ** (n - i - 1) for i in range(n)]]
    print_matrix(objective_func)

    constant_term = 0
    print(constant_term)


if __name__ == '__main__':
    main()
