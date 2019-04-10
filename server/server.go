package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mp2/blockchain"
	"mp2/config"
	"mp2/endpoints"
	"mp2/node_membership"
	"mp2/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Name                string
	tDetection          int64
	tSuspect            int64
	tFailure            int64
	pingNum             int
	TransactionCap      int
	IntroducerIpAddress string
	MembershipList      *node_membership.Membership
	MyAddress           string
	InitialTimeStamp    int64
	Bandwidth           float64
	BandwidthLock       sync.Mutex
	CurrBlock           blockchain.Block
	Transactions        *blockchain.TransactionList
	MessageReceive      int
	ServiceConn         net.Conn
	TransactionNumPerPing int
	BlockChain          blockchain.Tree
	VerifiedBlocks		*blockchain.BlockMap
}

func (s *Server) Constructor(name string, introducerIP string, myIP string, serviceConn net.Conn) {

	file, err := os.Open("config/config.json")
	utils.CheckError(err)
	decoder := json.NewDecoder(file)
	myConfig := config.Configuration{}
	err = decoder.Decode(&myConfig)
	utils.CheckError(err)

	currTimeStamp := time.Now().Unix()
	s.MembershipList = new(node_membership.Membership)
	s.ServiceConn = serviceConn
	s.MyAddress = myIP
	s.VerifiedBlocks = new(blockchain.BlockMap)
	s.IntroducerIpAddress = introducerIP
	s.InitialTimeStamp = currTimeStamp
	s.TransactionCap = myConfig.TransacCap
	s.TransactionNumPerPing = myConfig.TransactionNumPerPing
	s.tDetection = myConfig.DetectionTimeout
	s.tSuspect = myConfig.SuspiciousTimeout
	s.Transactions = new(blockchain.TransactionList)
	s.tFailure = myConfig.FailureTimeout
	s.pingNum = myConfig.PingNum
	s.Name = name
	var entry node_membership.Entry
	entry.Name = name
	entry.LastUpdatedTime = 0
	entry.EntryType = entry.EncodeEntryType("alive")
	entry.Incarnation = 0
	entry.InitialTimeStamp = currTimeStamp
	entry.IpAddress = myIP
	s.MembershipList.AddNewNode(entry)
	s.MessageReceive = 0
	s.BlockChain.Constructor()
}

func (s *Server) NodeInterCommunication(ServerConn net.Conn) {

	buf := make([]byte, 1024*1024)

	for {
		//wait for incoming response
		n, _ := ServerConn.Read(buf)

		var endpoint endpoints.Endpoint
		// parse resultMap to json format
		err := json.Unmarshal(buf[0:n], &endpoint)
		utils.CheckError(err)

		//log.Println("Data received:", resultMap.Record)
		for _, endpointType := range endpoint.GetEndpointTypes() {
			if endpointType == "FailureDetection" {
				//Customize different action
				resultMap := endpoint.FEndpoint
				if resultMap.Type == 1 {
					//received join
					s.MergeList(resultMap)
					s.Ack(resultMap.IpAddress)
				} else if resultMap.Type == 2 {
					//received ping
					//fmt.Println("Received Ping from ", resultMap.IpAddress)
					s.MergeList(resultMap)
					s.Ack(resultMap.IpAddress)
				} else if resultMap.Type == 3 {
					//received ack
					//fmt.Println("Received Ack from ", resultMap.IpAddress)
					for _, entry := range s.MembershipList.List {
						if entry.InitialTimeStamp == resultMap.InitialTimeStamp && entry.IpAddress == resultMap.IpAddress {
							s.MembershipList.UpdateNode2(resultMap.IpAddress, 0, 0)
							break
						}
					}
					s.MergeList(resultMap)
					//log.Println("After merging, server's membership list", myServer.MembershipList.List)
				} else if resultMap.Type == 4 {
					//fmt.Println("Received Quit from ", resultMap.IpAddress)
					//received leave
					//s.MembershipList.RemoveNode(incomingIP)
					s.MergeList(resultMap)
				}
			} else if endpointType == "Transaction" {
				//fmt.Println("Received new transaction: ", )
				transactionMeta := endpoint.TEndpoint
				//fmt.Println(transactionMeta)
				s.MergeTransactionList(transactionMeta)
			} else if endpointType == "Block" {
				fmt.Println("received block endpoint")
				receivedBlock := endpoint.BEndpoint.Block
				if !s.BlockChain.Has(receivedBlock) {
					// if has replica, drop the block
					s.BlockChain.PushToHoldBackQueue(receivedBlock)
					s.VerifyBlock(receivedBlock)
					s.SendBlock(receivedBlock)
				}
			} else if endpointType == "RequestMissingTransaction" {
				fmt.Println("Received RequestMissingTransaction.")
				fmt.Println("endpoint requester: ", endpoint.RMEndpoint.RequesterIPaddr)
				item, err := s.BlockChain.GetBlockByID(endpoint.RMEndpoint.MissingBlockID)
				if err != nil { // not found, disseminate to other nodes
					for _, index := range s.getPingTargets() {
						ipAddress := s.MembershipList.List[index].IpAddress
						if ipAddress != endpoint.RMEndpoint.RequesterIPaddr {
							s.sendMessageWithUDP(endpoint, ipAddress)
						}
					}
				} else {
					s.SendMissingBlockToNode(item, endpoint.RMEndpoint.RequesterIPaddr)
				}

			}
		}
	}
}

func (s *Server) ServiceServerCommunication(serviceConn net.Conn) {
	for {
		//parse incoming service server message
		message, _ := bufio.NewReader(serviceConn).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")

		fmt.Print("Message Received:", message, "\n")

		messageArr := strings.Split(message, " ")
		messageType := messageArr[0]
		if messageType == "INTRODUCE" {
			//received a introduce message from service server
			//fmt.Print("Message Received:", message, "\n")
			serverAddr := utils.Concatenate(messageArr[2], ":", messageArr[3])
			//fmt.Println("introducer serverAddr: ", serverAddr)
			s.Join(serverAddr)
		} else if messageType == "TRANSACTION" {
			s.MessageReceive += 1
			//received a transaction message from service server
			timeStamp, err := strconv.ParseFloat(messageArr[1], 64)
			utils.CheckError(err)
			transactionID := messageArr[2]
			sNum, err := strconv.Atoi(messageArr[3])
			utils.CheckError(err)
			dNum, err := strconv.Atoi(messageArr[4])
			utils.CheckError(err)
			amount, err := strconv.Atoi(messageArr[5])
			utils.CheckError(err)
			newTransaction := new(blockchain.Transaction)
			newTransaction.Timestamp = timeStamp
			newTransaction.ID = transactionID
			newTransaction.DNum = dNum
			newTransaction.SNum = sNum
			newTransaction.Amount = amount

			s.Transactions.Append(*newTransaction)
			log.Println(transactionID, time.Now().UnixNano())
		} else if messageType == "DIE" {
			//received a DIE message from service server
			//fmt.Println("Received a DIE message from service server.")
			os.Exit(6)
		} else if messageType == "SOLVED" {
			//received a solved puzzle solution
			fmt.Println("solved puzzle!")
			puzzleInput := messageArr[1]
			puzzleSol := messageArr[2]
			fmt.Println("puzzleInput: ", puzzleInput)
			fmt.Println("puzzleSol: ", puzzleSol)
			s.CurrBlock.Sol = puzzleSol
			s.BlockChain.InsertBlock(s.CurrBlock)
			s.updateTransactionCommitStatus()
			if s.CurrBlock.Term > 1 {
				s.SendBlock(s.CurrBlock)
			}

			go s.AskServiceToSolvePuzzle(0 * time.Second)
		} else if messageType == "VERIFY" {
			fmt.Println("Verified block!")
			status := messageArr[1]
			receivedBlock, _ := s.BlockChain.FindBlockInHoldBackQueueByPuzzle(messageArr[2])
			if status == "OK" {
				s.VerifiedBlocks.Set(receivedBlock.ID, receivedBlock)
				prevBlock, err := s.BlockChain.GetBlockByID(receivedBlock.PrevBlockID)
				if err != nil {
					//missing previous block(s), asking for other nodes to resend...
					fmt.Println("Verification failure: ", err)
					s.RequestMissingBlock(receivedBlock.PrevBlockID,s.MyAddress)
				} else {
					if s.IsBlockBalanceCorrect(prevBlock,receivedBlock) {
						//success
						s.BlockChain.RemoveBlockFromQueue(receivedBlock)
						s.BlockChain.InsertBlock(receivedBlock)
						if receivedBlock.Term > s.BlockChain.GetLeafBlockOfLongestChain().Term {
							go s.AskServiceToSolvePuzzle(0 * time.Second)
						}

					} else{
						fmt.Println("Verification failure: block has incorrect balance in it")
					}
				}
				s.AddBlocksFromHoldBackQueue()
				s.updateTransactionCommitStatus()
			} else {
				// verification failed ; report
				fmt.Println("Verification failure: service server fails to verify puzzle, ", messageArr[2])
			}
		}
	}
}

func (s *Server) sendMessageWithUDP(endpoint endpoints.Endpoint, ipAddress string) {
	//fmt.Println("ipAddress: ", ipAddress)
	arr := strings.Split(ipAddress, ":")
	myPort, err := strconv.Atoi(arr[1])
	utils.CheckError(err)

	iparr := utils.StringAddrToIntArr(ipAddress)
	Conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: iparr, Port: myPort, Zone: ""})
	utils.CheckError(err)
	defer Conn.Close()

	//fmt.Println("endpoint: ", endpoint)
	n, err := Conn.Write(endpoint.ToBytes())
	s.BandwidthLock.Lock()
	s.Bandwidth += float64(int(n) / 1024)
	s.BandwidthLock.Unlock()
	utils.CheckError(err)
}
