#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
int main(int argc, char* argv[]){
	if (argc != 3){
		puts("usage: ./test <num_process> <vm_num>");
		exit(1);
	}
	for (int i=0; i< atoi(argv[1]); i++){
		char cmd[512];
		int vm_num = atoi(argv[2]);
		sprintf(cmd ,"go run main.go node%d %d &", i+vm_num, 6000+(vm_num+i)*100);
		system(cmd);
		memset(cmd,0,512);
	}
	return 0;
}