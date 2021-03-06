package blockchain

import (
	"errors"
	"fmt"
	"os"
)

type Tree struct {
	blockmap      *BlockMap
	HoldbackQueue *BlockList
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
	t.HoldbackQueue = new(BlockList)
}

func (t *Tree) InsertBlock(b Block) {
	fmt.Println("------------------------")
	fmt.Println("Insert a new block: ")
	fmt.Println("Previous block id: ", b.PrevBlockID)
	prevBlock,err := t.GetBlockByID(b.PrevBlockID)
	if err != nil {
		fmt.Println("Prev block doesn't exist")
		os.Exit(15)
	}
	fmt.Println("Previous block balance: ", prevBlock.Balance)
	b.PrintContent()
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
	//fmt.Println("Longest leaf block:")

	res := t.blockmap.Get(id)
	//res.PrintContent()
	return res
}

func (t *Tree) PushToHoldBackQueue(b Block) {
	t.HoldbackQueue.Append(b)
}

func (t *Tree) FindBlockInHoldBackQueueByPuzzle(puzzle string) (Block, error) {
	for _, elem := range t.HoldbackQueue.GetAll() {
		if elem.GetPuzzle() == puzzle {
			return elem, nil
		}
	}
	return Block{}, errors.New("no block found")
}

func (t *Tree) Has(b Block) bool {
	for _, elem := range t.HoldbackQueue.GetAll() {
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
	t.HoldbackQueue.Delete(b)
}


func (t *Tree) GetCommittedTransaction(b Block) *TransactionList {
	txmap := make(map[string]int)
	ret := new(TransactionList)
	for {
		for _, elem := range b.TxList {
			if _, ok := txmap[elem.ID]; ok {
			} else {
				txmap[elem.ID] = 1
				ret.Append(elem)
			}
		}
		if b.PrevBlockID == "-1" {
			break
		}
		if t.blockmap.Has(b.PrevBlockID) {
			b = t.blockmap.Get(b.PrevBlockID)
		} else {
			break
		}

	}
	fmt.Println("=================== ")
	return ret
}

//func (t *Tree) GetBlockByPrevBlockInHoldBackQueue(id string)(Block,error){
//	for _,elem := range t.HoldbackQueue.GetAll(){
//		if id == elem.PrevBlockID{
//			return elem,nil
//		}
//	}
//	return Block{},errors.New("No satisfactory block has this prevId!")
//}

func (t *Tree) GetHoldBackQueue() [] Block {
	return t.HoldbackQueue.items
}

func (t *Tree) GetBlockInHoldBackQueueByID(blockID string) (Block,error) {
	for _,elem := range t.HoldbackQueue.GetAll(){
		if blockID == elem.ID {
			return elem,nil
		}
	}
	return Block{},errors.New("No satisfactory block has this prevId!")
}


func (t *Tree) CountSplitInChain()float64{
	res := make(map[string]int)
	for _,val := range t.blockmap.GetVals() {
		res[val.ID] = 0
	}
	for _,val:= range t.blockmap.GetVals() {
		res[val.PrevBlockID] += 1
	}
	count := 0
	for _,val := range res {
		if val > 1 {
			count += 1
		}
	}
	fmt.Println("split count = ",count," total block number: ",t.blockmap.Size())
	return float64(count)/float64(t.blockmap.Size())
}