import os
from typing import Iterator


def get_files(path, extension) -> Iterator[str]:
    import glob

    for root, _, _ in os.walk(path):
        yield from glob.iglob(root + '/*.' + extension)


def get_nums(file_name: str) -> Iterator[float]:
    with open(file_name, 'r') as file:
        for line in file:
            for num in line.split():
                yield float(num)


def run_test(input_file: str, output_file: str) -> bool:
    if not os.path.isfile(input_file):
        print(f'Missing file: {input_file}')
        return False
    if not os.path.isfile(output_file):
        print(f'Missing file: {output_file}')
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
