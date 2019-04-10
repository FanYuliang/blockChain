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
	"strings"
	"syscall"
	"time"
)



func main() {
	if len(os.Args) != 3 {
		fmt.Print("Usage: go run main.go [server name] [port] \n")
		return
	}

	//Parse input argument
	DEBUG := false
	name := os.Args[1]
	portNum, err := strconv.Atoi(os.Args[2])
	utils.CheckError(err)
	file, err := os.Open("config/config.json")
	utils.CheckError(err)
	decoder := json.NewDecoder(file)
	myConfig := config.Configuration{}
	err = decoder.Decode(&myConfig)
	utils.CheckError(err)

	if myConfig.ServiceIP == "127.0.0.1"{
		DEBUG = true
	}
	f := utils.SetupLog(name)
	defer f.Close()

	serviceAddr := utils.Concatenate(myConfig.ServiceIP, ":", myConfig.ServicePort)
	myAddr := utils.GetCurrentIP(DEBUG, portNum)



	//startTimestamp := time.Now().Second()
	fmt.Println(utils.Concatenate("Launching server ", name, " at ", myAddr))


	targetConn, err := net.Dial("tcp", serviceAddr)
	utils.CheckError(err)
	myAddrArr := strings.Split(myAddr, ":")
	_, err = fmt.Fprintf(targetConn, utils.Concatenate("CONNECT ", name, " ", myAddrArr[0], " ",myAddrArr[1], "\n"))
	utils.CheckError(err)

	iparr := utils.StringAddrToIntArr(myAddr)
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP:iparr,Port:portNum,Zone:""})
	defer ServerConn.Close()


	myServer := new(server.Server)
	myServer.Constructor(name, "",myAddr, targetConn)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-sigs
		//endTimestamp := time.Now().Second()

		//fmt.Println(utils.Concatenate("Ending server ", name, " at ", myServer.Bandwidth, endTimestamp-startTimestamp))
		fmt.Println("Received signal from user, about to gracefully terminate the server")
		myServer.Quit()
		//log.Printf(utils.Concatenate("Bandwidth: ", myServer.Bandwidth/float64(endTimestamp-startTimestamp)))
		log.Println(utils.Concatenate("Message received: ", myServer.MessageReceive))
		os.Exit(5)
	}()

	go myServer.ServiceServerCommunication(targetConn)
	go myServer.StartPing(time.Duration(myConfig.PingPeriod) * time.Second)
	//go myServer.AskServiceToSolvePuzzle(10 * time.Second)

	myServer.NodeInterCommunication(ServerConn)
}
