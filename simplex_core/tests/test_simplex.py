import os

def get_files(path, extension):
    import glob

    for root, _, _ in os.walk(path):
        for file_path in glob.iglob(root + "/*." + extension):
            yield file_path

def main() -> None:
    INPUT_EXTENSION = "in"
    OUTPUT_EXTENSION = "out"

    for input_file in get_files("./", INPUT_EXTENSION):
        output_file = input_file.removesuffix(INPUT_EXTENSION) + OUTPUT_EXTENSION

if __name__ == '__main__':
    main()
