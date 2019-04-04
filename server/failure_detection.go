package server

import (
	"fmt"
	"log"
	"mp2/endpoints"
	"mp2/node_membership"
	"mp2/utils"
	"time"
)


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

/*
	This function should ping to num processes. And at the same time, it should disseminate entries stored in the disseminateList
*/
func (s *Server) ping() {
	//fmt.Println("Start to ping...")
	targetIndices := s.getPingTargets()
	s.getNonFailureMembershipSize()
	//fmt.Println("membership list size: ", len(s.MembershipList.List))
	//fmt.Println("targetIndices", targetIndices)

	//blockToSend := blockchain.Block{}
	//if s.CurrBlock.IsReady {
	//	blockToSend = s.CurrBlock
	//	s.CurrBlock = blockchain.Block{}
	//}

	for _, index := range targetIndices {

		if s.MembershipList.List[index].LastUpdatedTime != 0 {
			continue
		}
		ipAddress := s.MembershipList.List[index].IpAddress

		var endpoint endpoints.Endpoint
		endpoint.TEndpoint = s.getTransactionEndpointMetadata()
		endpoint.FEndpoint = s.getFailureDetectionEndpointMetadata("Ping")
		endpoint.SetEndpointType("FailureDetection", "Transaction")
		s.sendMessageWithUDP(endpoint, ipAddress)
		s.MembershipList.List[index].LastUpdatedTime = time.Now().Unix()
	}


	var names []string
	for _, v := range s.MembershipList.List {
		names = append(names, v.Name)
	}
	//fmt.Println("server's membership list: ", names)
}



func (s *Server) getPingTargets() []int {
	selfInd := s.findSelfInMembershipList()
	tempArr := utils.Arange(selfInd, selfInd + int(len(s.MembershipList.List)/2) + 1, 1)
	var res []int
	for _, v := range tempArr {
		res = append(res, v%len(s.MembershipList.List))
	}
	return res
}

/*
	This function should reply to the ping from ipAddress, and disseminate its own disseminateList.
*/
func (s *Server) Ack(ipAddress string) {
	//fmt.Println("Sending ack")
	var endpoint endpoints.Endpoint
	endpoint.TEndpoint = s.getTransactionEndpointMetadata()
	endpoint.FEndpoint = s.getFailureDetectionEndpointMetadata("Ack")
	endpoint.SetEndpointType("FailureDetection", "Transaction")
	s.sendMessageWithUDP(endpoint, ipAddress)
}

/*
	This function invoke when it attempts to connect with the introducer node. If success, it should update its membership list
*/
func (s *Server) Join(ipAddress string) {
	//fmt.Println("Sending join request to ", introducerIPAddress)
	var endpoint endpoints.Endpoint
	endpoint.TEndpoint = s.getTransactionEndpointMetadata()
	endpoint.FEndpoint = s.getFailureDetectionEndpointMetadata("Join")
	endpoint.SetEndpointType("FailureDetection", "Transaction")
	s.sendMessageWithUDP(endpoint, ipAddress)
}

/*
	This function invoke when it quits the group
*/
func (s *Server) Quit() {
	fmt.Println("Sending Quit request")
	s.MembershipList.UpdateNode2(s.MyAddress, 2, 0)
	for _, entry := range s.MembershipList.List {
		s.MembershipList.ListMutex.Lock()
		ipAddress := entry.IpAddress
		s.MembershipList.ListMutex.Unlock()


		var endpoint endpoints.Endpoint
		endpoint.TEndpoint = s.getTransactionEndpointMetadata()
		endpoint.FEndpoint = s.getFailureDetectionEndpointMetadata("Quit")
		endpoint.SetEndpointType("FailureDetection", "Transaction")
		s.sendMessageWithUDP(endpoint, ipAddress)

	}
}

func (s *Server) MergeList(receivedRequest endpoints.FailureDetectionMeta) {
	//fmt.Println("Start to merge list...")
	for _, entry := range receivedRequest.Record {
		if entry.IpAddress != s.MyAddress {
			s.MembershipList.UpdateNode(entry)
		}
	}
}

func (s * Server)MergeTransactionList(receivedRequest endpoints.TransactionMeta) {
	for id, trans := range receivedRequest.Tx {
		if !s.Transactions.Has(id) {
			log.Println(id, time.Now().UnixNano())
			s.Transactions.Append(trans)
		}
	}
}


func (s *Server) checkMembershipList() {
	currTime := time.Now().Unix()
	//check if any process is MembershipList or failed
	for i := len(s.MembershipList.List) - 1; i >= 0; i-- {
		entry := s.MembershipList.List[i]
		if entry.EntryType == 0 && currTime-entry.LastUpdatedTime >= s.tDetection && entry.LastUpdatedTime != 0 {
			//alive now but passed detection timeout
			s.MembershipList.List[i].LastUpdatedTime = 0
			s.MembershipList.List[i].EntryType = 1
		} else if entry.EntryType == 1 && currTime-entry.LastUpdatedTime >= s.tSuspect && entry.LastUpdatedTime != 0 {
			//suspected now but passed suspected timeout
			s.MembershipList.List[i].EntryType = 2
		}
	}
}


func (s *Server) getMemebershipSubset(subsetNum int) []node_membership.Entry {
	tempArr := utils.Arange(0, len(s.MembershipList.List), 1)
	shuffledArr := utils.Shuffle(tempArr)
	var res [] node_membership.Entry
	for i, v := range shuffledArr {
		if i >= subsetNum {
			break
		}
		res = append(res, s.MembershipList.List[v])
	}
	return res
}


func (s *Server) findSelfInMembershipList() int {
	for ind, entry := range s.MembershipList.List {
		if s.MyAddress == entry.IpAddress {
			return ind
		}
	}

	fmt.Println("Fail to find self in membership list.")
	return -1
}


func (s *Server) getNonFailureMembershipSize() {
	size := 0
	for _, v := range s.MembershipList.List {
		if v.EntryType != 2 {
			size += 1
		}
	}
	fmt.Println("Non failure membership size: ", size)
}