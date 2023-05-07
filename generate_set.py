import os
import argparse
from string import ascii_uppercase
from tqdm import tqdm


parser = argparse.ArgumentParser('CSV set generator')
parser.add_argument('--files', '-f', type=int, default=10)
parser.add_argument('--lines', '-l', type=int, default=10)
parser.add_argument('--columns', '-c', type=int, default=20)

args = parser.parse_args()
NUMBER_OF_FILES = args.files
NUMBER_OF_LINES = args.lines
cols = args.columns

if not os.path.exists("csvset"):
    os.makedirs("csvset")

padding_file = len(str(NUMBER_OF_FILES))
padding_line = len(str(NUMBER_OF_LINES))

for i in tqdm(range(NUMBER_OF_FILES)):
    filename = str(i).zfill(padding_file)
    with open(f'csvset/{filename}.csv', 'w') as f:
        uppercase = list(ascii_uppercase[:cols])
        f.write(','.join(uppercase))
        f.write('\n')
        for l in range(NUMBER_OF_LINES):
            l_str = str(l).zfill(padding_line)
            f.write(','.join(f'{filename}{letter}{l_str}' for letter in uppercase))
            f.write('\n')
