import sys


def validateArgs():
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


def main():
    n = validateArgs()
    A = [[0] * n for _ in range(n)]
    for i in range(n):
        A[i][i] = 1
    
    for i in range(1, n):
        for j in range(n - i):
            A[j + i][j] = 2 ** (i + 1)
    
    print(A)

if __name__ == '__main__':
    main()
