package server

import (
	"fmt"
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
		//fmt.Println(tx.ID, time.Now().UnixNano())
		s.Transactions.Append(tx)
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
	fmt.Println("to verify block ", b.GetPuzzle())
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

func (s *Server) updateTransactionCommitStatus() {
	longestLeaf := s.BlockChain.GetLeafBlockOfLongestChain()
	totalTxlist := s.BlockChain.GetCommittedTransaction(longestLeaf)
	for _, tx := range s.Transactions.GetTransactionList() {
		if totalTxlist.Has(tx.ID) {
			s.Transactions.SetTransaction(tx.ID, "committed")
		} else {
			s.Transactions.SetTransaction(tx.ID, "uncommitted")
		}
	}
}

func (s *Server) AddBlocksFromHoldBackQueue(){
	for {
		isAnyBlockInQueueAddable := false
		for _, bInQ := range s.BlockChain.GetHoldBackQueue() {
			fmt.Println("CheckIfBlockCanAddFromHoldBackQueue")
			if s.CheckIfBlockCanAddFromHoldBackQueue(bInQ) {
				s.addBlocksFromHoldBackQueue(bInQ)
				isAnyBlockInQueueAddable = true
				break
			}
		}
		if !isAnyBlockInQueueAddable {
			break
		}
	}
}

func (s *Server) CheckIfBlockCanAddFromHoldBackQueue(currBlock blockchain.Block) bool {
	currBlock.PrintContent()
	if _, err := s.BlockChain.GetBlockByID(currBlock.ID); err == nil {
		return true
	} else {
		if !s.VerifiedBlocks.Has(currBlock.ID){
			return false
		}

		if b,err := s.BlockChain.GetBlockInHoldBackQueueByID(currBlock.PrevBlockID); err == nil { // found the block, continue put next block into chain
			if s.IsBlockBalanceCorrect(b,currBlock) {
				return s.CheckIfBlockCanAddFromHoldBackQueue(b)
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

func (s *Server) addBlocksFromHoldBackQueue(currBlock blockchain.Block) {
	if b,err := s.BlockChain.GetBlockInHoldBackQueueByID(currBlock.PrevBlockID); err == nil { // found the block, continue put next block into chain
		s.BlockChain.InsertBlock(currBlock)
		s.BlockChain.RemoveBlockFromQueue(currBlock)
		s.addBlocksFromHoldBackQueue(b)
	}
}
