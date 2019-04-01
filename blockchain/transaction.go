package blockchain

type Transaction struct {
	Timestamp float64
	ID        string
	SNum      int
	DNum      int
	Amount    int
	isCommit  bool
}

func (t *Transaction)  IsCommitted() bool {
	return  t.isCommit
}
