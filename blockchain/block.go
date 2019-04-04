package blockchain

import (
	"encoding/json"
	"time"
)

type Block struct {
	Term 			int
	Timestamp		int64
	TxList 			[] Transaction
	Sol				string
}


func (b *Block)  Constructor(term int, txList [] Transaction)  {
	b.Term = term
	b.Timestamp = time.Now().Unix()
	b.TxList = txList
	b.Sol = ""
}

func (b *Block)  ToBytes() []byte {
	res, _ := json.Marshal(b)
	return res
}