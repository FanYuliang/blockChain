package blockchain

import "encoding/json"

type Puzzle struct {
	PrevRef 		string
	TxList 			[] Transaction
}


func (p *Puzzle)  Constructor(prevRef string, txList [] Transaction) {
	p.PrevRef = prevRef
	p.TxList = txList
}

func (a *Puzzle)  ToBytes() []byte {
	res, _ := json.Marshal(a)
	return res
}