package server

type Transaction struct {
	Timestamp float64
	ID        string
	SNum      int
	DNum      int
	Amount    int
	sent	  bool
}
