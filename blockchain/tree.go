package blockchain
import("sync")

type Tree struct{
	Sentinel 		Block
	Leaf			[]Block
	lock 			sync.RWMutex
	holdbackQueue	[]Block
}

func (t *Tree)Constructor(){
	var tx = make([]Transaction,0)
	var m = make(map[string]int)
	var b = Block{"", "", tx, "", 0, m}
	t.Sentinel = b
	t.Leaf = make([]Block,0)
}


func (t *Tree)InsertBlock(b Block){

	for i,elem := range(t.Leaf){
		if elem.ID == b.PreviousBlockID{
			t.Leaf[i] = b
			return
		}
	}
	RequestBlock()

}

func (t* Tree)GetPreviousBlockID()Block{
	return Block{}
}