import os
import pprint

def process_file(filename, is_total=False):
	with open("logs/" + filename) as f:
		contents = f.readlines()

	contents = [content.split(' ') for content in contents]
	res = {}
	if is_total:
		for content in contents:
			res[content[0]] = float(content[1])
	else:
		for content in contents:
			res[content[2]] = float(content[3])/1e9
	return res

total_transaction = {}
nodes = {}
for filename in os.listdir("logs/"):
	if filename == "total.txt":
		total_transaction = process_file(filename, True)
	else:
		node_name = filename.split('.')[0]
		nodes[filename] = process_file(filename)

pprint.pprint(total_transaction)
pprint.pprint(nodes)
