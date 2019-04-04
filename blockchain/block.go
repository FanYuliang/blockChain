package blockchain

import (
	"encoding/json"
	"mp2/utils"
	"time"
)

type Block struct {
	ID 				string
	PrevBlockID		string
	TxList 			[] Transaction
	Sol				string
	balance 		map[string]int
	term	  		int
}

func (b *Block)  Constructor(txList [] Transaction)  {
	b.ID = utils.Concatenate(int(time.Now().Unix()))
	b.TxList = txList
}

func (b *Block)  ToBytes() []byte {
	res, _ := json.Marshal(b)
	return res
}