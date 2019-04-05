package blockchain

import (
	"encoding/json"
	"math/rand"
	"mp2/utils"
	"time"
)

type Block struct {
	ID 				string
	PrevBlockID		string
	TxList 			[] Transaction
	Sol				string
	balance 		map[int]int
	term	  		int
}

func (b *Block)  Constructor(prevBlockID string)  {
	b.ID = utils.Concatenate(rand.Intn(1000000), int(time.Now().Unix()))
	b.TxList = make([] Transaction, 0)
	b.PrevBlockID = prevBlockID
}

func (b *Block)  ToBytes() []byte {
	res, _ := json.Marshal(b)
	return res
}

func (b *Block) AddTransaction(transaction Transaction) bool {
	//not support for concurrency
	sourceBalance, ok1 := b.balance[transaction.SNum]
	_, ok2 := b.balance[transaction.DNum]
	if ok1 && sourceBalance >= transaction.Amount {
		b.balance[transaction.SNum] -= transaction.Amount
		if !ok2 {
			b.balance[transaction.DNum] = 0
		}
		b.balance[transaction.DNum] += transaction.Amount
		b.TxList = append(b.TxList, transaction)
		return true
	}
	return false
}

func (b *Block) GetPuzzle() string {
	currPuzzleHolder := new(Puzzle)
	currPuzzleHolder.Constructor(b.PrevBlockID, b.TxList)
	puzzleToSend := utils.GetSHA256(currPuzzleHolder.ToBytes())
	return puzzleToSend
}