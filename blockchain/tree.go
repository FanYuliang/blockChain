package blockchain

import (
	"errors"
	"fmt"
	"sync"
)

type Tree struct{
	blockmap		*BlockMap
	Leaf			*BlockMap
	holdbackQueue	*BlockList

}

func (t *Tree)Constructor(){

	var sentinelBlock Block

	initBalance := make(map[int]int)
	initBalance[0] = 0
	sentinelBlock.Constructor("-1",initBalance)
	t.blockmap = new(BlockMap)
	t.blockmap.Set("-1",sentinelBlock)
	t.Leaf = new(BlockMap)
	t.holdbackQueue = new(BlockList)
}

//func (t *Tree)InsertRoot(b Block){
//	var bl = Block{}
//	t.Sentinel = bl
//	t.Leaf = make([]Block,0)
//}

func (t *Tree) GetTermOfLongestChain()int{
	max := 0
	for _,elem := range(t.Leaf.GetVals()){
		if	elem.Term > max{
			max = elem.Term
		}
	}
	return max
}



func (t *Tree)InsertBlock(b Block){// Add the block into leaf; set it in blockmap
	t.blockmap.Set(b.ID,b)
	t.Leaf.Delete(b.PrevBlockID)
	t.Leaf.Set(b.PrevBlockID,b)
}


func (t* Tree)GetBlockByID(id string)(Block,error){

	//if val,ok := t.blockmap[id]; ok {
	//	//	return val,nil
	//	//}
	//	//return Block{},errors.New("No block with such id found")
	if t.blockmap.Has(id) {
		return t.blockmap.Get(id),nil
	}
	return Block{},errors.New("Not found")
}

func (t* Tree)GetBlockFromLeaf(id string)(Block,error){

	if t.Leaf.Has(id) {
		return t.Leaf.Get(id),nil
	}

	return Block{},errors.New("No such block")
}

func (t *Tree)GetPreviousBlockId()string{

	maxterm := 0
	id := ""
	for _,elem := range t.Leaf.GetVals(){
		if elem.Term > maxterm{
			maxterm = elem.Term
			id = elem.ID
		}
	}
	return id
}

func (t *Tree)PushToHoldBackQueue(b Block){
	t.holdbackQueue.Append(b)
}

func (t *Tree)FindBlockInHoldBackQueueByPuzzle(puzzle string)(Block,error){
	for _,elem := range(t.holdbackQueue.GetAll()){
		if elem.GetPuzzle() == puzzle {
			return elem,nil
		}
	}
	return Block{},errors.New("no block found")
}

func (t *Tree) Has(b Block)bool{
	for _,elem := range t.holdbackQueue.GetAll(){
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
	t.holdbackQueue.Delete(b)
}


func (t *Tree)GetCommitedTransaction(b Block)[]Transaction {
	var txmap= make(map[string]int)
	var ret= make([]Transaction,2000)
	for b.PrevBlockID != "-1" {
		for _, elem := range b.TxList {
			if _, ok := txmap[elem.ID]; ok {
				fmt.Println("repeated transaction!!!", elem.ID)
			} else {
				txmap[elem.ID] = 1
				ret = append(ret, elem)
			}
		}
	}
	return ret
}

func (t *Tree)GetBlockByPrevBlockInQueue(id string)(Block,error){
	for _,elem := range t.holdbackQueue.GetAll(){
		if id == elem.PrevBlockID{
			return elem,nil
		}
	}
	return Block{},errors.New("No satisfactory block has this prevId!")
}