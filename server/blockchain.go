package server

import (
	"fmt"
	"mp2/blockchain"
	"mp2/utils"
	"time"
)

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
