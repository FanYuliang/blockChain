package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mp2/config"
	"mp2/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	name                string
	tDetection          int64
	tSuspect            int64
	tFailure            int64
	tLeave              int64
	pingNum             int
	TransactionCap      int
	IntroducerIpAddress string
	MembershipList      *Membership
	MyAddress           string
	InitialTimeStamp    int64
	Transactions        map[string]*Transaction
	TransactionMutex    sync.Mutex
}

func (s * Server) Constructor(name string, introducerIP string, myIP string) {

	file, err := os.Open("config/config.json")
	utils.CheckError(err)
	decoder := json.NewDecoder(file)
	myConfig := config.Configuration{}
	err = decoder.Decode(&myConfig)
	utils.CheckError(err)

	currTimeStamp := time.Now().Unix()
	s.MembershipList = new(Membership)
	s.MyAddress = myIP
	s.IntroducerIpAddress = introducerIP
	s.InitialTimeStamp = currTimeStamp
	s.TransactionCap = myConfig.TransacCap
	s.tDetection = myConfig.DetectionTimeout
	s.tSuspect = myConfig.SuspiciousTimeout
	s.Transactions = make(map[string]*Transaction)
	s.tFailure = myConfig.FailureTimeout
	s.tLeave = myConfig.LeaveTimeout
	s.pingNum = myConfig.PingNum
	var entry Entry
	entry.Name = name
	entry.lastUpdatedTime = 0
	entry.EntryType = EncodeEntryType("alive")
	entry.Incarnation = 0
	entry.InitialTimeStamp = currTimeStamp
	entry.IpAddress = myIP
	s.MembershipList.AddNewNode(entry)
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
			serverAddr := utils.Concatenate(messageArr[2], ":", messageArr[3])
			fmt.Println("introducer serverAddr: ", serverAddr)
			s.Join(serverAddr)
		} else if messageType == "TRANSACTION" {
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
			newTransaction := new(Transaction)
			newTransaction.Timestamp = timeStamp
			newTransaction.ID = transactionID
			newTransaction.DNum = dNum
			newTransaction.SNum = sNum
			newTransaction.Amount = amount
			s.TransactionMutex.Lock()
			s.Transactions[transactionID] = newTransaction
			log.Println(transactionID, time.Now().UnixNano())
			s.TransactionMutex.Unlock()
		} else if messageType == "DIE" {
			//received a DIE message from service server
			fmt.Println("Received a DIE message from service server.")
			os.Exit(2)
		}
	}
}


func (s *Server) StartPing(duration time.Duration) {
	for {
		time.Sleep(duration)
		s.ping()
		s.checkMembershipList()
		fmt.Println("Transaction count: ", len(s.Transactions))
	}
}

/*
	This function should ping to num processes. And at the same time, it should disseminate entries stored in the disseminateList
 */
func (s *Server) ping() {
	fmt.Println("Start to ping...")
	targetIndices := s.getPingTargets()
	//fmt.Println("targetIndices", targetIndices)

	for _, index := range targetIndices {
		s.MembershipList.ListMutex.Lock()
		if s.MembershipList.List[index].lastUpdatedTime != 0 {
			s.MembershipList.ListMutex.Unlock()
			continue
		}
		ipAddress := s.MembershipList.List[index].IpAddress
		s.MembershipList.ListMutex.Unlock()

		s.sendMessageWithUDP("Ping", ipAddress, false)

		s.MembershipList.ListMutex.Lock()
		s.MembershipList.List[index].lastUpdatedTime = time.Now().Unix()
		s.MembershipList.ListMutex.Unlock()
	}

	s.MembershipList.BlacklistMutex.Lock()
	fmt.Println("server's Blacklist: ", s.MembershipList.Blacklist)
	s.MembershipList.BlacklistMutex.Unlock()


	var names []string
	for _, v := range s.MembershipList.List {
		names = append(names, v.Name)
	}
	fmt.Println("server's membership list: ", names)
}

/*
	This function should reply to the ping from ipAddress, and disseminate its own disseminateList.
 */
func (s *Server) Ack(ipAddress string, sendAll bool) {
	fmt.Println("Sending ack")
	s.sendMessageWithUDP("Ack", ipAddress, sendAll)
}


/*
	This function invoke when it attempts to connect with the introducer node. If success, it should update its membership list
 */
func (s *Server) Join(introducerIPAddress string) {
	fmt.Println("Sending join request to ", introducerIPAddress)
	s.sendMessageWithUDP("Join", introducerIPAddress, false)
}

/*
	This function invoke when it quits the group
 */
func (s *Server) Quit() {
	fmt.Println("Sending QUIT request")
	targetIndices := s.getPingTargets()
	s.MembershipList.UpdateNode2(s.MyAddress, 3, 0)
	//s.MembershipList.RemoveNode(s.MyAddress, s.InitialTimeStamp)

	for _, index := range targetIndices {
		s.MembershipList.ListMutex.Lock()
		ipAddress := s.MembershipList.List[index].IpAddress
		s.MembershipList.ListMutex.Unlock()
		s.sendMessageWithUDP("QUIT", ipAddress, false)
	}
}

func (s *Server) MergeList(receivedRequest Action) {
	fmt.Println("Start to merge list...")
	for _, entry := range receivedRequest.Record {
		if entry.IpAddress != s.MyAddress {
			index := s.MembershipList.UpdateNode(entry)
			if index != -1 {
				s.MembershipList.ListMutex.Lock()
				if s.MyAddress == s.MembershipList.List[index].IpAddress {
					//only process j can increase its own incarnation number
					s.MembershipList.List[index].Incarnation += 1
					s.MembershipList.List[index].EntryType = 0
				}
				s.MembershipList.ListMutex.Unlock()
			}
		}
	}

	s.TransactionMutex.Lock()
	for id, trans := range receivedRequest.Transactions {
		_, ok := s.Transactions[id]
		if !ok {
			log.Println(id, time.Now().UnixNano())
			s.Transactions[id] = &trans
		}

	}
	s.TransactionMutex.Unlock()

	s.MembershipList.ListMutex.Lock()

	s.MembershipList.ListMutex.Unlock()
}

func (s *Server) checkMembershipList() {
	s.MembershipList.ListMutex.Lock()
	defer s.MembershipList.ListMutex.Unlock()
	currTime := time.Now().Unix()
	//check if any process is MembershipList or failed
	for i:= len(s.MembershipList.List)-1; i>=0; i-- {
		entry := s.MembershipList.List[i]
		if entry.EntryType == 0 && currTime - entry.lastUpdatedTime >= s.tDetection && entry.lastUpdatedTime != 0  {
			//alive now but passed detection timeout
			s.MembershipList.List[i].EntryType = 1
			s.MembershipList.List[i].lastUpdatedTime = 0
		} else if entry.EntryType == 1 && currTime - entry.lastUpdatedTime >= s.tSuspect && entry.lastUpdatedTime != 0 {
			//suspected now but passed suspected timeout
			s.MembershipList.List[i].EntryType = 2
			s.MembershipList.List[i].lastUpdatedTime = 0
		} else if entry.EntryType == 2 && currTime - entry.lastUpdatedTime >= s.tFailure && entry.lastUpdatedTime != 0 {
			fmt.Println("failed now but passed failure timeout")
			s.MembershipList.List = append(s.MembershipList.List[:i], s.MembershipList.List[i+1:]...)
			s.MembershipList.AddToBlacklist(entry)
		}
	}
}

func (s *Server) sendMessageWithUDP (actionType string, ipAddress string, sendAll bool) {
	fmt.Println("ipAddress: ", ipAddress)
	arr := strings.Split(ipAddress, ":")

	myPort, err := strconv.Atoi(arr[1])
	utils.CheckError(err)
	Conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP:[]byte{127,0,0,1},Port:myPort,Zone:""})
	utils.CheckError(err)
	defer Conn.Close()
	var listToSend []Entry
	s.MembershipList.ListMutex.Lock()
	for _, v := range s.MembershipList.List {
		//if v.EntryType != 2 {
		listToSend = append(listToSend, v)
		//}
	}
	s.MembershipList.ListMutex.Unlock()

	transactionToSend := s.getTransactSubset()

	action := Action{EncodeActionType(actionType), listToSend, s.InitialTimeStamp, s.MyAddress, transactionToSend}
	//fmt.Println("actionToSend: ", action)
	_, err = Conn.Write(action.ToBytes())
	utils.CheckError(err)
}


func (s *Server) getTransactSubset() map[string] Transaction {
	s.TransactionMutex.Lock()
	defer s.TransactionMutex.Unlock()
	var orig []string
	for k, _ := range s.Transactions {
		orig = append(orig, k)
	}
	tempArr := utils.Arange(0,len(s.Transactions), 1)
	shuffledArr := utils.Shuffle(tempArr)


	res := make(map[string] Transaction)

	for _, v := range shuffledArr {
		if len(res) > s.TransactionCap {
			break
		}
		res[orig[v]] = *s.Transactions[orig[v]]
	}
	return res
}

func (s *Server) getPingTargets() []int {

	s.MembershipList.ListMutex.Lock()
	tempArr := utils.Arange(0,len(s.MembershipList.List), 1)
	s.MembershipList.ListMutex.Unlock()
	shuffledArr := utils.Shuffle(tempArr)
	var res [] int

	selfInd := s.findSelfInMembershipList()
	for _, v := range shuffledArr {
		if len(res) > s.pingNum {
			break
		}
		if v != selfInd {
			res = append(res, v)
		}
	}
	return res
}

func (s *Server) findSelfInMembershipList() int {
	s.MembershipList.ListMutex.Lock()
	defer s.MembershipList.ListMutex.Unlock()
	for ind, entry := range s.MembershipList.List {
		if s.MyAddress == entry.IpAddress {
			return ind
		}
	}

	fmt.Println("Fail to find self in membership list.")
	return -1
}
