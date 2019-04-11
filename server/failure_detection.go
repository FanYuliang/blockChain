package server

import (
	"fmt"
	"mp2/endpoints"
	"time"
)

func (s *Server) StartPing(duration time.Duration) {
	for {
		time.Sleep(duration)
		s.ping()
		s.MembershipList.CheckMembershipList(s.tDetection, s.tSuspect)
		//fmt.Println(s.Name, " Transaction count: ", s.Transactions.Size())
		//fmt.Println(s.Name, " Uncommitted transaction count: ", s.Transactions.UncommittedSize())
	}
}

/*
	This function should ping to num processes. And at the same time, it should disseminate entries stored in the disseminateList
*/
func (s *Server) ping() {
	fmt.Println("Start to ping...")
	targetIPs := s.getPingTargets()
	fmt.Println("target indices: ", targetIPs)
	s.MembershipList.GetNonFailureMembershipSize()
	s.MembershipList.PrintContent()

	//blockToSend := blockchain.Block{}
	//if s.CurrBlock.IsReady {
	//	blockToSend = s.CurrBlock
	//	s.CurrBlock = blockchain.Block{}
	//}

	for _, ipAddress := range targetIPs {

		if s.MembershipList.GetEntryByIpAddress(ipAddress).LastUpdatedTime != 0 {
			continue
		}
		//fmt.Println("Ping ", ipAddress)
		var endpoint endpoints.Endpoint
		endpoint.TEndpoint = s.getTransactionEndpointMetadata()
		endpoint.FEndpoint = s.getFailureDetectionEndpointMetadata("Ping")
		endpoint.SetEndpointType("FailureDetection", "Transaction")
		s.sendMessageWithUDP(endpoint, ipAddress)
		s.MembershipList.UpdateLastUpdatedTimeWithIP(ipAddress, time.Now().Unix())
	}

}

func (s *Server)  getPingTargets() []string {
	currMembershiplist := s.MembershipList.Copy()
	selfInd := currMembershiplist.FindEntryByIP(s.MyAddress)
	var required [] string
	var optional [] string
	v := selfInd + 1
	for {
		curr := v % currMembershiplist.Size()
		if curr != selfInd {
			if currMembershiplist.GetEntryByIndex(curr).EntryType == 1 {
				required = append(required, currMembershiplist.GetEntryByIndex(curr).IpAddress)
			} else if currMembershiplist.GetEntryByIndex(curr).EntryType == 0 {
				optional = append(optional, currMembershiplist.GetEntryByIndex(curr).IpAddress)
			}
		} else {
			break
		}
		v += 1
	}
	for _, v := range optional {
		if len(required) <= int(currMembershiplist.Size()/2) + 1 {
			required = append(required, v)
		}
	}
	return required
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
	for _, entry := range s.MembershipList.GetAll() {
		ipAddress := entry.IpAddress

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
		s.MembershipList.UpdateNode(entry)
	}
}

