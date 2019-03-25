import os
import pprint
import numpy as np

total_bandwidth = 0
def process_file(filename, is_total=False):
	global total_bandwidth
	with open("logs/" + filename) as f:
		contents = f.readlines()

	contents = [content.split(' ') for content in contents]
	res = {}
	if is_total:
		for content in contents:
			res[content[0]] = float(content[1])
	else:
		for i, content in enumerate(contents):
			if content[0] == "Bandwidth":
				total_bandwidth += int(content[1])
			else:
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

print("Total bandwidth: ", total_bandwidth)

def check_if_transaction_received(transaction_id, nodes):
	for filename, node_map in nodes.items():
		if transaction_id not in node_map:
			print(transaction_id, filename, " is missing")
			return False
	return True

def check_all_transaction_received_by_all_nodes(total_transaction, nodes):
	failed_transactions = []
	for transaction_id in total_transaction.keys():
		if not check_if_transaction_received(transaction_id, nodes):
			failed_transactions.append(transaction_id)
	return failed_transactions

def calculate_propagation_speed_of_a_transaction(transaction, total_transaction, nodes):
	start_time = total_transaction[transaction]
	end_time = start_time
	intervals = []
	for node_name, node_map in nodes.items():
		intervals.append(node_map[transaction] - start_time)
	intervals.sort()
	return  intervals[0], intervals[int(len(intervals)/2)], intervals[-1]

def calculate_avg_propagation_delay(total_transaction, nodes):
	delays1 = []
	delays2 = []
	delays3 = []
	for transactionID in total_transaction.keys():
		delay3, delay1, delay2 = calculate_propagation_speed_of_a_transaction(transactionID, total_transaction, nodes)
		delays1.append(delay1)
		delays2.append(delay2)
		delays3.append(delay3)
	print("median of one transaction: ")
	pprint.pprint(delays1)

	print("max of one transaction: ")
	pprint.pprint(delays2)

	print("min of one transaction: ")
	pprint.pprint(delays3)
	return np.mean(np.array(delays1)), np.mean(np.array(delays2))

failed_transactions = check_all_transaction_received_by_all_nodes(total_transaction, nodes)


if len(failed_transactions) == 0:
	print("TEST1: Successfully broadcasted all transactions")
else:
	print("TEST1: Failed to broadcast transactions for the following ids: ")
	pprint.pprint(failed_transactions)


print("Propagation delays: ")
avg_propagation_delay1, avg_propagation_delay2 = calculate_avg_propagation_delay(total_transaction, nodes)

print("Average propagation delays to reach half of the nodes: \n", avg_propagation_delay1)
print("Average propagation delays to reach all of the nodes: \n", avg_propagation_delay2)








