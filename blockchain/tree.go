package blockchain
import("sync")

type Tree struct{
	Sentinel 		Block
	Leaf			[]Block
	lock 			sync.RWMutex
	holdbackQueue	[]Block
}

func (t *Tree)Constructor(){
	var sentinelBlock Block
	sentinelBlock.Constructor([]Transaction{}, "-1")
	t.Sentinel.ID = string(0)
	t.Sentinel = sentinelBlock
	t.Leaf = make([]Block,0)
}

func (t *Tree)InsertBlock(b Block){

	for i,elem := range(t.Leaf){
		if elem.ID == b.PrevBlockID{
			t.Leaf[i] = b
			return
		}
	}
	//RequestBlock()
}

func (t* Tree)GetPreviousBlockID() string{
	return ""
}