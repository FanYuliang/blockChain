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



def check_if_transaction_received(transaction_id, nodes):
	for nodename, node_map in nodes.items():
		if transaction_id not in node_map:
			return False
	return True


def check_all_transaction_received_by_all_nodes(total_transaction, nodes):


	failed_transactions = []
	for transaction_id in total_transaction.keys():
		if not check_if_transaction_received(transaction_id, nodes):
			failed_transactions.append(transaction_id)
	return failed_transactions


failed_transactions = check_all_transaction_received_by_all_nodes(total_transaction, nodes)

if len(failed_transactions) == 0:
	print("TEST1: Successfully broadcasted all transactions")
else:
	print("TEST1: Failed to broadcast transactions for the following ids: ")
	pprint.pprint(failed_transactions)




