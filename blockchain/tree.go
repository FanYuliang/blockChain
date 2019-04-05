package blockchain
import(
	"errors"
	"sync")

type Tree struct{
	Sentinel 		Block
	Leaf			[]Block
	lock 			sync.RWMutex
	holdbackQueue	[]Block
}

func (t *Tree)Constructor(){

	var sentinelBlock Block
	sentinelBlock.Constructor("-1")
	t.Sentinel.ID = string(0)
	t.Sentinel = sentinelBlock
	t.Leaf = make([]Block,0)
}

//func (t *Tree)InsertRoot(b Block){
//	var bl = Block{}
//	t.Sentinel = bl
//	t.Leaf = make([]Block,0)
//}

func (t *Tree)GetLongestChainTerm()int{
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



func (t *Tree)InsertBlock(b Block)(error){
	t.lock.Lock()
	defer t.lock.Unlock()
	for i,elem := range t.Leaf{
		if elem.ID == b.PrevBlockID{
			t.Leaf[i] = b
			return nil
		}
	}
	//RequestBlock()
	return errors.New("missing block")
}


func (t* Tree)GetBlockByID(id string)(Block,error){
	visited := make(map[string]bool)
	for _,elem := range t.Leaf{
		for elem.PrevBlockID != t.Sentinel.ID {
			blockid := elem.ID
			if (visited[blockid]){
				break
			}else {
				visited[blockid] = true
				if blockid==id{
					return elem,nil
				}
			}

		}
	}
	return Block{},errors.New("No block with such id found")
}

func (t* Tree)GetPreviousBlock(id string)(Block,error){
	for i,elem := range t.Leaf{
		if elem.ID == id{
			return t.Leaf[i],nil
		}
	}
	return Block{},errors.New("No such block")
}

func (t *Tree)GetPreviousBlockId()string{
	maxterm := 0
	id := ""
	for _,elem := range t.Leaf{
		if elem.Term > maxterm{
			maxterm = elem.Term
			id = elem.ID
		}
	}
	return id
}

func (t *Tree)PushToHoldBackQueue(b Block){
	t.holdbackQueue = append(t.holdbackQueue, b)
}

func (t *Tree)PopFromHoldBackQueue()(Block,error){
	if len(t.holdbackQueue) == 0{
		return Block{},errors.New("nothing in hold back queue")
	}
	return t.holdbackQueue[0],nil
}