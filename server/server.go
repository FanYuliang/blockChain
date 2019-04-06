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
	CurrBlock 			blockchain.Block
	Transactions        *blockchain.TransactionList
	MessageReceive      int
	ServiceConn         net.Conn
	BlockChain 			blockchain.Tree
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
	s.IntroducerIpAddress = introducerIP
	s.InitialTimeStamp = currTimeStamp
	s.TransactionCap = myConfig.TransacCap
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
					fmt.Println("Received Join from ", resultMap.IpAddress)
					fmt.Println(endpoint.GetEndpointTypes())
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
				receivedBlock := endpoint.BEndpoint.Block
				fmt.Println("Received Block: ", receivedBlock)
				if !s.BlockChain.Has(receivedBlock){ // if has replica, drop the block
					s.BlockChain.PushToHoldBackQueue(receivedBlock)
					s.VerifyBlock(receivedBlock)
				}
			}else if endpointType == "HandleMissingTransaction"{
				if endpoint.REndpoint.Type == 0{// request missing transaction
					item,err := s.BlockChain.GetBlockByID(endpoint.REndpoint.MissingTransactionID)
					 if err != nil {// not found, disseminate to other nodes
					 	 s.ForwardMissingBlockToNode(endpoint.REndpoint.MissingTransactionID,endpoint.REndpoint.RequesterIPaddr)
					 }else{
						 s.SendMissingBlockToNode(item,endpoint.REndpoint.RequesterIPaddr)
					 }
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

		//fmt.Print("Message Received:", message, "\n")

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
			puzzleInput := messageArr[1]
			puzzleSol := messageArr[2]
			fmt.Println("puzzleInput: ", puzzleInput)
			fmt.Println("puzzleSol: ", puzzleSol)
			s.CurrBlock.Sol = puzzleSol
			s.BlockChain.InsertBlock(s.CurrBlock)
			s.SendBlock(s.CurrBlock)
			go s.AskServiceToSolvePuzzle(0 * time.Second)
		} else if messageType == "VERIFY" {
			status := messageArr[1]
			receivedBlock,_ := s.BlockChain.FindBlockInHoldBackQueueByPuzzle(messageArr[2])
			if receivedBlock.Term > s.BlockChain.GetTermOfLongestChain() { // block is latest
				if status == "ok" {
					prevBlock, err := s.BlockChain.GetBlockFromLeaf(receivedBlock.PrevBlockID)
					if err != nil { //missing previous block(s), asking for other nodes to resend...
						s.BlockChain.PushToHoldBackQueue(receivedBlock)
						s.RequestMissingBlockToNode(receivedBlock.PrevBlockID,s.MyAddress)
					}else{ // find parent of received block in my blockchain
						if (s.checkBlockBalance(prevBlock,receivedBlock)){// check whether final transaction sum is correct
							s.BlockChain.InsertBlock(receivedBlock)
							s.BlockChain.RemoveBlockFromQueue(receivedBlock)
							s.CommitTransactionInLongestChain(receivedBlock)// set all transactions in longest chain as committed
						} else{
							fmt.Println("block has incorrect sum in it")
						}
					}
					go s.AskServiceToSolvePuzzle(0 * time.Second)
				} else{ // verification failed ; report
					fmt.Println("this block is failed")
				}
			}else{ // not latest;
				prevblock,err := s.BlockChain.GetBlockByID(receivedBlock.PrevBlockID)
				if err != nil {// not found
					s.BlockChain.PushToHoldBackQueue(prevblock)
					s.RequestMissingBlockToNode(prevblock.ID,s.MyAddress)
				}else{
					s.BlockChain.InsertBlock(receivedBlock)
					s.AddBlockToChainFromQueue(receivedBlock)
				}
			}
		}
	}
}


func (s *Server) sendMessageWithUDP(endpoint endpoints.Endpoint, ipAddress string) {
	//fmt.Println("ipAddress: ", ipAddress)
	arr := strings.Split(ipAddress, ":")
	fmt.Println(arr)
	myPort, err := strconv.Atoi(arr[1])
	utils.CheckError(err)

	iparr := utils.StringAddrToIntArr(ipAddress)
	Conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: iparr, Port: myPort, Zone: ""})
	utils.CheckError(err)
	defer Conn.Close()


	//fmt.Println("endpoint: ", endpoint)
	n, err := Conn.Write(endpoint.ToBytes())
	s.BandwidthLock.Lock()
	s.Bandwidth += float64(int(n)/1024)
	s.BandwidthLock.Unlock()
	utils.CheckError(err)
}

func (s *Server)CommitTransactionInLongestChain(receivedBlock blockchain.Block){
	totalTxlist := s.BlockChain.GetCommitedTransaction(receivedBlock)
	for _,elem := range totalTxlist {
		if s.Transactions.Has(elem.ID) {
			s.Transactions.SetTransaction(elem,"committed")
		}
	}
}

func (s *Server)AddBlockToChainFromQueue(receivedBlock blockchain.Block){
	s.BlockChain.InsertBlock(receivedBlock)
	for {
		if b,err := s.BlockChain.GetBlockByPrevBlockInQueue(receivedBlock.ID);err==nil {// found the block, continue put next block into chain
			s.BlockChain.InsertBlock(b)
			s.BlockChain.RemoveBlockFromQueue(receivedBlock)
		}else{// can't find next block of the received block; break.
			break
		}
	}
}
