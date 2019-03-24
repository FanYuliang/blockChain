import subprocess
import sys
import signal

if (len(sys.argv) != 3):
	print("usage: aux.py <node_num> <vm_num>")
	sys.exit(1)

node_num = sys.argv[1]
vm_num = sys.argv[2]
arr = []

def signal_handler(sig,frame):
	for i,elem in enumerate(arr):
		elem.send_signal(signal.SIGINT)
		print("kill node %d",i)
	sys.exit(1)

for i in range(int(node_num)):
	num = "node"+str(int(vm_num)*10+i)
	port = str(6000+i)
	obj = subprocess.Popen(["go","run","main.go",num,port])
	arr.append(obj)

while (1):
	signal.signal(signal.SIGINT,signal_handler)