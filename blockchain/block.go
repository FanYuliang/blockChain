package blockchain

type Block struct {
	Term 			int
	Timestamp		int64
	TxList 			[] Transaction
	Puzzle			string
	Sol				string
}