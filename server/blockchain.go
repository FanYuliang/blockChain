package server

import (
	"fmt"
	"log"
	"mp2/blockchain"
	"mp2/endpoints"
	"mp2/utils"
	"time"
)

func (s * Server) AskServiceToSolvePuzzle() {
	time.Sleep(10 * time.Second)
	fmt.Println("Ask service to solve new puzzle")

	//prepare puzzle and current block
	s.CurrBlock = blockchain.Block{}
	//transactionToCommit := s.Transactions.GetTransactionToCommit(100)
	prevBlockID := s.BlockChain.GetPreviousBlockId()
	s.CurrBlock.Constructor(prevBlockID)

	for _, tx := range transactionToCommit {
		ok := s.CurrBlock.AddTransaction(tx)
		if ok {
			s.Transactions.SetTransaction(tx, "committed")
		} else {
			s.Transactions.SetTransaction(tx, "invalid")
		}
	}
	puzzleToSend := s.CurrBlock.GetPuzzle()
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("SOLVE ", puzzleToSend, "\n"))
	utils.CheckError(err)
}

func (s *Server) VerifyPuzzleSolution(block blockchain.Block) {
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("VERIFY ", block.ID, " ", block.Sol, "\n"))
	utils.CheckError(err)
}

func (s *Server) SolvePuzzle() {

}

func (s * Server) MergeTransactionList(receivedRequest endpoints.TransactionMeta) {
	for _, tx := range receivedRequest.Tx {
		if !s.Transactions.Has(tx.ID) {
			log.Println(tx.ID, time.Now().UnixNano())
			s.Transactions.Append(tx)
		}
	}
}


func (s *Server) SendBlock(b blockchain.Block) {
	targetIndices := s.getPingTargets()
	s.getNonFailureMembershipSize()
	for _, index := range targetIndices {

		//if s.MembershipList.List[index].LastUpdatedTime != 0 {
		//	continue
		//}
		ipAddress := s.MembershipList.List[index].IpAddress

		var endpoint endpoints.Endpoint
		endpoint.BEndpoint = s.getBlockMeta(b)
		endpoint.SetEndpointType( "Block")
		s.sendMessageWithUDP(endpoint, ipAddress)
	}

}

func (s *Server) RequestMissingBlockToNode(id string, myaddr string) {
	targetIndices := s.getPingTargets()
	s.getNonFailureMembershipSize()
	for _, index := range targetIndices {
		ipAddress := s.MembershipList.List[index].IpAddress

		var endpoint endpoints.Endpoint
		endpoint.REndpoint = s.getRequestMissingBlockMeta(id)
		endpoint.REndpoint.Type = 0 // request
		endpoint.REndpoint.RequesterIPaddr = myaddr
		endpoint.SetEndpointType("HandleMissingTransaction")
		s.sendMessageWithUDP(endpoint, ipAddress)
	}
}

func (s *Server) SendMissingBlockToNode(b blockchain.Block, ipaddr string) {

	var endpoint endpoints.Endpoint
	endpoint = s.getBlockMeta(b)
	endpoint.SetEndpointType("Block")
	s.sendMessageWithUDP(endpoint,ipaddr)

}

func (s* Server) ForwardMissingBlockToNode(missingblockid string, ipaddr string){
	targetIndices := s.getPingTargets()
	s.getNonFailureMembershipSize()
	for _, index := range targetIndices {

		if s.MembershipList.List[index].LastUpdatedTime != 0 {
			continue
		}
		ipAddress := s.MembershipList.List[index].IpAddress

		var endpoint endpoints.Endpoint
		endpoint.REndpoint = s.getRequestMissingBlockMeta(missingblockid)
		endpoint.REndpoint.Type = 0 // send
		endpoint.SetEndpointType( "Block")
		s.sendMessageWithUDP(endpoint, ipAddress)
	}
}

func (s* Server)VerifyBlock(b blockchain.Block){
	sol := b.Sol
	hash := b.ID
	_,err := fmt.Fprintf(s.ServiceConn,utils.Concatenate("VERIFY ",hash," ",sol,"\n"))
	utils.CheckError(err)
}

func (s *Server)checkBlockBalance(prevBlock blockchain.Block ,solBlock blockchain.Block)bool{
	currBalance := prevBlock.Balance
	for _,elem := range solBlock.TxList {
		amount := elem.Amount
		if elem.SNum-amount<0 {
			continue
		}else {
			currBalance[elem.SNum] -= amount
			currBalance[elem.DNum] += amount
		}
	}
	for k,_ := range currBalance {
		if currBalance[k] == currBalance[k]{
			continue
		}else{
			return false
		}
	}
	return true
}
