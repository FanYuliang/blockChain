package endpoints

type RequestMissingTransactionMeta struct{
	Type 			int // 0 receive, 1 reply
	RequesterIPaddr string
}
