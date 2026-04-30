import os
from typing import Iterator


def get_files(path, extension) -> Iterator[str]:
    import glob

    for root, _, _ in os.walk(path):
        yield from glob.iglob(root + '/*.' + extension)


def get_expected(file_name: str) -> Iterator[str]:
    with open(file_name, 'r') as file:
        yield from (expected for line in file for expected in line.split())


def get_output(file_name: str) -> Iterator[str]:
    import subprocess

    with open(file_name, 'r') as file:
        with subprocess.Popen(['../simplex_solver'], stdin=file, stdout=subprocess.PIPE, text=True) as p:
            yield from (output for line in p.stdout for output in line.split())


def compare_floats(f1: float, f2: float) -> bool:
    import sys
    return f1 - f2 <= float(sys.float_info.epsilon)


def run_test(input_file: str, output_file: str) -> bool:
    if not os.path.isfile(input_file):
        print(f'Missing file: {input_file}')
        return False
    if not os.path.isfile(output_file):
        print(f'Missing file: {output_file}')
        return False
    
    outputs = list(get_output(input_file))
    expecteds = list(get_expected(output_file))
    if len(outputs) != len(expecteds):
        print('len(outputs) != len(expecteds)')
        return False

    for out, expected in zip(get_output(input_file), get_expected(output_file)):
        can_convert_out = False
        try:
            out = float(out)
            can_convert_out = True
        except ValueError:
            pass

        can_convert_expected = False
        try:
            expected = float(expected)
            can_convert_expected = True
        except ValueError:
            pass

        if can_convert_out != can_convert_expected:
            print(f'Mismatched type: {out} vs. {expected}')
            return False
        if can_convert_expected and can_convert_out and not compare_floats(out, expected):
            print(f'Wrong output: {out} vs. {expected}')
            return False
        elif out != expected:
            print(f'Wrong output: {out} vs. {expected}')
            return False

    return True


def main() -> None:
    INPUT_EXTENSION = 'in'
    OUTPUT_EXTENSION = 'out'

    for input_file in get_files('./', INPUT_EXTENSION):
        output_file = input_file.removesuffix(INPUT_EXTENSION) + OUTPUT_EXTENSION
        if not run_test(input_file, output_file):
            print(f'Failed test: {input_file}')
            return

    print(f'All tests passed!')


if __name__ == '__main__':
    main()
