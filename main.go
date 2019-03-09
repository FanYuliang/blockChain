package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mp2/config"
	"mp2/utils"
	"net"
	"os"
	"strconv"
)


func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage: go run main.go [server name] [port] n. \n")
		return
	}
	//Parse input argument
	name := os.Args[1]
	portNum, _ := strconv.Atoi(os.Args[2])


	file, err := os.Open("config/config.json")
	utils.CheckError(err)
	decoder := json.NewDecoder(file)
	myConfig := config.Configuration{}
	err = decoder.Decode(&myConfig)
	utils.CheckError(err)

	serviceAddr := utils.Concatenate(myConfig.ServiceIP, ":", myConfig.ServicePort)
	myAddr := utils.Concatenate("127.0.0.1", ":", portNum)



	fmt.Println(utils.Concatenate("Launching server ", name, " at ", "127.0.0.1:", portNum))

	// listen on all interfaces

	ln, _ := net.Listen("tcp", myAddr)
	// accept connection on port

	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: portNum,
		},
	}
	targetConn, _ := dialer.Dial("tcp", serviceAddr)
	dialer.DialContext()

	fmt.Fprintf(targetConn, utils.Concatenate("CONNECT ", name, "127.0.0.1 ", portNum, "\n"))


	for {
		// will listen for message to process ending in newline (\n)
		fmt.Println("herereer")
		serverConn, err := ln.Accept()
		utils.CheckError(err)

		message, _ := bufio.NewReader(serverConn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message), "\n")
	}
}