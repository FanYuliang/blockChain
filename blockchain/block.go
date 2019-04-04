package blockchain

import (
	"encoding/json"
)

type Block struct {
	ID 				string
	PreviousBlockID string
	TxList 			[] Transaction
	Sol				string
	Term 	 		int
	Balance  		map[string]int
}


func (b *Block)  Constructor(term int, txList [] Transaction)  {
	b.Term = term
	b.TxList = txList
	b.Sol = ""
}

func (b *Block)  ToBytes() []byte {
	res, _ := json.Marshal(b)
	return res
}