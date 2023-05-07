import os
import argparse
from string import ascii_uppercase


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

for i in range(NUMBER_OF_FILES):
    with open(f'csvset/{i}.csv', 'w') as f:
        uppercase = list(ascii_uppercase[:cols])
        f.write(','.join(uppercase))
        f.write('\n')
        for l in range(NUMBER_OF_LINES):
            f.write(','.join(f'{letter}{l}' for letter in uppercase))
            f.write('\n')
