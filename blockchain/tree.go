package blockchain

import (
	"errors"
	"fmt"
)

type Tree struct {
	blockmap      				*BlockMap
	holdbackQueue 				*BlockList
}

func (t *Tree) Constructor() {
	var sentinelBlock Block
	initBalance := make(map[int]int)
	initBalance[0] = 0
	sentinelBlock.Constructor("0", initBalance, 0)
	sentinelBlock.ID = "-1"
	t.blockmap = new(BlockMap)
	t.blockmap.Set(sentinelBlock.ID, sentinelBlock)
	//t.leaf = new(BlockMap)
	//t.leaf.Set(sentinelBlock.ID, sentinelBlock)
	t.holdbackQueue = new(BlockList)
}

func (t *Tree) InsertBlock(b Block) {
	//fmt.Println("Insert a new block: ")
	//fmt.Println("Previous block id: ", b.PrevBlockID)
	b.PrintContent()
	//t.leaf.Delete(b.PrevBlockID)
	//t.leaf.Set(b.ID, b)
	t.blockmap.Set(b.ID, b)
}

func (t *Tree) GetBlockByID(id string) (Block, error) {
	if t.blockmap.Has(id) {
		return t.blockmap.Get(id), nil
	}
	return Block{}, errors.New("Not found")
}

func (t *Tree) GetLeafBlockOfLongestChain() Block {
	maxterm := 0
	id := ""
	for _, elem := range t.blockmap.GetVals() {
		if elem.Term >= maxterm {
			maxterm = elem.Term
			id = elem.ID
		}
	}
	return t.blockmap.Get(id)
}

func (t *Tree) PushToHoldBackQueue(b Block) {
	t.holdbackQueue.Append(b)
}

func (t *Tree) FindBlockInHoldBackQueueByPuzzle(puzzle string) (Block, error) {
	for _, elem := range t.holdbackQueue.GetAll() {
		if elem.GetPuzzle() == puzzle {
			return elem, nil
		}
	}
	return Block{}, errors.New("no block found")
}

func (t *Tree) Has(b Block) bool {
	for _, elem := range t.holdbackQueue.GetAll() {
		if b.ID == elem.ID {
			return true
		}
	}
	if _, err := t.GetBlockByID(b.ID); err != nil {
		return false
	} else {
		return true
	}

}

func (t *Tree) RemoveBlockFromQueue(b Block) {
	t.holdbackQueue.Delete(b)
}


func (t *Tree) GetCommittedTransaction(b Block) *TransactionList {
	txmap := make(map[string]int)
	ret := new(TransactionList)
	for {
		for _, elem := range b.TxList {
			if _, ok := txmap[elem.ID]; ok {
				fmt.Println("repeated transaction!!!", elem.ID)
			} else {
				txmap[elem.ID] = 1
				ret.Append(elem)
			}
		}
		if b.PrevBlockID == "-1" {
			break
		}
		b = t.blockmap.Get(b.PrevBlockID)
	}
	return ret
}

func (t *Tree) GetBlockByPrevBlockInHoldBackQueue(id string)(Block,error){
	for _,elem := range t.holdbackQueue.GetAll(){
		if id == elem.PrevBlockID{
			return elem,nil
		}
	}
	return Block{},errors.New("No satisfactory block has this prevId!")
}

func (t *Tree) GetHoldBackQueue() [] Block {
	return t.holdbackQueue.items
}