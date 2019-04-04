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
	"mp2/thread_safe_structures/cclist"
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
	Block               [] blockchain.Block
	CurrBlock 			blockchain.Block
	Transactions        *cclist.TransactionList
	MessageReceive      int
	ServiceConn         net.Conn
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
	s.Transactions = new(cclist.TransactionList)
	s.tFailure = myConfig.FailureTimeout
	s.pingNum = myConfig.PingNum
	s.Name = name
	var entry node_membership.Entry
	entry.Name = name
	entry.LastUpdatedTime = 0
	entry.EntryType = endpoints.EncodeEndpointType("alive")
	entry.Incarnation = 0
	entry.InitialTimeStamp = currTimeStamp
	entry.IpAddress = myIP
	s.MembershipList.AddNewNode(entry)
	s.MessageReceive = 0

	firstBlock := new(blockchain.Block)
	firstBlock.Constructor(0, [] blockchain.Transaction{},  "")
	s.Block = append(s.Block, *firstBlock)
}

func (s * Server) AskServiceToSolvePuzzle() {
	time.Sleep(10 * time.Second)
	fmt.Println("Ask service to solve new puzzle")
	if s.CurrBlock.IsReady {
		fmt.Println("This shouldn't happen: ready current block is about to be deleted.")
	}
	prevBlock := s.Block[len(s.Block)-1]
	//prepare puzzle and current block
	s.CurrBlock = blockchain.Block{}
	transactionToCommit := s.Transactions.Pop(100)
	s.CurrBlock.Constructor(prevBlock.Term + 1, transactionToCommit, "")
	prevRef := utils.Concatenate(prevBlock.Term, int(prevBlock.Timestamp))
	currPuzzleHolder := new(blockchain.Puzzle)
	currPuzzleHolder.Constructor(prevRef, s.CurrBlock.TxList)

	puzzleToSend := utils.GetSHA256(currPuzzleHolder.ToBytes())
	s.CurrBlock.Puzzle = puzzleToSend
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("SOLVE ", puzzleToSend, "\n"))
	utils.CheckError(err)
}

func (s *Server) VerifyPuzzleSolution(block blockchain.Block) {
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("VERIFY ", block.Puzzle, " ", block.Sol, "\n"))
	utils.CheckError(err)
}

func (s *Server) TalkWithServiceServer(serviceConn net.Conn) {
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
			//TODO: to add transaction through ISIS algorithm
			//s.Transactions.Set(transactionID, newTransaction)
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
			//1. add solution to the current block
			s.CurrBlock.Sol = puzzleSol

			//2. generate new puzzle
			go s.AskServiceToSolvePuzzle()
			//3. broadcast block
			s.CurrBlock.IsReady = true
		} else if messageType == "VERIFY" {
			status := messageArr[1]
			if status == "OK" {

			}
		}
	}
}

func (s *Server) StartPing(duration time.Duration) {
	for {
		time.Sleep(duration)
		s.MembershipList.ListMutex.Lock()
		s.ping()
		s.checkMembershipList()
		s.MembershipList.ListMutex.Unlock()
		fmt.Println(s.Name, " Transaction count: ", s.Transactions.Size())
	}
}



func (s *Server) SolvePuzzle() {

}

func (s *Server) getTransactSubset() [] blockchain.Transaction {
	/*
	orig := s.Transactions.GetKys()
	tempArr := utils.Arange(0, s.Transactions.Size(), 1)
	shuffledArr := utils.Shuffle(tempArr)

	res := make(map[string]blockchain.Transaction)

	for _, v := range shuffledArr {
		if len(res) > s.TransactionCap {
			break
		}
		res[orig[v]] = *s.Transactions.Get(orig[v])
	}
	return res
	*/
	var txList [] blockchain.Transaction
	return txList
}
