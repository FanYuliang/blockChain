package main

import (
	"bufio"
	//"bufio"
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


	// accept connection on port


	fmt.Println("serviceAddr: ", serviceAddr)
	fmt.Println("serverAddr", myAddr)
	targetConn, err := net.Dial("tcp", serviceAddr)
	utils.CheckError(err)

	_, err = fmt.Fprintf(targetConn, utils.Concatenate("CONNECT ", name, " 127.0.0.1 ", portNum, "\n"))
	utils.CheckError(err)
	//err = targetConn.Close()
	//utils.CheckError(err)

	for {
		// will listen for message to process ending in newline (\n)

		message, _ := bufio.NewReader(targetConn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message), "\n")
	}

}