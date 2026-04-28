def use_test(test: list[int], expected: list[int]) -> bool:
    import subprocess
    res = subprocess.run(['../simplex_solver', capture_output=True, text=True, check=True])

    if (res.returncode != 0):
        print("C++ error (nonzero return code):")
        return False
    
    # logic for checking (compare to .out)

def main() -> None:
    import glob
    in_files = glob.glob('*.in')
    out_files = glob.glob('*.out')

    for file in in_files:
        with open(file) as f:
            test = [[float(num) for num in line.split()] for line in f]
            if not use_test(test, expected):
                print(f'Failed test: {file}')

if __name__ == '__main__':
    main()
