# Introduction
This project is a dummy blockchain: instead of the real mining work, there's a server that solves the puzzle for each node. Each Node will collect the transactions it has received and hash them into a puzzle. The puzzle will be solved by service server and send back to it, if the service server solves. All the nodes will follow the longest chain principle to construct a tree-like blockchain where everyone has the same copy of work. 

# install requirement:
1. pull the repo.
2. install Golang and pyhton3.

# run the code:
1. rm logs/*
2. python3 mp2_service.py <port> [tx_rate] [solve_rate]
3. python aux.py <num_nodes> <vm_num> (if test locally, <vm_num> can be any int)
4. ctrl-c to kill all nodes.
5. python check_log [T/F] to get test result.