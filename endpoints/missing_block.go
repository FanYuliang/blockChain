package endpoints

type MissingBlockMeta struct{
	Type 					int // 0 request, 1 reply
	MissingTransactionID 	string
	RequesterIPaddr 		string
}