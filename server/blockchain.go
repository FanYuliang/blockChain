package server

import (
	"fmt"
	"log"
	"mp2/blockchain"
	"mp2/endpoints"
	"mp2/utils"
	"time"
)

func (s *Server) AskServiceToSolvePuzzle(waitTime time.Duration) {
	time.Sleep(waitTime)
	fmt.Println("Ask service to solve new puzzle")
	//prepare puzzle and current block
	s.CurrBlock = blockchain.Block{}
	transactionToCommit := s.Transactions.GetTransactionToCommit(10)
	leafBlock := s.BlockChain.GetLeafBlockOfLongestChain()
	s.CurrBlock.Constructor(leafBlock.ID, leafBlock.Balance, leafBlock.Term+1)

	for _, tx := range transactionToCommit {
		ok := s.CurrBlock.AddTransaction(tx)
		if !ok {
			s.Transactions.Delete(tx.ID)
		}
	}

	puzzleToSend := s.CurrBlock.GetPuzzle()
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("SOLVE ", puzzleToSend, "\n"))
	utils.CheckError(err)
}

func (s *Server) MergeTransactionList(receivedRequest endpoints.TransactionMeta) {
	for _, tx := range receivedRequest.Tx {
		if !s.Transactions.Has(tx.ID) {
			log.Println(tx.ID, time.Now().UnixNano())
			s.Transactions.Append(tx)
		}
	}
}

func (s *Server) SendBlock(b blockchain.Block) {
	//fmt.Println("Sending block: ", b)
	for _, index := range s.getPingTargets() {
		targetAddress := s.MembershipList.List[index].IpAddress
		var endpoint endpoints.Endpoint
		endpoint.BEndpoint = s.getBlockMeta(b)
		endpoint.SetEndpointType("Block")
		s.sendMessageWithUDP(endpoint, targetAddress)
	}
}

func (s *Server) RequestMissingBlock(missingBlockID string, requesterAddr string) {
	for _, index := range s.getPingTargets() {
		ipAddress := s.MembershipList.List[index].IpAddress
		var endpoint endpoints.Endpoint
		endpoint.RMEndpoint = s.getRequestMissingBlockMeta(missingBlockID)
		endpoint.RMEndpoint.RequesterIPaddr = requesterAddr
		endpoint.SetEndpointType("RequestMissingTransaction")
		s.sendMessageWithUDP(endpoint, ipAddress)
	}
}

func (s *Server) SendMissingBlockToNode(b blockchain.Block, ipAddr string) {
	fmt.Println("Send Missing Block To Node")
	var endpoint endpoints.Endpoint
	endpoint.BEndpoint = s.getBlockMeta(b)
	endpoint.SetEndpointType("Block")
	s.sendMessageWithUDP(endpoint, ipAddr)
}

func (s *Server) VerifyBlock(b blockchain.Block) {

	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("VERIFY ", b.GetPuzzle(), " ", b.Sol, "\n"))
	utils.CheckError(err)
}

func (s *Server) IsBlockBalanceCorrect(prevBlock blockchain.Block, solBlock blockchain.Block) bool {
	currBalance := prevBlock.Balance
	for _, elem := range solBlock.TxList {
		amount := elem.Amount
		if currBalance[elem.SNum] -amount < 0 && elem.SNum != 0{
			fmt.Println("invalid transaction!!!")
			return false
		} else {
			currBalance[elem.SNum] -= amount
			currBalance[elem.DNum] += amount
		}
	}
	for k, _ := range currBalance {
		if solBlock.Balance[k] == currBalance[k] || k == 0 {
			continue
		} else {
			return false
		}
	}
	return true
}

func (s *Server) updateTransactionCommitStatus(leafBlock blockchain.Block) {
	totalTxlist := s.BlockChain.GetCommittedTransaction(leafBlock)
	//fmt.Println("totalTxlist: ", totalTxlist)
	//fmt.Println("s.Transactions.GetTransactionList(): ", s.Transactions.GetTransactionList())
	for _, tx := range s.Transactions.GetTransactionList() {
		if totalTxlist.Has(tx.ID) {
			s.Transactions.SetTransaction(tx.ID, "committed")
		} else {
			s.Transactions.SetTransaction(tx.ID, "uncommitted")
		}
	}
}

func (s *Server)AddBlockToChainFromQueue(receivedBlock blockchain.Block){
	s.BlockChain.InsertBlock(receivedBlock)
	currBlock := receivedBlock
	for {
		if b,err := s.BlockChain.GetBlockByPrevBlockInHoldBackQueue(currBlock.ID);err==nil { // found the block, continue put next block into chain
			s.BlockChain.InsertBlock(b)
			s.BlockChain.RemoveBlockFromQueue(currBlock)
			currBlock = b // recurse forward
		}else{// can't find next block of the received block; break.
			break
		}
	}
}
