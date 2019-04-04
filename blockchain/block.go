package blockchain

import (
	"encoding/json"
	"time"
)

type Block struct {
	Term 			int
	Timestamp		int64
	TxList 			[] Transaction
	Puzzle			string
	Sol				string
	IsReady 		bool
}


func (b *Block)  Constructor(term int, txList [] Transaction, puzzle string)  {
	b.Term = term
	b.Timestamp = time.Now().Unix()
	b.TxList = txList
	b.Puzzle = puzzle
	b.Sol = ""
}

func (b *Block)  ToBytes() []byte {
	res, _ := json.Marshal(b)
	return res
}