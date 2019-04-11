package blockchain

import "encoding/json"

type Puzzle struct {
	PrevRef 		string
	TxList 			[] Transaction
	randNum 		int
}


func (p *Puzzle)  Constructor(prevRef string, txList [] Transaction, num int) {
	p.PrevRef = prevRef
	p.TxList = txList
	p.randNum = num
}

func (a *Puzzle)  ToBytes() []byte {
	res, _ := json.Marshal(a)
	return res
}