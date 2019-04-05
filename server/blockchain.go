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
	transactionToCommit := s.Transactions.GetTransactionToCommit(100)
	prevBlockID := s.BlockChain.GetPreviousBlockID()
	s.CurrBlock.Constructor(transactionToCommit, prevBlockID)

	currPuzzleHolder := new(blockchain.Puzzle)
	currPuzzleHolder.Constructor(prevBlockID, s.CurrBlock.TxList)

	puzzleToSend := utils.GetSHA256(currPuzzleHolder.ToBytes())
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("SOLVE ", puzzleToSend, "\n"))
	utils.CheckError(err)
}

func (s *Server) VerifyPuzzleSolution(block blockchain.Block) {
	_, err := fmt.Fprintf(s.ServiceConn, utils.Concatenate("VERIFY ", block.ID, " ", block.Sol, "\n"))
	utils.CheckError(err)
}

func (s *Server) SolvePuzzle() {

}

func (s * Server)MergeTransactionList(receivedRequest endpoints.TransactionMeta) {
	for _, tx := range receivedRequest.Tx {
		if !s.Transactions.Has(tx.ID) {
			log.Println(tx.ID, time.Now().UnixNano())
			s.Transactions.Append(tx)
		}
	}
}

