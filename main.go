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
	"os/signal"
	"strconv"
	"syscall"
	"time"
)



func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage: go run main.go [server name] [port] \n")
		return
	}
	//Parse input argument
	name := os.Args[1]
	portNum, err := strconv.Atoi(os.Args[2])
	utils.CheckError(err)
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

	targetConn, err := net.Dial("tcp", serviceAddr)
	utils.CheckError(err)

	_, err = fmt.Fprintf(targetConn, utils.Concatenate("CONNECT ", name, " 127.0.0.1 ", portNum, "\n"))
	utils.CheckError(err)


	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP:[]byte{127,0,0,1},Port:portNum,Zone:""})
	defer ServerConn.Close()


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-sigs

		fmt.Println("Received signal from user, about to gracefully terminate the server")
		myServer.Leave()
		os.Exit(1)
	}()

	go myServer.TalkWithServiceServer(targetConn)

	go myServer.StartPing(1 * time.Second)

	//wait for incoming response
	buf := make([]byte, 1024*1024)

	for {
		n, _ := ServerConn.Read(buf)
		var resultMap server.Action
		// parse resultMap to json format
		err = json.Unmarshal(buf[0:n], &resultMap)
		utils.CheckError(err)

		//log.Println("Data received:", resultMap.Record)

		//Customize different action
		if resultMap.ActionType == 0 {
			//received join
			log.Println("Received Join from ", resultMap.IpAddress)
			myServer.MergeList(resultMap)
			myServer.Ack(resultMap.IpAddress, true)
		} else if resultMap.ActionType == 1 {
			//received ping
			log.Println("Received Ping from ", resultMap.IpAddress)
			myServer.MergeList(resultMap)
			myServer.Ack(resultMap.IpAddress, false)
		} else if resultMap.ActionType == 2 {
			//received ack
			log.Println("Received Ack from ", resultMap.IpAddress)
			for _, entry := range myServer.MembershipList.List {
				if entry.InitialTimeStamp == resultMap.InitialTimeStamp && entry.IpAddress == resultMap.IpAddress {
					myServer.MembershipList.UpdateNode2(resultMap.InitialTimeStamp, resultMap.IpAddress, 0, 0)
					break
				}
			}
			myServer.MergeList(resultMap)
			//log.Println("After merging, server's membership list", myServer.MembershipList.List)
		} else if resultMap.ActionType == 3 {
			log.Println("Received Leave from ", resultMap.IpAddress)
			//received leave
			//s.MembershipList.RemoveNode(incomingIP)
			myServer.MergeList(resultMap)
		}

	}
}