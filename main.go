package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mp2/config"
	"mp2/server"
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


	myServer := new(server.Server)
	myServer.Constructor(name, "",myAddr)

	fmt.Println(utils.Concatenate("Launching server ", name, " at ", myAddr))

	// listen on all interfaces


	// accept connection on port


	fmt.Println("serviceAddr: ", serviceAddr)
	fmt.Println("serverAddr", myAddr)
	targetConn, err := net.Dial("tcp", serviceAddr)
	utils.CheckError(err)

	_, err = fmt.Fprintf(targetConn, utils.Concatenate("CONNECT ", name, " 127.0.0.1 ", portNum, "\n"))
	utils.CheckError(err)


	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP:[]byte{127,0,0,1},Port:portNum,Zone:""})
	defer ServerConn.Close()

	go myServer.TalkWithServiceServer(targetConn)

	//go myServer.StartPing(1 * time.Second)

	//wait for incoming response
	buf := make([]byte, 1024)

	for {
		n, _ := ServerConn.Read(buf)
		var resultMap server.Action
		// parse resultMap to json format
		json.Unmarshal(buf[0:n], &resultMap)

		//Customize different action
		if resultMap.ActionType == 0 {
			//received join
			log.Println("Received Join from ", resultMap.IpAddress)
			log.Println("Data received:", resultMap.Record)
			log.Println("server's membership list: ", myServer.MembershipList.List)
			myServer.MergeList(resultMap)
			log.Println("After merging, server's membership list", myServer.MembershipList.List)
			myServer.Ack(resultMap.IpAddress)
		} else if resultMap.ActionType == 1 {
			//received ping
			log.Println("Received Ping from ", resultMap.IpAddress)
			log.Println("Data received:", resultMap.Record)
			log.Println("server's membership list: ", myServer.MembershipList.List)
			myServer.MergeList(resultMap)
			log.Println("After merging, server's membership list", myServer.MembershipList.List)
			myServer.Ack(resultMap.IpAddress)
		} else if resultMap.ActionType == 2 {
			//received ack
			log.Println("Received Ack from ", resultMap.IpAddress)
			log.Println("Data received:", resultMap.Record)
			log.Println("server's membership list: ", myServer.MembershipList.List)
			for _, entry := range myServer.MembershipList.List {
				if entry.InitialTimeStamp == resultMap.InitialTimeStamp && entry.IpAddress == resultMap.IpAddress {
					myServer.MembershipList.UpdateNode2(resultMap.InitialTimeStamp, resultMap.IpAddress, 0, 0)
					break
				}
			}
			myServer.MergeList(resultMap)
			log.Println("After merging, server's membership list", myServer.MembershipList.List)
		} else if resultMap.ActionType == 3 {
			log.Println("Received Leave from ", resultMap.IpAddress)
			log.Println("Data received:", resultMap.Record)
			log.Println("server's membership list: ", myServer.MembershipList.List)
			//received leave
			//s.MembershipList.RemoveNode(incomingIP)
			myServer.MergeList(resultMap)
			log.Println("After merging, server's membership list", myServer.MembershipList.List)
		}

	}

}