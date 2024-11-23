# generate a random file system tree
# there is files and dirs in the tree
# there is no file in the root dir
# max depth of the tree is 7

import os
import random
import string

def random_string(length):
		return ''.join(random.choice(string.ascii_lowercase) for i in range(length))

def random_file():
		return "file_"+ random_string(8)

def random_dir():
		return "dir_"+ random_string(8)

file_path = []

def random_tree(depth, path: string):
		if depth == 0:
				return []
		else:
				num = random.randint(0, 10)
				mp = {
						0: 5,
						1: 5,
						2: 5,
						3: 5,
						4: 5,
						5: 4,
						6: 3,
						7: 2,
						8: 1,
						9: 1,
						10: 1,
				}
				num = mp[num]
				for i in range(num):
						if random.randint(0, 1) == 0 and depth < 5:
								file_path.append(os.path.join(path, random_file()))
						else:
								random_tree(depth - 1, os.path.join(path, random_dir()))
		
def print_tree():
		for i in file_path:
				print(i)

def make_file(file_path):
		for i in file_path:
				os.makedirs(os.path.dirname(i), exist_ok=True)
				with open(i, 'w') as f:
						f.write(random_string(1024))

def main():
		random.seed(0)
		random_tree(5, './fs')
		print_tree()
		make_file(file_path)

if __name__ == '__main__':
		main()
