package blockchain

import "fmt"

type Transaction struct {
	Timestamp float64
	ID        string
	SNum      int
	DNum      int
	Amount    int
}

func (t *Transaction) PrintContent() {
	fmt.Println("==========")
	fmt.Println("ID: ", t.ID)
	fmt.Println("Timestamp: ", t.Timestamp)
	fmt.Println("==========")
}