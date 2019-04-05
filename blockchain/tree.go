package blockchain

import (
	"errors"
	"sync"
)

type Tree struct{
	blockmap		map[string]Block
	Leaf			[]Block
	lock 			sync.RWMutex
	holdbackQueue	[]Block

}

func (t *Tree)Constructor(){
	var sentinelBlock Block
	initBalance := make(map[int]int)
	initBalance[0] = 0
	sentinelID := "-1"
	sentinelBlock.Constructor(sentinelID, initBalance)
	t.blockmap = make(map[string]Block)
	t.blockmap[sentinelID] = sentinelBlock
	t.Leaf = make([]Block,0)
}

//func (t *Tree)InsertRoot(b Block){
//	var bl = Block{}
//	t.Sentinel = bl
//	t.Leaf = make([]Block,0)
//}

func (t *Tree) GetTermOfLongestChain()int{
	t.lock.RLock()
	defer t.lock.RUnlock()
	max := 0
	for _,elem := range(t.Leaf){
		if	elem.Term > max{
			max = elem.Term
		}
	}
	return max
}



func (t *Tree)InsertBlock(b Block){
	t.lock.Lock()
	defer t.lock.Unlock()
	t.blockmap[b.ID] = b
	for i,elem := range t.Leaf {
		if elem.ID == b.PrevBlockID {
			t.Leaf[i] = b
			return
		}
	}
	t.Leaf = append(t.Leaf, b)

}


func (t* Tree)GetBlockByID(id string)(Block,error){
	t.lock.RLock()
	defer t.lock.RUnlock()
	if val,ok := t.blockmap[id]; ok {
		return val,nil
	}
	return Block{},errors.New("No block with such id found")
}

func (t* Tree)GetPreviousBlock(id string)(Block,error){
	t.lock.RLock()
	defer t.lock.RUnlock()
	for i,elem := range t.Leaf{
		if elem.ID == id{
			return t.Leaf[i],nil
		}
	}
	return Block{},errors.New("No such block")
}

func (t *Tree)GetPreviousBlockId()string{
	t.lock.RLock()
	defer t.lock.RUnlock()
	maxterm := 0
	id := "-1"
	for _,elem := range t.Leaf{
		if elem.Term > maxterm{
			maxterm = elem.Term
			id = elem.ID
		}
	}
	return id
}

func (t* Tree)GetBalance() map[int]int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	longestLeafID := t.GetPreviousBlockId()
	return t.blockmap[longestLeafID].Balance
}

func (t *Tree)PushToHoldBackQueue(b Block){
	t.holdbackQueue = append(t.holdbackQueue, b)
}

func (t *Tree)FindBlockInHoldBackQueueByPuzzle(puzzle string)(Block,error){
	for _,elem := range(t.holdbackQueue){
		if elem.GetPuzzle() == puzzle {
			return elem,nil
		}
	}
	return Block{},errors.New("no block found")
}

func (t *Tree) Has(b Block)bool{
	for _,elem := range t.holdbackQueue{
		if b.ID == elem.ID{
			return true
		}
	}
	if _,err := t.GetBlockByID(b.ID);err!=nil {
		return false
	} else{
		return true
	}

}

func (t *Tree) RemoveBlockFromQueue(b Block){
	for i,elem := range t.holdbackQueue{
		if elem.ID == b.ID {
			t.holdbackQueue[i], t.holdbackQueue[len(t.holdbackQueue)-1] = t.holdbackQueue[len(t.holdbackQueue)-1] ,t.holdbackQueue[i]
			t.holdbackQueue = t.holdbackQueue[:len(t.holdbackQueue)-1]
			return
		}
	}
}

func (t *Tree) GetCommittedTransaction(b Block)TransactionList {
	//var txmap = map[Transaction]int
	//var ret = TransactionList{}
	//for b.PrevBlockID != "-1"{
	//	for _,elem := range b.TxList {
	//		if v,ok:=txmap[elem]; ok{
	//			fmt.Println("repeated transaction!!!")
	//		}
	//		ret.Append(elem)
	//	}
	//}
	return TransactionList{}
}